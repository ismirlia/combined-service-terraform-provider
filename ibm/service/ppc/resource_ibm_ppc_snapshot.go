// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
	"github.com/IBM-Cloud/ppc-aas-go-client/ppc-aas/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
)

func ResourceIBMPPCSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPPCSnapshotCreate,
		ReadContext:   resourceIBMPPCSnapshotRead,
		UpdateContext: resourceIBMPPCSnapshotUpdate,
		DeleteContext: resourceIBMPPCSnapshotDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			helpers.PPCSnapshotName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique name of the snapshot",
			},
			helpers.PPCInstanceName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance name / id of the pvm",
			},
			helpers.PPCInstanceVolumeIds: {
				Type:             schema.TypeSet,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Set:              schema.HashString,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "List of PPC volumes",
			},
			helpers.PPCCloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: " Cloud Instance ID - This is the service_instance_id.",
			},
			"ppc_description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the PVM instance snapshot",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Snapshot description",
				Deprecated:  "This field is deprecated, use ppc_description instead",
			},

			// Computed Attributes
			helpers.PPCSnapshot: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id of the snapshot",
				Deprecated:  "This field is deprecated, use snapshot_id instead",
			},
			"snapshot_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the PVM instance snapshot",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_snapshots": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceIBMPPCSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)
	instanceid := d.Get(helpers.PPCInstanceName).(string)
	volids := flex.ExpandStringList((d.Get(helpers.PPCInstanceVolumeIds).(*schema.Set)).List())
	name := d.Get(helpers.PPCSnapshotName).(string)

	var description string
	if v, ok := d.GetOk("description"); ok {
		description = v.(string)
	}
	if v, ok := d.GetOk("ppc_description"); ok {
		description = v.(string)
	}

	client := st.NewIBMPPCInstanceClient(ctx, sess, cloudInstanceID)

	snapshotBody := &models.SnapshotCreate{Name: &name, Description: description}

	if len(volids) > 0 {
		snapshotBody.VolumeIDs = volids
	} else {
		log.Printf("no volumeids provided. Will snapshot the entire instance")
	}

	snapshotResponse, err := client.CreatePvmSnapShot(instanceid, snapshotBody)
	if err != nil {
		log.Printf("[DEBUG]  err %s", err)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *snapshotResponse.SnapshotID))

	pisnapclient := st.NewIBMPPCSnapshotClient(ctx, sess, cloudInstanceID)
	_, err = isWaitForPPCInstanceSnapshotAvailable(ctx, pisnapclient, *snapshotResponse.SnapshotID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIBMPPCSnapshotRead(ctx, d, meta)
}

func resourceIBMPPCSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("Calling the Snapshot Read function post create")
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, snapshotID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	snapshot := st.NewIBMPPCSnapshotClient(ctx, sess, cloudInstanceID)
	snapshotdata, err := snapshot.Get(snapshotID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(helpers.PPCSnapshotName, snapshotdata.Name)
	d.Set(helpers.PPCSnapshot, *snapshotdata.SnapshotID)
	d.Set("snapshot_id", *snapshotdata.SnapshotID)
	d.Set("status", snapshotdata.Status)
	d.Set("creation_date", snapshotdata.CreationDate.String())
	d.Set("volume_snapshots", snapshotdata.VolumeSnapshots)
	d.Set("last_update_date", snapshotdata.LastUpdateDate.String())

	return nil
}

func resourceIBMPPCSnapshotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log.Printf("Calling the IBM Power Snapshot  update call")
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, snapshotID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := st.NewIBMPPCSnapshotClient(ctx, sess, cloudInstanceID)

	if d.HasChange(helpers.PPCSnapshotName) || d.HasChange("description") {
		name := d.Get(helpers.PPCSnapshotName).(string)
		description := d.Get("description").(string)
		snapshotBody := &models.SnapshotUpdate{Name: name, Description: description}

		_, err := client.Update(snapshotID, snapshotBody)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = isWaitForPPCInstanceSnapshotAvailable(ctx, client, snapshotID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIBMPPCSnapshotRead(ctx, d, meta)
}

func resourceIBMPPCSnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, snapshotID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := st.NewIBMPPCSnapshotClient(ctx, sess, cloudInstanceID)
	snapshot, err := client.Get(snapshotID)
	if err != nil {
		// snapshot does not exist
		d.SetId("")
		return nil
	}

	log.Printf("The snapshot  to be deleted is in the following state .. %s", snapshot.Status)

	err = client.Delete(snapshotID)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = isWaitForPPCInstanceSnapshotDeleted(ctx, client, snapshotID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
func isWaitForPPCInstanceSnapshotAvailable(ctx context.Context, client *st.IBMPPCSnapshotClient, id string, timeout time.Duration) (interface{}, error) {

	log.Printf("Waiting for PPCInstance Snapshot (%s) to be available and active ", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"in_progress", "BUILD"},
		Target:     []string{"available", "ACTIVE"},
		Refresh:    isPPCInstanceSnapshotRefreshFunc(client, id),
		Delay:      30 * time.Second,
		MinTimeout: 2 * time.Minute,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPPCInstanceSnapshotRefreshFunc(client *st.IBMPPCSnapshotClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		snapshotInfo, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		//if pvm.Health.Status == helpers.PPCInstanceHealthOk {
		if snapshotInfo.Status == "available" && snapshotInfo.PercentComplete == 100 {
			log.Printf("The snapshot is now available")
			return snapshotInfo, "available", nil

		}
		return snapshotInfo, "in_progress", nil
	}
}

// Delete Snapshot

func isWaitForPPCInstanceSnapshotDeleted(ctx context.Context, client *st.IBMPPCSnapshotClient, id string, timeout time.Duration) (interface{}, error) {

	log.Printf("Waiting for (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PPCInstanceDeleting},
		Target:     []string{"Not Found"},
		Refresh:    isPPCInstanceSnapshotDeleteRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPPCInstanceSnapshotDeleteRefreshFunc(client *st.IBMPPCSnapshotClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		snapshot, err := client.Get(id)
		if err != nil {
			log.Printf("The snapshot is not found.")
			return snapshot, helpers.PPCInstanceNotFound, nil
		}
		return snapshot, helpers.PPCInstanceNotFound, nil

	}
}
