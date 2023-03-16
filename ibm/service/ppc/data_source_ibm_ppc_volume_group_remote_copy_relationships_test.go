// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMPPCVolumeGroupRemoteCopyRelationshipsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCVolumeGroupRemoteCopyRelationshipsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_ppc_volume_group_remote_copy_relationships.testacc_volume_group_remote_copy_relationships", "id"),
				),
			},
		},
	})
}

func testAccCheckIBMPPCVolumeGroupRemoteCopyRelationshipsDataSourceConfig() string {
	return fmt.Sprintf(`
data "ibm_ppc_volume_group_remote_copy_relationships" "testacc_volume_group_remote_copy_relationships" {
    ppc_cloud_instance_id = "%s"
    ppc_volume_group_id   = "%s"
}`, acc.Ppc_cloud_instance_id, acc.Ppc_volume_group_id)

}
