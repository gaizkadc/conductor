/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

 // The handler monitor collects information from deployment fragments and updates the status of services.

package monitor

import (
    pbConductor "github.com/nalej/grpc-conductor-go"
    pbCommon "github.com/nalej/grpc-common-go"
    "errors"
)

type ConductorMonitor struct {
    mng * Manager
}

func (m* ConductorMonitor) UpdateDeploymentFragmentStatus(request *pbConductor.DeploymentFragmentUpdateRequest) (*pbCommon.Success, error) {
    if request.FragmentId == "" {
        err := errors.New("non valid empty fragment id")
        return nil, err
    }
    err := m.mng.UpdateFragmentStatus(request)
    if err != nil {
        return nil, err
    }
    return &pbCommon.Success{}, nil
}




