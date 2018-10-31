/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package entities

import (
    pbConductor "github.com/nalej/grpc-conductor-go"
)

// List of requirements demanded by an app
type Requirements struct {
    List []Requirement `json:"list, omitempty"`
}

func NewRequirements() Requirements {
    return Requirements{List: make([]Requirement,0)}
}

// AddRequirement to the current list
func (r *Requirements) AddRequirement (req Requirement) {
    r.List = append(r.List, req)
}

func (r *Requirements) ToGRPC() []*pbConductor.Requirement {
    toReturn:= make([]*pbConductor.Requirement, len(r.List))
    for i, req := range r.List {
        toReturn[i] = &pbConductor.Requirement{
            Replicas: req.Replicas,
            Storage: req.Storage,
            Memory: req.Memory,
            Cpu: req.CPU,
            ServiceId: req.ServiceId,
        }
    }
    return toReturn
}

// Requirement for an app.
type Requirement struct {
    //Application id
    ServiceId string `json:"service_id, omitempty"`
    // Amount of CPU
    CPU int64 `json:"cpu, omitempty"`
    // Amount of memory
    Memory int64 `json:"memory, omitempty"`
    // Amount of storage
    Storage int64 `json:"storage, omitempty"`
    // Number of replicas
    Replicas int32 `json:"replicas, omitempty"`
}

func NewRequirement(appId string, cpu int64, memory int64, storage int64, replicas int32) Requirement {
    return Requirement{ServiceId: appId, CPU: cpu, Memory: memory, Storage: storage, Replicas: replicas}
}