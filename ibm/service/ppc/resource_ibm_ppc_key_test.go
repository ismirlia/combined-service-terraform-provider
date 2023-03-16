// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
)

func TestAccIBMPPCKey_basic(t *testing.T) {
	publicKey := strings.TrimSpace(`
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCKVmnMOlHKcZK8tpt3MP1lqOLAcqcJzhsvJcjscgVERRN7/9484SOBJ3HSKxxNG5JN8owAjy5f9yYwcUg+JaUVuytn5Pv3aeYROHGGg+5G346xaq3DAwX6Y5ykr2fvjObgncQBnuU5KHWCECO/4h8uWuwh/kfniXPVjFToc+gnkqA+3RKpAecZhFXwfalQ9mMuYGFxn+fwn8cYEApsJbsEmb0iJwPiZ5hjFC8wREuiTlhPHDgkBLOiycd20op2nXzDbHfCHInquEe/gYxEitALONxm0swBOwJZwlTDOB7C6y2dzlrtxr1L59m7pCkWI4EtTRLvleehBoj3u7jB4usR
`)
	name := fmt.Sprintf("tf-sp-sshkey-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCKeyConfig(publicKey, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCKeyExists("ibm_ppc_key.key"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_key.key", "ppc_key_name", name),
				),
			},
		},
	})
}
func testAccCheckIBMPPCKeyDestroy(s *terraform.State) error {

	sess, err := acc.TestAccProvider.Meta().(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_ppc_key" {
			continue
		}
		cloudInstanceID, key, err := splitID(rs.Primary.ID)
		if err != nil {
			return err
		}
		sshkeyC := st.NewIBMPPCKeyClient(context.Background(), sess, cloudInstanceID)
		_, err = sshkeyC.Get(key)
		if err == nil {
			return fmt.Errorf("SP key still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}
func testAccCheckIBMPPCKeyExists(n string) resource.TestCheckFunc {
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

		cloudInstanceID, key, err := splitID(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := st.NewIBMPPCKeyClient(context.Background(), sess, cloudInstanceID)
		_, err = client.Get(key)
		if err != nil {
			return err
		}
		return nil

	}
}

func testAccCheckIBMPPCKeyConfig(publicKey, name string) string {
	return fmt.Sprintf(`
		resource "ibm_ppc_key" "key" {
			ppc_cloud_instance_id = "%s"
			ppc_key_name          = "%s"
			ppc_ssh_key           = "%s"
		  }
	`, acc.Ppc_cloud_instance_id, name, publicKey)
}
