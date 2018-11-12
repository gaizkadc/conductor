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
)

type Manager struct{
    // Networking manager client
    NetClient pbNetwork.NetworksClient
}

func NewManager() (*Manager, error){
    // Network client
    netPool := conductor.GetNetworkingClients()
    if netPool != nil && len(netPool.GetConnections())==0{
        log.Panic().Msg("networking client was not started")
        return nil, errors.New("networking client was not started")
    }
    netClient := pbNetwork.NewNetworksClient(netPool.GetConnections()[0])

    return &Manager{netClient}, nil
}

func (m *Manager) AuthorizeNetworkMembership(organizationId string, networkId string, memberId string) error {
    req := pbNetwork.AuthorizeMemberRequest{
        OrganizationId: organizationId,
        NetworkId: networkId,
        MemberId: memberId,}
    _, err := m.NetClient.AuthorizeMember(context.Background(), &req)
    if err != nil {
        log.Error().Err(err).Msgf("AuthorizeNetworkMembership failed for %#v")
        return err
    }
    return nil
}

func (m *Manager) RegisterNetworkEntry() error {
    return nil
}