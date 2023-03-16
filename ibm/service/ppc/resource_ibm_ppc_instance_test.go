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

func testAccCheckIBMPPCInstanceConfig(name, instanceHealthStatus string) string {
	return fmt.Sprintf(`
	resource "ibm_ppc_key" "key" {
		ppc_cloud_instance_id = "%[1]s"
		ppc_key_name          = "%[2]s"
		ppc_ssh_key           = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCKVmnMOlHKcZK8tpt3MP1lqOLAcqcJzhsvJcjscgVERRN7/9484SOBJ3HSKxxNG5JN8owAjy5f9yYwcUg+JaUVuytn5Pv3aeYROHGGg+5G346xaq3DAwX6Y5ykr2fvjObgncQBnuU5KHWCECO/4h8uWuwh/kfniXPVjFToc+gnkqA+3RKpAecZhFXwfalQ9mMuYGFxn+fwn8cYEApsJbsEmb0iJwPiZ5hjFC8wREuiTlhPHDgkBLOiycd20op2nXzDbHfCHInquEe/gYxEitALONxm0swBOwJZwlTDOB7C6y2dzlrtxr1L59m7pCkWI4EtTRLvleehBoj3u7jB4usR"
	  }
	  data "ibm_ppc_image" "power_image" {
		ppc_image_name        = "%[3]s"
		ppc_cloud_instance_id = "%[1]s"
	  }
	  data "ibm_ppc_network" "power_networks" {
		ppc_cloud_instance_id = "%[1]s"
		ppc_network_name      = "%[4]s"
	  }
	  resource "ibm_ppc_volume" "power_volume" {
		ppc_volume_size       = 20
		ppc_volume_name       = "%[2]s"
		ppc_volume_shareable  = true
		ppc_volume_pool       = data.ibm_ppc_image.power_image.storage_pool
		ppc_cloud_instance_id = "%[1]s"
	  }
	  resource "ibm_ppc_instance" "power_instance" {
		ppc_memory             = "2"
		ppc_processors         = "0.25"
		ppc_instance_name      = "%[2]s"
		ppc_proc_type          = "shared"
		ppc_image_id           = data.ibm_ppc_image.power_image.id
		ppc_key_pair_name      = ibm_ppc_key.key.key_id
		ppc_sys_type           = "s922"
		ppc_cloud_instance_id  = "%[1]s"
		ppc_storage_pool       = data.ibm_ppc_image.power_image.storage_pool
		ppc_health_status      = "%[5]s"
		ppc_volume_ids         = [ibm_ppc_volume.power_volume.volume_id]
		ppc_network {
			network_id = data.ibm_ppc_network.power_networks.id
		}
	  }
	`, acc.Ppc_cloud_instance_id, name, acc.Ppc_image, acc.Ppc_network_name, instanceHealthStatus)
}

func testAccCheckIBMPPCInstanceDeploymentTypeConfig(name, instanceHealthStatus string) string {
	return fmt.Sprintf(`
	resource "ibm_ppc_key" "key" {
		ppc_cloud_instance_id = "%[1]s"
		ppc_key_name          = "%[2]s"
		ppc_ssh_key           = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCKVmnMOlHKcZK8tpt3MP1lqOLAcqcJzhsvJcjscgVERRN7/9484SOBJ3HSKxxNG5JN8owAjy5f9yYwcUg+JaUVuytn5Pv3aeYROHGGg+5G346xaq3DAwX6Y5ykr2fvjObgncQBnuU5KHWCECO/4h8uWuwh/kfniXPVjFToc+gnkqA+3RKpAecZhFXwfalQ9mMuYGFxn+fwn8cYEApsJbsEmb0iJwPiZ5hjFC8wREuiTlhPHDgkBLOiycd20op2nXzDbHfCHInquEe/gYxEitALONxm0swBOwJZwlTDOB7C6y2dzlrtxr1L59m7pCkWI4EtTRLvleehBoj3u7jB4usR"
	  }
	  data "ibm_ppc_image" "power_image" {
		ppc_image_name        = "%[3]s"
		ppc_cloud_instance_id = "%[1]s"
	  }
	  data "ibm_ppc_network" "power_networks" {
		ppc_cloud_instance_id = "%[1]s"
		ppc_network_name      = "%[4]s"
	  }
	  resource "ibm_ppc_instance" "power_instance" {
		ppc_memory             = "2"
		ppc_processors         = "1"
		ppc_instance_name      = "%[2]s"
		ppc_proc_type          = "dedicated"
		ppc_image_id           = data.ibm_ppc_image.power_image.id
		ppc_key_pair_name      = ibm_ppc_key.key.key_id
		ppc_sys_type           = "e980"
		ppc_cloud_instance_id  = "%[1]s"
		ppc_storage_type 	  = "tier1"
		ppc_health_status      = "%[5]s"
		ppc_network {
			network_id = data.ibm_ppc_network.power_networks.id
		}
		ppc_deployment_type          = "ESPC"
	  }
	`, acc.Ppc_cloud_instance_id, name, acc.Ppc_image, acc.Ppc_network_name, instanceHealthStatus)
}

