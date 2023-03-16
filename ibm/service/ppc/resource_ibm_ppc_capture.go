// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
	"github.com/IBM-Cloud/ppc-aas-go-client/ppc-aas/client/p_cloud_images"
	"github.com/IBM-Cloud/ppc-aas-go-client/ppc-aas/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const cloudStorageDestination string = "cloud-storage"
const imageCatalogDestination string = "image-catalog"

func ResourceIBMPPCCapture() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPPCCaptureCreate,
		ReadContext:   resourceIBMPPCCaptureRead,
		DeleteContext: resourceIBMPPCCaptureDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(75 * time.Minute),
			Delete: schema.DefaultTimeout(50 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			helpers.PPCCloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: " Cloud Instance ID - This is the service_instance_id.",
			},

			helpers.PPCInstanceName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance Name of the Power VM",
			},

			helpers.PPCInstanceCaptureName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the capture to create. Note : this must be unique",
			},

			helpers.PPCInstanceCaptureDestination: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Destination for the deployable image",
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"image-catalog", "cloud-storage", "both"}),
			},

			helpers.PPCInstanceCaptureVolumeIds: {
				Type:             schema.TypeSet,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Set:              schema.HashString,
				ForceNew:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "List of Data volume IDs",
			},

			helpers.PPCInstanceCaptureCloudStorageRegion: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "List of Regions to use",
			},

			helpers.PPCInstanceCaptureCloudStorageAccessKey: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "Name of Cloud Storage Access Key",
			},
			helpers.PPCInstanceCaptureCloudStorageSecretKey: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "Name of the Cloud Storage Secret Key",
			},
			helpers.PPCInstanceCaptureCloudStorageImagePath: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Cloud Storage Image Path (bucket-name [/folder/../..])",
			},
			// Computed Attribute
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image ID of Capture Instance",
			},
		},
	}
}

func resourceIBMPPCCaptureCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get(helpers.PPCInstanceName).(string)
	capturename := d.Get(helpers.PPCInstanceCaptureName).(string)
	capturedestination := d.Get(helpers.PPCInstanceCaptureDestination).(string)
	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)

	client := st.NewIBMPPCInstanceClient(context.Background(), sess, cloudInstanceID)

	captureBody := &models.PVMInstanceCapture{
		CaptureDestination: &capturedestination,
		CaptureName:        &capturename,
	}
	if capturedestination != imageCatalogDestination {
		if v, ok := d.GetOk(helpers.PPCInstanceCaptureCloudStorageRegion); ok {
			captureBody.CloudStorageRegion = v.(string)
		} else {
			return diag.Errorf("%s is required when capture destination is %s", helpers.PPCInstanceCaptureCloudStorageRegion, capturedestination)
		}
		if v, ok := d.GetOk(helpers.PPCInstanceCaptureCloudStorageAccessKey); ok {
			captureBody.CloudStorageAccessKey = v.(string)
		} else {
			return diag.Errorf("%s is required when capture destination is %s ", helpers.PPCInstanceCaptureCloudStorageAccessKey, capturedestination)
		}
		if v, ok := d.GetOk(helpers.PPCInstanceCaptureCloudStorageImagePath); ok {
			captureBody.CloudStorageImagePath = v.(string)
		} else {
			return diag.Errorf("%s is required when capture destination is %s ", helpers.PPCInstanceCaptureCloudStorageImagePath, capturedestination)
		}
		if v, ok := d.GetOk(helpers.PPCInstanceCaptureCloudStorageSecretKey); ok {
			captureBody.CloudStorageSecretKey = v.(string)
		} else {
			return diag.Errorf("%s is required when capture destination is %s ", helpers.PPCInstanceCaptureCloudStorageSecretKey, capturedestination)
		}
	}

	if v, ok := d.GetOk(helpers.PPCInstanceCaptureVolumeIds); ok {
		volids := flex.ExpandStringList((v.(*schema.Set)).List())
		if len(volids) > 0 {
			captureBody.CaptureVolumeIDs = volids
		}
	}

	captureResponse, err := client.CaptureInstanceToImageCatalogV2(name, captureBody)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", cloudInstanceID, capturename, capturedestination))
	jobClient := st.NewIBMPPCJobClient(ctx, sess, cloudInstanceID)
	_, err = waitForIBMPPCJobCompleted(ctx, jobClient, *captureResponse.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceIBMPPCCaptureRead(ctx, d, meta)
}

func resourceIBMPPCCaptureRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := parts[0]
	captureID := parts[1]
	capturedestination := parts[2]
	if capturedestination != cloudStorageDestination {
		imageClient := st.NewIBMPPCImageClient(ctx, sess, cloudInstanceID)
		imagedata, err := imageClient.Get(captureID)
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
	}
	d.Set(helpers.PPCCloudInstanceId, cloudInstanceID)
	return nil
}

func resourceIBMPPCCaptureDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := parts[0]
	captureID := parts[1]
	capturedestination := parts[2]
	if capturedestination != cloudStorageDestination {
		imageClient := st.NewIBMPPCImageClient(ctx, sess, cloudInstanceID)
		err = imageClient.Delete(captureID)
		if err != nil {
			uErr := errors.Unwrap(err)
			switch uErr.(type) {
			case *p_cloud_images.PcloudCloudinstancesImagesGetNotFound:
				log.Printf("[DEBUG] image does not exist while deleting %v", err)
				d.SetId("")
				return nil
			}
			log.Printf("[DEBUG] delete image failed %v", err)
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return nil
}
