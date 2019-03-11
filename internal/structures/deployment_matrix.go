/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package structures

import (
    "fmt"
    "github.com/nalej/conductor/internal/entities"
    "github.com/nalej/derrors"
    "github.com/rs/zerolog/log"
    "sort"
)

// The deployment matrix is a board containing information about deployments of groups across different clusters
// and their impact according to the scores collected by musicians.

type DeploymentMatrix struct {
    // Score for a cluster
    // Cluster id -> deployment score
    ClustersScore map[string]entities.ClusterDeploymentScore
    // Allocated scores during deployment analysis
    // Cluster -> allocated score
    AllocatedScore map[string]entities.ClusterDeploymentScore
    // Groups per cluster
    // cluster -> [groupA, groupB, groupC...]
    GroupsCluster map[string][]entities.ServiceGroup
}

// Build a deployment matrix using an existing DeploymentScore
func NewDeploymentMatrix(scores entities.DeploymentScore) *DeploymentMatrix {
    clusterScore := make(map[string]entities.ClusterDeploymentScore, 0)
    allocatedScore := make(map[string]entities.ClusterDeploymentScore, 0)

    // initialize data structures
    for _, ds := range scores.DeploymentsScore {
        clusterScore[ds.ClusterId] = ds
        allocatedScore[ds.ClusterId] = ds
    }

    return &DeploymentMatrix{
        ClustersScore: clusterScore,
        AllocatedScore: allocatedScore,
        GroupsCluster: make(map[string][]entities.ServiceGroup),
    }
}

func (dm *DeploymentMatrix) AddClusterScore(score entities.ClusterDeploymentScore) {
    dm.ClustersScore[score.ClusterId] = score
    dm.AllocatedScore[score.ClusterId] = entities.NewClusterDeploymentScore(score.ClusterId)
}


// Find the best targets for replicating a group across the deployment. The current solution checks
// that all the scores are larger than for the group. If not all the groups permit to allocate this
// group we return an error,
func (dm *DeploymentMatrix) FindBestTargetsForReplication(group entities.ServiceGroup) ([]string, derrors.Error) {
    // Iterate until we find the best solution to deploy as many replicas as required
    var desiredReplicas int

    if group.Specs.MultiClusterReplica{
        // This is a multiple cluster. Replicate as many times as available clusters we have.
        desiredReplicas = len(dm.AllocatedScore)
    } else {
        // Deploy as many replicas as mentioned in the deploy specs.
        desiredReplicas = int(group.Specs.NumReplicas)
    }

    targetClusters := make(map[string]float32, 0)
    // Allocate a replica in the cluster with the largest score
    // Greedy approach to find the best cluster with no allocated replica
    for i := 0; i < desiredReplicas; i++ {
        roundCandidate := ""
        candidateScore := float32(-1)
        for clusterId, clusterScore := range dm.AllocatedScore {
            _, usedCluster := targetClusters[clusterId]
            if !usedCluster {
                // We have not allocated anything in this cluster
                groupScoreInCluster, found := clusterScore.Scores[group.Name]
                if !found {
                    msg := fmt.Sprintf("cluster %s has no score for group %s", clusterScore.ClusterId, group.Name)
                    log.Warn().Msg(msg)
                } else if groupScoreInCluster >= 0 && groupScoreInCluster > candidateScore{
                    // Consider this cluster a potential candidate
                    roundCandidate = clusterId
                    candidateScore = groupScoreInCluster
                }
            }
        }
        if roundCandidate != "" {
            targetClusters[roundCandidate] = candidateScore
        } else {
            // if a multicluster replica set was chosen we made our best effort. If not we cannot allocate
            // the number of expected replicas.
            if !group.Specs.MultiClusterReplica {
                // It was impossible to allocate a remaining replica...
                msg := fmt.Sprintf("only %d replicas could be allocated out of the %d desired",len(targetClusters), desiredReplicas)
                return nil, derrors.NewUnavailableError(msg)
            }
        }
    }

    if len(targetClusters) == 0 {
        // no replicas were set
        return nil, derrors.NewUnavailableError("no replicas could be allocated")
    }

    // Allocate all the replicas we could find
    toReturn := make([]string,len(targetClusters))
    i := 0
    for clusterId, _ := range targetClusters {
        dm.allocateGroups(clusterId, group.Name,[]entities.ServiceGroup{group})
        toReturn[i] = clusterId
        i++
    }

    return toReturn, nil
}


// Allocate groups and update scores.
func (dm *DeploymentMatrix) allocateGroups(clusterId string, groupId string,groups []entities.ServiceGroup) {
    dm.GroupsCluster[clusterId] = groups
    scoreToRevisit := dm.AllocatedScore[clusterId]
    load := dm.AllocatedScore[clusterId].Scores[groupId]
    // substract scoring value for all the entries in this cluster
    for _,v := range scoreToRevisit.Scores {
        v = v - load
    }
    // Update the set of groups
    for _, g := range groups {
        dm.GroupsCluster[clusterId] = append(dm.GroupsCluster[clusterId], g)
    }
}

// Local function to generate the concatenated id for a group of group ids.
// E.G.:
// [groupA, groupC, groupB] -> groupAgroupBgroupC
func (dm *DeploymentMatrix) generateGroupId(groups []entities.ServiceGroup) string {
    groupNames := make([]string, len(groups))
    for i, g := range groups {
        groupNames[i] = g.Name
    }
    sort.Strings(groupNames)
    concatenated := ""
    for _, s := range groupNames {
        concatenated = concatenated + s
    }
    return concatenated
}
