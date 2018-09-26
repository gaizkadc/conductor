/*
 * Copyright (C) 2018  Nalej - All Rights Reserved
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
    //  return:
    //   requirements or error if any
    FindRequirements(appDescriptor *pbApplication.AppDescriptor) (*entities.Requirements, error)

}
