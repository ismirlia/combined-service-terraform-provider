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

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
	"github.com/IBM-Cloud/ppc-aas-go-client/ppc-aas/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMPPCNetworkPortAttach() *schema.Resource {
	return &schema.Resource{

		CreateContext: resourceIBMPPCNetworkPortAttachCreate,
		ReadContext:   resourceIBMPPCNetworkPortAttachRead,
		DeleteContext: resourceIBMPPCNetworkPortAttachDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			helpers.PPCCloudInstanceId: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			helpers.PPCInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance id to attach the network port to",
			},
			helpers.PPCNetworkName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Network Name - This is the subnet name  in the Cloud instance",
			},
			helpers.PPCNetworkPortDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A human readable description for this network Port",
				Default:     "Port Created via Terraform",
			},
			helpers.PPCNetworkPortIPAddress: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			//Computed Attributes
			"macaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "port_id attribute is deprecated, use network_port_id instead.",
			},
			"network_port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func resourceIBMPPCNetworkPortAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)
	networkname := d.Get(helpers.PPCNetworkName).(string)
	instanceID := d.Get(helpers.PPCInstanceId).(string)
	description := d.Get(helpers.PPCNetworkPortDescription).(string)
	nwportBody := &models.NetworkPortCreate{Description: description}

	if v, ok := d.GetOk(helpers.PPCNetworkPortIPAddress); ok {
		ipaddress := v.(string)
		nwportBody.IPAddress = ipaddress
	}

	nwportattachBody := &models.NetworkPortUpdate{
		Description:   &description,
		PvmInstanceID: &instanceID,
	}

	client := st.NewIBMPPCNetworkClient(ctx, sess, cloudInstanceID)

	networkPortResponse, err := client.CreatePort(networkname, nwportBody)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Printing the networkresponse %+v", &networkPortResponse)

	networkPortID := *networkPortResponse.PortID

	_, err = isWaitForIBMPPCNetworkportAvailable(ctx, client, networkPortID, networkname, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	networkPortResponse, err = client.UpdatePort(networkname, networkPortID, nwportattachBody)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = isWaitForIBMPPCNetworkPortAttachAvailable(ctx, client, networkPortID, networkname, instanceID, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", cloudInstanceID, networkname, networkPortID))

	return resourceIBMPPCNetworkPortAttachRead(ctx, d, meta)
}

func resourceIBMPPCNetworkPortAttachRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := parts[0]
	networkname := parts[1]
	portID := parts[2]

	networkC := st.NewIBMPPCNetworkClient(ctx, sess, cloudInstanceID)
	networkdata, err := networkC.GetPort(networkname, portID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(helpers.PPCNetworkPortIPAddress, networkdata.IPAddress)
	d.Set(helpers.PPCNetworkPortDescription, networkdata.Description)
	d.Set(helpers.PPCInstanceId, networkdata.PvmInstance.PvmInstanceID)
	d.Set("macaddress", networkdata.MacAddress)
	d.Set("status", networkdata.Status)
	d.Set("network_port_id", networkdata.PortID)
	d.Set("port_id", networkdata.PortID)
	d.Set("public_ip", networkdata.ExternalIP)

	return nil
}

func resourceIBMPPCNetworkPortAttachDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log.Printf("Calling the network delete functions. ")
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := parts[0]
	networkname := parts[1]
	portID := parts[2]

	client := st.NewIBMPPCNetworkClient(ctx, sess, cloudInstanceID)

	log.Printf("Calling the delete with the following params delete with cloud instance (%s) and networkid (%s) and portid (%s) ", cloudInstanceID, networkname, portID)
	err = client.DeletePort(networkname, portID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func isWaitForIBMPPCNetworkportAvailable(ctx context.Context, client *st.IBMPPCNetworkClient, id string, networkname string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Power Network (%s) that was created for Network Zone (%s) to be available.", id, networkname)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PPCNetworkProvisioning},
		Target:     []string{"DOWN"},
		Refresh:    isIBMPPCNetworkportRefreshFunc(client, id, networkname),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Minute,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPPCNetworkportRefreshFunc(client *st.IBMPPCNetworkClient, id, networkname string) resource.StateRefreshFunc {

	log.Printf("Calling the IsIBMPPCNetwork Refresh Function....with the following id (%s) for network port and following id (%s) for network name and waiting for network to be READY", id, networkname)
	return func() (interface{}, string, error) {
		network, err := client.GetPort(networkname, id)
		if err != nil {
			return nil, "", err
		}

		if *network.Status == "DOWN" {
			log.Printf(" The port has been created with the following ip address and attached to an instance ")
			return network, "DOWN", nil
		}

		return network, helpers.PPCNetworkProvisioning, nil
	}
}
func isWaitForIBMPPCNetworkPortAttachAvailable(ctx context.Context, client *st.IBMPPCNetworkClient, id, networkname, instanceid string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Power Network (%s) that was created for Network Zone (%s) to be available.", id, networkname)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PPCNetworkProvisioning},
		Target:     []string{"ACTIVE"},
		Refresh:    isIBMPPCNetworkPortAttachRefreshFunc(client, id, networkname, instanceid),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Minute,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPPCNetworkPortAttachRefreshFunc(client *st.IBMPPCNetworkClient, id, networkname, instanceid string) resource.StateRefreshFunc {

	log.Printf("Calling the IsIBMPPCNetwork Refresh Function....with the following id (%s) for network port and following id (%s) for network name and waiting for network to be READY", id, networkname)
	return func() (interface{}, string, error) {
		network, err := client.GetPort(networkname, id)
		if err != nil {
			return nil, "", err
		}

		if *network.Status == "ACTIVE" && network.PvmInstance.PvmInstanceID == instanceid {
			log.Printf(" The port has been created with the following ip address and attached to an instance ")
			return network, "ACTIVE", nil
		}

		return network, helpers.PPCNetworkProvisioning, nil
	}
}
