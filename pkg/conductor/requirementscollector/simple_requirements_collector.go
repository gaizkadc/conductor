/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
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


func (s *SimpleRequirementsCollector) FindRequirements(appInstance *pbApplication.AppInstance) (*entities.Requirements, error) {
    // TODO check non-sense requirements
    foundRequirements := entities.NewRequirements()
    for _, serv := range appInstance.Services {
        var totalStorage int64
        totalStorage = 0
        for _,st := range serv.Storage {
            totalStorage = totalStorage + st.Size
        }
        r := entities.NewRequirement(serv.AppInstanceId, serv.Specs.Cpu, serv.Specs.Memory, totalStorage, serv.Specs.Replicas)
        foundRequirements.AddRequirement(r)
    }

    return &foundRequirements, nil
}