/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/derrors"
	pbApplication "github.com/nalej/grpc-application-go"
	pbConductor "github.com/nalej/grpc-conductor-go"
	"sort"
	"strings"
	"time"
)

// DeploymentsScore for a set of potential deployments.
type DeploymentScore struct {
    // Score for every evaluated cluster
    DeploymentsScore [] ClusterDeploymentScore `json:"scoring,omitempty"`
    // Total number of evaluated clusters
    NumEvaluatedClusters int `json:"num_evaluated_clusters,omitempty"`
}

func NewClustersScore() DeploymentScore {
    return DeploymentScore{NumEvaluatedClusters: 0, DeploymentsScore: make([]ClusterDeploymentScore,0)}
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
		Scores: make(map[string]float32,0),
	}
}

// Add the score for a set of service groups.
func(cds *ClusterDeploymentScore) AddScore(serviceGroupIds []string, score float32) {
	sort.Strings(serviceGroupIds)
	// The key is the concatenation of the ids
	newKey := strings.Join(serviceGroupIds,"")
	cds.Scores[newKey] = score
}

// End of cluster deployment score ------



// Objects describing received deployment requests. These objects are designed to be stored into
// a storage structure such as a queue.
type DeploymentRequest struct {
	RequestId      string
	OrganizationId string
	ApplicationId  string
	InstanceId     string
	NumRetries     int32
	TimeRetry      *time.Time
}

// Fragment deployment Status definition

type DeploymentFragmentStatus int

const (
	FRAGMENT_WAITING DeploymentFragmentStatus = iota
	FRAGMENT_DEPLOYING
	FRAGMENT_DONE
	FRAGMENT_ERROR
	FRAGMENT_RETRYING
)

var DeploymentStatusToGRPC = map[pbConductor.DeploymentFragmentStatus]DeploymentFragmentStatus{
	pbConductor.DeploymentFragmentStatus_WAITING:   FRAGMENT_WAITING,
	pbConductor.DeploymentFragmentStatus_DEPLOYING: FRAGMENT_DEPLOYING,
	pbConductor.DeploymentFragmentStatus_DONE:      FRAGMENT_DONE,
	pbConductor.DeploymentFragmentStatus_ERROR:     FRAGMENT_ERROR,
	pbConductor.DeploymentFragmentStatus_RETRYING:  FRAGMENT_RETRYING,
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
	// AppInstanceId for the instance of the application to run
	AppInstanceId string `json:"app_instance_id,omitempty"`
	// AppNamed for the instance of the application to run
	AppName string `json:"app_name,omitempty"`
	// Identifier for this plan
	DeploymentId string `json:"deployment_id,omitempty"`
	// Fragment id
	FragmentId string `json:"fragment_id,omitempty"`
	// Group service id
	ServiceGroupId string `json:"service_group_id,omitempty"`
	// Group service instance id
	ServiceGroupInstanceId string `json:"service_group_instance_id,omitempty"`
	// Cluster id for the deployment target
	ClusterId string `json:"cluster_id,omitempty"`
	// Nalej variables
	NalejVariables map[string]string `json:"nalej_variables,omitempty"`
	// Deployment stages belonging to this fragment
	Stages []DeploymentStage `json:"stages,omitempty"`
}

func (df *DeploymentFragment) ToGRPC() *pbConductor.DeploymentFragment {
	convertedStages := make([]*pbConductor.DeploymentStage, len(df.Stages))
	for i, serv := range df.Stages {
		convertedStages[i] = serv.ToGRPC()
	}
	result := pbConductor.DeploymentFragment{
		OrganizationId: df.OrganizationId,
		FragmentId:     df.FragmentId,
		AppDescriptorId: df.AppDescriptorId,
		AppInstanceId:  df.AppInstanceId,
		ServiceGroupId: df.ServiceGroupId,
		ServiceGroupInstanceId: df.ServiceGroupInstanceId,
		OrganizationName: df.OrganizationName,
		AppName: 		df.AppName,
		DeploymentId:   df.DeploymentId,
		NalejVariables: df.NalejVariables,
		Stages:         convertedStages,
		ClusterId: df.ClusterId,
	}
	return &result
}

// ----


// Deployment fragment definition ----

type UndeployRequest struct {
	// OrganizationId this deployment belongs to
	OrganizationId string `json:"organization_id,omitempty"`
	// AppInstanceId for the instance of the application to run
	AppInstanceId string `json:"app_instance_id,omitempty"`
}



// Every deployment stage a fragment is made of.
type DeploymentStage struct {
	// Fragment id
	FragmentId string `json:"fragment_id,omitempty"`
	// Stage id
	StageId string `json:"stage_id,omitempty"`
	// Set of services
	Services []Service `json:"stage_id,omitempty"`
}

func (ds *DeploymentStage) ToGRPC() *pbConductor.DeploymentStage {
	convertedServices := make([]*pbApplication.Service, len(ds.Services))
	for i, serv := range ds.Services {
		convertedServices[i] = serv.ToGRPC()
	}
	result := pbConductor.DeploymentStage{
		FragmentId: ds.FragmentId,
		StageId:    ds.StageId,
		Services:   convertedServices,
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
	if request.AppId == nil {
		return derrors.NewInvalidArgumentError(emptyAppID)
	}
	if request.AppId.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationID)
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
