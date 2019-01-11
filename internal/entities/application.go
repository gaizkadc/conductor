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

type SecurityRule struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// RuleId with the security rule identifier.
	RuleId string `json:"rule_id,omitempty"`
	// Name of the security rule.
	Name string `json:"name,omitempty"`
	// SourceServiceId defining the service onto which the security rule is defined.
	SourceServiceId string `json:"source_service_id,omitempty"`
	// SourcePort defining the port that is affected by the current rule.
	SourcePort int32 `json:"source_port,omitempty"`
	// Access level to that port defining who can access the port.
	Access PortAccess `json:"access,omitempty"`
	// AuthServices defining a list of services that can access the port.
	AuthServices []string `json:"auth_services,omitempty"`
	// DeviceGroups defining a list of device groups that can access the port.
	DeviceGroups []string `json:"device_groups,omitempty"`
}

func (sr *SecurityRule) ToGRPC() *grpc_application_go.SecurityRule {
	access, _ := PortAccessToGRPC[sr.Access]
	return &grpc_application_go.SecurityRule{
		OrganizationId:  sr.OrganizationId,
		AppDescriptorId: sr.AppDescriptorId,
		RuleId:          sr.RuleId,
		Name:            sr.Name,
		SourceServiceId: sr.SourceServiceId, SourcePort: sr.SourcePort,
		Access:       access,
		AuthServices: sr.AuthServices,
		DeviceGroups: sr.DeviceGroups,
	}
}

type ServiceGroup struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// ServiceGroupId with the group identifier.
	ServiceGroupId string `json:"service_group_id,omitempty"`
	// Name of the service group.
	Name string `json:"name,omitempty"`
	// Description of the service group.
	Description string `json:"description,omitempty"`
	// Services defining a list of service identifiers that belong to the group.
	Services []string `json:"services,omitempty"`
	// Policy indicating the deployment collocation policy.
	Policy CollocationPolicy `json:"policy,omitempty"`
}


func (sg *ServiceGroup) ToGRPC() *grpc_application_go.ServiceGroup {
	policy, _ := CollocationPolicyToGRPC[sg.Policy]
	return &grpc_application_go.ServiceGroup{
		OrganizationId:  sg.OrganizationId,
		AppDescriptorId: sg.AppDescriptorId,
		ServiceGroupId:  sg.ServiceGroupId,
		Name:            sg.Name,
		Description:     sg.Description,
		Services:        sg.Services,
		Policy:          policy,
	}
}

func (sg * ServiceGroup) ToServiceGroupInstance(appInstID string) * ServiceGroupInstance {
	return &ServiceGroupInstance{
		OrganizationId:   sg.OrganizationId,
		AppDescriptorId:  sg.AppDescriptorId,
		AppInstanceId:    appInstID,
		ServiceGroupId:   sg.ServiceGroupId,
		Name:             sg.Name,
		Description:      sg.Description,
		ServiceInstances: sg.Services,
		Policy:           sg.Policy,
	}
}

type ServiceGroupInstance struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// AppInstanceId with the application instance identifier.
	AppInstanceId string `json:"app_instance_id,omitempty"`
	// ServiceGroupId with the group identifier.
	ServiceGroupId string `json:"service_group_id,omitempty"`
	// Name of the service group.
	Name string `json:"name,omitempty"`
	// Description of the service group.
	Description string `json:"description,omitempty"`
	// ServicesInstances defining a list of service identifiers that belong to the group.
	ServiceInstances []string `json:"service_instances,omitempty"`
	// Policy indicating the deployment collocation policy.
	Policy               CollocationPolicy `json:"policy,omitempty"`
}

