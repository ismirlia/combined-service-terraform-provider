// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
)

func TestAccIBMPPCVolumeAttachbasic(t *testing.T) {
	name := fmt.Sprintf("tf-sp-volume-attach-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCVolumeAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCVolumeAttachConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCVolumeAttachExists("ibm_ppc_volume_attach.power_attach_volume"),
					resource.TestCheckResourceAttrSet("ibm_ppc_volume_attach.power_attach_volume", "id"),
					resource.TestCheckResourceAttrSet("ibm_ppc_volume_attach.power_attach_volume", "status"),
				),
			},
		},
	})
}
func TestAccIBMPPCShareableVolumeAttachbasic(t *testing.T) {
	name := fmt.Sprintf("tf-sp-shareable-volume-attach-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCVolumeAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCShareableVolumeAttachConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCVolumeAttachExists("ibm_ppc_volume_attach.power_attach_volume"),
					resource.TestCheckResourceAttrSet("ibm_ppc_volume_attach.power_attach_volume", "id"),
					resource.TestCheckResourceAttrSet("ibm_ppc_volume_attach.power_attach_volume", "status"),
				),
			},
		},
	})
}
func testAccCheckIBMPPCVolumeAttachDestroy(s *terraform.State) error {
	sess, err := acc.TestAccProvider.Meta().(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_ppc_volume_attach" {
			continue
		}

		ids, err := flex.IdParts(rs.Primary.ID)
		if err != nil {
			return err
		}
		cloudInstanceID, pvmInstanceID, volumeID := ids[0], ids[1], ids[2]
		client := st.NewIBMPPCVolumeClient(context.Background(), sess, cloudInstanceID)
		volumeAttach, err := client.CheckVolumeAttach(pvmInstanceID, volumeID)
		if err == nil {
			log.Println("volume attach*****", volumeAttach.State)
			return fmt.Errorf("SP Volume Attach still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}
func testAccCheckIBMPPCVolumeAttachExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Record ID is set")
		}

		sess, err := acc.TestAccProvider.Meta().(conns.ClientSession).IBMPPCSession()
		if err != nil {
			return err
		}

		ids, err := flex.IdParts(rs.Primary.ID)
		if err != nil {
			return err
		}
		cloudInstanceID, pvmInstanceID, volumeID := ids[0], ids[1], ids[2]
		client := st.NewIBMPPCVolumeClient(context.Background(), sess, cloudInstanceID)

		_, err = client.CheckVolumeAttach(pvmInstanceID, volumeID)
		if err != nil {
			return err
		}
		return nil
	}
}
func testAccCheckIBMPPCVolumeAttachConfig(name string) string {
	return fmt.Sprintf(`
	  resource "ibm_ppc_volume" "power_volume" {
		ppc_volume_size       = 2
		ppc_volume_name       = "%[2]s"
		ppc_volume_shareable  = true
		ppc_volume_pool       = "Tier3-Flash-1"
		ppc_cloud_instance_id = "%[1]s"
	  }
	  resource "ibm_ppc_instance" "power_instance" {
		ppc_memory             = "2"
		ppc_processors         = "0.25"
		ppc_instance_name      = "%[2]s"
		ppc_proc_type          = "shared"
		ppc_image_id           = "%[3]s"
		ppc_sys_type           = "s922"
		ppc_cloud_instance_id  = "%[1]s"
		ppc_storage_pool       = "Tier3-Flash-1"
		ppc_network {
			network_id = "%[4]s"
		}
	  }
	  resource "ibm_ppc_volume_attach" "power_attach_volume"{
		ppc_cloud_instance_id 	= "%[1]s"
		ppc_volume_id			= ibm_ppc_volume.power_volume.volume_id
		ppc_instance_id 			= ibm_ppc_instance.power_instance.instance_id
	  }
	`, acc.Ppc_cloud_instance_id, name, acc.Ppc_image, acc.Ppc_network_name)
}

func testAccCheckIBMPPCShareableVolumeAttachConfig(name string) string {
	return fmt.Sprintf(`
	  resource "ibm_ppc_volume" "power_volume" {
		ppc_volume_size       = 2
		ppc_volume_name       = "%[2]s"
		ppc_volume_shareable  = true
		ppc_volume_pool       = "Tier3-Flash-1"
		ppc_cloud_instance_id = "%[1]s"
	  }
	  resource "ibm_ppc_instance" "power_instance" {
		count                 = 2
		ppc_memory             = "2"
		ppc_processors         = "0.25"
		ppc_instance_name      = "%[2]s-${count.index}"
		ppc_proc_type          = "shared"
		ppc_image_id           = "%[3]s"
		ppc_sys_type           = "s922"
		ppc_cloud_instance_id  = "%[1]s"
		ppc_storage_pool       = "Tier3-Flash-1"
		ppc_volume_ids         =  count.index == 0 ? [ibm_ppc_volume.power_volume.volume_id] : null
		ppc_network {
			network_id = "%[4]s"
		}
	  }
	  resource "ibm_ppc_volume_attach" "power_attach_volume"{
		ppc_cloud_instance_id 	= "%[1]s"
		ppc_volume_id 			= ibm_ppc_volume.power_volume.volume_id
		ppc_instance_id 			= ibm_ppc_instance.power_instance[1].instance_id
	  }
	`, acc.Ppc_cloud_instance_id, name, acc.Ppc_image, acc.Ppc_network_name)
}
