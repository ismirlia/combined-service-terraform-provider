---

subcategory: "Power Private Cloud"
layout: "ibm"
page_title: "IBM: ppc_volume"
description: |-
  Manages IBM Volume in the Power Private Cloud.
---

# ibm_ppc_volume
Create, update, or delete a volume to attach it to a Power Private Cloud instance. For more information, about managing volume, see [cloning a volume](https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-volume-snapshot-clone#cloning-volume).

## Example usage
The following example creates a 20 GB volume.

```terraform
resource "ibm_ppc_volume" "testacc_volume"{
  ppc_volume_size       = 20
  ppc_volume_name       = "test-volume"
  ppc_volume_type       = "ssd"
  ppc_volume_shareable  = true
  ppc_cloud_instance_id = "<value of the cloud_instance_id>"
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

ibm_ppc_volume provides the following [timeouts](https://www.terraform.io/docs/language/resources/syntax.html) configuration options:

- **create** - (Default 30 minutes) Used for creating volume.
- **update** - (Default 30 minutes) Used for updating volume.
- **delete** - (Default 10 minutes) Used for deleting volume.

## Argument reference 
Review the argument references that you can specify for your resource. 

- `ppc_affinity_instance` - (Optional, String) PVM Instance (ID or Name) to base volume affinity policy against; required if requesting `affinity` and `ppc_affinity_volume` is not provided.
- `ppc_affinity_policy` - (Optional, String) Affinity policy for data volume being created; ignored if `ppc_volume_pool` provided; for policy 'affinity' requires one of `ppc_affinity_instance` or `ppc_affinity_volume` to be specified; for policy 'anti-affinity' requires one of `ppc_anti_affinity_instances` or `ppc_anti_affinity_volumes` to be specified; Allowable values: `affinity`, `anti-affinity`
- `ppc_affinity_volume`- (Optional, String) Volume (ID or Name) to base volume affinity policy against; required if requesting `affinity` and `ppc_affinity_instance` is not provided.
- `ppc_anti_affinity_instances` - (Optional, String) List of pvmInstances to base volume anti-affinity policy against; required if requesting `anti-affinity` and `ppc_anti_affinity_volumes` is not provided.
- `ppc_anti_affinity_volumes`- (Optional, String) List of volumes to base volume anti-affinity policy against; required if requesting `anti-affinity` and `ppc_anti_affinity_instances` is not provided.
- `ppc_cloud_instance_id` - (Required, String) The GUID of the service instance associated with an account.
- `ppc_replication_enabled` - (Optional, Bool) Indicates if the volume should be replication enabled or not.
- `ppc_volume_name` - (Required, String) The name of the volume.
- `ppc_volume_pool` - (Optional, String) Volume pool where the volume will be created; if provided then `ppc_volume_type` and `ppc_affinity_policy` values will be ignored.
- `ppc_volume_shareable` - (Required, Bool) If set to **true**, the volume can be shared across Power Private Cloud instances. If set to **false**, you can attach it only to one instance. 
- `ppc_volume_size`  - (Required, Integer) The size of the volume in gigabytes. 
- `ppc_volume_type` - (Optional, String) Type of Disk, required if `ppc_affinity_policy` and `ppc_volume_pool` not provided, otherwise ignored. Supported values are `ssd`, `standard`, `tier1`, and `tier3`.

## Attribute reference
In addition to all argument reference list, you can access the following attribute reference after your resource is created.

- `auxiliary` - (Bool) Indicates if the volume is auxiliary or not.
- `auxiliary_volume_name` - (String) The auxiliary volume name.
- `consistency_group_name` - (String) The consistency group name if volume is a part of volume group.
- `delete_on_termination` - (Bool) Indicates if the volume should be deleted when the server terminates.
- `group_id` - (String) The volume group id to which volume belongs.
- `id` - (String) The unique identifier of the volume. The ID is composed of `<power_instance_id>/<volume_id>`.
- `master_volume_name` - (String) The master volume name.
- `mirroring_state` - (String) The mirroring state for replication enabled volume.
- `primary_role` - (String) Indicates whether `master`/`auxiliary` volume is playing the primary role.
- `replication_status` - (String) The replication status of the volume.
- `replication_type` - (String) The replication type of the volume `metro` or `global`.
- `status_description_errors` - (List of objects) - The status details of the volume group.

  Nested scheme for `status_description_errors`:
  - `key` - (String) The volume group error key.
  - `message` - (String) The failure message providing more details about the error key.
  - `volume_ids` - (List of strings) List of volume IDs, which failed to be added/removed to/from the volume group, with the given error.
- `volume_id` - (String) The unique identifier of the volume.
- `volume_status` - (String) The status of the volume.
- `wwn` - (String) The world wide name of the volume.

## Import

The `ibm_ppc_volume` resource can be imported by using `power_instance_id` and `volume_id`.

**Example**

```
$ terraform import ibm_ppc_volume.example d7bec597-4726-451f-8a63-e62e6f19c32c/cea6651a-bc0a-4438-9f8a-a0770bbf3ebb
```
