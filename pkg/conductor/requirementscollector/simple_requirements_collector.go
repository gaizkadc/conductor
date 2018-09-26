/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package requirementscollector

import (
    "github.com/nalej/conductor/internal/entities"
    pbApplication "github.com/nalej/grpc-application-go"
)

type SimpleRequirementsCollector struct {}

func NewSimpleRequirementsCollector() RequirementsCollector {
    return &SimpleRequirementsCollector{}
}


func (s *SimpleRequirementsCollector) FindRequirements(appDescriptor *pbApplication.AppDescriptor) (*entities.Requirements, error) {
    // TODO this is hardcoded until we have an available system model
    toReturn := entities.Requirements{CPU:0.5,Memory:100, Disk:100}
    return &toReturn, nil
}