func testAccIBMPPCInstanceNetworkConfig(name, privateNetIP string) string {
	return fmt.Sprintf(`
	resource "ibm_ppc_key" "key" {
		ppc_cloud_instance_id = "%[1]s"
		ppc_key_name          = "%[2]s"
		ppc_ssh_key           = "ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEArb2aK0mekAdbYdY9rwcmeNSxqVCwez3WZTYEq+1Nwju0x5/vQFPSD2Kp9LpKBbxx3OVLN4VffgGUJznz9DAr7veLkWaf3iwEil6U4rdrhBo32TuDtoBwiczkZ9gn1uJzfIaCJAJdnO80Kv9k0smbQFq5CSb9H+F5VGyFue/iVd5/b30MLYFAz6Jg1GGWgw8yzA4Gq+nO7HtyuA2FnvXdNA3yK/NmrTiPCdJAtEPZkGu9LcelkQ8y90ArlKfjtfzGzYDE4WhOufFxyWxciUePh425J2eZvElnXSdGha+FCfYjQcvqpCVoBAG70U4fJBGjB+HL/GpCXLyiYXPrSnzC9w=="
	}
	resource "ibm_ppc_network" "power_networks" {
		ppc_cloud_instance_id = "%[1]s"
		ppc_network_name      = "%[2]s"
		ppc_network_type      = "vlan"
		ppc_dns               = ["127.0.0.1"]
		ppc_gateway           = "192.168.17.2"
		ppc_cidr              = "192.168.17.0/24"
		ppc_ipaddress_range {
			ppc_ending_ip_address = "192.168.17.254"
			ppc_starting_ip_address = "192.168.17.3"
		}
	}
	resource "ibm_ppc_instance" "power_instance" {
		ppc_memory             = "2"
		ppc_processors         = "0.25"
		ppc_instance_name      = "%[2]s"
		ppc_proc_type          = "shared"
		ppc_image_id           = "f4501cad-d0f4-4517-9eea-85402309d90d"
		ppc_key_pair_name      = ibm_ppc_key.key.key_id
		ppc_sys_type           = "e980"
		ppc_storage_type 	  = "tier3"
		ppc_cloud_instance_id  = "%[1]s"
		ppc_network {
			network_id = resource.ibm_ppc_network.power_networks.id
			ip_address = "%[3]s"
		}
	}
	`, acc.Ppc_cloud_instance_id, name, privateNetIP)
}

func testAccIBMPPCInstanceVTLConfig(name string) string {
	return fmt.Sprintf(`
	resource "ibm_ppc_key" "vtl_key" {
		ppc_cloud_instance_id = "%[1]s"
		ppc_key_name          = "%[2]s"
		ppc_ssh_key           = "ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEArb2aK0mekAdbYdY9rwcmeNSxqVCwez3WZTYEq+1Nwju0x5/vQFPSD2Kp9LpKBbxx3OVLN4VffgGUJznz9DAr7veLkWaf3iwEil6U4rdrhBo32TuDtoBwiczkZ9gn1uJzfIaCJAJdnO80Kv9k0smbQFq5CSb9H+F5VGyFue/iVd5/b30MLYFAz6Jg1GGWgw8yzA4Gq+nO7HtyuA2FnvXdNA3yK/NmrTiPCdJAtEPZkGu9LcelkQ8y90ArlKfjtfzGzYDE4WhOufFxyWxciUePh425J2eZvElnXSdGha+FCfYjQcvqpCVoBAG70U4fJBGjB+HL/GpCXLyiYXPrSnzC9w=="
	}
	
	resource "ibm_ppc_network" "vtl_network" {
		ppc_cloud_instance_id = "%[1]s"
		ppc_network_name      = "%[2]s"
		ppc_network_type      = "pub-vlan"
	}

	resource "ibm_ppc_instance" "vtl_instance" {
		ppc_memory             = "22"
		ppc_processors         = "2"
		ppc_instance_name      = "%[2]s"
		ppc_license_repository_capacity = "3"
		ppc_proc_type          = "shared"
		ppc_image_id           = "%[3]s"
		ppc_key_pair_name      = ibm_ppc_key.vtl_key.key_id
		ppc_sys_type           = "s922"
		ppc_cloud_instance_id  = "%[1]s"
		ppc_storage_type 	  = "tier1"
		ppc_network {
			network_id = ibm_ppc_network.vtl_network.network_id
		}
	  }
	
	`, acc.Ppc_cloud_instance_id, name, acc.Ppc_image)
}

