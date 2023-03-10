---

subcategory: "Power Private Cloud"
layout: "ibm"
page_title: "IBM: ppc_network"
description: |-
  Manages a network in the IBM Power Private Cloud.
---

# ibm_ppc_network
Retrieve information about the network that your Power Private Cloud instance is connected to. For more information, about Power instance network, see [setting up an IBM i network install server](https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-preparing-install-server).

## Example usage

```terraform
data "ibm_ppc_network" "ds_network" {
  ppc_network_name = "APP"
  ppc_cloud_instance_id = "49fba6c9-23f8-40bc-9899-aca322ee7d5b"
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
  
## Argument reference
Review the argument references that you can specify for your data source. 

- `ppc_cloud_instance_id` - (Required, String) The GUID of the service instance associated with an account.
- `ppc_network_name` - (Required, String) The name of the network.

## Attribute reference
In addition to all argument reference list, you can access the following attribute references after your data source is created. 

- `available_ip_count` - (Float) The total number of IP addresses that you have in your network.
- `cidr` - (String) The CIDR of the network.
- `dns`- (Set of String) The DNS Servers for the network.
- `gateway` - (String) The network gateway that is attached to your network.
- `id` - (String) The ID of the network.
- `type` - (String) The type of network.
- `used_ip_count` - (Float) The number of used IP addresses.
- `used_ip_percent` - (Float) The percentage of IP addresses used.
- `vlan_id` - (String) The VLAN ID that the network is connected to.
- `mtu` - (Integer) Maximum Transmission Unit option of the network.
