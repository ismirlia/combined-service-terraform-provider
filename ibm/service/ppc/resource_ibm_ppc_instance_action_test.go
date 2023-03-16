// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMPPCInstanceAction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCInstanceActionConfig("stop"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_ppc_instance_action.example", "status", "SHUTOFF"),
				),
			},
			{
				Config: testAccCheckIBMPPCInstanceActionConfig("start"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_ppc_instance_action.example", "status", "ACTIVE"),
				),
			},
		},
	})
}

func TestAccIBMPPCInstanceActionHardReboot(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCInstanceActionWithHealthStatusConfig("hard-reboot"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_ppc_instance_action.example", "status", "ACTIVE"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_instance_action.example", "health_status", "WARNING"),
				),
			},
		},
	})
}

func TestAccIBMPPCInstanceActionResetState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCInstanceActionWithHealthStatusConfig("reset-state"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_ppc_instance_action.example", "status", "ACTIVE"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_instance_action.example", "health_status", "CRITICAL"),
				),
			},
		},
	})
}

func testAccCheckIBMPPCInstanceActionConfig(action string) string {
	return fmt.Sprintf(`
	resource "ibm_ppc_instance_action" "example" {
		ppc_cloud_instance_id	= "%s"
		ppc_instance_id			= "%s"
		ppc_action				= "%s"
	}
	`, acc.Ppc_cloud_instance_id, acc.Ppc_instance_name, action)
}

func testAccCheckIBMPPCInstanceActionWithHealthStatusConfig(action string) string {
	return fmt.Sprintf(`
	resource "ibm_ppc_instance_action" "example" {
		ppc_cloud_instance_id	= "%s"
		ppc_instance_id			= "%s"
		ppc_action				= "%s"
		ppc_health_status		= "WARNING"
	}
	`, acc.Ppc_cloud_instance_id, acc.Ppc_instance_name, action)
}
