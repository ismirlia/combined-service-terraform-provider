---

subcategory: "Power Private Cloud"
layout: "ibm"
page_title: "IBM: ppc_placement_group"
description: |-
  Manages a placement group in the Power Private Cloud.
---

# ibm_ppc_placement_group
Retrieve information about a placement group. For more information, about placement groups, see [Managing server placement groups](https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-placement-groups).

## Example Usage

```terraform
data "ibm_ppc_placement_group" "ds_placement_group" {
  ppc_placement_group_name   = "my-pg"
  ppc_cloud_instance_id = "49fba6c9-23f8-40bc-9899-aca322ee7d5b"
}
```

**Notes**
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
- `ppc_placement_group_name` - (Required, String) The name of the placement group.

## Attribute reference
In addition to all argument reference list, you can access the following attribute references after your data source is created. 

- `id` - (String) The ID of the placement group.
- `members` - (List of strings) The list of server instances IDs that are members of the placement group.
- `policy` - (String) The value of the group's affinity policy. Valid values are affinity and anti-affinity.
