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

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMPPCKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPPCKeysRead,
		Schema: map[string]*schema.Schema{

			// Arguments
			Arg_CloudInstanceID: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "PPC cloud instance ID",
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_Keys: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "SSH Keys",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_KeyName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User defined name for the SSH key",
						},
						Attr_Key: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SSH RSA key",
						},
						Attr_KeyCreationDate: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date of SSH key creation",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPPCKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// session
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	// arguments
	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)

	// get keys
	client := st.NewIBMPPCKeyClient(ctx, sess, cloudInstanceID)
	sshKeys, err := client.GetAll()
	if err != nil {
		log.Printf("[ERROR] get all keys failed %v", err)
		return diag.FromErr(err)
	}

	// set attributes
	result := make([]map[string]interface{}, 0, len(sshKeys.SSHKeys))
	for _, sshKey := range sshKeys.SSHKeys {
		key := map[string]interface{}{
			Attr_KeyName:         sshKey.Name,
			Attr_Key:             sshKey.SSHKey,
			Attr_KeyCreationDate: sshKey.CreationDate.String(),
		}
		result = append(result, key)
	}
	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(Attr_Keys, result)

	return nil
}
