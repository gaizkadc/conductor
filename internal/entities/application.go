/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package entities

import (
	"github.com/nalej/grpc-application-go"
)


type PortAccess int

const (
	AllAppServices PortAccess = iota + 1
	AppServices
	Public
	DeviceGroup
)

var PortAccessToGRPC = map[PortAccess]grpc_application_go.PortAccess{
	AllAppServices: grpc_application_go.PortAccess_ALL_APP_SERVICES,
	AppServices:    grpc_application_go.PortAccess_APP_SERVICES,
	Public:         grpc_application_go.PortAccess_PUBLIC,
	DeviceGroup:    grpc_application_go.PortAccess_DEVICE_GROUP,
}

var PortAccessFromGRPC = map[grpc_application_go.PortAccess]PortAccess{
	grpc_application_go.PortAccess_ALL_APP_SERVICES: AllAppServices,
	grpc_application_go.PortAccess_APP_SERVICES:     AppServices,
	grpc_application_go.PortAccess_PUBLIC:           Public,
	grpc_application_go.PortAccess_DEVICE_GROUP:     DeviceGroup,
}

type CollocationPolicy int

const (
	SameCluster CollocationPolicy = iota + 1
	SeparateClusters
)

var CollocationPolicyToGRPC = map[CollocationPolicy]grpc_application_go.CollocationPolicy{
	SameCluster:      grpc_application_go.CollocationPolicy_SAME_CLUSTER,
	SeparateClusters: grpc_application_go.CollocationPolicy_SEPARATE_CLUSTERS,
}

var CollocationPolicyFromGRPC = map[grpc_application_go.CollocationPolicy]CollocationPolicy{
	grpc_application_go.CollocationPolicy_SAME_CLUSTER:      SameCluster,
	grpc_application_go.CollocationPolicy_SEPARATE_CLUSTERS: SeparateClusters,
}


// InstanceMetadata----------


// Instance metadata
// This is a common metadata entity that collects information for a deployed instance. This instance can be a
// service instance or a service group instance.
type InstanceMetadata struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// AppInstanceId with the application instance identifier.
	AppInstanceId string `json:"app_instance_id,omitempty"`
	// ServiceGroupId with the id of the group this instance belongs to
	ServiceGroupId string `json:"service_group_id,omitempty"`
	// Identifier of the monitored entity
	MonitoredInstanceId string `json:"monitored_instance_id,omitempty"`
	// Type of instance this metadata refers to
	InstanceType InstanceType `json:"instance_type,omitempty"`
	// List of instances supervised byu this metadata structure
	InstancesId []string `json:"instances_id,omitempty"`
	// Number of desired replicas specified in the descriptor
	DesiredReplicas int32 `json:"desired_replicas,omitempty"`
	// Number of available replicas for this instance
	AvailableReplicas int32 `json:"available_replicas,omitempty"`
	// Number of unavaiable replicas for this descriptor
	UnavailableReplicas int32 `json:"unavailable_replicas,omitempty"`
	// Status of every item monitored by this metadata entry
	Status map[string]ServiceStatus `json:"Status,omitempty"`
	// Relevant information for every monitored instance
	Info map[string]string `json:"Info,omitempty"`
}

func (im *InstanceMetadata) ToGRPC() *grpc_application_go.InstanceMetadata {
	statuses := make(map[string]grpc_application_go.ServiceStatus,len(im.Status))
	for k,v := range im.Status {
		statuses[k] = ServiceStatusToGRPC[v]
	}
	return &grpc_application_go.InstanceMetadata{
		OrganizationId: im.OrganizationId,
		AppDescriptorId: im.AppDescriptorId,
		AppInstanceId: im.AppInstanceId,
		ServiceGroupId: im.ServiceGroupId,
		MonitoredInstanceId: im.MonitoredInstanceId,
		Type: InstanceTypeToGRPC[im.InstanceType],
		InstancesId: im.InstancesId,
		DesiredReplicas: im.DesiredReplicas,
		AvailableReplicas: im.AvailableReplicas,
		UnavailableReplicas: im.UnavailableReplicas,
		Status: statuses,
		Info: im.Info,
	}
}

func NewInstanceMetadataFromGRPC(ins *grpc_application_go.InstanceMetadata) InstanceMetadata {
	status := make(map[string]ServiceStatus,0)
	for k, v := range ins.Status{
		status[k] = ServiceStatusFromGRPC[v]
	}
	return InstanceMetadata{
		AppDescriptorId: ins.AppDescriptorId,
		AppInstanceId: ins.AppInstanceId,
		OrganizationId: ins.OrganizationId,
		ServiceGroupId: ins.ServiceGroupId,
		Info: ins.Info,
		UnavailableReplicas: ins.UnavailableReplicas,
		AvailableReplicas: ins.AvailableReplicas,
		DesiredReplicas: ins.DesiredReplicas,
		Status: status,
		InstanceType: InstanceTypeFromGRPC[ins.Type],
		InstancesId: ins.InstancesId,
		MonitoredInstanceId: ins.MonitoredInstanceId,
	}
}

// ----------

// InstanceType----------

type InstanceType int32

const (
	// A service
	ServiceInstanceType InstanceType = iota + 1
	// A service group that contains several running service instances
	ServiceGroupInstanceType
)

var InstanceTypeToGRPC = map[InstanceType]grpc_application_go.InstanceType{
	ServiceInstanceType: grpc_application_go.InstanceType_SERVICE_INSTANCE,
	ServiceGroupInstanceType: grpc_application_go.InstanceType_SERVICE_GROUP_INSTANCE,
}

