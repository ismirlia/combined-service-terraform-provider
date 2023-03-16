// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc

import (
	"context"
	"log"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPPCSAPProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPPCSAPProfilesRead,
		Schema: map[string]*schema.Schema{
			helpers.PPCCloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			// Computed Attributes
			PPCSAPProfiles: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						PPCSAPProfileID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SAP Profile ID",
						},
						PPCSAPProfileType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of profile",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPPCSAPProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)

	client := instance.NewIBMPPCSAPInstanceClient(ctx, sess, cloudInstanceID)
	sapProfiles, err := client.GetAllSAPProfiles(cloudInstanceID)
	if err != nil {
		log.Printf("[DEBUG] get all sap profiles failed %v", err)
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, 0, len(sapProfiles.Profiles))
	for _, sapProfile := range sapProfiles.Profiles {
		profile := map[string]interface{}{
			PPCSAPProfileCertified: *sapProfile.Certified,
			PPCSAPProfileCores:     *sapProfile.Cores,
			PPCSAPProfileMemory:    *sapProfile.Memory,
			PPCSAPProfileID:        *sapProfile.ProfileID,
			PPCSAPProfileType:      *sapProfile.Type,
		}
		result = append(result, profile)
	}

	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(PPCSAPProfiles, result)

	return nil
}
