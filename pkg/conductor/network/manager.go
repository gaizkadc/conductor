/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package network

import (
    pbNetwork "github.com/nalej/grpc-network-go"
    "github.com/nalej/conductor/pkg/conductor"
    "github.com/rs/zerolog/log"
    "context"
    "errors"
    "fmt"
)

type Manager struct{
    // Networking manager client
    NetClient pbNetwork.NetworksClient
    // DNS manager client
    DNSClient pbNetwork.DNSClient
}

func NewManager() (*Manager, error){
    // Network client
    netPool := conductor.GetNetworkingClients()
    if netPool != nil && len(netPool.GetConnections())==0{
        log.Panic().Msg("networking client was not started")
        return nil, errors.New("networking client was not started")
    }
    netClient := pbNetwork.NewNetworksClient(netPool.GetConnections()[0])
    dnsClient := pbNetwork.NewDNSClient(netPool.GetConnections()[0])

    return &Manager{netClient, dnsClient}, nil
}

func (m *Manager) AuthorizeNetworkMembership(organizationId string, networkId string, memberId string) error {
    req := pbNetwork.AuthorizeMemberRequest{
        OrganizationId: organizationId,
        NetworkId: networkId,
        MemberId: memberId,}
    _, err := m.NetClient.AuthorizeMember(context.Background(), &req)

    return err

}

func (m *Manager) RegisterNetworkEntry(organizationId string, networkId string, serviceName string, ip string) error {

    // Create the FQDN for this service
    fqdn := fmt.Sprintf("%s.%s",serviceName,organizationId)

    req := pbNetwork.AddDNSEntryRequest{
        NetworkId: networkId,
        OrganizationId: organizationId,
        Ip: ip,
        Fqdn: fqdn,
    }
    _, err := m.DNSClient.AddDNSEntry(context.Background(), &req)

    return err
}