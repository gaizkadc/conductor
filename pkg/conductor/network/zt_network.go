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
    "context"
    "github.com/nalej/conductor/internal/entities"
    "github.com/nalej/conductor/pkg/conductor"
    "github.com/nalej/conductor/pkg/conductor/plandesigner"
    "github.com/nalej/conductor/pkg/utils"
    "github.com/nalej/derrors"
    pbNetwork "github.com/nalej/grpc-network-go"
    pbApplication "github.com/nalej/grpc-application-go"
    "github.com/nalej/nalej-bus/pkg/queue/network/ops"
    "github.com/rs/zerolog/log"
    "net"
    "strings"
    "time"
)

// Manage network parameters when ZT networking is chosen.
const (
    // Initial address to use during the definition of VSA
    ConductorBaseVSA = "172.16.0.1"
    // Initial address to use during the definition of VSA
    ConductorOutboundVSA = "172.18.0.1"
    // Timeout when sending messages to the queue
    QueueTimeout = time.Second * 5
    // Nalej service suffix
    NalejServiceSuffix = "service.nalej"
)

// ZtNetworkingOperator is a wrapper that supports operations required to set a ZT VPN.
type ZtNetworkingOperator struct {
    AppClient pbApplication.ApplicationsClient
    NetworkOpsProducer *ops.NetworkOpsProducer
    NetClient pbNetwork.NetworksClient
}

func NewZtNetworkingOperator(connHelper *utils.ConnectionsHelper, networkOpsProducer *ops.NetworkOpsProducer) conductor.NetworkOperator {
    pool := connHelper.GetSystemModelClients()
    if pool != nil && len(pool.GetConnections()) == 0 {
        log.Panic().Msg("system model clients were not started")
        return nil
    }

    conn := pool.GetConnections()[0]
    // Create associated clients
    appClient := pbApplication.NewApplicationsClient(conn)

    // Network client
    netPool := connHelper.GetNetworkingClients()
    if netPool != nil && len(netPool.GetConnections()) == 0 {
        log.Panic().Msg("networking client was not started")
        return nil
    }

    netClient := pbNetwork.NewNetworksClient(netPool.GetConnections()[0])

    return &ZtNetworkingOperator{AppClient: appClient, NetworkOpsProducer: networkOpsProducer, NetClient: netClient}
}

func(zt *ZtNetworkingOperator) PrepareNetwork(appDescriptor *pbApplication.ParametrizedDescriptor,
    appInstance *pbApplication.AppInstance) (string, derrors.Error) {

    // Create VSA
    vsa, err := zt.CreateVSA(entities.NewParametrizedDescriptorFromGRPC(appDescriptor), appInstance.AppInstanceId)
    if err != nil {
        err := derrors.NewGenericError("impossible to create VAS", err)
        log.Error().Err(err).Str("appDescriptorId", appDescriptor.AppDescriptorId).Msg("impossible to create VAS")
        return "", err
    }

    // Create ZT-network with Network manager
    // we use the app instance id as the network id
    ztNetworkId, ztErr := zt.CreateZTNetwork(appInstance.AppInstanceId, appInstance.OrganizationId,
        appInstance.AppInstanceId, vsa)
    if ztErr != nil {
        err := derrors.NewGenericError("impossible to create zt network before deployment", ztErr)
        log.Error().Err(ztErr).Str("appDescriptorId", appDescriptor.AppDescriptorId)
        return "", err
    }

    return ztNetworkId,  nil
}

func(zt *ZtNetworkingOperator) GetNetworkId(appInstance *entities.AppInstance) (string, derrors.Error) {
    ctxNet, cancelNet := context.WithTimeout(context.Background(), QueueTimeout)
    defer cancelNet()
    networkId, err := zt.AppClient.GetAppZtNetwork(ctxNet,
        &pbApplication.GetAppZtNetworkRequest{
            AppInstanceId: appInstance.AppInstanceId, OrganizationId: appInstance.OrganizationId})

    if err != nil {
        log.Error().Err(err).Msg("service groups could not be deployed. The network id was not found")
        return "", derrors.NewInternalError("error retrieving zt network id", err)
    }
    return networkId.NetworkId, nil
}

// Generate the tuple key and value for a nalej service to be represented.
// params:
//  serviceName
//  appInstanceId
//  organizationId
// return:
//  variable name, variable value
func (zt *ZtNetworkingOperator) GetDeploymentVariableForService(serviceName string, appInstanceId string, organizationId string) (string, string) {
    key := fmt.Sprintf(conductor.NalejVariablePrefix, strings.ToUpper(serviceName))
    value := fmt.Sprintf("%s.%s", utils.GetVSAName(serviceName, organizationId, appInstanceId), NalejServiceSuffix)
    return key, value
}