var InstanceTypeFromGRPC = map[grpc_application_go.InstanceType]InstanceType {
	grpc_application_go.InstanceType_SERVICE_INSTANCE: ServiceInstanceType,
	grpc_application_go.InstanceType_SERVICE_GROUP_INSTANCE: ServiceGroupInstanceType,
}


// SecurityRule----------

type SecurityRule struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// RuleId with the security rule identifier.
	RuleId string `json:"rule_id,omitempty"`
	// Name of the security rule.
	Name string `json:"name,omitempty"`
	// TargetServiceGroupName defining the name of the target group
	TargetServiceGroupName string `json:"target_service_group_name,omitempty"`
	// TargetServiceName defining the name of the target service
	TargetServiceName string `json:"target_service_name,omitempty"`
	// TargetPort defining the port that is affected by the current rule.
	TargetPort int32 `json:"source_port,omitempty"`
	// Access level to that port defining who can access the port.
	Access PortAccess `json:"access,omitempty"`
	// Name of the authenticated group name
	AuthServiceGroupName string `json:"auth_service_group_name,omitempty"`
	// AuthServices defining a list of services that can access the port.
	AuthServices []string `json:"auth_services,omitempty"`
	// DeviceGroupNames defining a list of device groups that can access the port.
	DeviceGroupNames []string `json:"device_group_names,omitempty"`
	// DeviceGroupIds defining a list of device group ids that can access the port.
	DeviceGroupIds []string `json:"device_group_Ids,omitempty"`
}

func (sr *SecurityRule) ToGRPC() *grpc_application_go.SecurityRule {
	access, _ := PortAccessToGRPC[sr.Access]
	return &grpc_application_go.SecurityRule{
		OrganizationId:  sr.OrganizationId,
		AppDescriptorId: sr.AppDescriptorId,
		RuleId:          sr.RuleId,
		Name:            sr.Name,
		AuthServiceGroupName: sr.AuthServiceGroupName,
		TargetPort:      sr.TargetPort,
		TargetServiceGroupName: sr.TargetServiceGroupName,
		TargetServiceName: sr.TargetServiceName,
		Access:       access,
		AuthServices: sr.AuthServices,
		DeviceGroupNames: sr.DeviceGroupNames,
		DeviceGroupIds: sr.DeviceGroupIds,
	}
}

func NewSecurityRuleFromGRPC (s *grpc_application_go.SecurityRule) SecurityRule {
	return SecurityRule{
		OrganizationId: s.OrganizationId,
		AppDescriptorId: s.AppDescriptorId,
		RuleId: s.RuleId,
		Name: s.Name,
		TargetServiceGroupName: s.TargetServiceGroupName,
		TargetServiceName: s.TargetServiceName,
		TargetPort: s.TargetPort,
		Access: PortAccessFromGRPC[s.Access],
		AuthServiceGroupName: s.AuthServiceGroupName,
		AuthServices: s.AuthServices,
		DeviceGroupIds: s.DeviceGroupIds,
		DeviceGroupNames: s.DeviceGroupNames,
	}
}


// ----------

// Service Group ----------

type ServiceGroup struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// ServiceGroupId with the group identifier.
	ServiceGroupId string `json:"service_group_id,omitempty"`
	// Name of the service group.
	Name string `json:"name,omitempty"`
	// Services defining a list of service identifiers that belong to the group.
	Services []Service `json:"services,omitempty"`
	// Policy indicating the deployment collocation policy.
	Policy CollocationPolicy `json:"policy,omitempty"`
	// Particular deployment specs for this service
	Specs ServiceGroupDeploymentSpecs `json:"specs,omitempty"`
	// Labels defined by the user
	Labels map[string]string `json:"labels,omitempty"`
}


func (sg *ServiceGroup) ToGRPC() *grpc_application_go.ServiceGroup {
	policy, _ := CollocationPolicyToGRPC[sg.Policy]
	servs := make([]*grpc_application_go.Service, len(sg.Services))
	for i, s := range(sg.Services) {
		servs[i] = s.ToGRPC()
	}

	return &grpc_application_go.ServiceGroup{
		OrganizationId:  sg.OrganizationId,
		AppDescriptorId: sg.AppDescriptorId,
		ServiceGroupId:  sg.ServiceGroupId,
		Name:            sg.Name,
		Services:        servs,
		Policy:          policy,
		Labels:			 sg.Labels,
		Specs:           sg.Specs.ToGRPC(),
	}
}

func NewServiceGroupFromGRPC(g *grpc_application_go.ServiceGroup) ServiceGroup {
	servs := make([]Service, 0)
	for _, s := range g.Services {
		servs = append(servs, *NewServiceFromGRPC(g.AppDescriptorId, s))
	}
	return ServiceGroup{
		OrganizationId: g.OrganizationId,
		AppDescriptorId: g.AppDescriptorId,
		ServiceGroupId: g.ServiceGroupId,
		Name: g.Name,
		Labels: g.Labels,
		Policy: CollocationPolicyFromGRPC[g.Policy],
		Specs:  NewServiceGroupDeploymentSpecsFromGRPC(g.Specs),
		Services: servs,
	}
}

// ----------

// Service Group instance ----------

type ServiceGroupInstance struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// AppInstanceId with the application instance identifier.
	AppInstanceId string `json:"app_instance_id,omitempty"`
	// ServiceGroupId with the group identifier.
	ServiceGroupId string `json:"service_group_id,omitempty"`
	// ServiceGroupInstaceId with the instance identifier.
	ServiceGroupInstanceId string `json:"service_group_instance_id,omitempty"`
	// Name of the service group.
	Name string `json:"name,omitempty"`
	// ServicesInstances defining a list of service identifiers that belong to the group.
	ServiceInstances []ServiceInstance `json:"service_instances,omitempty"`
	// Policy indicating the deployment collocation policy.
	Policy               CollocationPolicy `json:"policy,omitempty"`
	// Service Status
	Status ServiceStatus `json:"service_status,omitempty"`
	// Instance metadata
	Metadata InstanceMetadata `json:"metadata,omitempty"`
	// Service group deployment specs
	Specs ServiceGroupDeploymentSpecs `json:"specs,omitempty"`
	// Labels defined by the user
	Labels map[string]string `json:"labels,omitempty"`
}



