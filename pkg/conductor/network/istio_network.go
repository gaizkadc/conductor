/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package network

import (
    "fmt"
    "github.com/nalej/conductor/internal/entities"
    "github.com/nalej/conductor/pkg/conductor"
    "github.com/nalej/derrors"
    pbApplication "github.com/nalej/grpc-application-go"
    "strings"
)

type IstioNetworkingOperator struct {

}


func NewIstioNetworkingOperator() conductor.NetworkOperator {
    return &IstioNetworkingOperator{}
}

// PrepareNetwork
// params:
//  descriptor to be deployed
//  instance of the application
// return:
//  error if any
func(io *IstioNetworkingOperator) PrepareNetwork(descriptor *pbApplication.ParametrizedDescriptor,
    instance *pbApplication.AppInstance) (string,derrors.Error) {
        return "", nil
}

// Get a network Id for an existing application instance
// params:
//  appInstance application instance to get the instance from
// return:
//  return the network Id, or error if any
func(io *IstioNetworkingOperator)  GetNetworkId(appInstance *entities.AppInstance) (string, derrors.Error) {
    return "", nil
}

// Generate the tuple key and value for a Nalej service to be represented in the networking modality.
// params:
//  serviceName
//  appInstanceId
//  organizationId
// return:
//  variable name, variable value
func(io *IstioNetworkingOperator) GetDeploymentVariableForService(serviceName string, appInstanceId string, organizationId string) (string, string) {
    key := fmt.Sprintf(conductor.NalejVariablePrefix, strings.ToUpper(serviceName))
    // The value for an Istio service is simply the name of the service.
    value := strings.ToLower(serviceName)
    return key,value
}