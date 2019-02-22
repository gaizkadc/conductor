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
    toReturn := make([]string,0)
    for _, clusterScore := range dm.AllocatedScore {
        // take a look at the score for this group
        groupScoreInCluster, found := clusterScore.Scores[group.Name]
        if !found {
            msg := fmt.Sprintf("cluster %s has no score for group %s", clusterScore.ClusterId, group.Name)
            err := derrors.NewFailedPreconditionError(msg)
            return nil, err
        }
        if groupScoreInCluster <= 0 {
            msg := fmt.Sprintf("cluster %s has negative score for group %s", clusterScore.ClusterId, group.Name)
            err := derrors.NewFailedPreconditionError(msg)
            return nil, err
        }
        // Positive score. Allocate and update matrix.
        dm.allocateGroups(clusterScore.ClusterId, group.Name,[]entities.ServiceGroup{group})
        toReturn = append(toReturn, clusterScore.ClusterId)
    }

    return toReturn, nil
}


// Find the best target to deploy a set of groups and update the matrix accordingly.
func (dm *DeploymentMatrix) FindBestTargetForGroups(groups []entities.ServiceGroup) string {
    log.Debug().Interface("matrix",dm).Msg("FindBestTargetForGroups")
    groupId := dm.generateGroupId(groups)
    // find the cluster with the largest score
    maxScore := float32(0)
    candidate := ""
    for _, clusterScore := range dm.AllocatedScore {
        score, found := clusterScore.Scores[groupId]
        if found {
            if score > maxScore {
                maxScore = score
                candidate = clusterScore.ClusterId
            }
        } else {
            log.Debug().Str("clusterId",clusterScore.ClusterId).Str("groupId", groupId).
                Msg("set of groups not found")
        }
    }

    if candidate != "" {
        dm.allocateGroups(candidate, groupId, groups)
    }

    return candidate
}

// Find the best target to deploy a single group and update the matrix accordingly.
func (dm *DeploymentMatrix) FindBestTargetForGroup(group entities.ServiceGroup) (string, derrors.Error) {
    log.Debug().Interface("matrix",dm).Msg("FindBestTargetForGroup")
    groupId := group.Name
    // find the cluster with the largest score
    maxScore := float32(0)
    candidate := ""
    for _, clusterScore := range dm.AllocatedScore {
        score, found := clusterScore.Scores[groupId]
        if found {
            if score > maxScore {
                maxScore = score
                candidate = clusterScore.ClusterId
            }
        } else {
            log.Debug().Str("clusterId",clusterScore.ClusterId).Str("groupId", groupId).
                Msg("set of groups not found")
        }
    }

    if candidate != "" {
        dm.allocateGroups(candidate, groupId, []entities.ServiceGroup{group})
        return candidate, nil
    }

    return "", derrors.NewGenericError("impossible to find cluster for single group deployment")


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

