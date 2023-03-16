// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
)

func TestAccIBMPPCVolumebasic(t *testing.T) {
	name := fmt.Sprintf("tf-sp-volume-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCVolumeConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCVolumeExists("ibm_ppc_volume.power_volume"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_volume.power_volume", "ppc_volume_name", name),
				),
			},
			{
				Config: testAccCheckIBMPPCVolumeSizeConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCVolumeExists("ibm_ppc_volume.power_volume"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_volume.power_volume", "ppc_volume_name", name),
					resource.TestCheckResourceAttr(
						"ibm_ppc_volume.power_volume", "ppc_volume_size", "30"),
				),
			},
		},
	})
}
func testAccCheckIBMPPCVolumeDestroy(s *terraform.State) error {

	sess, err := acc.TestAccProvider.Meta().(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_ppc_volume" {
			continue
		}
		cloudInstanceID, volumeID, err := splitID(rs.Primary.ID)
		if err != nil {
			return err
		}
		volumeC := st.NewIBMPPCVolumeClient(context.Background(), sess, cloudInstanceID)
		volume, err := volumeC.Get(volumeID)
		if err == nil {
			log.Println("volume*****", volume.State)
			return fmt.Errorf("SP Volume still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}
func testAccCheckIBMPPCVolumeExists(n string) resource.TestCheckFunc {
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
		cloudInstanceID, volumeID, err := splitID(rs.Primary.ID)
		if err != nil {
			return err
		}
		client := st.NewIBMPPCVolumeClient(context.Background(), sess, cloudInstanceID)

		_, err = client.Get(volumeID)
		if err != nil {
			return err
		}
		return nil

	}
}

func testAccCheckIBMPPCVolumeConfig(name string) string {
	return fmt.Sprintf(`
	resource "ibm_ppc_volume" "power_volume"{
		ppc_volume_size       = 20
		ppc_volume_name       = "%s"
		ppc_volume_type       = "tier1"
		ppc_volume_shareable  = true
		ppc_cloud_instance_id = "%s"
	  }
	`, name, acc.Ppc_cloud_instance_id)
}

func testAccCheckIBMPPCVolumeSizeConfig(name string) string {
	return fmt.Sprintf(`
	resource "ibm_ppc_volume" "power_volume"{
		ppc_volume_size       = 30
		ppc_volume_name       = "%s"
		ppc_volume_type       = "tier1"
		ppc_volume_shareable  = true
		ppc_cloud_instance_id = "%s"
	  }
	`, name, acc.Ppc_cloud_instance_id)
}

func TestAccIBMPPCVolumePool(t *testing.T) {
	name := fmt.Sprintf("tf-sp-volume-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCVolumePoolConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCVolumeExists("ibm_ppc_volume.power_volume"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_volume.power_volume", "ppc_volume_name", name),
					resource.TestCheckResourceAttr(
						"ibm_ppc_volume.power_volume", "ppc_volume_pool", "Tier3-Flash-1"),
				),
			},
		},
	})
}

func testAccCheckIBMPPCVolumePoolConfig(name string) string {
	return fmt.Sprintf(`
	resource "ibm_ppc_volume" "power_volume"{
		ppc_volume_size       = 20
		ppc_volume_name       = "%s"
		ppc_volume_pool       = "Tier3-Flash-1"
		ppc_volume_shareable  = true
		ppc_cloud_instance_id = "%s"
	  }
	`, name, acc.Ppc_cloud_instance_id)
}

// TestAccIBMPPCVolumeGRS test the volume replication feature which is part of global replication service(GRS)
func TestAccIBMPPCVolumeGRS(t *testing.T) {
	name := fmt.Sprintf("tf-sp-volume-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCVolumeGRSConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCVolumeExists("ibm_ppc_volume.power_volume"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_volume.power_volume", "ppc_volume_name", name),
					resource.TestCheckResourceAttr(
						"ibm_ppc_volume.power_volume", "ppc_replication_enabled", "true"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_volume.power_volume", "replication_type", "global"),
				),
			},
			{
				Config: testAccCheckIBMPPCVolumeGRSUpdateConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCVolumeExists("ibm_ppc_volume.power_volume"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_volume.power_volume", "ppc_volume_name", name),
					resource.TestCheckResourceAttr(
						"ibm_ppc_volume.power_volume", "ppc_replication_enabled", "false"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_volume.power_volume", "replication_type", ""),
				),
			},
		},
	})
}

func testAccCheckIBMPPCVolumeGRSConfig(name string) string {
	return testAccCheckIBMPPCVolumeGRSBasicConfig(name, acc.Ppc_cloud_instance_id, acc.PpcStoragePool, true)
}

func testAccCheckIBMPPCVolumeGRSUpdateConfig(name string) string {
	return testAccCheckIBMPPCVolumeGRSBasicConfig(name, acc.Ppc_cloud_instance_id, acc.PpcStoragePool, false)
}

func testAccCheckIBMPPCVolumeGRSBasicConfig(name, piCloudInstanceId, piStoragePool string, replicationEnabled bool) string {
	return fmt.Sprintf(`
	resource "ibm_ppc_volume" "power_volume"{
		ppc_volume_size         = 20
		ppc_volume_name         = "%[1]s"
		ppc_volume_pool         = "%[3]s"
		ppc_volume_shareable    = true
		ppc_cloud_instance_id   = "%[2]s"
		ppc_replication_enabled = %[4]v
	  }
	`, name, piCloudInstanceId, piStoragePool, replicationEnabled)
}
