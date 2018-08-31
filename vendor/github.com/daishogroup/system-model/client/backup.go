//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Definition of the Backup operations.

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
)

// Dump is the client interface for import/export information of the System Model.
type Backup interface {

    // Export all or specific  information from system model into a Backup structure.
    //   returns:
    //     A dump structure with the system model information.
    //     An error if the data cannot be obtained.
    Export(component string)(* entities.BackupRestore, derrors.DaishoError)

    // Restore all or specific  information into system model provided by backup.
    //   returns:
    //     Success or failure.
    //     if restore fails.
    // Currently there is no mechanism to rollback partly restored data.
    Import(component string, entity * entities.BackupRestore)( derrors.DaishoError)

}
