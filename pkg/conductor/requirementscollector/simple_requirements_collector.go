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

package requirementscollector

import (
	"fmt"
	"github.com/nalej/conductor/internal/entities"
	"github.com/nalej/derrors"
	pbApplication "github.com/nalej/grpc-application-go"
	"github.com/rs/zerolog/log"
)

type SimpleRequirementsCollector struct{}

func NewSimpleRequirementsCollector() RequirementsCollector {
	return &SimpleRequirementsCollector{}
}

func (s *SimpleRequirementsCollector) FindRequirements(appDescriptor *pbApplication.ParametrizedDescriptor, appInstanceId string) (*entities.Requirements, error) {
	// Check if we have any service group to deploy
	if len(appDescriptor.Groups) == 0 {
		return nil, derrors.NewFailedPreconditionError("no groups available for this application")
	}

	foundRequirements := entities.NewRequirements()

	// Generate one set of requirements per service group
	for _, g := range appDescriptor.Groups {

		req, err := s.FindRequirementForGroup(g.ServiceGroupId, appInstanceId, appDescriptor)

		if err != nil {
			log.Error().Str("appInstanceId", appInstanceId).Str("serviceGroupName", g.Name).
				Msg("impossible to find requirements for the group")
			return nil, derrors.NewFailedPreconditionError(fmt.Sprintf("impossible to find requirements for group %s", g.Name))
		}

		foundRequirements.AddRequirement(*req)
	}

	return &foundRequirements, nil
}

func (s *SimpleRequirementsCollector) FindRequirementsForGroups(serviceGroupsIds []string, appInstanceId string, appDescriptor *pbApplication.ParametrizedDescriptor) (*entities.Requirements, error) {
	foundRequirements := entities.NewRequirements()

	for _, sg := range serviceGroupsIds {
		req, err := s.FindRequirementForGroup(sg, appInstanceId, appDescriptor)
		if err != nil {
			return nil, err
		}
		foundRequirements.AddRequirement(*req)
	}

	return &foundRequirements, nil
}

func (s *SimpleRequirementsCollector) FindRequirementForGroup(serviceGroupId string, appInstanceId string, appDescriptor *pbApplication.ParametrizedDescriptor) (*entities.Requirement, error) {

	var g *pbApplication.ServiceGroup = nil

	for _, aux := range appDescriptor.Groups {
		if aux.ServiceGroupId == serviceGroupId {
			g = aux
			break
		}
	}

	if g == nil {
		return nil, derrors.NewFailedPreconditionError("service group not found inside descriptor")
	}

	if len(g.Services) == 0 {
		return nil, derrors.NewFailedPreconditionError("no services specified for the application")
	}

	var totalStorage int64 = 0
	var totalCPU int64 = 0
	var totalMemory int64 = 0

	for _, serv := range g.Services {

		numServReplicas := int64(1)
		if serv.Specs != nil && serv.Specs.Replicas > 0 {
			numServReplicas = int64(serv.Specs.Replicas)
		}

		totalCPU = totalCPU + (serv.Specs.Cpu * numServReplicas)
		totalMemory = totalMemory + (serv.Specs.Memory * numServReplicas)
		// accumulate requested provider
		for _, st := range serv.Storage {
			totalStorage = totalStorage + (st.Size * numServReplicas)
		}
	}

	selectors := map[string]string{}
	if g.Specs != nil {
		if g.Specs.DeploymentSelectors != nil {
			selectors = g.Specs.DeploymentSelectors
		}
	}

	// TODO: requirements for every fragment only permit one replica per requirement. Requirements are for a single service group
	toReturn := entities.NewRequirement(appInstanceId, g.Name, totalCPU, totalMemory, totalStorage, 1, selectors)
	return &toReturn, nil

}
