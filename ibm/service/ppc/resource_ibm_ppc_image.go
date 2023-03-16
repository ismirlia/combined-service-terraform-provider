// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/ppc-aas-go-client/errors"
	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
	"github.com/IBM-Cloud/ppc-aas-go-client/ppc-aas/client/p_cloud_images"
	"github.com/IBM-Cloud/ppc-aas-go-client/ppc-aas/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

func ResourceIBMPPCImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPPCImageCreate,
		ReadContext:   resourceIBMPPCImageRead,
		DeleteContext: resourceIBMPPCImageDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			helpers.PPCCloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PPC cloud instance ID",
				ForceNew:    true,
			},
			helpers.PPCImageName: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Image name",
				DiffSuppressFunc: flex.ApplyOnce,
				ForceNew:         true,
			},
			helpers.PPCImageId: {
				Type:             schema.TypeString,
				Optional:         true,
				ExactlyOneOf:     []string{helpers.PPCImageId, helpers.PPCImageBucketName},
				Description:      "Instance image id",
				DiffSuppressFunc: flex.ApplyOnce,
				ConflictsWith:    []string{helpers.PPCImageBucketName},
				ForceNew:         true,
			},

			// COS import variables
			helpers.PPCImageBucketName: {
				Type:          schema.TypeString,
				Optional:      true,
				ExactlyOneOf:  []string{helpers.PPCImageId, helpers.PPCImageBucketName},
				Description:   "Cloud Object Storage bucket name; bucket-name[/optional/folder]",
				ConflictsWith: []string{helpers.PPCImageId},
				RequiredWith:  []string{helpers.PPCImageBucketRegion, helpers.PPCImageBucketFileName},
				ForceNew:      true,
			},
			helpers.PPCImageBucketAccess: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Indicates if the bucket has public or private access",
				Default:       "public",
				ValidateFunc:  validate.ValidateAllowedStringValues([]string{"public", "private"}),
				ConflictsWith: []string{helpers.PPCImageId},
				ForceNew:      true,
			},
			helpers.PPCImageAccessKey: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Cloud Object Storage access key; required for buckets with private access",
				ForceNew:     true,
				Sensitive:    true,
				RequiredWith: []string{helpers.PPCImageSecretKey},
			},
			helpers.PPCImageSecretKey: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Cloud Object Storage secret key; required for buckets with private access",
				ForceNew:     true,
				Sensitive:    true,
				RequiredWith: []string{helpers.PPCImageAccessKey},
			},
			helpers.PPCImageBucketRegion: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Cloud Object Storage region",
				ConflictsWith: []string{helpers.PPCImageId},
				RequiredWith:  []string{helpers.PPCImageBucketName},
				ForceNew:      true,
			},
			helpers.PPCImageBucketFileName: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Cloud Object Storage image filename",
				ConflictsWith: []string{helpers.PPCImageId},
				RequiredWith:  []string{helpers.PPCImageBucketName},
				ForceNew:      true,
			},
			helpers.PPCImageStorageType: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of storage",
				ForceNew:    true,
			},
			helpers.PPCImageStoragePool: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Storage pool where the image will be loaded, if provided then ppc_image_storage_type and ppc_affinity_policy will be ignored",
				ForceNew:    true,
			},
			PPCAffinityPolicy: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Affinity policy for image; ignored if ppc_image_storage_pool provided; for policy affinity requires one of ppc_affinity_instance or ppc_affinity_volume to be specified; for policy anti-affinity requires one of ppc_anti_affinity_instances or ppc_anti_affinity_volumes to be specified",
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"affinity", "anti-affinity"}),
				ForceNew:     true,
			},
			PPCAffinityVolume: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Volume (ID or Name) to base storage affinity policy against; required if requesting affinity and ppc_affinity_instance is not provided",
				ConflictsWith: []string{PPCAffinityInstance},
				ForceNew:      true,
			},
			PPCAffinityInstance: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "PVM Instance (ID or Name) to base storage affinity policy against; required if requesting storage affinity and ppc_affinity_volume is not provided",
				ConflictsWith: []string{PPCAffinityVolume},
				ForceNew:      true,
			},
			PPCAntiAffinityVolumes: {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "List of volumes to base storage anti-affinity policy against; required if requesting anti-affinity and ppc_anti_affinity_instances is not provided",
				ConflictsWith: []string{PPCAntiAffinityInstances},
				ForceNew:      true,
			},
			PPCAntiAffinityInstances: {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "List of pvmInstances to base storage anti-affinity policy against; required if requesting anti-affinity and ppc_anti_affinity_volumes is not provided",
				ConflictsWith: []string{PPCAntiAffinityVolumes},
				ForceNew:      true,
			},

			// Computed Attribute
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image ID",
			},
		},
	}
}

func resourceIBMPPCImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		log.Printf("Failed to get the session")
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)
	imageName := d.Get(helpers.PPCImageName).(string)

	client := st.NewIBMPPCImageClient(ctx, sess, cloudInstanceID)
	// image copy
	if v, ok := d.GetOk(helpers.PPCImageId); ok {
		imageid := v.(string)
		source := "root-project"
		var body = &models.CreateImage{
			ImageName: imageName,
			ImageID:   imageid,
			Source:    &source,
		}
		imageResponse, err := client.Create(body)
		if err != nil {
			return diag.FromErr(err)
		}

		IBMPPCImageID := imageResponse.ImageID
		d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *IBMPPCImageID))

		_, err = isWaitForIBMPPCImageAvailable(ctx, client, *IBMPPCImageID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			log.Printf("[DEBUG]  err %s", err)
			return diag.FromErr(err)
		}
	}

	// COS image import
	if v, ok := d.GetOk(helpers.PPCImageBucketName); ok {
		bucketName := v.(string)
		bucketImageFileName := d.Get(helpers.PPCImageBucketFileName).(string)
		bucketRegion := d.Get(helpers.PPCImageBucketRegion).(string)
		bucketAccess := d.Get(helpers.PPCImageBucketAccess).(string)

		body := &models.CreateCosImageImportJob{
			ImageName:     &imageName,
			BucketName:    &bucketName,
			BucketAccess:  &bucketAccess,
			ImageFilename: &bucketImageFileName,
			Region:        &bucketRegion,
		}

		if v, ok := d.GetOk(helpers.PPCImageAccessKey); ok {
			body.AccessKey = v.(string)
		}
		if v, ok := d.GetOk(helpers.PPCImageSecretKey); ok {
			body.SecretKey = v.(string)
		}

		if v, ok := d.GetOk(helpers.PPCImageStorageType); ok {
			body.StorageType = v.(string)
		}
		if v, ok := d.GetOk(helpers.PPCImageStoragePool); ok {
			body.StoragePool = v.(string)
		}
		if ap, ok := d.GetOk(PPCAffinityPolicy); ok {
			policy := ap.(string)
			affinity := &models.StorageAffinity{
				AffinityPolicy: &policy,
			}

			if policy == "affinity" {
				if av, ok := d.GetOk(PPCAffinityVolume); ok {
					afvol := av.(string)
					affinity.AffinityVolume = &afvol
				}
				if ai, ok := d.GetOk(PPCAffinityInstance); ok {
					afins := ai.(string)
					affinity.AffinityPVMInstance = &afins
				}
			} else {
				if avs, ok := d.GetOk(PPCAntiAffinityVolumes); ok {
					afvols := flex.ExpandStringList(avs.([]interface{}))
					affinity.AntiAffinityVolumes = afvols
				}
				if ais, ok := d.GetOk(PPCAntiAffinityInstances); ok {
					afinss := flex.ExpandStringList(ais.([]interface{}))
					affinity.AntiAffinityPVMInstances = afinss
				}
			}
			body.StorageAffinity = affinity
		}
		imageResponse, err := client.CreateCosImage(body)
		if err != nil {
			return diag.FromErr(err)
		}

		jobClient := st.NewIBMPPCJobClient(ctx, sess, cloudInstanceID)
		_, err = waitForIBMPPCJobCompleted(ctx, jobClient, *imageResponse.ID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}

		// Once the job is completed find by name
		image, err := client.Get(imageName)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *image.ImageID))
	}

	return resourceIBMPPCImageRead(ctx, d, meta)
}

func resourceIBMPPCImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, imageID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	imageC := st.NewIBMPPCImageClient(ctx, sess, cloudInstanceID)
	imagedata, err := imageC.Get(imageID)
	if err != nil {
		uErr := errors.Unwrap(err)
		switch uErr.(type) {
		case *p_cloud_images.PcloudCloudinstancesImagesGetNotFound:
			log.Printf("[DEBUG] image does not exist %v", err)
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] get image failed %v", err)
		return diag.FromErr(err)
	}

	imageid := *imagedata.ImageID
	d.Set("image_id", imageid)
	d.Set(helpers.PPCCloudInstanceId, cloudInstanceID)

	return nil
}

func resourceIBMPPCImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, imageID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	imageC := st.NewIBMPPCImageClient(ctx, sess, cloudInstanceID)
	err = imageC.Delete(imageID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func isWaitForIBMPPCImageAvailable(ctx context.Context, client *st.IBMPPCImageClient, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Power Image (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PPCImageQueStatus},
		Target:     []string{helpers.PPCImageActiveStatus},
		Refresh:    isIBMPPCImageRefreshFunc(ctx, client, id),
		Timeout:    timeout,
		Delay:      20 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPPCImageRefreshFunc(ctx context.Context, client *st.IBMPPCImageClient, id string) resource.StateRefreshFunc {

	log.Printf("Calling the isIBMPPCImageRefreshFunc Refresh Function....")
	return func() (interface{}, string, error) {
		image, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if image.State == "active" {
			return image, helpers.PPCImageActiveStatus, nil
		}

		return image, helpers.PPCImageQueStatus, nil
	}
}

func waitForIBMPPCJobCompleted(ctx context.Context, client *st.IBMPPCJobClient, jobID string, timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{helpers.JobStatusQueued, helpers.JobStatusReadyForProcessing, helpers.JobStatusInProgress, helpers.JobStatusRunning, helpers.JobStatusWaiting},
		Target:  []string{helpers.JobStatusCompleted, helpers.JobStatusFailed},
		Refresh: func() (interface{}, string, error) {
			job, err := client.Get(jobID)
			if err != nil {
				log.Printf("[DEBUG] get job failed %v", err)
				return nil, "", fmt.Errorf(errors.GetJobOperationFailed, jobID, err)
			}
			if job == nil || job.Status == nil {
				log.Printf("[DEBUG] get job failed with empty response")
				return nil, "", fmt.Errorf("failed to get job status for job id %s", jobID)
			}
			if *job.Status.State == helpers.JobStatusFailed {
				log.Printf("[DEBUG] job status failed with message: %v", job.Status.Message)
				return nil, helpers.JobStatusFailed, fmt.Errorf("job status failed for job id %s with message: %v", jobID, job.Status.Message)
			}
			return job, *job.Status.State, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}
