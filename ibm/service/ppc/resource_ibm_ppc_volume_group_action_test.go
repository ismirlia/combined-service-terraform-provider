// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
)

func TestAccIBMPPCVolumeGroupActionbasic(t *testing.T) {
	name := fmt.Sprintf("tf-sp-volume-group-action-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCVolumeGroupStopActionConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCVolumeGroupActionExists("ibm_ppc_volume_group_action.power_volume_group_action"),
					resource.TestCheckResourceAttrSet("ibm_ppc_volume_group_action.power_volume_group_action", "id"),
					resource.TestCheckResourceAttrSet("ibm_ppc_volume_group_action.power_volume_group_action", "volume_group_status"),
				),
			},
			{
				Config: testAccCheckIBMPPCVolumeGroupStartActionConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCVolumeGroupActionExists("ibm_ppc_volume_group_action.power_volume_group_action"),
					resource.TestCheckResourceAttrSet("ibm_ppc_volume_group_action.power_volume_group_action", "id"),
					resource.TestCheckResourceAttrSet("ibm_ppc_volume_group_action.power_volume_group_action", "volume_group_status"),
				),
			},
		},
	})
}

func testAccCheckIBMPPCVolumeGroupActionExists(n string) resource.TestCheckFunc {
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

		ids, err := flex.IdParts(rs.Primary.ID)
		if err != nil {
			return err
		}
		cloudInstanceID, vgID := ids[0], ids[1]
		client := st.NewIBMPPCVolumeGroupClient(context.Background(), sess, cloudInstanceID)

		_, err = client.Get(vgID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckIBMPPCVolumeGroupStopActionConfig(name string) string {
	return testAccCheckIBMPPCVolumeGroupConfig(name) + fmt.Sprintf(`
	  resource "ibm_ppc_volume_group_action" "power_volume_group_action" {
		ppc_cloud_instance_id   = "%[1]s"
		ppc_volume_group_id     = ibm_ppc_volume_group.power_volume_group.volume_group_id
		ppc_volume_group_action {
			stop {
				access = true
			}
		}
	  }
	`, acc.Ppc_cloud_instance_id)
}

func testAccCheckIBMPPCVolumeGroupStartActionConfig(name string) string {
	return testAccCheckIBMPPCVolumeGroupConfig(name) + fmt.Sprintf(`
	  resource "ibm_ppc_volume_group_action" "power_volume_group_action" {
		ppc_cloud_instance_id   = "%[1]s"
		ppc_volume_group_id     = ibm_ppc_volume_group.power_volume_group.volume_group_id
		ppc_volume_group_action {
			start {
				source = "master"
			}
		}
	  }
	`, acc.Ppc_cloud_instance_id)
}
