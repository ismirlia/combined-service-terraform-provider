// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPPCSAPProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPPCSAPProfileRead,
		Schema: map[string]*schema.Schema{
			helpers.PPCCloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			PPCSAPInstanceProfileID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SAP Profile ID",
			},
			// Computed Attributes
			PPCSAPProfileCertified: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Has certification been performed on profile",
			},
			PPCSAPProfileCores: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Amount of cores",
			},
			PPCSAPProfileMemory: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Amount of memory (in GB)",
			},
			PPCSAPProfileType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of profile",
			},
		},
	}
}

func dataSourceIBMPPCSAPProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)
	profileID := d.Get(PPCSAPInstanceProfileID).(string)

	client := instance.NewIBMPPCSAPInstanceClient(ctx, sess, cloudInstanceID)
	sapProfile, err := client.GetSAPProfile(profileID)
	if err != nil {
		log.Printf("[DEBUG] get sap profile failed %v", err)
		return diag.FromErr(err)
	}

	d.SetId(*sapProfile.ProfileID)
	d.Set(PPCSAPProfileCertified, *sapProfile.Certified)
	d.Set(PPCSAPProfileCores, *sapProfile.Cores)
	d.Set(PPCSAPProfileMemory, *sapProfile.Memory)
	d.Set(PPCSAPProfileType, *sapProfile.Type)

	return nil
}
