// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMPPCStorageTypesCapacityDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCStorageTypesCapacityDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_ppc_storage_types_capacity.types", "id"),
				),
			},
		},
	})
}

func testAccCheckIBMPPCStorageTypesCapacityDataSourceConfig() string {
	return fmt.Sprintf(`
		data "ibm_ppc_storage_types_capacity" "types" {
			ppc_cloud_instance_id = "%s"
		}
	`, acc.Ppc_cloud_instance_id)
}
