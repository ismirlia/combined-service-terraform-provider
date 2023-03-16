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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
)

func TestAccIBMPPCDhcpbasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCDhcpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCDhcpConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCDhcpExists("ibm_ppc_dhcp.dhcp_service"),
					resource.TestCheckResourceAttrSet(
						"ibm_ppc_dhcp.dhcp_service", "dhcp_id"),
				),
			},
		},
	})
}

func testAccCheckIBMPPCDhcpDestroy(s *terraform.State) error {
	sess, err := acc.TestAccProvider.Meta().(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_ppc_dhcp" {
			continue
		}

		cloudInstanceID, dhcpID, err := splitID(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := st.NewIBMPPCDhcpClient(context.Background(), sess, cloudInstanceID)
		_, err = client.Get(dhcpID)
		if err == nil {
			return fmt.Errorf("SP DHCP still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}
func testAccCheckIBMPPCDhcpExists(n string) resource.TestCheckFunc {
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

		cloudInstanceID, dhcpID, err := splitID(rs.Primary.ID)
		if err != nil {
			return err
		}
		client := st.NewIBMPPCDhcpClient(context.Background(), sess, cloudInstanceID)

		_, err = client.Get(dhcpID)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckIBMPPCDhcpConfig() string {
	return fmt.Sprintf(`
	resource "ibm_ppc_dhcp" "dhcp_service" {
		ppc_cloud_instance_id = "%s"
	}
	`, acc.Ppc_cloud_instance_id)
}

func TestAccIBMPPCDhcpWithCidrName(t *testing.T) {
	name := fmt.Sprintf("tf-dhcp-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCDhcpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCDhcpWithCidrNameConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCDhcpExists("ibm_ppc_dhcp.dhcp_service"),
					resource.TestCheckResourceAttrSet(
						"ibm_ppc_dhcp.dhcp_service", "dhcp_id"),
					resource.TestCheckResourceAttrSet(
						"ibm_ppc_dhcp.dhcp_service", "status"),
					resource.TestCheckResourceAttrSet(
						"ibm_ppc_dhcp.dhcp_service", "network_id"),
					resource.TestCheckResourceAttrSet(
						"ibm_ppc_dhcp.dhcp_service", "network_name"),
				),
			},
		},
	})
}

func testAccCheckIBMPPCDhcpWithCidrNameConfig(name string) string {
	return fmt.Sprintf(`
		resource "ibm_ppc_dhcp" "dhcp_service" {
			ppc_cloud_instance_id 	= "%[1]s"
			ppc_dhcp_name = "%[2]s"
			ppc_cidr = "192.168.103.0/24"
		}
	`, acc.Ppc_cloud_instance_id, name)
}

func TestAccIBMPPCDhcpSNAT(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCDhcpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCDhcpConfigWithSNATDisabled(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCDhcpExists("ibm_ppc_dhcp.dhcp_service"),
					resource.TestCheckResourceAttrSet(
						"ibm_ppc_dhcp.dhcp_service", "dhcp_id"),
				),
			},
		},
	})
}

func testAccCheckIBMPPCDhcpConfigWithSNATDisabled() string {
	return fmt.Sprintf(`
	resource "ibm_ppc_dhcp" "dhcp_service" {
		ppc_cloud_instance_id = "%s"
		ppc_dhcp_snat_enabled = false
	}
	`, acc.Ppc_cloud_instance_id)
}

func splitID(id string) (id1, id2 string, err error) {
	parts, err := flex.IdParts(id)
	if err != nil {
		return
	}
	id1 = parts[0]
	id2 = parts[1]
	return
}
