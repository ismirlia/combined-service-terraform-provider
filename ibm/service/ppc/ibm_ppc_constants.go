package ppc

import "time"

const (
	// used by all
	Arg_CloudInstanceID = "ppc_cloud_instance_id"

	// Keys
	Arg_KeyName = "ppc_key_name"
	Arg_Key     = "ppc_ssh_key"

	Attr_KeyID           = "key_id"
	Attr_Keys            = "keys"
	Attr_KeyCreationDate = "creation_date"
	Attr_Key             = "ssh_key"
	Attr_KeyName         = "name"

	// SAP Profile
	PPCSAPProfiles         = "profiles"
	PPCSAPProfileCertified = "certified"
	PPCSAPProfileCores     = "cores"
	PPCSAPProfileMemory    = "memory"
	PPCSAPProfileID        = "profile_id"
	PPCSAPProfileType      = "type"

	// DHCP
	Arg_DhcpCidr              = "ppc_cidr"
	Arg_DhcpID                = "ppc_dhcp_id"
	Arg_DhcpCloudConnectionID = "ppc_cloud_connection_id"
	Arg_DhcpDnsServer         = "ppc_dns_server"
	Arg_DhcpName              = "ppc_dhcp_name"
	Arg_DhcpSnatEnabled       = "ppc_dhcp_snat_enabled"

	Attr_DhcpServers           = "servers"
	Attr_DhcpID                = "dhcp_id"
	Attr_DhcpLeases            = "leases"
	Attr_DhcpLeaseInstanceIP   = "instance_ip"
	Attr_DhcpLeaseInstanceMac  = "instance_mac"
	Attr_DhcpNetworkDeprecated = "network" // to deprecate
	Attr_DhcpNetworkID         = "network_id"
	Attr_DhcpNetworkName       = "network_name"
	Attr_DhcpStatus            = "status"

	// Instance
	Arg_PVMInstanceId           = "ppc_instance_id"
	Arg_PVMInstanceActionType   = "ppc_action"
	Arg_PVMInstanceHealthStatus = "ppc_health_status"

	Attr_Status       = "status"
	Attr_Progress     = "progress"
	Attr_HealthStatus = "health_status"

	PVMInstanceHealthOk      = "OK"
	PVMInstanceHealthWarning = "WARNING"

	//Added timeout values for warning  and active status
	warningTimeOut = 60 * time.Second
	activeTimeOut  = 2 * time.Minute
	// power service instance capabilities
	CUSTOM_VIRTUAL_CORES                  = "custom-virtualcores"
	PPCInstanceDeploymentType             = "ppc_deployment_type"
	PPCInstanceNetwork                    = "ppc_network"
	PPCInstanceStoragePool                = "ppc_storage_pool"
	PPCSAPInstanceProfileID               = "ppc_sap_profile_id"
	PPCSAPInstanceDeploymentType          = "ppc_sap_deployment_type"
	PPCInstanceStoragePoolAffinity        = "ppc_storage_pool_affinity"
	Arg_PPCInstanceSharedProcessorPool    = "ppc_shared_processor_pool"
	Attr_PPCInstanceSharedProcessorPool   = "shared_processor_pool"
	Attr_PPCInstanceSharedProcessorPoolID = "shared_processor_pool_id"

	// Placement Group
	PPCPlacementGroupID      = "placement_group_id"
	PPCPlacementGroupMembers = "members"

	// Volume
	PPCAffinityPolicy        = "ppc_affinity_policy"
	PPCAffinityVolume        = "ppc_affinity_volume"
	PPCAffinityInstance      = "ppc_affinity_instance"
	PPCAntiAffinityInstances = "ppc_anti_affinity_instances"
	PPCAntiAffinityVolumes   = "ppc_anti_affinity_volumes"

	// IBM PPC Volume Group
	PPCVolumeGroupName                 = "ppc_volume_group_name"
	PPCVolumeGroupsVolumeIds           = "ppc_volume_ids"
	PPCVolumeGroupConsistencyGroupName = "ppc_consistency_group_name"
	PPCVolumeGroupID                   = "ppc_volume_group_id"
	PPCVolumeGroupAction               = "ppc_volume_group_action"
	PPCVolumeOnboardingID              = "ppc_volume_onboarding_id"

	// Disaster Recovery Location
	PPCDRLocation = "location"

	// VPN
	PPCVPNConnectionId                         = "connection_id"
	PPCVPNConnectionStatus                     = "connection_status"
	PPCVPNConnectionDeadPeerDetection          = "dead_peer_detections"
	PPCVPNConnectionDeadPeerDetectionAction    = "action"
	PPCVPNConnectionDeadPeerDetectionInterval  = "interval"
	PPCVPNConnectionDeadPeerDetectionThreshold = "threshold"
	PPCVPNConnectionLocalGatewayAddress        = "local_gateway_address"
	PPCVPNConnectionVpnGatewayAddress          = "gateway_address"

	// Cloud Connections
	PPCCloudConnectionTransitEnabled = "ppc_cloud_connection_transit_enabled"

	// Shared Processor Pool
	Arg_SharedProcessorPoolName                      = "ppc_shared_processor_pool_name"
	Arg_SharedProcessorPoolHostGroup                 = "ppc_shared_processor_pool_host_group"
	Arg_SharedProcessorPoolPlacementGroupID          = "ppc_shared_processor_pool_placement_group_id"
	Arg_SharedProcessorPoolReservedCores             = "ppc_shared_processor_pool_reserved_cores"
	Arg_SharedProcessorPoolID                        = "ppc_shared_processor_pool_id"
	Attr_SharedProcessorPoolID                       = "shared_processor_pool_id"
	Attr_SharedProcessorPoolName                     = "name"
	Attr_SharedProcessorPoolReservedCores            = "reserved_cores"
	Attr_SharedProcessorPoolAvailableCores           = "available_cores"
	Attr_SharedProcessorPoolAllocatedCores           = "allocated_cores"
	Attr_SharedProcessorPoolHostID                   = "host_id"
	Attr_SharedProcessorPoolStatus                   = "status"
	Attr_SharedProcessorPoolStatusDetail             = "status_detail"
	Attr_SharedProcessorPoolPlacementGroups          = "spp_placement_groups"
	Attr_SharedProcessorPoolInstances                = "instances"
	Attr_SharedProcessorPoolInstanceCpus             = "cpus"
	Attr_SharedProcessorPoolInstanceUncapped         = "uncapped"
	Attr_SharedProcessorPoolInstanceAvailabilityZone = "availability_zone"
	Attr_SharedProcessorPoolInstanceId               = "id"
	Attr_SharedProcessorPoolInstanceMemory           = "memory"
	Attr_SharedProcessorPoolInstanceName             = "name"
	Attr_SharedProcessorPoolInstanceStatus           = "status"
	Attr_SharedProcessorPoolInstanceVcpus            = "vcpus"

	// PPCP Placement Group
	Arg_PPCPPlacementGroupName     = "ppc_spp_placement_group_name"
	Arg_PPCPPlacementGroupPolicy   = "ppc_spp_placement_group_policy"
	Attr_PPCPPlacementGroupID      = "spp_placement_group_id"
	Attr_PPCPPlacementGroupMembers = "members"
	Arg_PPCPPlacementGroupID       = "ppc_spp_placement_group_id"
	Attr_PPCPPlacementGroupPolicy  = "policy"
	Attr_PPCPPlacementGroupName    = "name"

	// status
	// common status states
	StatusShutoff = "SHUTOFF"
	StatusActive  = "ACTIVE"
	StatusResize  = "RESIZE"
	StatusError   = "ERROR"
	StatusBuild   = "BUILD"
	StatusPending = "PENDING"
	SctionStart   = "start"
	SctionStop    = "stop"
)
