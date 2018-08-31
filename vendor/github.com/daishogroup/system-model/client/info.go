//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
)

// Info is the client interface for extract summary info of the system model.
type Info interface {
    // ReducedInfo get all the essential information in the system model.
    //   returns:
    //     The essential information of the system model.
    //     An error if the data cannot be obtained.
    ReducedInfo() (*entities.ReducedInfo, derrors.DaishoError)

    // SummaryInfo exports the basic information of the system model.
    //   returns:
    //     A summary with the counters of each entity.
    //     An error if the data cannot be obtained.
    SummaryInfo() (*entities.SummaryInfo, derrors.DaishoError)

    // ReducedInfoByNetwork basic information in the system model filter by networkID.
    //   params:
    //     networkID The selected network.
    //   returns:
    //     A dump structure with the system model information.
    //     An error if the data cannot be obtained.
    ReducedInfoByNetwork(networkID string) (*entities.ReducedInfo, derrors.DaishoError)
}
