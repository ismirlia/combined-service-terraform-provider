// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIBMPPCCaptureBasic(t *testing.T) {
	captureRes := "ibm_ppc_capture.capture_instance"
	name := fmt.Sprintf("tf-sp-capture-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCCaptureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCCaptureConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCCaptureExists(captureRes),
					resource.TestCheckResourceAttr(captureRes, "ppc_capture_name", name),
					resource.TestCheckResourceAttrSet(captureRes, "image_id"),
				),
			},
		},
	})
}
func TestAccIBMPPCCaptureWithVolume(t *testing.T) {
	captureRes := "ibm_ppc_capture.capture_instance"
	name := fmt.Sprintf("tf-sp-capture-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCCaptureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCCaptureWithVolumeConfig(name, helpers.PPCInstanceHealthOk),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCCaptureExists(captureRes),
					resource.TestCheckResourceAttr(captureRes, "ppc_capture_name", name),
					resource.TestCheckResourceAttrSet(captureRes, "image_id"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccIBMPPCCaptureCloudStorage(t *testing.T) {
	captureRes := "ibm_ppc_capture.capture_instance"
	name := fmt.Sprintf("tf-sp-capture-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCCaptureCloudStorageConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(captureRes, "ppc_capture_name", name),
				),
			},
		},
	})
}

func TestAccIBMPPCCaptureBoth(t *testing.T) {
	captureRes := "ibm_ppc_capture.capture_instance"
	name := fmt.Sprintf("tf-sp-capture-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCCaptureBothConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(captureRes, "ppc_capture_name", name),
					resource.TestCheckResourceAttrSet(captureRes, "image_id"),
				),
			},
		},
	})
}

func testAccCheckIBMPPCCaptureExists(n string) resource.TestCheckFunc {
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
		parts, err := flex.IdParts(rs.Primary.ID)
		cloudInstanceID := parts[0]
		captureID := parts[1]
		if err != nil {
			return err
		}
		client := st.NewIBMPPCImageClient(context.Background(), sess, cloudInstanceID)

		_, err = client.Get(captureID)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckIBMPPCCaptureDestroy(s *terraform.State) error {
	sess, err := acc.TestAccProvider.Meta().(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_ppc_capture" {
			continue
		}
		parts, err := flex.IdParts(rs.Primary.ID)
		cloudInstanceID := parts[0]
		captureID := parts[1]
		if err != nil {
			return err
		}
		imageClient := st.NewIBMPPCImageClient(context.Background(), sess, cloudInstanceID)
		_, err = imageClient.Get(captureID)
		if err == nil {
			return fmt.Errorf("PPC Image still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckIBMPPCCaptureWithVolumeConfig(name string, healthStatus string) string {
	return testAccCheckIBMPPCInstanceConfig(name, healthStatus) + fmt.Sprintf(`
	resource "ibm_ppc_capture" "capture_instance" {
		depends_on=[ibm_ppc_instance.power_instance]
		ppc_cloud_instance_id="%[1]s"
		ppc_capture_name  = "%[2]s"
		ppc_instance_name = ibm_ppc_instance.power_instance.ppc_instance_name
		ppc_capture_destination = "image-catalog"
		ppc_capture_volume_ids = [ibm_ppc_volume.power_volume.volume_id]
	}
	`, acc.Ppc_cloud_instance_id, name)
}

func testAccCheckIBMPPCCaptureConfigBasic(name string) string {
	return fmt.Sprintf(`
	resource "ibm_ppc_capture" "capture_instance" {
		ppc_cloud_instance_id="%[1]s"
		ppc_capture_name = "%s"
		ppc_instance_name = "%s"
		ppc_capture_destination = "image-catalog"
	}
	`, acc.Ppc_cloud_instance_id, name, acc.Ppc_instance_name)
}

func testAccCheckIBMPPCCaptureCloudStorageConfig(name string) string {
	return fmt.Sprintf(`
	resource "ibm_ppc_capture" "capture_instance" {
		ppc_cloud_instance_id="%[1]s"
		ppc_capture_name  = "%s"
		ppc_instance_name = "%s"
		ppc_capture_destination = "cloud-storage"
		ppc_capture_cloud_storage_region = "us-east"
		ppc_capture_cloud_storage_access_key = "%s"
		ppc_capture_cloud_storage_secret_key = "%s"
		ppc_capture_storage_image_path = "%s"
	}
	`, acc.Ppc_cloud_instance_id, name, acc.Ppc_instance_name, acc.Ppc_capture_cloud_storage_access_key, acc.Ppc_capture_cloud_storage_secret_key, acc.Ppc_capture_storage_image_path)
}

func testAccCheckIBMPPCCaptureBothConfig(name string) string {
	return fmt.Sprintf(`
	resource "ibm_ppc_capture" "capture_instance" {
		ppc_cloud_instance_id="%[1]s"
		ppc_capture_name = "%s"
		ppc_instance_name = "%s"
		ppc_capture_destination  = "both"
		ppc_capture_cloud_storage_region = "us-east"
		ppc_capture_cloud_storage_access_key = "%s"
		ppc_capture_cloud_storage_secret_key = "%s"
		ppc_capture_storage_image_path = "%s"
	}
	`, acc.Ppc_cloud_instance_id, name, acc.Ppc_instance_name, acc.Ppc_capture_cloud_storage_access_key, acc.Ppc_capture_cloud_storage_secret_key, acc.Ppc_capture_storage_image_path)
}
