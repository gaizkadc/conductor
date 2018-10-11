/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package entities

import (
    "time"
    pbConductor "github.com/nalej/grpc-conductor-go"
)
// System status representation.
type Status struct {
    Timestamp time.Time `json:"timestamp"`
    Mem float64 `json: "mem"`
    CPU float64 `json: "cpu"`
    Disk float64 `json: "disk"`
}


type Requirements struct {
    CPU float32
    Memory float32
    Disk float32
}

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

// Service status definition

type ServiceStatus int

const (
    SERVICE_SCHEDULED = iota
    SERVICE_WAITING
    SERVICE_DEPLOYING
    SERVICE_RUNNING
    SERVICE_ERROR
)

var ServiceStatusToGRPC = map[pbConductor.ServiceStatus] ServiceStatus {
    pbConductor.ServiceStatus_SERVICE_SCHEDULED : SERVICE_SCHEDULED,
    pbConductor.ServiceStatus_SERVICE_WAITING : SERVICE_WAITING,
    pbConductor.ServiceStatus_SERVICE_DEPLOYING : SERVICE_DEPLOYING,
    pbConductor.ServiceStatus_SERVICE_RUNNING : SERVICE_RUNNING,
    pbConductor.ServiceStatus_SERVICE_ERROR : SERVICE_ERROR,
}