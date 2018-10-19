/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 */


package entities

import (
    pbConductor "github.com/nalej/grpc-conductor-go"
    pbApplication "github.com/nalej/grpc-application-go"
)


// Representation of the score for a potential deployment candidate.
type ClusterScore struct {
    RequestID string
    ClusterID string
    Score float32
    TotalEvaluated int
}

// Objects describing received deployment requests. These objects are designed to be stored into
// a storage structure such as a queue.
type DeploymentRequest struct {
    RequestID string
    OrganizationID string
    ApplicationID string
    InstanceID string
}

// Fragment deployment status definition

type DeploymentFragmentStatus int

const (
    FRAGMENT_WAITING   DeploymentFragmentStatus = iota
    FRAGMENT_DEPLOYING
    FRAGMENT_DONE
    FRAGMENT_ERROR
    FRAGMENT_RETRYING
)

var DeploymentStatusToGRPC = map[pbConductor.DeploymentFragmentStatus] DeploymentFragmentStatus {
    pbConductor.DeploymentFragmentStatus_WAITING : FRAGMENT_WAITING,
    pbConductor.DeploymentFragmentStatus_DEPLOYING : FRAGMENT_DEPLOYING,
    pbConductor.DeploymentFragmentStatus_DONE : FRAGMENT_DONE,
    pbConductor.DeploymentFragmentStatus_ERROR : FRAGMENT_ERROR,
    pbConductor.DeploymentFragmentStatus_RETRYING : FRAGMENT_RETRYING,
}


// Plan related entries

// Data structure defining a deployment plan for Nalej applications.
type DeploymentPlan struct{
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
type DeploymentFragment struct{
    // OrganizationId this deployment belongs to
    OrganizationId string `json:"organization_id,omitempty"`
    // AppInstanceId for the instance of the application to run
    AppInstanceId string `json:"app_instance_id,omitempty"`
    // Identifier for this plan
    DeploymentId string `json:"deployment_id,omitempty"`
    // Fragment id
    FragmentId string `json:"fragment_id,omitempty"`
    // Cluster id
    ClusterId string `json:"cluster_id,omitempty"`
    // Cluster Ip
    ClusterIp string `json:"cluster_ip,omitempty"`
    // Deployment stages belonging to this fragment
    Stages []DeploymentStage `json:"stages,omitempty"`
}


func(df *DeploymentFragment) ToGRPC() *pbConductor.DeploymentFragment {
    convertedStages := make([]*pbConductor.DeploymentStage,len(df.Stages))
    for i,serv := range df.Stages {
        convertedStages[i] = serv.ToGRPC()
    }
    result := pbConductor.DeploymentFragment{
        OrganizationId: df.OrganizationId,
        FragmentId: df.FragmentId,
        AppInstanceId: df.AppInstanceId,
        DeploymentId: df.DeploymentId,
        Stages: convertedStages,
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

func (ds *DeploymentStage) ToGRPC() *pbConductor.DeploymentStage{
    convertedServices := make([]*pbApplication.Service,len(ds.Services))
    for i,serv := range ds.Services {
        convertedServices[i] = serv.ToGRPC()
    }
    result := pbConductor.DeploymentStage{
        FragmentId: ds.FragmentId,
        StageId: ds.StageId,
        Services: convertedServices,
    }
    return &result
}
