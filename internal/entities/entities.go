/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package entities

import "time"

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