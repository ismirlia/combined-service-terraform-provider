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

func ResourceIBMPPCNetworkPort() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPPCNetworkPortCreate,
		ReadContext:   resourceIBMPPCNetworkPortRead,
		UpdateContext: resourceIBMPPCNetworkPortUpdate,
		DeleteContext: resourceIBMPPCNetworkPortDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			helpers.PPCNetworkName: {
				Type:     schema.TypeString,
				Required: true,
			},
			helpers.PPCCloudInstanceId: {
				Type:     schema.TypeString,
				Required: true,
			},
			helpers.PPCNetworkPortDescription: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			helpers.PPCNetworkPortIPAddress: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			//Computed Attributes
			"macaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"portid": {
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
		DeprecationMessage: "Resource ibm_ppc_network_port is deprecated. Use ibm_ppc_network_port_attach to create & attach a network port to a pvm instance",
	}
}

func resourceIBMPPCNetworkPortCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)
	networkname := d.Get(helpers.PPCNetworkName).(string)
	description := d.Get(helpers.PPCNetworkPortDescription).(string)

	ipaddress := d.Get(helpers.PPCNetworkPortIPAddress).(string)

	nwportBody := &models.NetworkPortCreate{Description: description}

	if ipaddress != "" {
		log.Printf("IP address provided. ")
		nwportBody.IPAddress = ipaddress
	}

	client := st.NewIBMPPCNetworkClient(ctx, sess, cloudInstanceID)

	networkPortResponse, err := client.CreatePort(networkname, nwportBody)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Printing the networkresponse %+v", &networkPortResponse)

	IBMPPCNetworkPortID := *networkPortResponse.PortID

	d.SetId(fmt.Sprintf("%s/%s/%s", cloudInstanceID, networkname, IBMPPCNetworkPortID))

	_, err = isWaitForIBMPPCNetworkPortAvailable(ctx, client, IBMPPCNetworkPortID, networkname, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIBMPPCNetworkPortRead(ctx, d, meta)
}

func resourceIBMPPCNetworkPortRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	d.Set("macaddress", networkdata.MacAddress)
	d.Set("status", networkdata.Status)
	d.Set("portid", networkdata.PortID)
	d.Set("public_ip", networkdata.ExternalIP)

	return nil
}

func resourceIBMPPCNetworkPortUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceIBMPPCNetworkPortDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

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

func isWaitForIBMPPCNetworkPortAvailable(ctx context.Context, client *st.IBMPPCNetworkClient, id string, networkname string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for Power Network (%s) that was created for Network Zone (%s) to be available.", id, networkname)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PPCNetworkProvisioning},
		Target:     []string{"DOWN"},
		Refresh:    isIBMPPCNetworkPortRefreshFunc(client, id, networkname),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Minute,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPPCNetworkPortRefreshFunc(client *st.IBMPPCNetworkClient, id, networkname string) resource.StateRefreshFunc {

	log.Printf("Calling the IsIBMPPCNetwork Refresh Function....with the following id (%s) for network port and following id (%s) for network name and waiting for network to be READY", id, networkname)
	return func() (interface{}, string, error) {
		network, err := client.GetPort(networkname, id)
		if err != nil {
			return nil, "", err
		}

		if &network.PortID != nil {
			//if network.State == "available" {
			log.Printf(" The port has been created with the following ip address and attached to an instance ")
			return network, "DOWN", nil
		}

		return network, helpers.PPCNetworkProvisioning, nil
	}
}
