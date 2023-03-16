// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
	"github.com/IBM-Cloud/ppc-aas-go-client/ppc-aas/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
)

func ResourceIBMPPCImageExport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPPCImageExportCreate,
		ReadContext:   resourceIBMPPCImageExportRead,
		DeleteContext: resourceIBMPPCImageExportDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			//required attributes
			helpers.PPCCloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PPC cloud instance ID",
				ForceNew:    true,
			},
			helpers.PPCImageId: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Instance image id",
				DiffSuppressFunc: flex.ApplyOnce,
				ForceNew:         true,
			},
			helpers.PPCImageBucketName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloud Object Storage bucket name; bucket-name[/optional/folder]",
				ForceNew:    true,
			},
			helpers.PPCImageAccessKey: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloud Object Storage access key; required for buckets with private access",
				Sensitive:   true,
				ForceNew:    true,
			},

			helpers.PPCImageSecretKey: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloud Object Storage secret key; required for buckets with private access",
				Sensitive:   true,
				ForceNew:    true,
			},
			helpers.PPCImageBucketRegion: {
				Type:        schema.TypeString,
				Description: "Cloud Object Storage region",
				ForceNew:    true,
				Required:    true,
			},
		},
	}
}

func resourceIBMPPCImageExportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		log.Printf("Failed to get the session")
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)
	imageid := d.Get(helpers.PPCImageId).(string)
	bucketName := d.Get(helpers.PPCImageBucketName).(string)
	accessKey := d.Get(helpers.PPCImageAccessKey).(string)

	client := st.NewIBMPPCImageClient(ctx, sess, cloudInstanceID)

	// image export
	var body = &models.ExportImage{
		BucketName: &bucketName,
		AccessKey:  &accessKey,
		Region:     d.Get(helpers.PPCImageBucketRegion).(string),
		SecretKey:  d.Get(helpers.PPCImageSecretKey).(string),
	}

	imageResponse, err := client.ExportImage(imageid, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", imageid, bucketName, d.Get(helpers.PPCImageBucketRegion).(string)))

	jobClient := st.NewIBMPPCJobClient(ctx, sess, cloudInstanceID)
	_, err = waitForIBMPPCJobCompleted(ctx, jobClient, *imageResponse.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceIBMPPCImageExportRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceIBMPPCImageExportDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
