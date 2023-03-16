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

func TestAccIBMPPCNetworkPortAttachbasic(t *testing.T) {
	name := fmt.Sprintf("tf-sp-network-port-attach-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCNetworkPortAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCNetworkPortAttachConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCNetworkPortAttachExists("ibm_ppc_network_port_attach.power_network_port_attach"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_network_port_attach.power_network_port_attach", "ppc_network_name", name),
					resource.TestCheckResourceAttrSet("ibm_ppc_network_port_attach.power_network_port_attach", "id"),
					resource.TestCheckResourceAttrSet("ibm_ppc_network_port_attach.power_network_port_attach", "network_port_id"),
					resource.TestCheckResourceAttrSet("ibm_ppc_network_port_attach.power_network_port_attach", "public_ip"),
				),
			},
		},
	})
}

func TestAccIBMPPCNetworkPortAttachVlanbasic(t *testing.T) {
	name := fmt.Sprintf("tf-sp-network-port-attach-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCNetworkPortAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCNetworkPortAttachVlanConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCNetworkPortAttachExists("ibm_ppc_network_port_attach.power_network_port_attach"),
					resource.TestCheckResourceAttr(
						"ibm_ppc_network_port_attach.power_network_port_attach", "ppc_network_name", name),
					resource.TestCheckResourceAttrSet("ibm_ppc_network_port_attach.power_network_port_attach", "id"),
					resource.TestCheckResourceAttrSet("ibm_ppc_network_port_attach.power_network_port_attach", "network_port_id"),
				),
			},
		},
	})
}
func testAccCheckIBMPPCNetworkPortAttachDestroy(s *terraform.State) error {
	sess, err := acc.TestAccProvider.Meta().(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_ppc_network_port_attach" {
			continue
		}
		parts, err := flex.IdParts(rs.Primary.ID)
		if err != nil {
			return err
		}
		cloudInstanceID := parts[0]
		networkname := parts[1]
		portID := parts[2]
		networkC := st.NewIBMPPCNetworkClient(context.Background(), sess, cloudInstanceID)
		_, err = networkC.GetPort(networkname, portID)
		if err == nil {
			return fmt.Errorf("SP Network Port still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}
func testAccCheckIBMPPCNetworkPortAttachExists(n string) resource.TestCheckFunc {
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
		parts, err := flex.IdParts(rs.Primary.ID)
		if err != nil {
			return err
		}
		cloudInstanceID := parts[0]
		networkname := parts[1]
		portID := parts[2]
		client := st.NewIBMPPCNetworkClient(context.Background(), sess, cloudInstanceID)

		_, err = client.GetPort(networkname, portID)
		if err != nil {
			return err
		}
		return nil

	}
}

func testAccCheckIBMPPCNetworkPortAttachConfig(name string) string {
	return testAccCheckIBMPPCNetworkConfig(name) + fmt.Sprintf(`
	resource "ibm_ppc_network_port_attach" "power_network_port_attach" {
		ppc_cloud_instance_id  = "%s"
		ppc_network_name       = ibm_ppc_network.power_networks.ppc_network_name
		ppc_network_port_description = "IP Reserved for Test UAT"
		ppc_instance_id = "%s"
	}
	`, acc.Ppc_cloud_instance_id, acc.Ppc_instance_name)
}

func testAccCheckIBMPPCNetworkPortAttachVlanConfig(name string) string {
	return testAccCheckIBMPPCNetworkGatewayConfig(name) + fmt.Sprintf(`
	resource "ibm_ppc_network_port_attach" "power_network_port_attach" {
		ppc_cloud_instance_id  = "%s"
		ppc_network_name       = ibm_ppc_network.power_networks.ppc_network_name
		ppc_network_port_description = "IP Reserved for Test UAT"
		ppc_instance_id = "%s"
	}
	`, acc.Ppc_cloud_instance_id, acc.Ppc_instance_name)
}
