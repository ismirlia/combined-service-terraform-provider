// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
	models "github.com/IBM-Cloud/ppc-aas-go-client/ppc-aas/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func ResourceIBMPPCPlacementGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPPCPlacementGroupCreate,
		ReadContext:   resourceIBMPPCPlacementGroupRead,
		UpdateContext: resourceIBMPPCPlacementGroupUpdate,
		DeleteContext: resourceIBMPPCPlacementGroupDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			helpers.PPCPlacementGroupName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the placement group",
			},

			helpers.PPCPlacementGroupPolicy: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"affinity", "anti-affinity"}),
				Description:  "Policy of the placement group",
			},

			helpers.PPCCloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PPC cloud instance ID",
			},

			PPCPlacementGroupMembers: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Server IDs that are the placement group members",
			},

			PPCPlacementGroupID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "PPC placement group ID",
			},
		},
	}
}

func resourceIBMPPCPlacementGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)
	name := d.Get(helpers.PPCPlacementGroupName).(string)
	policy := d.Get(helpers.PPCPlacementGroupPolicy).(string)
	client := st.NewIBMPPCPlacementGroupClient(ctx, sess, cloudInstanceID)
	body := &models.PlacementGroupCreate{
		Name:   &name,
		Policy: &policy,
	}

	response, err := client.Create(body)
	if err != nil || response == nil {
		return diag.FromErr(fmt.Errorf("error creating the shared processor pool: %s", err))
	}

	log.Printf("Printing the placement group %+v", &response)

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *response.ID))
	return resourceIBMPPCPlacementGroupRead(ctx, d, meta)
}

func resourceIBMPPCPlacementGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := parts[0]
	client := st.NewIBMPPCPlacementGroupClient(ctx, sess, cloudInstanceID)

	response, err := client.Get(parts[1])
	if err != nil {
		log.Printf("[DEBUG]  err %s", err)
		return diag.FromErr(err)
	}

	d.Set(helpers.PPCPlacementGroupName, response.Name)
	d.Set(PPCPlacementGroupID, response.ID)
	d.Set(helpers.PPCPlacementGroupPolicy, response.Policy)
	d.Set(PPCPlacementGroupMembers, response.Members)

	return nil

}

func resourceIBMPPCPlacementGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceIBMPPCPlacementGroupRead(ctx, d, meta)
}

func resourceIBMPPCPlacementGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := parts[0]
	client := st.NewIBMPPCPlacementGroupClient(ctx, sess, cloudInstanceID)
	err = client.Delete(parts[1])

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
