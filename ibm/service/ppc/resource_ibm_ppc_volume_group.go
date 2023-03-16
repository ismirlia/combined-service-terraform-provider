// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
	"github.com/IBM-Cloud/ppc-aas-go-client/ppc-aas/client/p_cloud_volume_groups"
	"github.com/IBM-Cloud/ppc-aas-go-client/ppc-aas/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMPPCVolumeGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPPCVolumeGroupCreate,
		ReadContext:   resourceIBMPPCVolumeGroupRead,
		UpdateContext: resourceIBMPPCVolumeGroupUpdate,
		DeleteContext: resourceIBMPPCVolumeGroupDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			helpers.PPCCloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloud Instance ID - This is the service_instance_id.",
			},
			PPCVolumeGroupName: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Volume Group Name to create",
				ConflictsWith: []string{PPCVolumeGroupConsistencyGroupName},
			},
			PPCVolumeGroupConsistencyGroupName: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The name of consistency group at storage controller level",
				ConflictsWith: []string{PPCVolumeGroupName},
			},
			PPCVolumeGroupsVolumeIds: {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of volumes to add in volume group",
			},

			// Computed Attributes
			"volume_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume Group ID",
			},
			"volume_group_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume Group Status",
			},
			"replication_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume Group Replication Status",
			},
			"status_description_errors": vgStatusDescriptionErrors(),
			"consistency_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Consistency Group Name if volume is a part of volume group",
			},
		},
	}
}

func resourceIBMPPCVolumeGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	vgName := d.Get(PPCVolumeGroupName).(string)
	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)
	body := &models.VolumeGroupCreate{
		Name: vgName,
	}

	volids := flex.ExpandStringList((d.Get(PPCVolumeGroupsVolumeIds).(*schema.Set)).List())
	body.VolumeIDs = volids

	if v, ok := d.GetOk(PPCVolumeGroupConsistencyGroupName); ok {
		body.ConsistencyGroupName = v.(string)
	}

	client := st.NewIBMPPCVolumeGroupClient(ctx, sess, cloudInstanceID)
	vg, err := client.CreateVolumeGroup(body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *vg.ID))

	_, err = isWaitForIBMPPCVolumeGroupAvailable(ctx, client, *vg.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIBMPPCVolumeGroupRead(ctx, d, meta)
}

func resourceIBMPPCVolumeGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, vgID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := st.NewIBMPPCVolumeGroupClient(ctx, sess, cloudInstanceID)

	vg, err := client.GetDetails(vgID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("volume_group_id", vg.ID)
	d.Set("volume_group_status", vg.Status)
	d.Set("consistency_group_name", vg.ConsistencyGroupName)
	d.Set("replication_status", vg.ReplicationStatus)
	d.Set(PPCVolumeGroupName, vg.Name)
	d.Set(PPCVolumeGroupsVolumeIds, vg.VolumeIDs)
	d.Set("status_description_errors", flattenVolumeGroupStatusDescription(vg.StatusDescription.Errors))

	return nil
}

func resourceIBMPPCVolumeGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, vgID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := st.NewIBMPPCVolumeGroupClient(ctx, sess, cloudInstanceID)
	if d.HasChanges(PPCVolumeGroupsVolumeIds) {
		old, new := d.GetChange(PPCVolumeGroupsVolumeIds)
		oldList := old.(*schema.Set)
		newList := new.(*schema.Set)
		body := &models.VolumeGroupUpdate{
			AddVolumes:    flex.ExpandStringList(newList.Difference(oldList).List()),
			RemoveVolumes: flex.ExpandStringList(oldList.Difference(newList).List()),
		}
		err := client.UpdateVolumeGroup(vgID, body)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = isWaitForIBMPPCVolumeGroupAvailable(ctx, client, vgID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIBMPPCVolumeGroupRead(ctx, d, meta)
}
func resourceIBMPPCVolumeGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, vgID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := st.NewIBMPPCVolumeGroupClient(ctx, sess, cloudInstanceID)

	volids := flex.ExpandStringList((d.Get(PPCVolumeGroupsVolumeIds).(*schema.Set)).List())
	if len(volids) > 0 {
		body := &models.VolumeGroupUpdate{
			RemoveVolumes: volids,
		}
		err = client.UpdateVolumeGroup(vgID, body)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = isWaitForIBMPPCVolumeGroupAvailable(ctx, client, vgID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = client.DeleteVolumeGroup(vgID)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = isWaitForIBMPPCVolumeGroupDeleted(ctx, client, vgID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
func isWaitForIBMPPCVolumeGroupAvailable(ctx context.Context, client *st.IBMPPCVolumeGroupClient, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Volume Group (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PPCVolumeProvisioning},
		Target:     []string{helpers.PPCVolumeProvisioningDone},
		Refresh:    isIBMPPCVolumeGroupRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 2 * time.Minute,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPPCVolumeGroupRefreshFunc(client *st.IBMPPCVolumeGroupClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vg, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if vg.Status == "available" {
			return vg, helpers.PPCVolumeProvisioningDone, nil
		}

		return vg, helpers.PPCVolumeProvisioning, nil
	}
}

func isWaitForIBMPPCVolumeGroupDeleted(ctx context.Context, client *st.IBMPPCVolumeGroupClient, id string, timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting", "updating"},
		Target:     []string{"deleted"},
		Refresh:    isIBMPPCVolumeGroupDeleteRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 2 * time.Minute,
		Timeout:    timeout,
	}
	return stateConf.WaitForStateContext(ctx)
}

func isIBMPPCVolumeGroupDeleteRefreshFunc(client *st.IBMPPCVolumeGroupClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vg, err := client.Get(id)
		if err != nil {
			uErr := errors.Unwrap(err)
			switch uErr.(type) {
			case *p_cloud_volume_groups.PcloudVolumegroupsGetNotFound:
				log.Printf("[DEBUG] volume-group does not exist while deleteing %v", err)
				return vg, "deleted", nil
			}
			return nil, "", err
		}
		if vg == nil {
			return vg, "deleted", nil
		}
		return vg, "deleting", nil
	}
}
