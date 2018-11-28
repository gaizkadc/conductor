/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/derrors"
	pbApplication "github.com/nalej/grpc-application-go"
	pbConductor "github.com/nalej/grpc-conductor-go"
)

// Representation of the score for a potential deployment candidate.
type ClustersScore struct {
    // RequestId for this score request
    Scoring  [] ClusterScore `json:"scoring,omitempty"`
    TotalEvaluated int `json:"total_evaluated,omitempty"`
}

func NewClustersScore() ClustersScore {
    return ClustersScore{TotalEvaluated: 0, Scoring: make([]ClusterScore,0)}
}

// AddClusterScore appends a cluster score and updates the set of total evaluated clusters
func (c *ClustersScore) AddClusterScore(score ClusterScore) {
    c.Scoring = append(c.Scoring, score)
    c.TotalEvaluated = c.TotalEvaluated + 1
}

// Scoring returned by a certain cluster
type ClusterScore struct {
    // ClusterId for the queried cluster
    ClusterId string `json:"cluster_id,omitempty"`
    // Score returned by this cluster
    Score float32 `json:"score,omitempty"`
}

// Objects describing received deployment requests. These objects are designed to be stored into
// a storage structure such as a queue.
type DeploymentRequest struct {
	RequestId      string
	OrganizationId string
	ApplicationId  string
	InstanceId     string
}

// Fragment deployment status definition

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
}

// Data structure representing the components of a plan that will be
// spread across the available clusters.
type DeploymentFragment struct {
	// OrganizationId this deployment belongs to
	OrganizationId string `json:"organization_id,omitempty"`
	// OrganizationName this deployment belongs to
	OrganizationName string `json:"organization_name,omitempty"`
	// AppInstanceId for the instance of the application to run
	AppInstanceId string `json:"app_instance_id,omitempty"`
	// AppNamed for the instance of the application to run
	AppName string `json:"app_name,omitempty"`
	// Identifier for this plan
	DeploymentId string `json:"deployment_id,omitempty"`
	// Fragment id
	FragmentId string `json:"fragment_id,omitempty"`
	// Cluster id
	ClusterId string `json:"cluster_id,omitempty"`
	// Deployment stages belonging to this fragment
	Stages []DeploymentStage `json:"stages,omitempty"`
}

type UndeployRequest struct {
	// OrganizationId this deployment belongs to
	OrganizationId string `json:"organization_id,omitempty"`
	// AppInstanceId for the instance of the application to run
	AppInstanceId string `json:"app_instance_id,omitempty"`
}

func (df *DeploymentFragment) ToGRPC() *pbConductor.DeploymentFragment {
	convertedStages := make([]*pbConductor.DeploymentStage, len(df.Stages))
	for i, serv := range df.Stages {
		convertedStages[i] = serv.ToGRPC()
	}
	result := pbConductor.DeploymentFragment{
		OrganizationId: df.OrganizationId,
		FragmentId:     df.FragmentId,
		AppInstanceId:  df.AppInstanceId,
		OrganizationName: df.OrganizationName,
		AppName: 		df.AppName,
		DeploymentId:   df.DeploymentId,
		Stages:         convertedStages,
	}
	return &result
}

// Every deployment stage a frament is made of.
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
