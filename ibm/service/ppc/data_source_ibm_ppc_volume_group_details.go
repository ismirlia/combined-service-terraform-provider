// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPPCVolumeGroupDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPPCVolumeGroupDetailsRead,
		Schema: map[string]*schema.Schema{
			PPCVolumeGroupID: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the volume group",
				ValidateFunc: validation.NoZeroValues,
			},
			helpers.PPCCloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes
			"volume_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume group name",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replication_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"consistency_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_description_errors": vgStatusDescriptionErrors(),
			"volume_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceIBMPPCVolumeGroupDetailsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)
	vgClient := instance.NewIBMPPCVolumeGroupClient(ctx, sess, cloudInstanceID)
	vgData, err := vgClient.GetDetails(d.Get(PPCVolumeGroupID).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*vgData.ID)
	d.Set("status", vgData.Status)
	d.Set("consistency_group_name", vgData.ConsistencyGroupName)
	d.Set("volume_group_name", vgData.Name)
	d.Set("replication_status", vgData.ReplicationStatus)
	d.Set("volume_ids", vgData.VolumeIDs)
	if vgData.StatusDescription != nil {
		d.Set("status_description_errors", flattenVolumeGroupStatusDescription(vgData.StatusDescription.Errors))
	}

	return nil
}
