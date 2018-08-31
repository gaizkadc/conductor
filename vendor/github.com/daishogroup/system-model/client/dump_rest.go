//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// REST implementation of the dump client.

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "fmt"
    "github.com/daishogroup/dhttp"
)

const DumpExportURI = "/api/v0/dump/export"

// Client Rest for Network resources.
type DumpRest struct {
    client dhttp.Client
}

//Deprecated: Use NewDumpClientRest
func NewDumpRest(basePath string) Dump {
    return NewDumpClientRest(ParseHostPort(basePath))
}

func NewDumpClientRest(host string, port int) Dump {
    rest := dhttp.NewClientSling(dhttp.NewRestBasicConfig(host, port))
    return &DumpRest{rest}
}

// Export all the information in the system model into a Dump structure.
//   returns:
//     A dump structure with the system model information.
//     An error if the data cannot be obtained.
func (rest *DumpRest) Export() (*entities.Dump, derrors.DaishoError) {
    response := rest.client.Get(fmt.Sprintf(DumpExportURI), new(entities.Dump))
    if response.Error != nil {
        return nil, response.Error
    } else {
        n := response.Result.(*entities.Dump)
        return n, nil
    }
}
