/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package entities

import (
	"github.com/nalej/derrors"
	pbApplication "github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-application-network-go"
	pbConductor "github.com/nalej/grpc-conductor-go"
	"sort"
	"strings"
	"time"
)

// DeploymentsScore for a set of potential deployments.
type DeploymentScore struct {
	// Score for every evaluated cluster
	DeploymentsScore []ClusterDeploymentScore `json:"scoring,omitempty"`
	// Total number of evaluated clusters
	NumEvaluatedClusters int `json:"num_evaluated_clusters,omitempty"`
}

func NewClustersScore() DeploymentScore {
	return DeploymentScore{NumEvaluatedClusters: 0, DeploymentsScore: make([]ClusterDeploymentScore, 0)}
}

// AddClusterScore appends a cluster score and updates the set of total evaluated clusters
func (c *DeploymentScore) AddClusterScore(score ClusterDeploymentScore) {
	c.DeploymentsScore = append(c.DeploymentsScore, score)
	c.NumEvaluatedClusters = c.NumEvaluatedClusters + 1
}

// Cluster deployment score ------

// Combinations of deployments evaluated by a cluster. Every combination is a set of service groups
// and their scoring.
type ClusterDeploymentScore struct {
	// Cluster that returned this score
	ClusterId string `json:"cluster_id,omitempty"`
	// Score of different service group combinations
	// The key for this map is a concatenation of service group ids. The result is the concatenation of the service
	// group ids after sorting.
	Scores map[string]float32 `json: "scores,omitempty"`
}

func NewClusterDeploymentScore(clusterId string) ClusterDeploymentScore {
	return ClusterDeploymentScore{
		ClusterId: clusterId,
		Scores:    make(map[string]float32, 0),
	}
}

// Add the score for a set of service groups.
func (cds *ClusterDeploymentScore) AddScore(serviceGroupIds []string, score float32) {
	sort.Strings(serviceGroupIds)
	// The key is the concatenation of the ids
	newKey := strings.Join(serviceGroupIds, "")
	cds.Scores[newKey] = score
}

// End of cluster deployment score ------

// Objects describing received deployment requests. These objects are designed to be stored into
// a provider structure such as a queue.
type DeploymentRequest struct {
	RequestId      string     `json:"request_id,omitempty"`
	OrganizationId string     `json:"organization_id,omitempty"`
	ApplicationId  string     `json:"application_id,omitempty"`
	InstanceId     string     `json:"instance_id,omitempty"`
	NumRetries     int32      `json:"num_retries,omitempty"`
	TimeRetry      *time.Time `json:"time_retry,omitempty"`
	// The AppInstanceId is internally used to link this request with a certain instance
	AppInstanceId string                                            `json:"app_instance_id,omitempty"`
	Connections   []*grpc_application_network_go.ConnectionInstance `json:"connections,omitempty"`
}

// Fragment deployment Status definition

type DeploymentFragmentStatus int

const (
	FRAGMENT_WAITING DeploymentFragmentStatus = iota
	FRAGMENT_DEPLOYING
	FRAGMENT_DONE
	FRAGMENT_ERROR
	FRAGMENT_RETRYING
	FRAGMENT_TERMINATING
)

var DeploymentStatusToGRPC = map[pbConductor.DeploymentFragmentStatus]DeploymentFragmentStatus{
	pbConductor.DeploymentFragmentStatus_WAITING:     FRAGMENT_WAITING,
	pbConductor.DeploymentFragmentStatus_DEPLOYING:   FRAGMENT_DEPLOYING,
	pbConductor.DeploymentFragmentStatus_DONE:        FRAGMENT_DONE,
	pbConductor.DeploymentFragmentStatus_ERROR:       FRAGMENT_ERROR,
	pbConductor.DeploymentFragmentStatus_RETRYING:    FRAGMENT_RETRYING,
	pbConductor.DeploymentFragmentStatus_TERMINATING: FRAGMENT_TERMINATING,
}

// Plan related entries

