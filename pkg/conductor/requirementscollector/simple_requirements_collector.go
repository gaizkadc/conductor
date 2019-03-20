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


func (s *SimpleRequirementsCollector) FindRequirements(appDescriptor *pbApplication.AppDescriptor, appInstanceId string) (*entities.Requirements, error) {
    // Check if we have any service group to deploy
    if len(appDescriptor.Groups) == 0 {
        return nil, derrors.NewFailedPreconditionError("no groups available for this application")
    }


    foundRequirements := entities.NewRequirements()

    // Generate one set of requirements per service group
    for _, g := range appDescriptor.Groups {

        // Check if there are any services to be analysed
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
            // accumulate requested storage
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
        r := entities.NewRequirement(appInstanceId, g.Name, totalCPU, totalMemory,
            totalStorage, 1, selectors)
        foundRequirements.AddRequirement(r)
    }

    return &foundRequirements, nil
}