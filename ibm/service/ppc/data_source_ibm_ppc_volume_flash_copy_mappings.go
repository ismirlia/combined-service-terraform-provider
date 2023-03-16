// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPPCVolumeFlashCopyMappings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPPCVolumeFlashCopyMappings,
		Schema: map[string]*schema.Schema{
			helpers.PPCVolumeId: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Volume ID",
				ValidateFunc: validation.NoZeroValues,
			},
			helpers.PPCCloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes
			"flash_copy_mappings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"copy_rate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicates the rate of flash copy operation of a volume",
						},
						"flash_copy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates flash copy name of the volume",
						},
						"progress": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicates the progress of flash copy operation",
						},
						"source_volume_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates name of the source volume",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the start time of flash copy operation",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Copy status of a volume",
						},
						"target_volume_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates name of the target volume",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPPCVolumeFlashCopyMappings(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)
	volClient := instance.NewIBMPPCVolumeClient(ctx, sess, cloudInstanceID)
	volData, err := volClient.GetVolumeFlashCopyMappings(d.Get(helpers.PPCVolumeId).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	results := make([]map[string]interface{}, 0, len(volData))
	for _, i := range volData {
		if i != nil {
			l := map[string]interface{}{
				"copy_rate":          i.CopyRate,
				"progress":           i.Progress,
				"source_volume_name": i.SourceVolumeName,
				"start_time":         i.StartTime.String(),
				"status":             i.Status,
				"target_volume_name": i.TargetVolumeName,
			}
			if i.FlashCopyName != nil {
				l["flash_copy_name"] = i.FlashCopyName
			}
			results = append(results, l)
		}
	}
	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set("flash_copy_mappings", results)

	return nil
}