func (sgi *ServiceGroupInstance) ToGRPC() *grpc_application_go.ServiceGroupInstance {
	policy, _ := CollocationPolicyToGRPC[sgi.Policy]
	instances := make([]*grpc_application_go.ServiceInstance,len(sgi.ServiceInstances))
	for i, ins := range sgi.ServiceInstances {
		instances[i] = ins.ToGRPC()
	}
	return &grpc_application_go.ServiceGroupInstance{
		OrganizationId:       sgi.OrganizationId,
		AppDescriptorId:      sgi.AppDescriptorId,
		AppInstanceId:        sgi.AppInstanceId,
		ServiceGroupId:       sgi.ServiceGroupId,
		ServiceGroupInstanceId: sgi.ServiceGroupInstanceId,
		Name:                 sgi.Name,
		ServiceInstances:     instances,
		Policy:               policy,
		Status:               ServiceStatusToGRPC[sgi.Status],
		Metadata:             sgi.Metadata.ToGRPC(),
		Specs:                sgi.Specs.ToGRPC(),
		Labels:               sgi.Labels,
	}
}

func NewServiceGroupInstanceFromGRPC(group *grpc_application_go.ServiceGroupInstance) ServiceGroupInstance {
	serviceInstances := make([]ServiceInstance,0)
	for _, serv := range group.ServiceInstances {
		serviceInstances = append(serviceInstances, NewServiceInstanceFromGRPC(serv))
	}
	metadata := InstanceMetadata{}
	if group.Metadata != nil {
		metadata = NewInstanceMetadataFromGRPC(group.Metadata)
	}
	return ServiceGroupInstance{
		Name: group.Name,
		Status: ServiceStatusFromGRPC[group.Status],
		OrganizationId: group.OrganizationId,
		Labels: group.Labels,
		AppInstanceId: group.AppInstanceId,
		AppDescriptorId: group.AppDescriptorId,
		Specs: NewServiceGroupDeploymentSpecsFromGRPC(group.Specs),
		ServiceGroupInstanceId: group.ServiceGroupInstanceId,
		ServiceGroupId: group.ServiceGroupId,
		Policy: CollocationPolicyFromGRPC[group.Policy],
		ServiceInstances: serviceInstances,
		Metadata: metadata,
	}
}

// ServiceGroupInstanceList
type ServiceGroupInstancesList struct {
	ServiceGroupInstances []ServiceGroupInstance
}

func NewServiceGroupInstanceListFromGRPC(list *grpc_application_go.ServiceGroupInstancesList) ServiceGroupInstancesList {
	l := make([]ServiceGroupInstance,len(list.ServiceGroupInstances))
	for i, g := range list.ServiceGroupInstances {
		l[i] = NewServiceGroupInstanceFromGRPC(g)
	}
	return ServiceGroupInstancesList{ServiceGroupInstances: l}
}

// ----------

// ServiceGroupDeploymentSpecs----------

type ServiceGroupDeploymentSpecs struct {
	NumReplicas int32 `json:"num_replicas,omitempty"`
	MultiClusterReplica bool `json:"multi_cluster_replica,omitempty"`
	DeploymentSelectors map[string]string `json:"deployment_selectors,omitempty"`
}

func(sgds *ServiceGroupDeploymentSpecs) ToGRPC() *grpc_application_go.ServiceGroupDeploymentSpecs {
	return &grpc_application_go.ServiceGroupDeploymentSpecs{
		NumReplicas: sgds.NumReplicas,
		MultiClusterReplica: sgds.MultiClusterReplica,
		DeploymentSelectors: sgds.DeploymentSelectors,
	}
}

func NewServiceGroupDeploymentSpecsFromGRPC(s *grpc_application_go.ServiceGroupDeploymentSpecs) ServiceGroupDeploymentSpecs {
	return ServiceGroupDeploymentSpecs{
		NumReplicas: s.NumReplicas,
		MultiClusterReplica: s.MultiClusterReplica,
		DeploymentSelectors: s.DeploymentSelectors,
	}
}

// ----------

type ServiceType int32

const (
	DockerService ServiceType = iota + 1
)

var ServiceTypeToGRPC = map[ServiceType]grpc_application_go.ServiceType{
	DockerService: grpc_application_go.ServiceType_DOCKER,
}

var ServiceTypeFromGRPC = map[grpc_application_go.ServiceType]ServiceType{
	grpc_application_go.ServiceType_DOCKER: DockerService,
}

type ImageCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	DockerRepository     string   `json:"docker_repository,omitempty"`
}

func NewImageCredentialsFromGRPC(credentials * grpc_application_go.ImageCredentials) *ImageCredentials {
	if credentials == nil {
		return nil
	}
	return &ImageCredentials{
		Username: credentials.Username,
		Password: credentials.Password,
		Email:    credentials.Email,
		DockerRepository: credentials.DockerRepository,
	}
}

func (ic *ImageCredentials) ToGRPC() *grpc_application_go.ImageCredentials {
	if ic == nil {
		return nil
	}
	return &grpc_application_go.ImageCredentials{
		Username: ic.Username,
		Password: ic.Password,
		Email:    ic.Email,
		DockerRepository: ic.DockerRepository,
	}
}

