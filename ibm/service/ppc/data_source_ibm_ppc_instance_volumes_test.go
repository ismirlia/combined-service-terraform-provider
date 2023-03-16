// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMPPCVolumesDataSource_basic(t *testing.T) {
	name := fmt.Sprintf("tf-pi-volume-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCVolumesDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_ppc_instance_volumes.testacc_ds_volumes", "id"),
				),
			},
		},
	})
}

func testAccCheckIBMPPCVolumesDataSourceConfig(name string) string {
	return fmt.Sprintf(`
data "ibm_ppc_instance_volumes" "testacc_ds_volumes" {
    ppc_instance_name = "%s"
    ppc_cloud_instance_id = "%s"
}`, acc.Ppc_instance_name, acc.Ppc_cloud_instance_id)

}