// Data structure defining a deployment plan for Nalej applications.
type DeploymentPlan struct {
	// Identifier for this plan
	DeploymentId string `json:"deployment_id,omitempty"`
	// OrganizationId this deployment belongs to
	OrganizationId string `json:"organization_id,omitempty"`
	// AppInstanceId for the instance of the application to run
	AppInstanceId string `json:"app_instance_id,omitempty"`
	// Fragments this plan is made of
	Fragments []DeploymentFragment `json:"fragments,omitempty"`
	// Associated deployment request
	DeploymentRequest *DeploymentRequest `json:"deployment_request,omitempty"`
}

// Start deployment fragment definition ----

// Data structure representing the components of a plan that will be
// spread across the available clusters.
type DeploymentFragment struct {
	// OrganizationId this deployment belongs to
	OrganizationId string `json:"organization_id,omitempty"`
	// OrganizationName this deployment belongs to
	OrganizationName string `json:"organization_name,omitempty"`
	// AppDescriptorId
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// AppDescriptorName with the name of the descriptor
	AppDescriptorName string `json:"app_descriptor_name,omitempty"`
	// AppInstanceId for the instance of the application to run
	AppInstanceId string `json:"app_instance_id,omitempty"`
	// AppNamed for the instance of the application to run
	AppName string `json:"app_name,omitempty"`
	// Identifier for this plan
	DeploymentId string `json:"deployment_id,omitempty"`
	// Fragment id
	FragmentId string `json:"fragment_id,omitempty"`
	// Cluster id for the deployment target
	ClusterId string `json:"cluster_id,omitempty"`
	// Nalej variables
	NalejVariables map[string]string `json:"nalej_variables,omitempty"`
	// Deployment stages belonging to this fragment
	Stages []DeploymentStage `json:"stages,omitempty"`
	// Identifier for the ZtNetworkID. This is a value only used by conductor.
	ZtNetworkID string `json:"stages,omitempty"`
	// Status for this deployment fragment
	Status DeploymentFragmentStatus `json:"status"`
}

func (df *DeploymentFragment) ToGRPC() *pbConductor.DeploymentFragment {
	convertedStages := make([]*pbConductor.DeploymentStage, len(df.Stages))
	for i, serv := range df.Stages {
		convertedStages[i] = serv.ToGRPC()
	}
	result := pbConductor.DeploymentFragment{
		OrganizationId:   df.OrganizationId,
		FragmentId:       df.FragmentId,
		AppDescriptorId:  df.AppDescriptorId,
		AppDescriptorName:df.AppDescriptorName,
		AppInstanceId:    df.AppInstanceId,
		OrganizationName: df.OrganizationName,
		AppInstanceName:  df.AppName,
		DeploymentId:     df.DeploymentId,
		NalejVariables:   df.NalejVariables,
		Stages:           convertedStages,
		ClusterId:        df.ClusterId,
	}
	return &result
}

// ----

type UndeployRequest struct {
	// OrganizationId this deployment belongs to
	OrganizationId string `json:"organization_id,omitempty"`
	// AppInstanceId for the instance of the application to run
	AppInstanceId string `json:"app_instance_id,omitempty"`
}

// ------------------------------------ //
// -- Deployment fragment definition -- //
// ------------------------------------ //
// DeviceGroupSecurityRuleInstance
type DeviceGroupSecurityRuleInstance struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// RuleId with the security rule identifier.
	RuleId string `json:"rule_id,omitempty"`
	// TargetServiceGroupId with the group identifier as provided by the user.
	TargetServiceGroupId string `json:"target_service_group_id,omitempty"`
	// TargetServiceGroupInstanceId with the group identifier provided by the system.
	TargetServiceGroupInstanceId string `json:"target_service_group_instance_id,omitempty"`
	// TargetServiceId with the service identifier as provided by the user.
	TargetServiceId string `json:"target_service_id,omitempty"`
	// TargetServiceInstanceId with the service identifier as provided by the system.
	TargetServiceInstanceId string `json:"target_service_instance_id,omitempty"`
	// TargetPort defining the port that is affected by the current rule.
	TargetPort int32 `json:"target_port,omitempty"`
	// DeviceGroupIds with the identifiers of the device groups that have access to the service.
	DeviceGroupIds []string `json:"device_group_ids,omitempty"`
	// DeviceGroupJWTSecrets with the secrets of those groups so that JWT can be enforced by the apps.
	DeviceGroupJwtSecrets []string `json:"device_group_jwt_secrets,omitempty"`
}

