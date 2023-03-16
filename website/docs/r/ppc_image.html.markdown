---

subcategory: "Power Private Cloud"
layout: "ibm"
page_title: "IBM: ppc_image"
description: |-
  Manages IBM Image in the Power Private Cloud.
---

# ibm_ppc_image
Create, update, or delete for a Power Private Cloud image. For more information, about IBM Power Private Cloud, see [getting started with IBM Power Private Cloud Virtual Servers](https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-getting-started).

## Example usage
The following example enables you to create a image:

```terraform
resource "ibm_ppc_image" "testacc_image  "{
  ppc_image_name       = "7200-03-02"
  ppc_image_id         = <"image id obtained from the datasource">
  ppc_cloud_instance_id = "<value of the cloud_instance_id>"
}
```

```terraform
resource "ibm_ppc_image" "testacc_image  "{
  ppc_image_name       = "test_image"
  ppc_cloud_instance_id = "<value of the cloud_instance_id>"
  ppc_image_bucket_name = "images-public-bucket"
  ppc_image_bucket_access = "public"
  ppc_image_bucket_region = "us-south"
  ppc_image_bucket_file_name = "rhcos-48-07222021.ova.gz"
  ppc_image_storage_type = "tier1"
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

The   ibm_ppc_image   provides the following [timeouts](https://www.terraform.io/docs/language/resources/syntax.html) configuration options:

- **Create** The creation of the image is considered failed if no response is received for 60 minutes. 
- **Delete** The deletion of the image is considered failed if no response is received for 60 minutes. 

## Argument reference
Review the argument references that you can specify for your resource. 

- `ppc_affinity_instance` - (Optional, String) PVM Instance (ID or Name) to base storage affinity policy against; required if requesting `affinity` and `ppc_affinity_volume` is not provided.
- `ppc_affinity_policy` - (Optional, String) Affinity policy for image; ignored if `ppc_image_storage_pool` provided; for policy affinity requires one of `ppc_affinity_instance` or `ppc_affinity_volume` to be specified; for policy anti-affinity requires one of `ppc_anti_affinity_instances` or `ppc_anti_affinity_volumes` to be specified; Allowable values: `affinity`, `anti-affinity`
- `ppc_affinity_volume`- (Optional, String) Volume (ID or Name) to base storage affinity policy against; required if requesting `affinity` and `ppc_affinity_instance` is not provided.
- `ppc_anti_affinity_instances` - (Optional, String) List of pvmInstances to base storage anti-affinity policy against; required if requesting `anti-affinity` and `ppc_anti_affinity_volumes` is not provided.
- `ppc_anti_affinity_volumes`- (Optional, String) List of volumes to base storage anti-affinity policy against; required if requesting `anti-affinity` and `ppc_anti_affinity_instances` is not provided.
- `ppc_cloud_instance_id` - (Required, String) The GUID of the service instance associated with an account.
- `ppc_image_name` - (Required, String) The name of an image.
- `ppc_image_id` - (Optional, String) Image ID of existing source image; required for copy image.
  - Either `ppc_image_id` or `ppc_image_bucket_name` is required.
- `ppc_image_bucket_name` - (Optional, String) Cloud Object Storage bucket name; `bucket-name[/optional/folder]`
  - Either `ppc_image_bucket_name` or `ppc_image_id` is required.
- `ppc_image_access_key` - (Optional, String, Sensitive) Cloud Object Storage access key; required for buckets with private access.
  - `ppc_image_access_key` is required with `ppc_image_secret_key`
- `ppc_image_bucket_access` - (Optional, String) Indicates if the bucket has public or private access. The default value is `public`.
- `ppc_image_bucket_file_name` - (Optional, String) Cloud Object Storage image filename
  - `ppc_image_bucket_file_name` is required with `ppc_image_bucket_name`
- `ppc_image_bucket_region` - (Optional, String) Cloud Object Storage region
  - `ppc_image_bucket_region` is required with `ppc_image_bucket_name`
- `ppc_image_secret_key` - (Optional, String, Sensitive) Cloud Object Storage secret key; required for buckets with private access.
  - `ppc_image_secret_key` is required with `ppc_image_access_key`
- `ppc_image_storage_pool` - (Optional, String) Storage pool where the image will be loaded, if provided then `ppc_image_storage_type` and `ppc_affinity_policy` will be ignored.
- `ppc_image_storage_type` - (Optional, String) Type of storage. Will be ignored if `ppc_image_storage_pool` or `ppc_affinity_policy` is provided. If only using `ppc_image_storage_type` for storage selection then the storage pool with the most available space will be selected.


## Attribute reference
In addition to all argument reference list, you can access the following attribute reference after your resource is created.

- `id` - (String) The unique identifier of an image. The ID is composed of `<ppc_cloud_instance_id>/<image_id>`. 
- `image_id` - (String) The unique identifier of an image.

## Import

The `ibm_ppc_image` can be imported by using `ppc_cloud_instance_id` and `image_id`.

**Example**

```
$ terraform import ibm_ppc_image.example d7bec597-4726-451f-8a63-e62e6f19c32c/cea6651a-bc0a-4438-9f8a-a0770bbf3ebb
```