func testAccCheckIBMPPCInstanceDestroy(s *terraform.State) error {
	sess, err := acc.TestAccProvider.Meta().(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_ppc_instance" {
			continue
		}
		cloudInstanceID, instanceID, err := splitID(rs.Primary.ID)
		if err == nil {
			return err
		}
		client := st.NewIBMPPCInstanceClient(context.Background(), sess, cloudInstanceID)
		_, err = client.Get(instanceID)
		if err == nil {
			return fmt.Errorf("SP Instance still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}
func testAccCheckIBMPPCInstanceExists(n string) resource.TestCheckFunc {
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

		cloudInstanceID, instanceID, err := splitID(rs.Primary.ID)
		if err == nil {
			return err
		}
		client := st.NewIBMPPCInstanceClient(context.Background(), sess, cloudInstanceID)

		_, err = client.Get(instanceID)
		if err != nil {
			return err
		}

		return nil
	}
}

func TestAccIBMPPCInstanceBasic(t *testing.T) {
	instanceRes := "ibm_ppc_instance.power_instance"
	name := fmt.Sprintf("tf-pi-instance-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCInstanceConfig(name, helpers.PPCInstanceHealthWarning),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCInstanceExists(instanceRes),
					resource.TestCheckResourceAttr(instanceRes, "ppc_instance_name", name),
				),
			},
		},
	})
}

func TestAccIBMPPCInstanceDeploymentType(t *testing.T) {
	instanceRes := "ibm_ppc_instance.power_instance"
	name := fmt.Sprintf("tf-pi-instance-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPPCInstanceDeploymentTypeConfig(name, helpers.PPCInstanceHealthWarning),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCInstanceExists(instanceRes),
					resource.TestCheckResourceAttr(instanceRes, "ppc_instance_name", name),
				),
			},
		},
	})
}

func TestAccIBMPPCInstanceNetwork(t *testing.T) {
	instanceRes := "ibm_ppc_instance.power_instance"
	name := fmt.Sprintf("tf-pi-instance-%d", acctest.RandIntRange(10, 100))
	privateNetIP := "192.168.17.253"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIBMPPCInstanceNetworkConfig(name, privateNetIP),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCInstanceExists(instanceRes),
					resource.TestCheckResourceAttr(instanceRes, "ppc_instance_name", name),
					resource.TestCheckResourceAttrSet(instanceRes, "ppc_network.0.network_id"),
					resource.TestCheckResourceAttrSet(instanceRes, "ppc_network.0.mac_address"),
					resource.TestCheckResourceAttr(instanceRes, "ppc_network.0.ip_address", privateNetIP),
				),
			},
		},
	})
}

func TestAccIBMPPCInstanceVTL(t *testing.T) {
	instanceRes := "ibm_ppc_instance.vtl_instance"
	name := fmt.Sprintf("tf-pi-vtl-instance-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIBMPPCInstanceVTLConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCInstanceExists(instanceRes),
					resource.TestCheckResourceAttr(instanceRes, "ppc_instance_name", name),
					resource.TestCheckResourceAttr(instanceRes, "ppc_license_repository_capacity", "3"),
				),
			},
		},
	})
}