func NewDeviceGroupSecurityRuleInstance(service pbApplication.ServiceInstance, rule SecurityRule, jwtSecrets []string) *DeviceGroupSecurityRuleInstance {
	return &DeviceGroupSecurityRuleInstance{
		OrganizationId:               service.OrganizationId,
		AppDescriptorId:              service.AppDescriptorId,
		RuleId:                       rule.RuleId,
		TargetServiceGroupId:         service.ServiceGroupId,
		TargetServiceGroupInstanceId: service.ServiceGroupInstanceId,
		TargetServiceId:              service.ServiceId,
		TargetServiceInstanceId:      service.ServiceInstanceId,
		TargetPort:                   rule.TargetPort,
		DeviceGroupIds:               rule.DeviceGroupIds,
		DeviceGroupJwtSecrets:        jwtSecrets,
	}
}

func (dg *DeviceGroupSecurityRuleInstance) ToGRPC() *pbConductor.DeviceGroupSecurityRuleInstance {
	return &pbConductor.DeviceGroupSecurityRuleInstance{
		OrganizationId:               dg.OrganizationId,
		AppDescriptorId:              dg.AppDescriptorId,
		RuleId:                       dg.RuleId,
		TargetServiceGroupId:         dg.TargetServiceGroupId,
		TargetServiceGroupInstanceId: dg.TargetServiceGroupInstanceId,
		TargetServiceId:              dg.TargetServiceId,
		TargetServiceInstanceId:      dg.TargetServiceInstanceId,
		TargetPort:                   dg.TargetPort,
		DeviceGroupIds:               dg.DeviceGroupIds,
		DeviceGroupJwtSecrets:        dg.DeviceGroupJwtSecrets,
	}
}

// PublicSecurityRuleInstance
type PublicSecurityRuleInstance struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// RuleId with the security rule identifier.
	RuleId string `json:"rule_id,omitempty"`
	// TargetServiceGroupId with the group identifier as provided by the user.
	TargetServiceGroupId string `json:"target_service_group_id,omitempty"`
	// TargetServiceGroupInstanceId with the group identifier provided by the system.
	TargetServiceGroupInstanceId string `json:"target_service_group_instance_id,omitempty"`
	// TargetServiceId with the service identifier as provided by the user.
	TargetServiceId string `json:"target_service_id,omitempty"`
	// TargetServiceInstanceId with the service identifier as provided by the system.
	TargetServiceInstanceId string `json:"target_service_instance_id,omitempty"`
	// TargetPort defining the port that is affected by the current rule.
	TargetPort int32 `json:"target_port,omitempty"`
}

func NewPublicSercurityRuleInstance(service pbApplication.ServiceInstance, rule SecurityRule) *PublicSecurityRuleInstance {

	return &PublicSecurityRuleInstance{
		OrganizationId:               service.OrganizationId,
		AppDescriptorId:              service.AppDescriptorId,
		RuleId:                       rule.RuleId,
		TargetServiceGroupId:         service.ServiceGroupId,
		TargetServiceGroupInstanceId: service.ServiceGroupInstanceId,
		TargetServiceId:              service.ServiceId,
		TargetServiceInstanceId:      service.ServiceInstanceId,
		TargetPort:                   rule.TargetPort,
	}
}

func (pr *PublicSecurityRuleInstance) ToGRPC() *pbConductor.PublicSecurityRuleInstance {
	return &pbConductor.PublicSecurityRuleInstance{
		OrganizationId:               pr.OrganizationId,
		AppDescriptorId:              pr.AppDescriptorId,
		RuleId:                       pr.RuleId,
		TargetServiceGroupId:         pr.TargetServiceGroupId,
		TargetServiceGroupInstanceId: pr.TargetServiceGroupInstanceId,
		TargetServiceId:              pr.TargetServiceId,
		TargetServiceInstanceId:      pr.TargetServiceInstanceId,
		TargetPort:                   pr.TargetPort,
	}
}

