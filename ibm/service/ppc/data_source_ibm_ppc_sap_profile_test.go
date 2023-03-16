// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMPPCSAPProfileDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCSAPProfileDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_ppc_sap_profile.test", "id"),
					resource.TestCheckResourceAttr("data.ibm_ppc_sap_profile.test", "id", acc.PiSAPProfileID),
				),
			},
		},
	})
}

func testAccCheckIBMPPCSAPProfileDataSourceConfig() string {
	return fmt.Sprintf(`
		data "ibm_ppc_sap_profile" "test" {
			ppc_cloud_instance_id = "%s"
			ppc_sap_profile_id = "%s"
		}`, acc.Ppc_cloud_instance_id, acc.PiSAPProfileID)
}
