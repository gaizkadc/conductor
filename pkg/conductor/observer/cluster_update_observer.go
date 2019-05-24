/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package observer

import (
    "context"
    "github.com/nalej/conductor/internal/entities"
    "github.com/nalej/conductor/pkg/conductor/baton"
    pbApplication "github.com/nalej/grpc-application-go"
    pbInfrastructure "github.com/nalej/grpc-infrastructure-go"
    "github.com/rs/zerolog/log"
    "time"
)

// The cluster update observer can be invoked when a cluster update is detected.
// It will decide when changes in cluster definitions must trigger operations such as scheduling operations.


const ClusterUpdateObserverReallocationTimeout = time.Minute * 5


type ClusterUpdateObserver struct{
    // reference to conductor's baton to perform scheduling if needed
    baton *baton.Manager
}

func NewClusterUpdateObserver (baton *baton.Manager) ClusterUpdateObserver {
    return ClusterUpdateObserver{
        baton: baton,
    }
}

func (cuo *ClusterUpdateObserver) ObserveChanges(update *pbInfrastructure.UpdateClusterRequest) {
    log.Info().Str("clusterId", update.ClusterId).Msg("started cluster update changes observation...")

    // Create a local infrastructure client
    conn := cuo.baton.ConnHelper.GetSystemModelClients().GetConnections()[0]
    infrClient := pbInfrastructure.NewClustersClient(conn)

    // get the final cluster definition stored in the system
    clusterEntry, err := infrClient.GetCluster(ctx,&pbInfrastructure.ClusterId{
        OrganizationId: update.OrganizationId, ClusterId: update.ClusterId})
    if err!=nil{
        log.Error().Err(err).Msg("impossible to retrieve clsuter status to observe changes")
        return
    }

    // get deployment fragments in the cluster
    dfs, err := cuo.baton.AppClusterDB.GetFragmentsInCluster(update.ClusterId)
    if err != nil {
        log.Error().Err(err).Msg("skip potential reallocation due to error retrieving deployment fragments")
        return
    }

    if dfs == nil || len(dfs) == 0 {
        log.Debug().Msg("no deployment fragments allocated in cluster. Skip potential reallocation")
        return
    }

    // check if the current cluster definition can allocate this deployment group

    // list of ids of descriptor fragments to be reallocated
    toReallocate := make([]ObservableDeploymentFragment,0)
    // store in this map those descriptors already requested
    descriptors := make(map[string]entities.AppDescriptor,0)
    for _, df := range dfs {
        // get the descriptor for this fragment when required
        if _, found := descriptors[df.AppDescriptorId]; !found {
            ctx, cancel := context.WithTimeout(context.Background(),time.Second*10)
            parametrizedDesc,err := cuo.baton.AppClient.GetParametrizedDescriptor(ctx,&pbApplication.AppInstanceId{
                OrganizationId: update.OrganizationId,
                AppInstanceId: df.AppInstanceId})
            cancel()
            if err!=nil{
                log.Error().Err(err).Str("fragmentId",df.FragmentId).Str("appInstanceId",df.AppInstanceId).
                    Msg("error retrieving ")
                continue
            }
            descriptors[df.AppDescriptorId] = entities.NewParametrizedDescriptorFromGRPC(parametrizedDesc)
        }

        if cuo.reallocationRequired(df, appInstance, clusterEntry) {
            toReallocate = append(toReallocate, ObservableDeploymentFragment{ClusterId:df.ClusterId,
                FragmentId: df.FragmentId, AppInstanceId: df.AppInstanceId})
        }
    }

    log.Info().Interface("toReallocate",toReallocate).
        Msgf("there is a total of %d deployment fragments to be reallocated",len(toReallocate))

    if len(toReallocate) == 0 {
        log.Info().Msg("no deployment fragments to reallocate. Exit")
        return
    }

    observer := NewDeploymentFragmentsObserver(toReallocate, cuo.baton.AppClusterDB)
    // Run an observer in a separated thread to send the schedule to the queue when is terminating
    go observer.Observe(ClusterUpdateObserverReallocationTimeout,entities.FRAGMENT_TERMINATING, cuo.baton.scheduleDeploymentFragment)

    log.Info().Str("clusterId",update.ClusterId).Msg("scheduled fragments reallocation")

    // Drain the whole cluster
    for _, fragment := range toReallocate {
        c.undeployClustersInstance(update.OrganizationId, fragment.AppInstanceId, []string{fragment.ClusterId})
    }

    log.Info().Str("clusterId", update.ClusterId).Msg("cluster update changes observation done")

}

// Determine if a deployment fragment has to be reallocated in the context of a new cluster definition for its
// current application descriptor.
func(cuo *ClusterUpdateObserver) reallocationRequired(df entities.DeploymentFragment,
    descriptor entities.AppDescriptor, cluster *pbInfrastructure.Cluster) bool {

        // Find the service group definition as stated in the application descriptor
        var serviceGroup *entities.ServiceGroup = nil
        for _, sg := range descriptor.Groups {
            // all the services in the deployment fragment belong to the same service group
            if df.Stages[0].Services[0].ServiceGroupId == sg.ServiceGroupId {
                serviceGroup = &sg
                break
            }
        }

        if serviceGroup == nil {
            // this is really strange to happen. This means inconsistency in the database
            log.Error().Interface("deploymentFragment",df).
                Msg("the service group definition could not be found for a deployment fragment")
            return true
        }

    // check if this cluster has all the required labels by the service group definition
    for key, expectedValue := range serviceGroup.Labels {
        clusterValue, found := cluster.Labels[key]
        if !found {
            log.Debug().Interface("groupLabels",serviceGroup.Labels).
                Interface("clusterLabels", cluster.Labels).Msgf("service group expects %s label", key)
            return true
        }
        if clusterValue!=expectedValue {
            log.Debug().Interface("groupLabels",serviceGroup.Labels).
                Interface("clusterLabels", cluster.Labels).
                Msgf("service group expects %s:%s not %s:%s",key,expectedValue,key,clusterValue)
            return true
        }
    }

    // everything was correct
    return false
}