// Every deployment stage a fragment is made of.
type DeploymentStage struct {
	// Fragment id
	FragmentId string `json:"fragment_id,omitempty"`
	// Stage id
	StageId string `json:"stage_id,omitempty"`
	// Set of services
	Services []ServiceInstance `json:"services,omitempty"`
	// DeviceGroupRules with the security rules affecting device group access.
	DeviceGroupRules []DeviceGroupSecurityRuleInstance `json:"device_group_rules,omitempty"`
	// PublicRules with the security rules related to public access of a given endpoint.
	PublicRules []PublicSecurityRuleInstance `json:"public_rules,omitempty"`
}

func (ds *DeploymentStage) ToGRPC() *pbConductor.DeploymentStage {
	convertedServices := make([]*pbConductor.ServiceInstance, len(ds.Services))
	for i, serv := range ds.Services {
		convertedServices[i] = serv.toConductorGRPC()
	}
	publicRules := make([]*pbConductor.PublicSecurityRuleInstance, len(ds.PublicRules))
	for i, rules := range ds.PublicRules {
		publicRules[i] = rules.ToGRPC()
	}
	deviceGroupRules := make([]*pbConductor.DeviceGroupSecurityRuleInstance, len(ds.DeviceGroupRules))
	for i, deviceRules := range ds.DeviceGroupRules {
		deviceGroupRules[i] = deviceRules.ToGRPC()
	}
	result := pbConductor.DeploymentStage{
		FragmentId:       ds.FragmentId,
		StageId:          ds.StageId,
		Services:         convertedServices,
		PublicRules:      publicRules,
		DeviceGroupRules: deviceGroupRules,
	}
	return &result
}

//ValidDeploymentRequest validates request data before executing a deployment
func ValidDeploymentRequest(request *pbConductor.DeploymentRequest) derrors.Error {
	if request == nil {
		return derrors.NewInvalidArgumentError(invalidRequest)
	}
	if request.RequestId == "" {
		return derrors.NewInvalidArgumentError(emptyRequestID)
	}
	if request.Name == "" {
		return derrors.NewInvalidArgumentError(emptyName)
	}
	if request.AppInstanceId == nil {
		return derrors.NewInvalidArgumentError(emptyAppID)
	}
	if request.AppInstanceId.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationID)
	}
	if request.AppInstanceId.AppInstanceId == "" {
		return derrors.NewInvalidArgumentError(emptyAppInstanceID)
	}
	return nil
}

//ValidDeploymentFragmentID validates request fragment ID is not empty
func ValidDeploymentFragmentID(fragmentID string) derrors.Error {
	if fragmentID == "" {
		return derrors.NewInvalidArgumentError(emptyFragmentID)
	}
	return nil
}

//ValidDeploymentFragmentUpdateRequest validates request data before updating a deployment fragment
func ValidDeploymentFragmentUpdateRequest(request *pbConductor.DeploymentFragmentUpdateRequest) derrors.Error {
	if request == nil {
		return derrors.NewInvalidArgumentError(invalidRequest)
	}
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationID)
	}
	if err := ValidDeploymentFragmentID(request.FragmentId); err != nil {
		return err
	}
	if request.AppInstanceId == "" {
		return derrors.NewInvalidArgumentError(emptyAppInstanceID)
	}
	if request.ClusterId == "" {
		return derrors.NewInvalidArgumentError(emptyClusterID)
	}
	return nil
}

//ValidUndeployRequest validates request data before executing an undeployment
func ValidUndeployRequest(request *pbConductor.UndeployRequest) derrors.Error {
	if request == nil {
		return derrors.NewInvalidArgumentError(invalidRequest)
	}
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationID)
	}
	if request.AppInstanceId == "" {
		return derrors.NewInvalidArgumentError(emptyAppInstanceID)
	}
	return nil
}
