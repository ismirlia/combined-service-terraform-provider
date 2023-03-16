// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPPCPlacementGroup() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceIBMPPCPlacementGroupRead,
		Schema: map[string]*schema.Schema{
			helpers.PPCPlacementGroupName: {
				Type:     schema.TypeString,
				Required: true,
			},

			"policy": {
				Type:     schema.TypeString,
				Computed: true,
			},

			helpers.PPCCloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			PPCPlacementGroupMembers: {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceIBMPPCPlacementGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)
	placementGroupName := d.Get(helpers.PPCPlacementGroupName).(string)
	client := st.NewIBMPPCPlacementGroupClient(ctx, sess, cloudInstanceID)

	response, err := client.Get(placementGroupName)
	if err != nil {
		log.Printf("[DEBUG]  err %s", err)
		return diag.FromErr(err)
	}

	d.SetId(*response.ID)
	d.Set("policy", response.Policy)
	d.Set(PPCPlacementGroupMembers, response.Members)

	return nil
}
