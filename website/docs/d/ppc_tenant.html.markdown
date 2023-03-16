---

subcategory: "Power Private Cloud"
layout: "ibm"
page_title: "IBM: ppc_tenant"
description: |-
  Manages a tenant in the IBM Power Private Cloud.
---

# ibm_ppc_tenant
Retrieve information about the tenants that are configured for your Power Private Cloud instance. For more information, about Power tenants, see [network security](https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-network-security).

## Example usage
The following example retrieves all tenants for the Power Private Cloud instance with the ID.

```terraform
data "ibm_ppc_tenant" "ds_tenant" {
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

## Attribute reference
In addition to all argument reference list, you can access the following attribute references after your data source is created. 

- `creation_date` - (Timestamp) The timestamp when the tenant was created.
- `cloud_instances` - (List) A list with the regions and Power Private Cloud instance IDs that the tenant owns.

  Nested scheme for `cloud_instances`:
	- `cloud_instance_id` - (String) The unique identifier of the cloud instance.
	- `region` - (String) The region of the cloud instance.
- `enabled` - (Bool) Indicates if the tenant is enabled for the Power Private Cloud instance ID.
- `id` - (String) The ID of the tenant.
- `tenant_name` -  (String) The name of the tenant.
