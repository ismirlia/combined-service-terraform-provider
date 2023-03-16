// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
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
	"github.com/IBM-Cloud/ppc-aas-go-client/ppc-aas/client/p_cloud_volumes"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMPPCVolumeAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPPCVolumeAttachCreate,
		ReadContext:   resourceIBMPPCVolumeAttachRead,
		DeleteContext: resourceIBMPPCVolumeAttachDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			helpers.PPCCloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: " Cloud Instance ID - This is the service_instance_id.",
			},

			helpers.PPCVolumeId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the volume to attach. Note these volumes should have been created",
			},

			helpers.PPCInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "PPC Instance Id",
			},

			// Computed Attribute
			helpers.PPCVolumeAttachStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIBMPPCVolumeAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	volumeID := d.Get(helpers.PPCVolumeId).(string)
	pvmInstanceID := d.Get(helpers.PPCInstanceId).(string)
	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)

	volClient := st.NewIBMPPCVolumeClient(ctx, sess, cloudInstanceID)
	volinfo, err := volClient.Get(volumeID)
	if err != nil {
		return diag.FromErr(err)
	}

	if volinfo.State == "available" || *volinfo.Shareable {
		log.Printf(" In the current state the volume can be attached to the instance ")
	}

	if volinfo.State == "in-use" && *volinfo.Shareable {

		log.Printf("Volume State /Status is  permitted and hence attaching the volume to the instance")
	}

	if volinfo.State == helpers.PPCVolumeAllowableAttachStatus && !*volinfo.Shareable {
		return diag.Errorf("the volume cannot be attached in the current state. The volume must be in the *available* state. No other states are permissible")
	}

	err = volClient.Attach(pvmInstanceID, volumeID)
	if err != nil {
		log.Printf("[DEBUG]  err %s", err)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", cloudInstanceID, pvmInstanceID, *volinfo.VolumeID))

	_, err = isWaitForIBMPPCVolumeAttachAvailable(ctx, volClient, *volinfo.VolumeID, cloudInstanceID, pvmInstanceID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIBMPPCVolumeAttachRead(ctx, d, meta)
}

func resourceIBMPPCVolumeAttachRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	ids, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID, pvmInstanceID, volumeID := ids[0], ids[1], ids[2]

	client := st.NewIBMPPCVolumeClient(ctx, sess, cloudInstanceID)

	vol, err := client.CheckVolumeAttach(pvmInstanceID, volumeID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(helpers.PPCVolumeAttachStatus, vol.State)
	return nil
}

func resourceIBMPPCVolumeAttachDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}
	ids, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID, pvmInstanceID, volumeID := ids[0], ids[1], ids[2]
	client := st.NewIBMPPCVolumeClient(ctx, sess, cloudInstanceID)

	log.Printf("the id of the volume to detach is %s ", volumeID)

	err = client.Detach(pvmInstanceID, volumeID)
	if err != nil {
		uErr := errors.Unwrap(err)
		switch uErr.(type) {
		case *p_cloud_volumes.PcloudCloudinstancesVolumesGetNotFound:
			log.Printf("[DEBUG] volume does not exist while detaching %v", err)
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] volume detach failed %v", err)
		return diag.FromErr(err)
	}

	_, err = isWaitForIBMPPCVolumeDetach(ctx, client, volumeID, cloudInstanceID, pvmInstanceID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	// wait for power volume states to be back as available. if it's attached it will be in-use
	d.SetId("")
	return nil
}

func isWaitForIBMPPCVolumeAttachAvailable(ctx context.Context, client *st.IBMPPCVolumeClient, id, cloudInstanceID, pvmInstanceID string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Volume (%s) to be available for attachment", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PPCVolumeProvisioning},
		Target:     []string{helpers.PPCVolumeAllowableAttachStatus},
		Refresh:    isIBMPPCVolumeAttachRefreshFunc(client, id, cloudInstanceID, pvmInstanceID),
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPPCVolumeAttachRefreshFunc(client *st.IBMPPCVolumeClient, id, cloudInstanceID, pvmInstanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vol, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if vol.State == "in-use" && flex.StringContains(vol.PvmInstanceIDs, pvmInstanceID) {
			return vol, helpers.PPCVolumeAllowableAttachStatus, nil
		}

		return vol, helpers.PPCVolumeProvisioning, nil
	}
}

func isWaitForIBMPPCVolumeDetach(ctx context.Context, client *st.IBMPPCVolumeClient, id, cloudInstanceID, pvmInstanceID string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Volume (%s) to be available after detachment", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"detaching", helpers.PowerVolumeAttachDeleting},
		Target:     []string{helpers.PPCVolumeProvisioningDone},
		Refresh:    isIBMPPCVolumeDetachRefreshFunc(client, id, cloudInstanceID, pvmInstanceID),
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPPCVolumeDetachRefreshFunc(client *st.IBMPPCVolumeClient, id, cloudInstanceID, pvmInstanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vol, err := client.Get(id)
		if err != nil {
			uErr := errors.Unwrap(err)
			switch uErr.(type) {
			case *p_cloud_volumes.PcloudCloudinstancesVolumesGetNotFound:
				log.Printf("[DEBUG] volume does not exist while detaching %v", err)
				return vol, helpers.PPCVolumeProvisioningDone, nil
			}
			return nil, "", err
		}

		// Check if Instance ID is in the Volume's Instance list
		// Also validate the Volume state is 'available' when it is not Sharable
		// In case of Sharable Volume it can be `in-use` state
		if !flex.StringContains(vol.PvmInstanceIDs, pvmInstanceID) &&
			(*vol.Shareable || (!*vol.Shareable && vol.State == "available")) {
			return vol, helpers.PPCVolumeProvisioningDone, nil
		}

		return vol, "detaching", nil
	}
}
