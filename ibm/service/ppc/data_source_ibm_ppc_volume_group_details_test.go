// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMPPCVolumeGroupDetailsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCVolumeGroupDetailsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_ppc_volume_group_details.testacc_volume_group_details", "id"),
				),
			},
		},
	})
}

func testAccCheckIBMPPCVolumeGroupDetailsDataSourceConfig() string {
	return fmt.Sprintf(`
data "ibm_ppc_volume_group_details" "testacc_volume_group_details" {
    ppc_volume_group_id   = "%s"
    ppc_cloud_instance_id = "%s"
}`, acc.Ppc_volume_group_id, acc.Ppc_cloud_instance_id)

}
