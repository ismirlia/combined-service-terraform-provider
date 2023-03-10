// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ppc

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/ppc-aas-go-client/clients/instance"
	"github.com/IBM-Cloud/ppc-aas-go-client/helpers"
	"github.com/IBM-Cloud/ppc-aas-go-client/ppc-aas/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

const (
	piEndingIPAaddress   = "ppc_ending_ip_address"
	piStartingIPAaddress = "ppc_starting_ip_address"
)

func ResourceIBMPPCNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPPCNetworkCreate,
		ReadContext:   resourceIBMPPCNetworkRead,
		UpdateContext: resourceIBMPPCNetworkUpdate,
		DeleteContext: resourceIBMPPCNetworkDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			helpers.PPCNetworkType: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"vlan", "pub-vlan"}),
				Description:  "PPC network type",
			},
			helpers.PPCNetworkName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PPC network name",
			},
			helpers.PPCNetworkDNS: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of PPC network DNS name",
			},
			helpers.PPCNetworkCidr: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "PPC network CIDR",
			},
			helpers.PPCNetworkGateway: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "PPC network gateway",
			},
			helpers.PPCNetworkMtu: {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "PPC Maximum Transmission Unit option",
			},
			helpers.PPCCloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PPC cloud instance ID",
			},
			helpers.PPCNetworkIPAddressRange: {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of one or more ip address range(s)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						piEndingIPAaddress: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Ending ip address",
						},
						piStartingIPAaddress: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Starting ip address",
						},
					},
				},
			},

			//Computed Attributes
			"network_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "PPC network ID",
			},
			"vlan_id": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "VLAN Id value",
			},
		},
	}
}

func resourceIBMPPCNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := d.Get(helpers.PPCCloudInstanceId).(string)
	networkname := d.Get(helpers.PPCNetworkName).(string)
	networktype := d.Get(helpers.PPCNetworkType).(string)

	client := st.NewIBMPPCNetworkClient(ctx, sess, cloudInstanceID)
	var body = &models.NetworkCreate{
		Type: &networktype,
		Name: networkname,
	}
	if v, ok := d.GetOk(helpers.PPCNetworkDNS); ok {
		networkdns := flex.ExpandStringList((v.(*schema.Set)).List())
		if len(networkdns) > 0 {
			body.DNSServers = networkdns
		}
	}

	if v, ok := d.GetOk(helpers.PPCNetworkMtu); ok {
		body.Mtu = v.(int64)
	}

	if networktype == "vlan" {
		var networkcidr string
		var ipBodyRanges []*models.IPAddressRange
		if v, ok := d.GetOk(helpers.PPCNetworkCidr); ok {
			networkcidr = v.(string)
		} else {
			return diag.Errorf("%s is required when %s is vlan", helpers.PPCNetworkCidr, helpers.PPCNetworkType)
		}

		gateway, firstip, lastip, err := generateIPData(networkcidr)
		if err != nil {
			return diag.FromErr(err)
		}

		ipBodyRanges = []*models.IPAddressRange{{EndingIPAddress: &lastip, StartingIPAddress: &firstip}}

		if g, ok := d.GetOk(helpers.PPCNetworkGateway); ok {
			gateway = g.(string)
		}

		if ips, ok := d.GetOk(helpers.PPCNetworkIPAddressRange); ok {
			ipBodyRanges = getIPAddressRanges(ips.([]interface{}))
		}

		body.IPAddressRanges = ipBodyRanges
		body.Gateway = gateway
		body.Cidr = networkcidr
	}

	networkResponse, err := client.Create(body)
	if err != nil {
		return diag.FromErr(err)
	}

	networkID := *networkResponse.NetworkID

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, networkID))

	_, err = isWaitForIBMPPCNetworkAvailable(ctx, client, networkID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIBMPPCNetworkRead(ctx, d, meta)
}

func resourceIBMPPCNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, networkID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	networkC := st.NewIBMPPCNetworkClient(ctx, sess, cloudInstanceID)
	networkdata, err := networkC.Get(networkID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("network_id", networkdata.NetworkID)
	d.Set(helpers.PPCNetworkCidr, networkdata.Cidr)
	d.Set(helpers.PPCNetworkDNS, networkdata.DNSServers)
	d.Set("vlan_id", networkdata.VlanID)
	d.Set(helpers.PPCNetworkName, networkdata.Name)
	d.Set(helpers.PPCNetworkType, networkdata.Type)
	d.Set(helpers.PPCNetworkMtu, networkdata.Mtu)
	d.Set(helpers.PPCNetworkGateway, networkdata.Gateway)
	ipRangesMap := []map[string]interface{}{}
	if networkdata.IPAddressRanges != nil {
		for _, n := range networkdata.IPAddressRanges {
			if n != nil {
				v := map[string]interface{}{
					piEndingIPAaddress:   n.EndingIPAddress,
					piStartingIPAaddress: n.StartingIPAddress,
				}
				ipRangesMap = append(ipRangesMap, v)
			}
		}
	}
	d.Set(helpers.PPCNetworkIPAddressRange, ipRangesMap)

	return nil

}

func resourceIBMPPCNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, networkID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges(helpers.PPCNetworkName, helpers.PPCNetworkDNS, helpers.PPCNetworkGateway, helpers.PPCNetworkIPAddressRange) {
		networkC := st.NewIBMPPCNetworkClient(ctx, sess, cloudInstanceID)
		body := &models.NetworkUpdate{
			DNSServers: flex.ExpandStringList((d.Get(helpers.PPCNetworkDNS).(*schema.Set)).List()),
		}
		if d.Get(helpers.PPCNetworkType).(string) == "vlan" {
			body.Gateway = flex.PtrToString(d.Get(helpers.PPCNetworkGateway).(string))
			body.IPAddressRanges = getIPAddressRanges(d.Get(helpers.PPCNetworkIPAddressRange).([]interface{}))
		}

		if d.HasChange(helpers.PPCNetworkName) {
			body.Name = flex.PtrToString(d.Get(helpers.PPCNetworkName).(string))
		}

		_, err = networkC.Update(networkID, body)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIBMPPCNetworkRead(ctx, d, meta)
}

func resourceIBMPPCNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log.Printf("Calling the network delete functions. ")
	sess, err := meta.(conns.ClientSession).IBMPPCSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, networkID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	networkC := st.NewIBMPPCNetworkClient(ctx, sess, cloudInstanceID)
	err = networkC.Delete(networkID)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func isWaitForIBMPPCNetworkAvailable(ctx context.Context, client *st.IBMPPCNetworkClient, id string, timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PPCNetworkProvisioning},
		Target:     []string{"NETWORK_READY"},
		Refresh:    isIBMPPCNetworkRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPPCNetworkRefreshFunc(client *st.IBMPPCNetworkClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		network, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if network.VlanID != nil {
			return network, "NETWORK_READY", nil
		}

		return network, helpers.PPCNetworkProvisioning, nil
	}
}

func generateIPData(cdir string) (gway, firstip, lastip string, err error) {
	_, ipv4Net, err := net.ParseCIDR(cdir)

	if err != nil {
		return "", "", "", err
	}

	var subnetToSize = map[string]int{
		"21": 2048,
		"22": 1024,
		"23": 512,
		"24": 256,
		"25": 128,
		"26": 64,
		"27": 32,
		"28": 16,
		"29": 8,
		"30": 4,
		"31": 2,
	}

	//subnetsize, _ := ipv4Net.Mask.Size()

	gateway, err := cidr.Host(ipv4Net, 1)
	if err != nil {
		log.Printf("Failed to get the gateway for this cidr passed in %s", cdir)
		return "", "", "", err
	}
	ad := cidr.AddressCount(ipv4Net)

	convertedad := strconv.FormatUint(ad, 10)
	// Powervc in wdc04 has to reserve 3 ip address hence we start from the 4th. This will be the default behaviour
	firstusable, err := cidr.Host(ipv4Net, 4)
	if err != nil {
		log.Print(err)
		return "", "", "", err
	}
	lastusable, err := cidr.Host(ipv4Net, subnetToSize[convertedad]-2)
	if err != nil {
		log.Print(err)
		return "", "", "", err
	}
	return gateway.String(), firstusable.String(), lastusable.String(), nil

}

func getIPAddressRanges(ipAddressRanges []interface{}) []*models.IPAddressRange {
	ipRanges := make([]*models.IPAddressRange, 0, len(ipAddressRanges))
	for _, v := range ipAddressRanges {
		if v != nil {
			ipAddressRange := v.(map[string]interface{})
			ipRange := &models.IPAddressRange{
				EndingIPAddress:   flex.PtrToString(ipAddressRange[piEndingIPAaddress].(string)),
				StartingIPAddress: flex.PtrToString(ipAddressRange[piStartingIPAaddress].(string)),
			}
			ipRanges = append(ipRanges, ipRange)
		}
	}
	return ipRanges
}
