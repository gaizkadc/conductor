//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//

package client

import (
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/entities"
    "fmt"
    "github.com/daishogroup/dhttp"
)

// ReducedInfoURI is the uri of the ReducedInfo method.
const ReducedInfoURI = "/api/v0/info/reduced"

// SummaryInfoURI is the uri of the SummaryInfo method.
const SummaryInfoURI = "/api/v0/info/summary"

// ReducedInfoByNetworkURI is the uri of the ReducedInfoByNetwork method.
const ReducedInfoByNetworkURI = "/api/v0/info/%s/reduced"


// InfoRest is the client Rest for Info resources.
type InfoRest struct {
    client dhttp.Client
}
// Deprecated: Use NewInfoClientRest
func NewInfoRest(basePath string) Info {
    return NewInfoClientRest(ParseHostPort(basePath))
}

// NewInfoClientRest is the basic constructor of InfoRest.
func NewInfoClientRest(host string, port int) Info {
    rest := dhttp.NewClientSling(dhttp.NewRestBasicConfig(host,port))
    return &InfoRest{rest}
}

// ReducedInfo get all the essential information in the system model.
//   returns:
//     The essential information of the system model.
//     An error if the data cannot be obtained.
func (rest *InfoRest) ReducedInfo() (*entities.ReducedInfo, derrors.DaishoError) {
    response := rest.client.Get(fmt.Sprintf(ReducedInfoURI), new(entities.ReducedInfo))
    if response.Error != nil {
        return nil, response.Error
    }
    n := response.Result.(*entities.ReducedInfo)
    return n, nil
}

// SummaryInfo exports the basic information of the system model.
//   returns:
//     A summary with the counters of each entity.
//     An error if the data cannot be obtained.
func (rest *InfoRest) SummaryInfo() (*entities.SummaryInfo, derrors.DaishoError) {
    response := rest.client.Get(fmt.Sprintf(SummaryInfoURI), new(entities.SummaryInfo))
    if response.Error != nil {
        return nil, response.Error
    }
    n := response.Result.(*entities.SummaryInfo)
    return n, nil
}

// ReducedInfoByNetwork basic information in the system model filter by networkID.
//   params:
//     networkID The selected network.
//   returns:
//     A dump structure with the system model information.
//     An error if the data cannot be obtained.
func (rest *InfoRest) ReducedInfoByNetwork(networkID string) (*entities.ReducedInfo, derrors.DaishoError) {
    response := rest.client.Get(fmt.Sprintf(ReducedInfoByNetworkURI,networkID), new(entities.ReducedInfo))
    if response.Error != nil {
        return nil, response.Error
    }
    n := response.Result.(*entities.ReducedInfo)
    return n, nil
}