type DeploySpecs struct {
	Cpu      int64 `json:"cpu,omitempty"`
	Memory   int64 `json:"memory,omitempty"`
	Replicas int32 `json:"replicas,omitempty"`
}

func NewDeploySpecsFromGRPC(specs * grpc_application_go.DeploySpecs) * DeploySpecs {
	if specs == nil {
		return nil
	}
	return &DeploySpecs{
		Cpu:      specs.Cpu,
		Memory:   specs.Memory,
		Replicas: specs.Replicas,
	}
}

func (ds *DeploySpecs) ToGRPC() *grpc_application_go.DeploySpecs {
	return &grpc_application_go.DeploySpecs{
		Cpu:      ds.Cpu,
		Memory:   ds.Memory,
		Replicas: ds.Replicas,
	}
}

type StorageType int32

const (
	Ephemeral StorageType = iota + 1
	ClusterLocal
	ClusterReplica
	CloudPersistent
)

var StorageTypeToGRPC = map[StorageType]grpc_application_go.StorageType{
	Ephemeral:       grpc_application_go.StorageType_EPHEMERAL,
	ClusterLocal:    grpc_application_go.StorageType_CLUSTER_LOCAL,
	ClusterReplica:  grpc_application_go.StorageType_CLUSTER_REPLICA,
	CloudPersistent: grpc_application_go.StorageType_CLOUD_PERSISTENT,
}

var StorageTypeFromGRPC = map[grpc_application_go.StorageType]StorageType{
	grpc_application_go.StorageType_EPHEMERAL:        Ephemeral,
	grpc_application_go.StorageType_CLUSTER_LOCAL:    ClusterLocal,
	grpc_application_go.StorageType_CLUSTER_REPLICA:  ClusterReplica,
	grpc_application_go.StorageType_CLOUD_PERSISTENT: CloudPersistent,
}

type Storage struct {
	Size      int64       `json:"size,omitempty"`
	MountPath string      `json:"mount_path,omitempty"`
	Type      StorageType `json:"type,omitempty"`
}

func NewStorageFromGRPC(storage * grpc_application_go.Storage) * Storage{
	if storage == nil {
		return nil
	}
	storageType, _ := StorageTypeFromGRPC[storage.Type]
	return &Storage{
		Size:      storage.Size,
		MountPath: storage.MountPath,
		Type:      storageType,
	}
}

func (s *Storage) ToGRPC() *grpc_application_go.Storage {
	convertedType, _ := StorageTypeToGRPC[s.Type]
	return &grpc_application_go.Storage{
		Size:      s.Size,
		MountPath: s.MountPath,
		Type:      convertedType,
	}
}

type EndpointType int

const (
	IsAlive EndpointType = iota + 1
	Rest
	Web
	Prometheus
)

var EndpointTypeToGRPC = map[EndpointType]grpc_application_go.EndpointType{
	IsAlive:    grpc_application_go.EndpointType_IS_ALIVE,
	Rest:       grpc_application_go.EndpointType_REST,
	Web:        grpc_application_go.EndpointType_WEB,
	Prometheus: grpc_application_go.EndpointType_PROMETHEUS,
}

var EndpointTypeFromGRPC = map[grpc_application_go.EndpointType]EndpointType{
	grpc_application_go.EndpointType_IS_ALIVE:   IsAlive,
	grpc_application_go.EndpointType_REST:       Rest,
	grpc_application_go.EndpointType_WEB:        Web,
	grpc_application_go.EndpointType_PROMETHEUS: Prometheus,
}

type Endpoint struct {
	Type EndpointType `json:"type,omitempty"`
	Path string       `json:"path,omitempty"`
}

func NewEndpointFromGRPC( endpoint * grpc_application_go.Endpoint) * Endpoint {
	if endpoint == nil {
		return nil
	}
	endpointType, _ := EndpointTypeFromGRPC[endpoint.Type]
	return &Endpoint{
		Type: endpointType,
		Path: endpoint.Path,
	}
}

func (e *Endpoint) ToGRPC() *grpc_application_go.Endpoint {
	convertedType, _ := EndpointTypeToGRPC[e.Type]
	return &grpc_application_go.Endpoint{
		Type: convertedType,
		Path: e.Path,
	}
}

type Port struct {
	Name         string     `json:"name,omitempty"`
	InternalPort int32      `json:"internal_port,omitempty"`
	ExposedPort  int32      `json:"exposed_port,omitempty"`
	Endpoints    []Endpoint `json:"endpoints,omitempty"`
}

func NewPortFromGRPC(port *grpc_application_go.Port) * Port {
	if port == nil {
		return nil
	}
	endpoints := make([]Endpoint, 0)
	for _, e := range port.Endpoints{
		endpoints = append(endpoints, *NewEndpointFromGRPC(e))
	}
	return &Port{
		Name:         port.Name,
		InternalPort: port.InternalPort,
		ExposedPort:  port.ExposedPort,
		Endpoints:    endpoints,
	}
}

func (p *Port) ToGRPC() *grpc_application_go.Port {
	endpoints := make([]*grpc_application_go.Endpoint, 0)

	for _, ep := range p.Endpoints {
		endpoints = append(endpoints, ep.ToGRPC())
	}

	return &grpc_application_go.Port{
		Name:         p.Name,
		InternalPort: p.InternalPort,
		ExposedPort:  p.ExposedPort,
		Endpoints:    endpoints,
	}
}

type ConfigFile struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// ConfigFileId with the config file identifier.
	ConfigFileId string `json:"config_file_id,omitempty"`
	// Name for the config file
	Name string `json:"name,omitempty"`
	// Content of the configuration file.
	Content []byte `json:"content,omitempty"`
	// MountPath of the configuration file in the service instance.
	MountPath string `json:"mount_path,omitempty"`
}


