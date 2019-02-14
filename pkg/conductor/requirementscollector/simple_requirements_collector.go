/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package requirementscollector

import (
    "github.com/nalej/conductor/internal/entities"
    "github.com/nalej/derrors"
    pbApplication "github.com/nalej/grpc-application-go"

)

type SimpleRequirementsCollector struct {}

func NewSimpleRequirementsCollector() RequirementsCollector {
    return &SimpleRequirementsCollector{}
}


func (s *SimpleRequirementsCollector) FindRequirements(appInstance *pbApplication.AppInstance) (*entities.Requirements, error) {
    // Check if we have any service group to deploy
    if len(appInstance.Groups) == 0 {
        return nil, derrors.NewFailedPreconditionError("no groups available for this application")
    }


    foundRequirements := entities.NewRequirements()

    // Generate one set of requirements per service group
    for _, g := range appInstance.Groups {

        // Check if there are any services to be analysed
        if len(g.ServiceInstances) == 0 {
            return nil, derrors.NewFailedPreconditionError("no services specified for the application")
        }

        var totalStorage int64
        var totalCPU int64
        var totalMemory int64

        for _, serv := range g.ServiceInstances{


            totalCPU = totalCPU + serv.Specs.Cpu
            totalMemory = totalMemory + serv.Specs.Memory
            // accumulate requested storage
            for _, st := range serv.Storage {
                totalStorage = totalStorage + st.Size
            }
        }

        r := entities.NewRequirement(appInstance.AppInstanceId, g.Name, totalCPU, totalMemory,
            totalStorage, g.Specs.NumReplicas)
        foundRequirements.AddRequirement(r)
    }

    return &foundRequirements, nil
}