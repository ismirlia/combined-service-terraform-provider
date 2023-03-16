// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMPPCSnapshotsDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCSnapshotsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_ppc_instance_snapshots.testacc_ds_snapshots", "id"),
				),
			},
		},
	})
}

func testAccCheckIBMPPCSnapshotsDataSourceConfig() string {
	return fmt.Sprintf(`
	
data "ibm_ppc_instance_snapshots" "testacc_ds_snapshots" {
    ppc_cloud_instance_id = "%s"
}`, acc.Ppc_cloud_instance_id)

}
