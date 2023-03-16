---

subcategory: "Power Private Cloud"
layout: "ibm"
page_title: "IBM: ppc_instance_action"
description: |-
  Performs an action start, stop, reboot, immediate-shutdown, reset on a p VM instance.
---

# ibm_ppc_instance_action
Performs an action on a [Power Private Cloud instance](https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-creating-power-virtual-server).

## Example usage
The following example perform an action "hard-reboot" on a Power Private Cloud instance.

```terraform
resource "ibm_ppc_instance_action" "example" {
  ppc_cloud_instance_id  = "d7bec597-4726-451f-8a63-e62e6f19c32c"
  ppc_instance_id        = "cea6651a-bc0a-4438-9f8a-a0770b112ebb"
  ppc_action             = "hard-reboot"
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

The `ibm_ppc_instance_action` provides the following [timeouts](https://www.terraform.io/docs/language/resources/syntax.html) configuration options:

- **create** - The action on the instance is considered failed if no response is received for 15 minutes.
- **Update** The update action on the instance is considered failed if no response is received for 15 minutes.


## Argument reference
Review the argument references that you can specify for your resource.

- `ppc_action` - (Required, String) Name of the action to take. Allowed values are `start`, `stop`, `hard-reboot`, `soft-reboot`, `immediate-shutdown`, `reset-state`.
- `ppc_cloud_instance_id` - (Required, String) The GUID of the service instance associated with an account.
- `ppc_health_status` - (Optional, String) Specifies if Terraform should poll for the health status to be `OK` or `WARNING`. The default value is `OK`. Ignored for `ppc_action = "reset-state"`.
- `ppc_instance_id` - (Required, String) Custom deployment type information (For Internal Use Only).

## Attribute reference
In addition to all argument reference list, you can access the following attribute reference after your resource is created.

- `health_status` - (String) The health status of the VM.
- `id` - (String) The unique identifier of the instance. The ID is composed of `<cloud_instance_id>/<instance_id>`.
- `progress` - (Float) - The progress of the instance.
- `status` - (String) - The status of the instance.

## Import

The `ibm_ppc_instance_action` can be imported using `cloud_instance_id` and `instance_id`.

**Example**

```
$ terraform import ibm_ppc_instance_action.example d7bec597-4726-451f-8a63-e62e6f19c32c/cea6651a-bc0a-4438-9f8a-a0770b112ebb
```
