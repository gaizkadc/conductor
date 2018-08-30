//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Interface to be implemented by any appmgr client

package asm

import (
    log "github.com/sirupsen/logrus"
    "github.com/daishogroup/system-model/entities"
    entitiesConductor "github.com/daishogroup/conductor/entities"
    "github.com/daishogroup/derrors"
)

var logger = log.WithField("package", "asm_client")

type Client interface {
    // Start an application
    // params:
    //   descriptor Information about the application descriptor
    //   appRequest Related application request with additional information
    // return:
    //   Error if any
    Start(descriptor entities.AppDescriptor, appRequest entitiesConductor.DeployAppRequest) derrors.DaishoError

    // Stop an application
    // params:
    //   instanceName name of the instance we want to stop
    // return:
    //   boolean indicating wheter the application is stopped or not
    //   error if any
    Stop(instanceName string) (bool, derrors.DaishoError)

    // Pods get the pod list of the application.
    // params:
    //   instanceName name of the instance we want to stop
    // return:
    //   AppPodsResponse the list of pods
    //   error if any
    Pods(instanceName string) (*AppPodsResponse, derrors.DaishoError)

    // List all running applications.
    List() (*AppListResponse, derrors.DaishoError)
}

// This client factory encapsulates the creation of asmclients
// ASM clients lifetime are expected to be very short. Serveral kinds of clients can coexist pointing to different
// urls
type ClientFactory interface {
    // Create a new client
    // params:
    //    host name of the machine pointint at
    //    port number
    // returns:
    //    instance of a client
    CreateClient(host string, port int) Client

    // Create a new client with Address
    // params:
    //    url URL address to point at
    // returns:
    //    instance of a client
    //CreateClientWithAddress(address string) Client
}
