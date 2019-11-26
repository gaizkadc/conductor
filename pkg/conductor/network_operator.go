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

package conductor

import (
    "github.com/nalej/conductor/internal/entities"
    pbApplication "github.com/nalej/grpc-application-go"
    "github.com/nalej/derrors"
)

const (
    // Definition of the prefix to be used when defining any Nalej service variable
    NalejVariablePrefix = "NALEJ_SERV_%s"
    // Definition of the prefix to be used when defining any Nalej outbound variable
    NalejVariableOutboundPrefix = "NALEJ_OUTBOUND_%s"
    // Suffix to add to the outbound FQDNs
    OutboundSuffix = "-OUT-"
)


type NetworkOperator interface {

    // PrepareNetwork
    // params:
    //  descriptor to be deployed
    //  instance of the application
    // return:
    //  error if any
    PrepareNetwork(descriptor *pbApplication.ParametrizedDescriptor, instance *pbApplication.AppInstance) (string,derrors.Error)

    // Get a network Id for an existing application instance
    // params:
    //  appInstance application instance to get the instance from
    // return:
    //  return the network Id, or error if any
    GetNetworkId(appInstance *entities.AppInstance) (string, derrors.Error)

    // Generate the tuple key and value for a Nalej service to be represented in the networking modality.
    // params:
    //  serviceName
    //  appInstanceId
    //  organizationId
    // return:
    //  variable name, variable value
    GetDeploymentVariableForService(serviceName string, appInstanceId string, organizationId string) (string, string)

}
