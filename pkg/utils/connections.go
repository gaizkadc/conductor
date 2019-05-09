/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package utils

// Set of common routines for conductor components. A pool of already opened client connections is maintained
// for the components below and implemented in a singleton instance accessible by all the elements in this package.
// When running tests, this pool uses listening buffers.

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "errors"
    "fmt"
    "github.com/nalej/derrors"
    pbInfrastructure "github.com/nalej/grpc-infrastructure-go"
    pbOrganization "github.com/nalej/grpc-organization-go"
    "github.com/nalej/grpc-utils/pkg/tools"
    "github.com/rs/zerolog/log"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "io/ioutil"
    "sync"
)

// Internal struct to store information about a cluster connection. This struct can be used to query the latest
// known status of a cluster.
type ClusterEntry struct {
    // Hostname for this cluster
    Hostname string
    // Cordon true if this cluster is in a cordon status
    Cordon bool
}


type ConnectionsHelper struct {
    // Singleton instance of connections with cluster clients clients
    ClusterClients *tools.ConnectionsMap
    onceClusters   sync.Once
    // Singleton instance of connections with the system model
    SMClients *tools.ConnectionsMap
    onceSM sync.Once
    // Singleton instance of connections with the network client
    NetworkingClients *tools.ConnectionsMap
    onceNC sync.Once
    // Translation map between cluster ids and their ip addresses
    ClusterReference map[string]ClusterEntry
    // useTLS connections
    useTLS bool
    // path for the CA
    caCertPath string
    // skip CA validation
    skipCAValidation bool
    // Singleton instance of connections with the Authx
    AuthxClients *tools.ConnectionsMap
    onceAuthx sync.Once
    // Singleton instance of connections with the UnifiedLogging
    UnifiedLoggingClients *tools.ConnectionsMap
    onceUnifiedLogging sync.Once

}

func NewConnectionsHelper(useTLS bool, caCertPath string, skipCAValidation bool) *ConnectionsHelper {

    return &ConnectionsHelper{
        ClusterReference: make(map[string]ClusterEntry, 0),
        useTLS: useTLS,
        caCertPath: caCertPath,
        skipCAValidation: skipCAValidation,
    }
}

func(h *ConnectionsHelper) GetSystemModelClients() *tools.ConnectionsMap {
    h.onceSM.Do(func(){
        // reuse the conductor factory
        h.SMClients = tools.NewConnectionsMap(systemModelClientFactory)
    })
    return h.SMClients
}

func(h *ConnectionsHelper) GetAuthxClients() *tools.ConnectionsMap {
    h.onceAuthx.Do(func(){
        // reuse the conductor factory
        h.AuthxClients = tools.NewConnectionsMap(authxClientFactory)
    })
    return h.AuthxClients
}

func(h *ConnectionsHelper) GetUnifiedLoggingClients() *tools.ConnectionsMap {
    h.onceUnifiedLogging.Do(func(){
        // reuse the conductor factory
        h.UnifiedLoggingClients = tools.NewConnectionsMap(unifiedLoggingClientFactory)
    })
    return h.UnifiedLoggingClients
}

func(h *ConnectionsHelper) GetClusterClients() *tools.ConnectionsMap {
    h.onceClusters.Do(func(){
        h.ClusterClients = tools.NewConnectionsMap(clusterClientFactory)
        if h.ClusterReference == nil {
            h.ClusterReference = make(map[string]ClusterEntry, 0)
        }
    })
    return h.ClusterClients
}


func(h *ConnectionsHelper) GetNetworkingClients() *tools.ConnectionsMap {
    h.onceNC.Do(func() {
        h.NetworkingClients = tools.NewConnectionsMap(networkingClientFactory)
        if h.NetworkingClients == nil {
            h.NetworkingClients = tools.NewConnectionsMap(networkingClientFactory)
        }
    })
    return h.NetworkingClients
}


// Factory in charge of generating new basic connections with a grpc server.
//  params:
//   hostname
//   port
//   params
//  return:
//   grpc connection and error if any
func basicClientFactory(hostname string, port int, params...interface{}) (*grpc.ClientConn, error) {
    address := fmt.Sprintf("%s:%d",hostname,port)
    log.Debug().Str("address", address).Msg("basicClientFactory being created")
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Fatal().Msgf("Failed to start gRPC connection: %v", err)
    }
    log.Info().Msgf("Connected to address at %s", address)
    return conn, err
}

// Factory in charge of generation a secure connection with a grpc server.
//  params:
//   hostname of the target server
//   port of the target server
//   useTLS flag indicating whether to use the TLS security
//   caCert path of the CA certificate
//   skipCAValidation skip the validation of the CA
//  return:
//   grpc connection and error if any
func secureClientFactory(hostname string, port int, useTLS bool, caCertPath string, skipCAValidation bool) (*grpc.ClientConn, error) {
    rootCAs := x509.NewCertPool()
    tlsConfig := &tls.Config{
        ServerName:   hostname,
    }

    if caCertPath != "" {
        log.Debug().Str("caCertPath", caCertPath).Msg("loading CA cert")
        caCert, err := ioutil.ReadFile(caCertPath)
        if err != nil {
            return nil, derrors.NewInternalError("Error loading CA certificate")
        }
        added := rootCAs.AppendCertsFromPEM(caCert)
        if !added {
            return nil, derrors.NewInternalError("cannot add CA certificate to the pool")
        }
        tlsConfig.RootCAs = rootCAs
    }

    targetAddress := fmt.Sprintf("%s:%d", hostname, port)
    log.Debug().Str("address", targetAddress).Bool("useTLS", useTLS).Str("caCertPath", caCertPath).Bool("skipCAValidation", skipCAValidation).Msg("creating secure connection")

    if skipCAValidation {
        tlsConfig.InsecureSkipVerify = true
    }

    creds := credentials.NewTLS(tlsConfig)

    log.Debug().Interface("creds", creds.Info()).Msg("Secure credentials")
    sConn, dErr := grpc.Dial(targetAddress, grpc.WithTransportCredentials(creds))
    if dErr != nil {
        log.Error().Err(dErr).Msg("impossible to create secure client factory connection")
        return nil, derrors.AsError(dErr, "cannot create connection with the signup service")
    }

    return sConn, nil

}