func NewConfigFileFromGRPC(appDescriptorID string, config * grpc_application_go.ConfigFile) * ConfigFile {
	if config == nil {
		return nil
	}
	return &ConfigFile{
		OrganizationId:  config.OrganizationId,
		AppDescriptorId: appDescriptorID,
		ConfigFileId:    config.ConfigFileId,
		Name:            config.Name,
		Content:         config.Content,
		MountPath:       config.MountPath,
	}
}


func (cf *ConfigFile) ToGRPC() *grpc_application_go.ConfigFile {
	return &grpc_application_go.ConfigFile{
		OrganizationId:  cf.OrganizationId,
		AppDescriptorId: cf.AppDescriptorId,
		ConfigFileId:    cf.ConfigFileId,
		Content:         cf.Content,
		MountPath:       cf.MountPath,
	}
}

type Service struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// ServiceGroupId with the service identifier.
	ServiceGroupId string `json:"service_id,omitempty"`
	// GroupServiceId with the service identifier.
	ServiceId string `json:"service_id,omitempty"`
	// Name of the service.
	Name string `json:"name,omitempty"`
	// ServiceType represents the underlying technology of the service to be launched.
	Type ServiceType `json:"type,omitempty"`
	// Image contains the URL/name of the image to be executed.
	Image string `json:"image,omitempty"`
	// ImageCredentials with the data required to access the repository the image is available at.
	Credentials * ImageCredentials `json:"credentials,omitempty"`
	// DeploySpecs with the resource specs required by the service.
	Specs * DeploySpecs `json:"specs,omitempty"`
	// Storage restrictions
	Storage []Storage `json:"storage,omitempty"`
	// ExposedPorts contains the list of ports exposed by the current service.
	ExposedPorts []Port `json:"exposed_ports,omitempty"`
	// EnvironmentVariables defines a key-value map of environment variables and values that will be passed to all
	// running services.
	EnvironmentVariables map[string]string `json:"environment_variables,omitempty"`
	// Configs contains the configuration files required by the service.
	Configs []ConfigFile `json:"configs,omitempty"`
	// Labels with the user defined labels.
	Labels map[string]string `json:"labels,omitempty"`
	// DeployAfter contains the list of services that must be running before launching a service.
	DeployAfter []string `json:"deploy_after,omitempty"`
	// RunArguments containts a list of arguments
	RunArguments [] string `json:"run_arguments" omitempty"`

}


func (s *Service) ToGRPC() *grpc_application_go.Service {
	serviceType, _ := ServiceTypeToGRPC[s.Type]
	storage := make([]*grpc_application_go.Storage, 0)
	for _, s := range s.Storage {
		storage = append(storage, s.ToGRPC())
	}
	exposedPorts := make([]*grpc_application_go.Port, 0)
	for _, ep := range s.ExposedPorts {
		exposedPorts = append(exposedPorts, ep.ToGRPC())
	}
	configs := make([]*grpc_application_go.ConfigFile, 0)
	for _, c := range s.Configs {
		configs = append(configs, c.ToGRPC())
	}
	return &grpc_application_go.Service{
		OrganizationId:       s.OrganizationId,
		AppDescriptorId:      s.AppDescriptorId,
		ServiceId:            s.ServiceId,
		Name:                 s.Name,
		ServiceGroupId:       s.ServiceGroupId,
		Type:                 serviceType,
		Image:                s.Image,
		Credentials:          s.Credentials.ToGRPC(),
		Specs:                s.Specs.ToGRPC(),
		Storage:              storage,
		ExposedPorts:         exposedPorts,
		EnvironmentVariables: s.EnvironmentVariables,
		Configs:              configs,
		Labels:               s.Labels,
		DeployAfter:          s.DeployAfter,
		RunArguments:         s.RunArguments,
	}
}

func NewServiceFromGRPC(appDescriptorID string, service *grpc_application_go.Service) * Service {
	if service == nil{
		return nil
	}

	storage := make([]Storage, 0)
	for _, s := range service.Storage {
		storage = append(storage, *NewStorageFromGRPC(s))
	}
	ports := make([]Port, 0)
	for _, p := range service.ExposedPorts {
		ports = append(ports, *NewPortFromGRPC(p))
	}
	configs := make([]ConfigFile, 0)
	for _, cf := range service.Configs {
		configs = append(configs, *NewConfigFileFromGRPC(appDescriptorID, cf))
	}

	serviceType, _ := ServiceTypeFromGRPC[service.Type]
	return &Service{
		OrganizationId:       service.OrganizationId,
		AppDescriptorId:      appDescriptorID,
 		ServiceGroupId:       service.ServiceGroupId,
		ServiceId:            service.ServiceId,
		Name:                 service.Name,
		Type:                 serviceType,
		Image:                service.Image,
		Credentials:          NewImageCredentialsFromGRPC(service.Credentials),
		Specs:                NewDeploySpecsFromGRPC(service.Specs),
		Storage:              storage,
		ExposedPorts:         ports,
		EnvironmentVariables: service.EnvironmentVariables,
		Configs:              configs,
		Labels:               service.Labels,
		DeployAfter:          service.DeployAfter,
		RunArguments:         service.RunArguments,
	}
}

func (s * Service) ToServiceInstance(appInstanceID string) * ServiceInstance {

	return &ServiceInstance{
		OrganizationId:       s.OrganizationId,
		AppDescriptorId:      s.AppDescriptorId,
		AppInstanceId:        appInstanceID,
		ServiceId:            s.ServiceId,
		Name:                 s.Name,
		Type:                 s.Type,
		Image:                s.Image,
		Credentials:          s.Credentials,
		Specs:                s.Specs,
		Storage:              s.Storage,
		ExposedPorts:         s.ExposedPorts,
		EnvironmentVariables: s.EnvironmentVariables,
		Configs:              s.Configs,
		Labels:               s.Labels,
		DeployAfter:          s.DeployAfter,
		Status:               ServiceWaiting,
		RunArguments:         s.RunArguments,
	}
}

