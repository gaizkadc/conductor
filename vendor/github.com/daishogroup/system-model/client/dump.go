//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Definition of the dump operations.

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
)

// Dump is the client interface for import/export information of the System Model.
type Dump interface {

    // Export all the information in the system model into a Dump structure.
    //   returns:
    //     A dump structure with the system model information.
    //     An error if the data cannot be obtained.
    Export()(* entities.Dump, derrors.DaishoError)

}