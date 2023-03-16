// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIBMPPCVolumeGroupUpdate(t *testing.T) {
	name := fmt.Sprintf("tf-sp-volume-group-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCVolumeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCVolumeGroupConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCVolumeGroupExists("ibm_ppc_volume_group.power_volume_group"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_volume_group.power_volume_group", "ppc_volume_group_name", name),
				),
			},
			{
				Config: testAccCheckIBMPPCVolumeGroupUpdateConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCVolumeGroupExists("ibm_ppc_volume_group.power_volume_group"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_volume_group.power_volume_group", "ppc_volume_group_name", name),
				),
			},
			{
				Config: testAccCheckIBMPPCVolumeGroupEmptyVolumeConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCVolumeGroupExists("ibm_ppc_volume_group.power_volume_group"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_volume_group.power_volume_group", "ppc_volume_group_name", name),
				),
			},
		},
	})
}

func testAccCheckIBMPPCVolumeGroupDestroy(s *terraform.State) error {

	sess, err := acc.TestAccProvider.Meta().(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_ppc_volume_group" {
			continue
		}
		cloudInstanceID, vgID, err := splitID(rs.Primary.ID)
		if err != nil {
			return err
		}
		vgC := st.NewIBMPPCVolumeGroupClient(context.Background(), sess, cloudInstanceID)
		vg, err := vgC.Get(vgID)
		if err == nil {
			log.Println("volume-group*****", vg.Status)
			return fmt.Errorf("SP Volume Group still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}
func testAccCheckIBMPPCVolumeGroupExists(n string) resource.TestCheckFunc {
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
		cloudInstanceID, vgID, err := splitID(rs.Primary.ID)
		if err != nil {
			return err
		}
		client := st.NewIBMPPCVolumeGroupClient(context.Background(), sess, cloudInstanceID)

		_, err = client.Get(vgID)
		if err != nil {
			return err
		}
		return nil

	}
}

func testAccCheckIBMPPCVolumeGroupConfig(name string) string {
	return volumeConfig(name, acc.Ppc_cloud_instance_id) + fmt.Sprintf(`
	resource "ibm_ppc_volume_group" "power_volume_group"{
		ppc_volume_group_name       = "%[1]s"
		ppc_cloud_instance_id 	   = "%[2]s"
		ppc_volume_ids              = [ibm_ppc_volume.power_volume[0].volume_id,ibm_ppc_volume.power_volume[1].volume_id]
	  }
	`, name, acc.Ppc_cloud_instance_id)
}

func testAccCheckIBMPPCVolumeGroupUpdateConfig(name string) string {
	return volumeConfig(name, acc.Ppc_cloud_instance_id) + fmt.Sprintf(`
	resource "ibm_ppc_volume_group" "power_volume_group"{
		ppc_volume_group_name       = "%[1]s"
		ppc_cloud_instance_id 	   = "%[2]s"
		ppc_volume_ids              = [ibm_ppc_volume.power_volume[2].volume_id]
	  }
	`, name, acc.Ppc_cloud_instance_id)
}

func testAccCheckIBMPPCVolumeGroupEmptyVolumeConfig(name string) string {
	return volumeConfig(name, acc.Ppc_cloud_instance_id) + fmt.Sprintf(`
	resource "ibm_ppc_volume_group" "power_volume_group"{
		ppc_volume_group_name       = "%[1]s"
		ppc_cloud_instance_id 	   = "%[2]s"
		ppc_volume_ids              = []
	  }
	`, name, acc.Ppc_cloud_instance_id)
}

func volumeConfig(name, cloud_instance_id string) string {
	return fmt.Sprintf(`
	resource "ibm_ppc_volume" "power_volume" {
	count = 3
	ppc_volume_size         = 2
	ppc_volume_name         = "%[1]s-${count.index}"
	ppc_volume_shareable    = true
	ppc_volume_pool         = "Tier1-Flash-1"
	ppc_cloud_instance_id   = "%[2]s"
	ppc_replication_enabled = true
	 }
	`, name, cloud_instance_id)
}
