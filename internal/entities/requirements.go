/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package entities

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

// Requirement for an app.
type Requirement struct {
    //Application id
    AppId string `json:"app_id, omitempty"`
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
    return Requirement{AppId: appId, CPU: cpu, Memory: memory, Storage: storage, Replicas: replicas}
}