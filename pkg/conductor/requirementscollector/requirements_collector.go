/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package requirementscollector

import (
    pbApplication "github.com/nalej/grpc-application-go"
    "github.com/nalej/conductor/internal/entities"
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
