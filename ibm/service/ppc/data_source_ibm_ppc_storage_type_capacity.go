// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc

import (
	"context"
	"fmt"

	"log"

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"

	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	PPCTypeName = "ppc_storage_type"
)

func DataSourceIBMPPCStorageTypeCapacity() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPPCStorageTypeCapacityRead,
		Schema: map[string]*schema.Schema{
			helpers.PPCCloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			PPCTypeName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "Storage type name",
			},
			// Computed Attributes
			MaximumStorageAllocation: {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Maximum storage allocation",
			},
			StoragePoolsCapacity: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Storage pools capacity",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						MaxAllocationSize: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum allocation storage size (GB)",
						},
						PoolName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pool name",
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
				},
			},
		},
	}
}

func dataSourceIBMPPCStorageTypeCapacityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)
	storageType := d.Get(PPCTypeName).(string)

	client := st.NewIBMPPCStorageCapacityClient(ctx, sess, cloudInstanceID)
	stc, err := client.GetStorageTypeCapacity(storageType)
	if err != nil {
		log.Printf("[ERROR] get storage type capacity failed %v", err)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, storageType))

	if stc.MaximumStorageAllocation != nil {
		msa := stc.MaximumStorageAllocation
		data := map[string]interface{}{
			MaxAllocationSize: *msa.MaxAllocationSize,
			StoragePool:       *msa.StoragePool,
			StorageType:       *msa.StorageType,
		}
		d.Set(MaximumStorageAllocation, flex.Flatten(data))
	}

	result := make([]map[string]interface{}, 0, len(stc.StoragePoolsCapacity))
	for _, sp := range stc.StoragePoolsCapacity {
		data := map[string]interface{}{
			MaxAllocationSize: *sp.MaxAllocationSize,
			PoolName:          sp.PoolName,
			StorageType:       sp.StorageType,
			TotalCapacity:     sp.TotalCapacity,
		}
		result = append(result, data)
	}
	d.Set(StoragePoolsCapacity, result)

	return nil
}
