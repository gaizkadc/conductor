//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Mockup client emulating a connection with the appmgr.

package asm

import (
    "github.com/daishogroup/system-model/entities"
    entitiesConductor "github.com/daishogroup/conductor/entities"
    "github.com/daishogroup/derrors"
)

// Factory in charge of creating new rest clients.
type MockupClientFactory struct{}

func (f *MockupClientFactory) CreateClient(host string, port int) Client {
    return NewMockupClient()
}

func (f *MockupClientFactory) CreateClientWithAddress(address string) Client {
    return NewMockupClient()
}

func NewMockupClientFactory() ClientFactory {
    return &MockupClientFactory{}
}

type MockupClient struct {
    // Emtpy now
}

func NewMockupClient() Client {

    return &MockupClient{}
}

func (c *MockupClient) Start(descriptor entities.AppDescriptor,
    appRequest entitiesConductor.DeployAppRequest) derrors.DaishoError {
    // TODO
    return nil
}

func (c *MockupClient) Stop(instanceName string) (bool, derrors.DaishoError) {
    // TODO
    return true, nil
}

// Pods get the pod list of the application.
// params:
//   instanceName name of the instance we want to stop
// return:
//   AppPodsResponse the list of pods
//   error if any
func (c *MockupClient) Pods(instanceName string) (*AppPodsResponse, derrors.DaishoError) {
    // TODO
    response := &AppPodsResponse{Pods: []string{"pod1", "pod2"}}
    return response, nil
}

func (c * MockupClient) List() (*AppListResponse, derrors.DaishoError) {
    return NewAppListResponse("none", 0, "", 0, make([]DaishoApplicationInfo, 0)), nil
}