func TestAccIBMPPCSASPnstance(t *testing.T) {
	instanceRes := "ibm_ppc_instance.sap"
	name := fmt.Sprintf("tf-pi-sap-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIBMPPCSASPnstanceConfig(name, "tinytest-1x4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCInstanceExists(instanceRes),
					resource.TestCheckResourceAttr(instanceRes, "ppc_instance_name", name),
					resource.TestCheckResourceAttr(instanceRes, "ppc_sap_profile_id", "tinytest-1x4"),
				),
			},
			{
				Config: testAccIBMPPCSASPnstanceConfig(name, "tinytest-1x8"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCInstanceExists(instanceRes),
					resource.TestCheckResourceAttr(instanceRes, "ppc_instance_name", name),
					resource.TestCheckResourceAttr(instanceRes, "ppc_sap_profile_id", "tinytest-1x8"),
				),
			},
		},
	})
}
func testAccIBMPPCSASPnstanceConfig(name, sapProfile string) string {
	return fmt.Sprintf(`
	resource "ibm_ppc_network" "power_network" {
		ppc_cloud_instance_id	= "%[1]s"
		ppc_network_name			= "%[2]s"
		ppc_network_type			= "pub-vlan"
	}
	resource "ibm_ppc_instance" "sap" {
		ppc_cloud_instance_id  	= "%[1]s"
		ppc_instance_name      	= "%[2]s"
		ppc_sap_profile_id       = "%[3]s"
		ppc_image_id           	= "%[4]s"
		ppc_storage_type			= "tier1"
		ppc_network {
			network_id = ibm_ppc_network.power_network.network_id
		}
		ppc_health_status		= "OK"
	}
	`, acc.Ppc_cloud_instance_id, name, sapProfile, acc.Ppc_sap_image)
}

func TestAccIBMPPCInstanceMixedStorage(t *testing.T) {
	instanceRes := "ibm_ppc_instance.instance"
	name := fmt.Sprintf("tf-pi-mixedstorage-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPPCInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIBMPPCInstanceMixedStorage(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMPPCInstanceExists(instanceRes),
					resource.TestCheckResourceAttr(instanceRes, "ppc_instance_name", name),
					resource.TestCheckResourceAttr(instanceRes, "ppc_storage_pool_affinity", "false"),
				),
			},
		},
	})
}

func testAccIBMPPCInstanceMixedStorage(name string) string {
	return fmt.Sprintf(`
	resource "ibm_ppc_key" "key" {
		ppc_cloud_instance_id = "%[1]s"
		ppc_key_name          = "%[2]s"
		ppc_ssh_key           = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCKVmnMOlHKcZK8tpt3MP1lqOLAcqcJzhsvJcjscgVERRN7/9484SOBJ3HSKxxNG5JN8owAjy5f9yYwcUg+JaUVuytn5Pv3aeYROHGGg+5G346xaq3DAwX6Y5ykr2fvjObgncQBnuU5KHWCECO/4h8uWuwh/kfniXPVjFToc+gnkqA+3RKpAecZhFXwfalQ9mMuYGFxn+fwn8cYEApsJbsEmb0iJwPiZ5hjFC8wREuiTlhPHDgkBLOiycd20op2nXzDbHfCHInquEe/gYxEitALONxm0swBOwJZwlTDOB7C6y2dzlrtxr1L59m7pCkWI4EtTRLvleehBoj3u7jB4usR"
	}
	resource "ibm_ppc_network" "power_network" {
		ppc_cloud_instance_id = "%[1]s"
		ppc_network_name      = "%[2]s"
		ppc_network_type      = "pub-vlan"
	}
	resource "ibm_ppc_volume" "power_volume" {
		ppc_cloud_instance_id = "%[1]s"
		ppc_volume_size       = 20
		ppc_volume_name       = "%[2]s"
		ppc_volume_shareable  = true
		ppc_volume_type       = "tier3"
	}
	resource "ibm_ppc_instance" "instance" {
		ppc_cloud_instance_id     = "%[1]s"
		ppc_memory                = "2"
		ppc_processors            = "0.25"
		ppc_instance_name         = "%[2]s"
		ppc_proc_type             = "shared"
		ppc_image_id              = "ca4ea55f-b329-4cf5-bdce-d2f38cfc6da3"
		ppc_key_pair_name         = ibm_ppc_key.key.key_id
		ppc_sys_type              = "s922"
		ppc_storage_type          = "tier1"
		ppc_storage_pool_affinity = false
		ppc_network {
			network_id = ibm_ppc_network.power_network.network_id
		}
	}
	resource "ibm_ppc_volume_attach" "power_attach_volume"{
		ppc_cloud_instance_id = "%[1]s"
		ppc_volume_id         = ibm_ppc_volume.power_volume.volume_id
		ppc_instance_id       = ibm_ppc_instance.instance.instance_id
	}
	`, acc.Ppc_cloud_instance_id, name)
}
