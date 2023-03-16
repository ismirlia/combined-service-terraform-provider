---

subcategory: "Power Private Cloud"
layout: "ibm"
page_title: "IBM: ppc_key"
description: |-
  Manages IBM SSH keys in the Power Private Cloud.
---

# ibm_ppc_key
Create, update, or delete an SSH key for your Power Private Cloud instance. The SSH key is used to access the instance after it is created. For more information, about SSH keys in Power, see [getting started with IBM Power Private Cloud Virtual Servers](https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-getting-started).

## Example usage
The following example enables you to create a SSH key to be used during creation of a Power instance:

```terraform
resource "ibm_ppc_key" "testacc_sshkey" {
  ppc_key_name          = "testkey"
  ppc_ssh_key           = "ssh-rsa <value>"
  ppc_cloud_instance_id = "<value of the cloud_instance_id>"
}
```

## Argument reference
Review the argument references that you can specify for your resource. 

- `ppc_cloud_instance_id` - (Required, String) Cloud Instance ID of a PCloud Instance.
- `ppc_key_name`  - (Required, String) User defined name for the SSH key. 
- `ppc_ssh_key` - (Required, String) SSH RSA key. 

## Attribute reference
 In addition to all argument reference list, you can access the following attribute reference after your resource is created.

- `id` - (String) The unique identifier of the key. The ID is composed of `<ppc_cloud_instance_id>/<ppc_key_name>`.
- `key_id` - (String) User defined name for the SSH key (deprecated - replaced by `name`).
- `name` - (String) User defined name for the SSH key
- `creation_date` - (String) Date of SSH Key creation. 
- `ssh_key` - (String) SSH RSA key.

## Timeouts

ibm_ppc_key provides the following [timeouts](https://www.terraform.io/docs/language/resources/syntax.html) configuration options:

- **create** - (Default 60 minutes) Used for creating a SSH key.
- **delete** - (Default 60 minutes) Used for deleting a SSH key.

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
## Import

The `ibm_ppc_key` resource can be imported by using `ppc_cloud_instance_id` and `ppc_key_name`.

**Example**

```
$ terraform import ibm_ppc_key.example d7bec597-4726-451f-8a63-e62e6f19c32c/mykey
```