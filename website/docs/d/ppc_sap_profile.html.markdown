---

subcategory: "Power Private Cloud"
layout: "ibm"
page_title: "IBM: ppc_sap_profile"
description: |-
  Manages SAP profile in the Power Private Cloud.
---

# ibm_ppc_sap_profile
Retrieve information about a SAP profile. For more information, see [getting started with IBM Power Private Cloud](https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-getting-started).

## Example usage

```terraform
data "ibm_ppc_sap_profile" "example" {
  ppc_cloud_instance_id = "<value of the cloud_instance_id>"
  ppc_sap_profile_id    = "tinytest-1x4"
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
- `ppc_sap_profile_id` - (Required, String) SAP Profile ID.

## Attribute reference
In addition to all argument reference list, you can access the following attribute references after your data source is created.

- `certified` - (Boolean) Has certification been performed on profile.
- `cores` - (Integer) Amount of cores.
- `memory` - (Integer) Amount of memory (in GB).
- `type` - (String) Type of profile.
