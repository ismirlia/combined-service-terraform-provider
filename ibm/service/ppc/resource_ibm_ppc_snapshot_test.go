// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
)

func TestAccIBMPPCInstanceSnapshotbasic(t *testing.T) {

	name := fmt.Sprintf("tf-sp-instance-snapshot-%d", acctest.RandIntRange(10, 100))
	snapshotRes := "ibm_ppc_snapshot.power_snapshot"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCInstanceSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCInstanceSnapshotConfig(name, helpers.PPCInstanceHealthOk),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCInstanceSnapshotExists(snapshotRes),
					resource.TestCheckResourceAttr(snapshotRes, "ppc_snap_shot_name", name),
					resource.TestCheckResourceAttr(snapshotRes, "status", "available"),
					resource.TestCheckResourceAttrSet(snapshotRes, "id"),
				),
			},
		},
	})
}
func testAccCheckIBMPPCInstanceSnapshotDestroy(s *terraform.State) error {

	sess, err := acc.TestAccProvider.Meta().(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_ppc_snapshot" {
			continue
		}
		cloudInstanceID, snapshotID, err := splitID(rs.Primary.ID)
		if err != nil {
			return err
		}
		snapshotC := st.NewIBMPPCSnapshotClient(context.Background(), sess, cloudInstanceID)
		_, err = snapshotC.Get(snapshotID)
		if err == nil {
			return fmt.Errorf("PPC Instance Snapshot still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}
func testAccCheckIBMPPCInstanceSnapshotExists(n string) resource.TestCheckFunc {
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
		cloudInstanceID, snapshotID, err := splitID(rs.Primary.ID)
		if err != nil {
			return err
		}
		client := st.NewIBMPPCSnapshotClient(context.Background(), sess, cloudInstanceID)

		_, err = client.Get(snapshotID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckIBMPPCInstanceSnapshotConfig(name, healthStatus string) string {
	return testAccCheckIBMPPCInstanceConfig(name, healthStatus) + fmt.Sprintf(`
	resource "ibm_ppc_snapshot" "power_snapshot"{
		depends_on=[ibm_ppc_instance.power_instance]
		ppc_instance_name       = ibm_ppc_instance.power_instance.ppc_instance_name
		ppc_cloud_instance_id = "%s"
		ppc_snap_shot_name       = "%s"
		ppc_volume_ids         = [ibm_ppc_volume.power_volume.volume_id]
	  }
	`, acc.Ppc_cloud_instance_id, name)
}