type ServiceStatus int

const (
	ServiceScheduled ServiceStatus = iota + 1
	ServiceWaiting
	ServiceDeploying
	ServiceRunning
	ServiceError
)

var ServiceStatusToGRPC = map[ServiceStatus]grpc_application_go.ServiceStatus{
	ServiceScheduled:    grpc_application_go.ServiceStatus_SERVICE_SCHEDULED,
	ServiceWaiting: grpc_application_go.ServiceStatus_SERVICE_WAITING,
	ServiceDeploying:       grpc_application_go.ServiceStatus_SERVICE_DEPLOYING,
	ServiceRunning:        grpc_application_go.ServiceStatus_SERVICE_RUNNING,
	ServiceError: grpc_application_go.ServiceStatus_SERVICE_ERROR,
}

var ServiceStatusFromGRPC = map[grpc_application_go.ServiceStatus]ServiceStatus{
	grpc_application_go.ServiceStatus_SERVICE_SCHEDULED : ServiceScheduled,
	grpc_application_go.ServiceStatus_SERVICE_WAITING : ServiceWaiting,
	grpc_application_go.ServiceStatus_SERVICE_DEPLOYING : ServiceDeploying,
	grpc_application_go.ServiceStatus_SERVICE_RUNNING : ServiceRunning,
	grpc_application_go.ServiceStatus_SERVICE_ERROR : ServiceError,
}

// ServiceInstance----------

type ServiceInstance struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// AppInstanceId with the application instance identifier.
	AppInstanceId string `json:"app_instance_id,omitempty"`
	// ServiceInstanceId with the service instance identifier.
	ServiceInstanceId string `json:"service_instance_id,omitempty"`
	// GroupServiceId with the service identifier.
	ServiceId string `json:"service_id,omitempty"`
	// ServiceGroupInstanceId with the service group identifier.
	ServiceGroupInstanceId string `json:"service_group_instance_id,omitempty"`
	// ServiceGroupId with the service group identifier.
	ServiceGroupId string `json:"service_group_id,omitempty"`
	// Name of the service.
	Name string `json:"name,omitempty"`
	// ServiceType represents the underlying technology of the service to be launched.
	Type ServiceType `json:"type,omitempty"`
	// Image contains the URL/name of the image to be executed.
	Image string `json:"image,omitempty"`
	// ImageCredentials with the data required to access the repository the image is available at.
	Credentials * ImageCredentials `json:"credentials,omitempty"`
	// DeploySpecs with the resource specs required by the service.
	Specs * DeploySpecs `json:"specs,omitempty"`
	// Storage restrictions
	Storage []Storage `json:"storage,omitempty"`
	// ExposedPorts contains the list of ports exposed by the current service.
	ExposedPorts []Port `json:"exposed_ports,omitempty"`
	// EnvironmentVariables defines a key-value map of environment variables and values that will be passed to all
	// running services.
	EnvironmentVariables map[string]string `json:"environment_variables,omitempty"`
	// Configs contains the configuration files required by the service.
	Configs []ConfigFile `json:"configs,omitempty"`
	// Labels with the user defined labels.
	Labels map[string]string `json:"labels,omitempty"`
	// DeployAfter contains the list of services that must be running before launching a service.
	DeployAfter []string `json:"deploy_after,omitempty"`
	// Status of the deployed service
	Status ServiceStatus `json:"status,omitempty"`
	// Endpoints exposed to the users by the service.
	Endpoints []EndpointInstance `json:"endpoints,omitempty"`
	// DeployedOnClusterId specifies which is the cluster where the service is running.
	DeployedOnClusterID string `json:"deployed_on_cluster_id,omitempty"`
	// RunArguments containts a list of arguments
	RunArguments [] string `json:"run_arguments" omitempty"`
	// Information about new service instance
	Info string `json:"info,omitempty"`

}

func (si *ServiceInstance) ToGRPC() *grpc_application_go.ServiceInstance {
	serviceType, _ := ServiceTypeToGRPC[si.Type]
	serviceStatus, _ := ServiceStatusToGRPC[si.Status]
	storage := make([]*grpc_application_go.Storage, 0)
	for _, s := range si.Storage {
		storage = append(storage, s.ToGRPC())
	}
	exposedPorts := make([]*grpc_application_go.Port, 0)
	for _, ep := range si.ExposedPorts {
		exposedPorts = append(exposedPorts, ep.ToGRPC())
	}
	configs := make([]*grpc_application_go.ConfigFile, 0)
	for _, c := range si.Configs {
		configs = append(configs, c.ToGRPC())
	}
	endpoints := make([]*grpc_application_go.EndpointInstance,0)
	for _, e := range si.Endpoints {
		endpoints = append(endpoints, e.ToGRPC())
	}
	return &grpc_application_go.ServiceInstance{
		OrganizationId:       si.OrganizationId,
		AppDescriptorId:      si.AppDescriptorId,
		AppInstanceId:        si.AppInstanceId,
		ServiceId:           si.ServiceId,
		Name:                 si.Name,
		Type:                 serviceType,
		Image:                si.Image,
		Credentials:          si.Credentials.ToGRPC(),
		Specs:                si.Specs.ToGRPC(),
		Storage:              storage,
		ExposedPorts:         exposedPorts,
		EnvironmentVariables: si.EnvironmentVariables,
		Configs:              configs,
		Labels:               si.Labels,
		DeployAfter:          si.DeployAfter,
		Status:               serviceStatus,
		RunArguments:         si.RunArguments,
		ServiceGroupInstanceId: si.ServiceGroupInstanceId,
		Info:                 si.Info,
		Endpoints:            endpoints,
		ServiceInstanceId:    si.ServiceInstanceId,
		ServiceGroupId:       si.ServiceGroupId,
		DeployedOnClusterId: si.DeployedOnClusterID,
	}

}

