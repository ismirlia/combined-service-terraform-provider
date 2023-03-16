// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMPPCImageEport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCImageExportConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("ibm_ppc_image_export.power_image_export", "id"),
				),
			},
		},
	})
}

func testAccCheckIBMPPCImageExportConfig() string {
	return fmt.Sprintf(`
	data "ibm_ppc_image" "power_image" {
		ppc_image_name        = "%[6]s"
		ppc_cloud_instance_id = "%[1]s"
	  }
	resource "ibm_ppc_image_export" "power_image_export" {
		ppc_image_id         = data.ibm_ppc_image.power_image.id
		ppc_cloud_instance_id = "%[1]s"
		ppc_image_bucket_name = "%[2]s"
		ppc_image_access_key = "%[3]s"
		ppc_image_secret_key = "%[4]s"
		ppc_image_bucket_region = "%[5]s"
	  }
	`, acc.Ppc_cloud_instance_id, acc.Ppc_image_bucket_name, acc.Ppc_image_bucket_access_key, acc.Ppc_image_bucket_secret_key, acc.Ppc_image_bucket_region, acc.Ppc_image)
}
