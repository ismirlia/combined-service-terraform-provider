// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
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

func DataSourceIBMPPCKey() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceIBMPPCKeyRead,
		Schema: map[string]*schema.Schema{

			// Arguments
			Arg_KeyName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "SSH key name for a pcloud tenant",
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_CloudInstanceID: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_KeyCreationDate: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date of sshkey creation",
			},
			Attr_Key: {
				Type:        schema.TypeString,
				Sensitive:   true,
				Computed:    true,
				Description: "SSH RSA key",
			},
			"sshkey": {
				Type:       schema.TypeString,
				Sensitive:  true,
				Computed:   true,
				Deprecated: "This field is deprecated, use ssh_key instead",
			},
		},
	}
}

func dataSourceIBMPPCKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// session
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	// arguments
	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)

	// get key
	sshkeyC := instance.NewIBMPPCKeyClient(ctx, sess, cloudInstanceID)
	sshkeydata, err := sshkeyC.Get(d.Get(helpers.PPCKeyName).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	// set attributes
	d.SetId(*sshkeydata.Name)
	d.Set(Attr_KeyCreationDate, sshkeydata.CreationDate.String())
	d.Set(Attr_Key, sshkeydata.SSHKey)
	d.Set("sshkey", sshkeydata.SSHKey) // TODO: deprecated, to remove

	return nil
}
