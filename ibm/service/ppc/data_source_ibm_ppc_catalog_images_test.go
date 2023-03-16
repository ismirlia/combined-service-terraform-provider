// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccCheckIBMPPCCatalogImagesDataSourceBasicConfig() string {
	return fmt.Sprintf(`
	data "ibm_ppc_catalog_images" "power_catalog_images_basic" {
		ppc_cloud_instance_id = "%s"
	}
	`, acc.Ppc_cloud_instance_id)
}

func testAccCheckIBMPPCCatalogImagesDataSourceSAPConfig() string {
	return fmt.Sprintf(`
	data "ibm_ppc_catalog_images" "power_catalog_images_sap" {
		ppc_cloud_instance_id = "%s"
		sap = "true"
	}
	`, acc.Ppc_cloud_instance_id)
}

func testAccCheckIBMPPCCatalogImagesDataSourceVTLConfig() string {
	return fmt.Sprintf(`
	data "ibm_ppc_catalog_images" "power_catalog_images_vtl" {
		ppc_cloud_instance_id = "%s"
		vtl = "true"
	}
	`, acc.Ppc_cloud_instance_id)
}

func testAccCheckIBMPPCCatalogImagesDataSourceSAP_And_VTLConfig() string {
	return fmt.Sprintf(`
	data "ibm_ppc_catalog_images" "power_catalog_images_sap_and_vtl" {
		ppc_cloud_instance_id = "%s"
		sap = "true"
		vtl = "true"
	}
	`, acc.Ppc_cloud_instance_id)
}

func TestAccIBMPPCCatalogImagesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCCatalogImagesDataSourceBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_ppc_catalog_images.power_catalog_images_basic", "id"),
				),
			},
		},
	})
}

func TestAccIBMPPCCatalogImagesDataSourceSAP(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCCatalogImagesDataSourceSAPConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_ppc_catalog_images.power_catalog_images_sap", "id"),
				),
			},
		},
	})
}

func TestAccIBMPPCCatalogImagesDataSourceVTL(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCCatalogImagesDataSourceVTLConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_ppc_catalog_images.power_catalog_images_vtl", "id"),
				),
			},
		},
	})
}

func TestAccIBMPPCCatalogImagesDataSourceSAPAndVTL(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCCatalogImagesDataSourceSAP_And_VTLConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_ppc_catalog_images.power_catalog_images_sap_and_vtl", "id"),
				),
			},
		},
	})
}
