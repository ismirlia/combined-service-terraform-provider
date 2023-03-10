---

subcategory: "Power Private Cloud"
layout: "ibm"
page_title: "IBM: ppc_network"
description: |-
  Manages networks in the IBM Power Private Cloud.
---

# ibm_ppc_network
Create, update, or delete a network connection for your Power Private Cloud instance. For more information, about Power instance network, see [setting up an IBM i network install server](https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-preparing-install-server).

## Example usage
The following example creates a network connection for your Power Private Cloud instance.

```terraform
resource "ibm_ppc_network" "power_networks" {
  count                = 1
  ppc_network_name      = "power-network"
  ppc_cloud_instance_id = "<value of the cloud_instance_id>"
  ppc_network_type      = "vlan"
  ppc_cidr              = "<Network in CIDR notation (192.168.0.0/24)>"
  ppc_dns               = [<"DNS Servers">]
  ppc_gateway           = "192.168.0.1"
  ppc_ipaddress_range {
    ppc_starting_ip_address  = "192.168.0.2"
    ppc_ending_ip_address    = "192.168.0.254"
  }
}
```

**Note**
* Please find [supported Regions](https://cloud.ibm.com/apidocs/power-cloud#endpoint) for endpoints.
* If a Power Private Cloud instance is provisioned at `lon04`, The provider level attributes should be as follows:
  * `region` - `lon`
  * `zone` - `lon04`
  
  Example usage:

  ```terraform
    provider "ibm" {
      region    =   "lon"
      zone      =   "lon04"
    }
  ```

## Timeouts

The `ibm_ppc_network` provides the following [Timeouts](https://www.terraform.io/docs/language/resources/syntax.html) configuration options:

- **create** - (Default 60 minutes) Used for creating a network.
- **update** - (Default 60 minutes) Used for updating a network.
- **delete** - (Default 60 minutes) Used for deleting a network.

## Argument reference 
Review the argument references that you can specify for your resource. 

- `ppc_cloud_instance_id` - (Required, String) The GUID of the service instance associated with an account.
- `ppc_network_name` - (Required, String) The name of the network.
- `ppc_network_type` - (Required, String) The type of network that you want to create, such as `pub-vlan` or `vlan`.
- `ppc_dns` - (Optional, Set of String) The DNS Servers for the network. Required for `vlan` network type.
- `ppc_cidr` - (Optional, String) The network CIDR. Required for `vlan` network type.
- `ppc_gateway` - (Optional, String) The gateway ip address.
- `ppc_ipaddress_range` - (Optional, List of Map) List of one or more ip address range. The `ppc_ipaddress_range` object structure is documented below. 
  The `ppc_ipaddress_range` block supports:
  - `ppc_ending_ip_address` - (Required, String) The ending ip address.
  - `ppc_starting_ip_address` - (Required, String) The staring ip address. **Note** if the `ppc_gateway` or `ppc_ipaddress_range` is not provided, it will calculate the value based on CIDR respectively.
- `ppc_network_mtu` - (Optional, Integer) Maximum Transmission Unit option of the network.

## Attribute reference
In addition to all argument reference list, you can access the following attribute reference after your resource is created.

- `id` - (String) The unique identifier of the network. The ID is composed of `<power_instance_id>/<network_id>`.
- `network_id` - (String) The unique identifier of the network.
- `vlan_id` - (Integer) The ID of the VLAN that your network is attached to. 

## Import
The `ibm_ppc_network` resource can be imported by using `power_instance_id` and `network_id`.

**Example**

```
$ terraform import ibm_ppc_network.example d7bec597-4726-451f-8a63-e62e6f19c32c/cea6651a-bc0a-4438-9f8a-a0770bbf3ebb
```
