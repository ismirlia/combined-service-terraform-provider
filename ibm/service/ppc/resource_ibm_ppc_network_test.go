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
)

func TestAccIBMPPCNetworkbasic(t *testing.T) {
	name := fmt.Sprintf("tf-sp-network-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCNetworkConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCNetworkExists("ibm_ppc_network.power_networks"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_network.power_networks", "ppc_network_name", name),
					resource.TestCheckResourceAttrSet("ibm_ppc_network.power_networks", "id"),
					resource.TestCheckResourceAttrSet("ibm_ppc_network.power_networks", "ppc_gateway"),
					resource.TestCheckResourceAttrSet("ibm_ppc_network.power_networks", "ppc_ipaddress_range.#"),
				),
			},
			{
				Config: testAccCheckIBMPPCNetworkConfigUpdateDNS(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCNetworkExists("ibm_ppc_network.power_networks"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_network.power_networks", "ppc_network_name", name),
					resource.TestCheckResourceAttr(
						"ibm_ppc_network.power_networks", "ppc_dns.#", "1"),
					resource.TestCheckResourceAttrSet("ibm_ppc_network.power_networks", "id"),
					resource.TestCheckResourceAttrSet("ibm_ppc_network.power_networks", "ppc_gateway"),
					resource.TestCheckResourceAttrSet("ibm_ppc_network.power_networks", "ppc_ipaddress_range.#"),
				),
			},
		},
	})
}
func TestAccIBMPPCNetworkGatewaybasic(t *testing.T) {
	name := fmt.Sprintf("tf-sp-network-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCNetworkGatewayConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCNetworkExists("ibm_ppc_network.power_networks"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_network.power_networks", "ppc_network_name", name),
					resource.TestCheckResourceAttrSet("ibm_ppc_network.power_networks", "ppc_gateway"),
					resource.TestCheckResourceAttrSet("ibm_ppc_network.power_networks", "id"),
					resource.TestCheckResourceAttrSet("ibm_ppc_network.power_networks", "ppc_ipaddress_range.#"),
				),
			},
			{
				Config: testAccCheckIBMPPCNetworkConfigGatewayUpdateDNS(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCNetworkExists("ibm_ppc_network.power_networks"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_network.power_networks", "ppc_network_name", name),
					resource.TestCheckResourceAttr(
						"ibm_ppc_network.power_networks", "ppc_gateway", "192.168.17.2"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_network.power_networks", "ppc_ipaddress_range.0.ppc_ending_ip_address", "192.168.17.254"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_network.power_networks", "ppc_ipaddress_range.0.ppc_starting_ip_address", "192.168.17.3"),
				),
			},
		},
	})
}
func testAccCheckIBMPPCNetworkDestroy(s *terraform.State) error {

	sess, err := acc.TestAccProvider.Meta().(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_ppc_network" {
			continue
		}
		cloudInstanceID, networkID, err := splitID(rs.Primary.ID)
		if err != nil {
			return err
		}
		networkC := st.NewIBMPPCNetworkClient(context.Background(), sess, cloudInstanceID)
		_, err = networkC.Get(networkID)
		if err == nil {
			return fmt.Errorf("SP Network still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}
func testAccCheckIBMPPCNetworkExists(n string) resource.TestCheckFunc {
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
		cloudInstanceID, networkID, err := splitID(rs.Primary.ID)
		if err != nil {
			return err
		}
		client := st.NewIBMPPCNetworkClient(context.Background(), sess, cloudInstanceID)

		_, err = client.Get(networkID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckIBMPPCNetworkConfig(name string) string {
	return fmt.Sprintf(`
		resource "ibm_ppc_network" "power_networks" {
			ppc_cloud_instance_id = "%s"
			ppc_network_name      = "%s"
			ppc_network_type      = "pub-vlan"
		}
	`, acc.Ppc_cloud_instance_id, name)
}

func testAccCheckIBMPPCNetworkConfigUpdateDNS(name string) string {
	return fmt.Sprintf(`
		resource "ibm_ppc_network" "power_networks" {
			ppc_cloud_instance_id = "%s"
			ppc_network_name      = "%s"
			ppc_network_type      = "pub-vlan"
			ppc_dns               = ["127.0.0.1"]
		}
	`, acc.Ppc_cloud_instance_id, name)
}

func testAccCheckIBMPPCNetworkGatewayConfig(name string) string {
	return fmt.Sprintf(`
		resource "ibm_ppc_network" "power_networks" {
			ppc_cloud_instance_id = "%s"
			ppc_network_name      = "%s"
			ppc_network_type      = "vlan"
			ppc_cidr              = "192.168.17.0/24"
		}
	`, acc.Ppc_cloud_instance_id, name)
}

func testAccCheckIBMPPCNetworkConfigGatewayUpdateDNS(name string) string {
	return fmt.Sprintf(`
		resource "ibm_ppc_network" "power_networks" {
			ppc_cloud_instance_id = "%s"
			ppc_network_name      = "%s"
			ppc_network_type      = "vlan"
			ppc_dns               = ["127.0.0.1"]
			ppc_gateway           = "192.168.17.2"
			ppc_cidr              = "192.168.17.0/24"
			ppc_ipaddress_range {
				ppc_ending_ip_address = "192.168.17.254"
				ppc_starting_ip_address = "192.168.17.3"
			}
		}
	`, acc.Ppc_cloud_instance_id, name)
}
