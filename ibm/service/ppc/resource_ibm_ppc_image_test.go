// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
)

func TestAccIBMPPCImagebasic(t *testing.T) {

	name := fmt.Sprintf("tf-sp-image-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCImageConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCImageExists("ibm_ppc_image.power_image"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_image.power_image", "ppc_image_name", name),
				),
			},
		},
	})
}

func testAccCheckIBMPPCImageDestroy(s *terraform.State) error {
	sess, err := acc.TestAccProvider.Meta().(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_ppc_image" {
			continue
		}
		cloudInstanceID, imageID, err := splitID(rs.Primary.ID)
		if err != nil {
			return err
		}
		imageC := st.NewIBMPPCImageClient(context.Background(), sess, cloudInstanceID)
		_, err = imageC.Get(imageID)
		if err == nil {
			return fmt.Errorf("SP Image still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckIBMPPCImageExists(n string) resource.TestCheckFunc {
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
		cloudInstanceID, imageID, err := splitID(rs.Primary.ID)
		if err != nil {
			return err
		}
		client := st.NewIBMPPCImageClient(context.Background(), sess, cloudInstanceID)

		_, err = client.Get(imageID)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckIBMPPCImageConfig(name string) string {
	return fmt.Sprintf(`
	resource "ibm_ppc_image" "power_image" {
		ppc_image_name       = "%s"
		ppc_image_id         = "IBMi-74-01-001"
		ppc_cloud_instance_id = "%s"
	  }
	`, name, acc.Ppc_cloud_instance_id)
}

func TestAccIBMPPCImageCOSPublicImport(t *testing.T) {
	imageRes := "ibm_ppc_image.cos_image"
	name := fmt.Sprintf("tf-sp-image-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCImageCOSPublicConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCImageExists(imageRes),
					resource.TestCheckResourceAttr(imageRes, "ppc_image_name", name),
					resource.TestCheckResourceAttrSet(imageRes, "image_id"),
				),
			},
		},
	})
}

func testAccCheckIBMPPCImageCOSPublicConfig(name string) string {
	return fmt.Sprintf(`
	resource "ibm_ppc_image" "cos_image" {
		ppc_image_name       = "%[1]s"
		ppc_cloud_instance_id = "%[2]s"
		ppc_image_bucket_name = "%[3]s"
		ppc_image_bucket_access = "public"
		ppc_image_bucket_region = "us-south"
		ppc_image_bucket_file_name = "%[4]s"
		ppc_image_storage_type = "tier1"
	}
	`, name, acc.Ppc_cloud_instance_id, acc.Ppc_image_bucket_name, acc.Ppc_image_bucket_file_name)
}