// Create the virtual application addresses for a given application descriptor
// params:
//  appDescriptor requiring the VAS entries
//  appInstanceId to work with
// return:
//  map with the VSA list
//  error if the operation failed
func (zt *ZtNetworkingOperator) CreateVSA(appDescriptor entities.AppDescriptor, appInstanceId string) (map[string]string, derrors.Error) {
    currentIp := net.ParseIP(ConductorBaseVSA).To4()
    // store the generated vsa
    vsa := make(map[string]string, 0)
    servicesByName := make(map[string]entities.Service)

    for _, sg := range appDescriptor.Groups {
        for _, serv := range sg.Services {
            servicesByName[serv.Name] = serv
            fqdn := utils.GetVSAName(serv.Name, appDescriptor.OrganizationId, appInstanceId)
            dnsRequest := pbNetwork.AddDNSEntryRequest{
                OrganizationId: serv.OrganizationId,
                ServiceName:    serv.Name,
                Fqdn:           fqdn,
                Ip:             currentIp.String(),
                Tags: []string{
                    fmt.Sprintf("appInstanceId:%s", appInstanceId),
                    fmt.Sprintf("organizationId:%s", appDescriptor.OrganizationId),
                    fmt.Sprintf("descriptorId:%s", appDescriptor.AppDescriptorId),
                    fmt.Sprintf("serviceGroupId:%s", sg.ServiceGroupId),
                    fmt.Sprintf("serviceId:%s", serv.ServiceId),
                },
            }
            ctx, cancel := context.WithTimeout(context.Background(), QueueTimeout)
            err := zt.NetworkOpsProducer.Send(ctx, &dnsRequest)
            cancel()
            if err != nil {
                log.Error().Err(err).Interface("request", dnsRequest).Msg("impossible to send a dns entry request")
                return nil, err
            }
            vsa[fqdn] = currentIp.String()
            // Increase the IP
            currentIp = utils.NextIP(currentIp, 1)
        }
    }

    currentOutboundIp := net.ParseIP(ConductorOutboundVSA).To4()
    for _, securityRule := range appDescriptor.Rules {
        if securityRule.OutboundNetInterfaceName != "" {
            // TODO beware to not overflow the maximum length for a DNS name (253 chars)
            fqdn := utils.GetVSAName(securityRule.TargetServiceName, appDescriptor.OrganizationId, appInstanceId)+
                plandesigner.OutboundSuffix + securityRule.OutboundNetInterfaceName
            targetService := servicesByName[securityRule.TargetServiceName]
            dnsRequest := pbNetwork.AddDNSEntryRequest{
                OrganizationId: securityRule.OrganizationId,
                ServiceName:    targetService.Name,
                AppInstanceId:  appInstanceId,
                Fqdn:           fqdn,
                Ip:             currentOutboundIp.String(),
                Tags: []string{
                    fmt.Sprintf("appInstanceId:%s", appInstanceId),
                    fmt.Sprintf("organizationId:%s", appDescriptor.OrganizationId),
                    fmt.Sprintf("descriptorId:%s", appDescriptor.AppDescriptorId),
                    fmt.Sprintf("serviceGroupId:%s", targetService.ServiceGroupId),
                    fmt.Sprintf("serviceId:%s", targetService.ServiceId),
                    fmt.Sprintf("securityRule:%s", securityRule.Name),
                },
            }
            ctx, cancel := context.WithTimeout(context.Background(), QueueTimeout)
            err := zt.NetworkOpsProducer.Send(ctx, &dnsRequest)
            cancel()
            if err != nil {
                log.Error().Err(err).Interface("request", dnsRequest).Msg("impossible to send a dns entry request")
                return nil, err
            }
            vsa[fqdn] = currentOutboundIp.String()
            // Increase the IP
            currentOutboundIp = utils.NextIP(currentOutboundIp, 1)
        }
    }

    return vsa, nil
}


// Create a new zero tier network and return the corresponding network id.
// params:
//  name of the network
//  organizationId for this network
// returns:
//  networkId or error otherwise
func (zt *ZtNetworkingOperator) CreateZTNetwork(name string, organizationId string, appInstanceId string,
    vsa map[string]string) (string, error) {

    request := pbNetwork.AddNetworkRequest{
        Name:           name,
        OrganizationId: organizationId,
        AppInstanceId:  appInstanceId,
        Vsa:            vsa}

    log.Debug().Interface("addNetworkRequest", request).Msgf("create a network request")

    timeout, cancel := context.WithTimeout(context.Background(), QueueTimeout)
    defer cancel()
    ztNetworkId, err := zt.NetClient.AddNetwork(timeout, &request)

    if err != nil {
        log.Error().Err(err).Msgf("there was a problem when creating network for name: %s with org: %s", name, organizationId)
        return "", err
    }
    return ztNetworkId.NetworkId, err
}