func (sgi *ServiceGroupInstance) ToGRPC() *grpc_application_go.ServiceGroupInstance {
	policy, _ := CollocationPolicyToGRPC[sgi.Policy]
	return &grpc_application_go.ServiceGroupInstance{
		OrganizationId:       sgi.OrganizationId,
		AppDescriptorId:      sgi.AppDescriptorId,
		AppInstanceId:        sgi.AppInstanceId,
		ServiceGroupId:       sgi.ServiceGroupId,
		Name:                 sgi.Name,
		Description:          sgi.Description,
		ServiceInstances:     sgi.ServiceInstances,
		Policy:               policy,
	}
}


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
	// ServiceId with the service identifier.
	ServiceId string `json:"service_id,omitempty"`
	// Name of the service.
	Name string `json:"name,omitempty"`
	// Description of the service.
	Description string `json:"description,omitempty"`
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
		Description:          s.Description,
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
		ServiceId:            service.ServiceId,
		Name:                 service.Name,
		Description:          service.Description,
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
		Description:          s.Description,
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

type ServiceInstance struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// AppInstanceId with the application instance identifier.
	AppInstanceId string `json:"app_instance_id,omitempty"`
	// ServiceId with the service identifier.
	ServiceId string `json:"service_id,omitempty"`
	// Name of the service.
	Name string `json:"name,omitempty"`
	// Description of the service.
	Description string `json:"description,omitempty"`
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
	// RunArguments containts a list of arguments
	RunArguments [] string `json:"run_arguments" omitempty"`

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
	return &grpc_application_go.ServiceInstance{
		OrganizationId:       si.OrganizationId,
		AppDescriptorId:      si.AppDescriptorId,
		AppInstanceId:        si.AppInstanceId,
		ServiceId:           si.ServiceId,
		Name:                 si.Name,
		Description:          si.Description,
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
	}

}

type AppDescriptor struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// Name of the application.
	Name string `json:"name,omitempty"`
	// Description of the application.
	Description string `json:"description,omitempty"`
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
	// Services of the application.
	Services []Service `json:"services,omitempty"`
}

func NewAppDescriptor(organizationID string, appDescriptorID string, name string, description string,
	configOptions map[string]string, envVars map[string]string,
	labels map[string]string,
	rules []SecurityRule, groups []ServiceGroup, services []Service) *AppDescriptor {
	return &AppDescriptor{
		organizationID, appDescriptorID,
		name, description,
		configOptions,
		envVars,
		labels,
		rules,
		groups,
		services,
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
	services := make([]*grpc_application_go.Service, 0)
	for _, s := range d.Services {
		services = append(services, s.ToGRPC())
	}

	return &grpc_application_go.AppDescriptor{
		OrganizationId:       d.OrganizationId,
		AppDescriptorId:      d.AppDescriptorId,
		Name:                 d.Name,
		Description:          d.Description,
		ConfigurationOptions: d.ConfigurationOptions,
		EnvironmentVariables: d.EnvironmentVariables,
		Labels:               d.Labels,
		Rules:                rules,
		Groups:               groups,
		Services:             services,
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

type AppInstance struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// AppInstanceId with the application instance identifier.
	AppInstanceId string `json:"app_instance_id,omitempty"`
	// Name of the application.
	Name string `json:"name,omitempty"`
	// Description of the application.
	Description string `json:"description,omitempty"`
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
	// Services of the application.
	Services []ServiceInstance `json:"services,omitempty"`
	// Status of the deployed instance.
	Status  ApplicationStatus `json:"status,omitempty"`
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
	services := make([]*grpc_application_go.ServiceInstance, 0)
	for _, s := range i.Services {
		services = append(services, s.ToGRPC())
	}

	status, _ := AppStatusToGRPC[i.Status]

	return &grpc_application_go.AppInstance{
		OrganizationId:       i.OrganizationId,
		AppDescriptorId:      i.AppDescriptorId,
		AppInstanceId:        i.AppInstanceId,
		Name:                 i.Name,
		Description:          i.Description,
		ConfigurationOptions: i.ConfigurationOptions,
		EnvironmentVariables: i.EnvironmentVariables,
		Labels:               i.Labels,
		Rules:                rules,
		Groups:               groups,
		Services:             services,
		Status:               status,
	}
}