func NewServiceInstanceFromGRPC(ins *grpc_application_go.ServiceInstance) ServiceInstance {

	endpoints := make([]EndpointInstance, len(ins.Endpoints))
	for i, e := range ins.Endpoints {
		endpoints[i] = NewEndpointInstanceFromGRPC(e)
	}
	storages := make([]Storage, len(ins.Storage))
	for i, s := range ins.Storage {
		storages[i] = *NewStorageFromGRPC(s)
	}
	configs := make([]ConfigFile, len(ins.Configs))
	for i, c := range ins.Configs {
		configs[i] = *NewConfigFileFromGRPC(ins.AppDescriptorId, c)
	}
	exposedPorts := make([]Port, len(ins.ExposedPorts))
	for i, p := range ins.ExposedPorts {
		exposedPorts[i] = *NewPortFromGRPC(p)
	}

	return ServiceInstance{
		Status: ServiceStatusFromGRPC[ins.Status],
		Type: ServiceTypeFromGRPC[ins.Type],
		Info: ins.Info,
		OrganizationId: ins.OrganizationId,
		AppInstanceId: ins.AppInstanceId,
		AppDescriptorId: ins.AppDescriptorId,
		ServiceGroupInstanceId: ins.ServiceGroupInstanceId,
		ServiceGroupId: ins.ServiceGroupId,
		Labels: ins.Labels,
		Name: ins.Name,
		Specs: NewDeploySpecsFromGRPC(ins.Specs),
		EnvironmentVariables: ins.EnvironmentVariables,
		ServiceInstanceId: ins.ServiceInstanceId,
		ServiceId: ins.ServiceId,
		DeployedOnClusterID: ins.DeployedOnClusterId,
		RunArguments: ins.RunArguments,
		DeployAfter: ins.DeployAfter,
		Image: ins.Image,
		Endpoints: endpoints,
		Storage: storages,
		Configs: configs,
		Credentials: NewImageCredentialsFromGRPC(ins.Credentials),
		ExposedPorts: exposedPorts,
	}
}

// EndpointInstance----------

type EndpointInstance struct {
	// EndpointInstanceId
	EnpointInstanceId string `json:"endpoint_instance_id,omitempty"`
	// EndpointType
	EndpointType EndpointType `json:"endpoint_type,omitempty"`
	// FQDN
	FQDN string `json:"fqdn,omitempty"`
}

func(e *EndpointInstance) ToGRPC() * grpc_application_go.EndpointInstance {
	return &grpc_application_go.EndpointInstance{
		Type: EndpointTypeToGRPC[e.EndpointType],
		EndpointInstanceId: e.EnpointInstanceId,
		Fqdn: e.FQDN,
	}
}

func NewEndpointInstanceFromGRPC(endpoint *grpc_application_go.EndpointInstance) EndpointInstance {
	return EndpointInstance{
		FQDN: endpoint.Fqdn,
		EnpointInstanceId: endpoint.EndpointInstanceId,
		EndpointType: EndpointTypeFromGRPC[endpoint.Type],
	}
}

// ----------

// ----------

type AppDescriptor struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// Name of the application.
	Name string `json:"name,omitempty"`
	// ConfigurationOptions defines a key-value map of configuration options.
	ConfigurationOptions map[string]string `json:"configuration_options,omitempty"`
	// EnvironmentVariables defines a key-value map of environment variables and values that will be passed to all
	// running services.
	EnvironmentVariables map[string]string `json:"environment_variables,omitempty"`
	// Labels defined by the user.
	Labels map[string]string `json:"labels,omitempty"`
	// Rules that define the connectivity between the elements of an application.
	Rules []SecurityRule `json:"rules,omitempty"`
	// Groups with the Service collocation strategies.
	Groups []ServiceGroup `json:"groups,omitempty"`
}

func NewAppDescriptor(organizationID string, appDescriptorID string, name string,
	configOptions map[string]string, envVars map[string]string,
	labels map[string]string, rules []SecurityRule, groups []ServiceGroup) *AppDescriptor {
	return &AppDescriptor{
		OrganizationId: organizationID,
		AppDescriptorId: appDescriptorID,
		Name: name,
		ConfigurationOptions: configOptions,
		EnvironmentVariables: envVars,
		Labels: labels,
		Rules: rules,
		Groups: groups,
	}
}

func (d *AppDescriptor) ToGRPC() *grpc_application_go.AppDescriptor {

	rules := make([]*grpc_application_go.SecurityRule, 0)
	for _, r := range d.Rules {
		rules = append(rules, r.ToGRPC())
	}
	groups := make([]*grpc_application_go.ServiceGroup, 0)
	for _, g := range d.Groups {
		groups = append(groups, g.ToGRPC())
	}

	return &grpc_application_go.AppDescriptor{
		OrganizationId:       d.OrganizationId,
		AppDescriptorId:      d.AppDescriptorId,
		Name:                 d.Name,
		ConfigurationOptions: d.ConfigurationOptions,
		EnvironmentVariables: d.EnvironmentVariables,
		Labels:               d.Labels,
		Rules:                rules,
		Groups:               groups,
	}
}

