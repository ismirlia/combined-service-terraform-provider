// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc

import (
	"context"
	"fmt"

	"log"

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"

	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	PPCPoolName = "ppc_storage_pool"
)

func DataSourceIBMPPCStoragePoolCapacity() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPPCStoragePoolCapacityRead,
		Schema: map[string]*schema.Schema{
			helpers.PPCCloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			PPCPoolName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Storage pool name",
			},
			// Computed Attributes
			MaxAllocationSize: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum allocation storage size (GB)",
			},
			StorageType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Storage type of the storage pool",
			},
			TotalCapacity: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total pool capacity (GB)",
			},
		},
	}
}

func dataSourceIBMPPCStoragePoolCapacityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)
	storagePool := d.Get(PPCPoolName).(string)

	client := st.NewIBMPPCStorageCapacityClient(ctx, sess, cloudInstanceID)
	sp, err := client.GetStoragePoolCapacity(storagePool)
	if err != nil {
		log.Printf("[ERROR] get storage pool capacity failed %v", err)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, storagePool))

	d.Set(MaxAllocationSize, *sp.MaxAllocationSize)
	d.Set(StorageType, sp.StorageType)
	d.Set(TotalCapacity, sp.TotalCapacity)

	return nil
}
