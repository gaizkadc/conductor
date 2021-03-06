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
	pbConductor "github.com/nalej/grpc-conductor-go"
)

// List of requirements demanded by an app
type Requirements struct {
	List []Requirement `json:"list, omitempty"`
}

func NewRequirements() Requirements {
	return Requirements{List: make([]Requirement, 0)}
}

// AddRequirement to the current list
func (r *Requirements) AddRequirement(req Requirement) {
	r.List = append(r.List, req)
}

func (r *Requirements) ToGRPC() []*pbConductor.Requirement {
	toReturn := make([]*pbConductor.Requirement, len(r.List))

	for i, req := range r.List {
		toReturn[i] = &pbConductor.Requirement{
			Replicas:               req.Replicas,
			Storage:                req.Storage,
			Memory:                 req.Memory,
			Cpu:                    req.CPU,
			AppInstanceId:          req.AppInstanceId,
			GroupServiceInstanceId: req.GroupServiceId,
			RequestId:              "",
		}
	}
	return toReturn
}

// Requirement for an app.
type Requirement struct {
	// Application instance id
	AppInstanceId string `json: "app_instance_id, omitempty"`
	//Groupo service id
	GroupServiceId string `json:"service_id, omitempty"`
	// Amount of CPUNum
	CPU int64 `json:"cpu, omitempty"`
	// Amount of memory
	Memory int64 `json:"memory, omitempty"`
	// Amount of storage
	Storage int64 `json:"storage, omitempty"`
	// Number of replicas
	Replicas int32 `json:"replicas, omitempty"`
	// Cluster selection labels
	DeploymentSelectors map[string]string `json:"deployment_selectors, omitempty"`
}

func NewRequirement(appInstanceId string, groupServiceId string, cpu int64, memory int64, storage int64,
	replicas int32, deploymentSelectors map[string]string) Requirement {
	return Requirement{AppInstanceId: appInstanceId, GroupServiceId: groupServiceId, CPU: cpu, Memory: memory,
		Storage: storage, Replicas: replicas, DeploymentSelectors: deploymentSelectors}
}
