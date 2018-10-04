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
    //   appInstance application to be analyzed
    //  return:
    //   requirements or error if any
    FindRequirements(appInstance *pbApplication.AppInstance) (*entities.Requirements, error)

}