func NewAppDescriptorFromGRPC(app *grpc_application_go.AppDescriptor) AppDescriptor {
	groups := make([]ServiceGroup,0)
	for _, g := range app.Groups {
		groups = append(groups, NewServiceGroupFromGRPC(g))
	}
	rules := make([]SecurityRule,0)
	for _, r := range app.Rules {
		rules = append(rules, NewSecurityRuleFromGRPC(r))
	}
	return AppDescriptor{
		OrganizationId: app.OrganizationId,
		AppDescriptorId: app.AppDescriptorId,
		Labels: app.Labels,
		EnvironmentVariables: app.EnvironmentVariables,
		ConfigurationOptions: app.ConfigurationOptions,
		Name: app.Name,
		Groups: groups,
		Rules: rules,
	}
}


type ApplicationStatus int

const (
	Queued ApplicationStatus = iota +1
	Planning
	Scheduled
	Deploying
	Running
	Incomplete
	PlanningError
	DeploymentError
	Error
)

var AppStatusToGRPC = map[ApplicationStatus]grpc_application_go.ApplicationStatus{
	Queued: grpc_application_go.ApplicationStatus_QUEUED,
	Planning: grpc_application_go.ApplicationStatus_PLANNING,
	Scheduled: grpc_application_go.ApplicationStatus_SCHEDULED,
	Deploying: grpc_application_go.ApplicationStatus_DEPLOYING,
	Running: grpc_application_go.ApplicationStatus_RUNNING,
	Incomplete: grpc_application_go.ApplicationStatus_INCOMPLETE,
	PlanningError: grpc_application_go.ApplicationStatus_PLANNING_ERROR,
	DeploymentError: grpc_application_go.ApplicationStatus_DEPLOYMENT_ERROR,
	Error: grpc_application_go.ApplicationStatus_ERROR,
}

var AppStatusFromGRPC = map[grpc_application_go.ApplicationStatus]ApplicationStatus{
	grpc_application_go.ApplicationStatus_QUEUED:Queued,
	grpc_application_go.ApplicationStatus_PLANNING:Planning,
	grpc_application_go.ApplicationStatus_SCHEDULED:Scheduled,
	grpc_application_go.ApplicationStatus_DEPLOYING:Deploying,
	grpc_application_go.ApplicationStatus_RUNNING:Running,
	grpc_application_go.ApplicationStatus_INCOMPLETE:Incomplete,
	grpc_application_go.ApplicationStatus_PLANNING_ERROR:PlanningError,
	grpc_application_go.ApplicationStatus_DEPLOYMENT_ERROR:DeploymentError,
	grpc_application_go.ApplicationStatus_ERROR:Error,
}

// AppInstance----------

type AppInstance struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// AppInstanceId with the application instance identifier.
	AppInstanceId string `json:"app_instance_id,omitempty"`
	// Name of the application.
	Name string `json:"name,omitempty"`
	// ConfigurationOptions defines a key-value map of configuration options.
	ConfigurationOptions map[string]string `json:"configuration_options,omitempty"`
	// EnvironmentVariables defines a key-value map of environment variables and values that will be passed to all
	// running services.
	EnvironmentVariables map[string]string `json:"environment_variables,omitempty"`
	// Labels defined by the user.
	Labels map[string]string `json:"labels,omitempty"`
	// Rules that define the connectivity between the elements of an application.
	Rules []SecurityRule `json:"rules,omitempty"`
	// Groups with the Service collocation strategies.
	Groups []ServiceGroupInstance `json:"groups,omitempty"`
	// Status of the deployed instance.
	Status  ApplicationStatus `json:"status,omitempty"`
	// Metadata
	Metadata []InstanceMetadata  `json: "metadata, omitempty"`
}


func (i *AppInstance) ToGRPC() *grpc_application_go.AppInstance {
	rules := make([]*grpc_application_go.SecurityRule, 0)
	for _, r := range i.Rules {
		rules = append(rules, r.ToGRPC())
	}
	groups := make([]*grpc_application_go.ServiceGroupInstance, 0)
	for _, g := range i.Groups {
		groups = append(groups, g.ToGRPC())
	}

	metadatas := make([]*grpc_application_go.InstanceMetadata,0)
	for _, m := range i.Metadata {
		metadatas = append(metadatas, m.ToGRPC())
	}

	status, _ := AppStatusToGRPC[i.Status]

	return &grpc_application_go.AppInstance{
		OrganizationId:       i.OrganizationId,
		AppDescriptorId:      i.AppDescriptorId,
		AppInstanceId:        i.AppInstanceId,
		Name:                 i.Name,
		ConfigurationOptions: i.ConfigurationOptions,
		EnvironmentVariables: i.EnvironmentVariables,
		Labels:               i.Labels,
		Rules:                rules,
		Groups:               groups,
		Status:               status,
		Metadata:             metadatas,
	}
}

func NewAppInstanceFromGRPC(app *grpc_application_go.AppInstance) AppInstance {
	groups := make([]ServiceGroupInstance,0)
	for _, g := range app.Groups {
		groups = append(groups, NewServiceGroupInstanceFromGRPC(g))
	}
	rules := make([]SecurityRule,0)
	for _, r := range app.Rules {
		rules = append(rules, NewSecurityRuleFromGRPC(r))
	}
	metadata := make([]InstanceMetadata,0)

	for _, m := range app.Metadata {
		metadata = append(metadata, NewInstanceMetadataFromGRPC(m))
	}

	return AppInstance{
		AppDescriptorId: app.AppDescriptorId,
		Name: app.Name,
		AppInstanceId: app.AppInstanceId,
		Labels: app.Labels,
		OrganizationId: app.OrganizationId,
		ConfigurationOptions: app.ConfigurationOptions,
		EnvironmentVariables: app.EnvironmentVariables,
		Status: AppStatusFromGRPC[app.Status],
		Groups: groups,
		Rules: rules,
		Metadata: metadata,
	}
}

// ----------