// Factory in charge of generating new connections for Conductor->system model.
//  params:
//   hostname
//   port
//   params
//  return:
//   client and error if any
func systemModelClientFactory(hostname string, port int, params...interface{}) (*grpc.ClientConn, error) {
    return basicClientFactory(hostname, port)
}

// Factory in charge of generating new connections for Conductor->authx.
//  params:
//   hostname
//   port
//   params
//  return:
//   client and error if any
func authxClientFactory(hostname string, port int, params...interface{}) (*grpc.ClientConn, error) {
    return basicClientFactory(hostname, port)
}

// Factory in charge of generating new connections for Conductor->authx.
//  params:
//   hostname
//   port
//   params
//  return:
//   client and error if any
func unifiedLoggingClientFactory(hostname string, port int, params...interface{}) (*grpc.ClientConn, error) {
    return basicClientFactory(hostname, port)
}

// Factory in charge of generating new connections for Conductor->cluster communication.
//  params:
//   hostname of the target server
//   port of the target server
//   useTLS flag indicating whether to use the TLS security
//   caCert path of the CA certificate
//   skipCAValidation skip the validation of the CA
//  return:
//   client and error if any
func clusterClientFactory(hostname string, port int, params...interface{}) (*grpc.ClientConn, error) {
    log.Debug().Str("hostname", hostname).Int("port", port).Int("len", len(params)).Interface("params", params).Msg("calling cluster client factory")
    if len(params) != 3 {
        log.Fatal().Interface("params",params).Msg("cluster client factory called with not enough parameters")
    }
    useTLS := params[0].(bool)
    caCertPath := params[1].(string)
    skipCAValidation := params[2].(bool)
    return secureClientFactory(hostname, port, useTLS, caCertPath, skipCAValidation)
}


// Factory in charge of generating new connections for Conductor->Networkcommunication.
//  params:
//   hostname
//   port
//   params
//  return:
//   client and error if any
func networkingClientFactory(hostname string, port int, params...interface{}) (*grpc.ClientConn, error) {
    return basicClientFactory(hostname, port)
}



// This is a common sharing function to check the system model and update the available clusters.
// Additionally, the function updates the available connections for musicians and deployment managers.
// The common ClusterReference object is updated with the cluster ids and the corresponding ip.
//  params:
//   organizationId
func(h *ConnectionsHelper) UpdateClusterConnections(organizationId string) error{
    log.Debug().Msg("update cluster connections...")
    // Rebuild the map
    h.ClusterReference = make(map[string]ClusterEntry,0)

    cmClients := h.GetSystemModelClients()
    // no available system model client
    if cmClients.NumConnections() == 0 {
        log.Error().Msg("there are no available system model clients")
        return errors.New("there are no available system model clients")
    }

    // Get an infrastructure client and check the available clusters
    client := pbInfrastructure.NewClustersClient(cmClients.GetConnections()[0])
    req := pbOrganization.OrganizationId{OrganizationId:organizationId}
    clusterList, err := client.ListClusters(context.Background(), &req)
    if err != nil {
        msg := fmt.Sprintf("there was a problem getting the list of " +
            "available cluster for org %s",organizationId)
        log.Error().Err(err).Msg(msg)
        return errors.New(msg)
    }

    toReturn := make([]string,0)
    clusters := h.GetClusterClients()

    for _, cluster := range clusterList.Clusters {
        // The cluster is running and is not in cordon status
        if h.isClusterAvailable(cluster){
            targetHostname := fmt.Sprintf("appcluster.%s", cluster.Hostname)
            log.Debug().Str("clusterId", cluster.ClusterId).Str("hostname", cluster.Hostname).Str("targetHostname", targetHostname).Msg("add connection to cluster")
            h.ClusterReference[cluster.ClusterId] = ClusterEntry{Hostname: targetHostname, Cordon: cluster.Cordon}
            targetPort := int(APP_CLUSTER_API_PORT)
            params := make([]interface{}, 0)
            params = append(params, h.useTLS)
            params = append(params, h.caCertPath)
            params = append(params, h.skipCAValidation)

            clusters.AddConnection(targetHostname, targetPort, params ... )
            toReturn = append(toReturn, targetHostname)
        }
    }
    return nil
}

// Internal function to check if a cluster meets all the conditions to be added to the list of available clusters.
func (h * ConnectionsHelper) isClusterAvailable(cluster *pbInfrastructure.Cluster) bool {
    if cluster.Status != pbInfrastructure.InfraStatus_RUNNING {
        log.Debug().Str("clusterID", cluster.ClusterId).Msg("cluster ignored because it is not running")
        return false
    }
    // Others...
    return true
}


