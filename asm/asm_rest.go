//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Rest application client to connect with an existing application manager.

package asm

import (
    "fmt"
    "github.com/daishogroup/system-model/entities"
    "strings"
    "github.com/daishogroup/dhttp"
    entitiesConductor "github.com/nalej/conductor/entities"
    "github.com/daishogroup/derrors"
)

const (
    PackageUploadEndpoint   = "/appmgr/packages"
    ApplicationEndpoint     = "/appmgr/apps"
    ApplicationStopEndpoint = "/appmgr/apps/%s"
    ApplicationPodsEndpoint = "/appmgr/resources/pods/%s"
    PackageUploadFormKey    = "package"

    ApplicationTypeLabel = "daishoApplicationType"
    // Port where slaves are listening to
    SlavePort = 30000
)

// Factory in charge of creating new rest clients.
type RestClientFactory struct{}

func (f *RestClientFactory) CreateClient(host string, port int) Client {
    return NewASMRestClient(host, port)
}

/*
func (f *RestClientFactory) CreateClientWithAddress(host string, port int) Client {
    targetAddress := fmt.Sprintf("https://%s:30000", address)
    return NewASMRestClient(targetAddress)
}
*/

func NewRestClientFactory() ClientFactory {
    return &RestClientFactory{}
}

type RestClient struct {
    client dhttp.Client
}

// Constructor of an REST client to connect with an appmgr instance
// params:
//   host string with the hostname
//   port number indicating the target port
// return:
//   instance of the rest client
func NewASMRestClient(host string, port int) Client {
    logger.Debugf("Create https client pointing at https://%s:%d", host, port)
    conf := dhttp.NewRestBasicHTTPS(host, port)
    rest := dhttp.NewClientSling(conf)
    return &RestClient{rest}
}

// Start an application
// params:
//   descriptor Application descriptor
// return:
//   Error if any
func (c *RestClient) Start(descriptor entities.AppDescriptor,
    appRequest entitiesConductor.DeployAppRequest) derrors.DaishoError {
    logger.Debug("Called application start")
    // Can only start if we've uploaded it somewhere
    if c.client == nil {
        logger.Error("impossible to run Deploy if the client is null")
        return derrors.NewGenericError("impossible to run Deploy if the client is null")
    }

    var startRequest entitiesConductor.AppStartRequest
    if len(appRequest.Arguments) != 0 {
        splitArgs := strings.Split(appRequest.Arguments, " ")
        startRequest = *entitiesConductor.NewAppStartRequest(
            appRequest.Name, descriptor.ServiceName, descriptor.ServiceVersion, splitArgs, appRequest.Labels)
    } else {
        startRequest = *entitiesConductor.NewAppStartRequest(
            appRequest.Name, descriptor.ServiceName, descriptor.ServiceVersion, make([]string, 0), appRequest.Labels)
    }

    /*
    splitArgs := strings.Split(appRequest.Arguments, " ")
    startRequest := entitiesConductor.NewAppStartRequest(appRequest.Name, descriptor.ServiceName, descriptor.ServiceVersion, splitArgs)
    */
    logger.Debugf("Sending start request to appmgr: %#v", startRequest)
    startResponse := &entitiesConductor.AppStartResponse{}
    response := c.client.Post(ApplicationEndpoint, startRequest, startResponse)
    if response.Error!= nil {
        logger.Errorf("There was an error requesting an application to start [%s]", response.Error)
        return derrors.NewOperationError("failed to start application", response.Error)
    }

    logger.Debugf("Start response \n%s", response)
    logger.Debugf("Result from start application %s", response.Result)

    return nil
}

// Stop an application
// params:
//   instanceName name of the instance we want to stop
// return:
//   error if any
func (c *RestClient) Stop(instanceName string) (bool, derrors.DaishoError) {
    // TODO appmgr is using an HTTP delete operation expecting a body. This should be fixed.
    // In this code we have to extend the sling rest client with a delete with body operation
    logger.Debug("Called application stop: " + instanceName)
    // Create a stop request
    if c.client == nil {
        logger.Error("impossible to undeploy. The client is null.")
        return false, derrors.NewGenericError("impossible to undeploy. The client is null")
    }

    stopRequest := entitiesConductor.NewAppStopRequest(instanceName)
    stopResponse := entitiesConductor.NewAppStopResponse()
    targetUrl := fmt.Sprintf(ApplicationStopEndpoint, instanceName)
    response := c.client.DeleteWithBody(targetUrl, stopRequest, stopResponse)
    if response.Error!=nil {
        logger.Errorf("there was an error requesting an application to stop [%s]", response.Error)
        return false, derrors.NewOperationError(stopResponse.Info,response.Error)
    }

    return stopResponse.Stopped, nil
}

// Pods get the pod list of the application.
// params:
//   instanceName name of the instance we want to get logs
// return:
//   AppPodsResponse the list of pods
//   error if any
func (c *RestClient) Pods(instanceName string) (*AppPodsResponse, derrors.DaishoError) {

    logger.Debug("Called list of pods: " + instanceName)

    if c.client == nil {
        logger.Error("impossible to list of pods. The client is null.")
        return nil, derrors.NewGenericError("impossible to list of pods. The client is null")
    }

    podsRequest := NewAppPodsRequest()
    podsResponse := &AppPodsResponse{}
    targetUrl := fmt.Sprintf(ApplicationPodsEndpoint, instanceName)

    response := c.client.GetWithBody(targetUrl, podsRequest, podsResponse)
    if !response.Ok() {
        return nil, derrors.NewOperationError("error returning the list of pods",response.Error)
    }
    return podsResponse, nil
}

func (c* RestClient) List() (*AppListResponse, derrors.DaishoError) {
    logger.Debug("Called application list")

    if c.client == nil {
        logger.Error("impossible to list apps. The client is null.")
        return nil, derrors.NewGenericError("impossible to list apps. The client is null")
    }

    request := NewAppListRequest()
    listResponse := &AppListResponse{}
    response := c.client.GetWithBody(ApplicationEndpoint, request, listResponse)
    if response.Error != nil {
        logger.Error("Error retrieving applicaiton list")
        return nil, response.Error
    }
    return listResponse, nil
}
