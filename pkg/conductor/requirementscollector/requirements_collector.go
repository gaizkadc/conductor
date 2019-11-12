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
	"github.com/nalej/conductor/internal/entities"
	pbApplication "github.com/nalej/grpc-application-go"
)

// Interface to be implemented by any requirements collector

type RequirementsCollector interface {

	// Find the set of requirements demanded by the application in an internally processable format.
	//  params:
	//   appDescriptor application to be analyzed
	//   appInstanceId id of the referenced application instance
	//  return:
	//   requirements or error if any
	FindRequirements(appDescriptor *pbApplication.ParametrizedDescriptor, appInstanceId string) (*entities.Requirements, error)

	// Find requirements for the service group which id is contained in the application descriptor
	//  params:
	//   serviceGroupId service group identifier
	//   appInstanceId
	//   appDescriptor
	//  return:
	//   requirement for the service group or error if any
	FindRequirementForGroup(serviceGroupId string, appInstanceId string, appDescriptor *pbApplication.ParametrizedDescriptor) (*entities.Requirement, error)

	// Find the set of requirements for a list of service groups
	//  params:
	//   serviceGroupId service group identifier
	//   appInstanceId
	//   appDescriptor
	//  return:
	//   requirement for the service group or error if any
	FindRequirementsForGroups(serviceGroupsIds []string, appInstanceId string, appDescriptor *pbApplication.ParametrizedDescriptor) (*entities.Requirements, error)
}
