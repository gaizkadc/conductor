/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package structures

import (
    "github.com/nalej/conductor/internal/entities"
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
    GroupsCluster map[string][]string
}

func NewDeploymentMatrix() DeploymentMatrix {
    return DeploymentMatrix{
        ClustersScore: make(map[string]entities.ClusterDeploymentScore, 0),
        AllocatedScore: make(map[string]entities.ClusterDeploymentScore, 0),
        GroupsCluster: make(map[string][]string),
    }
}

func (dm *DeploymentMatrix) AddClusterScore(score entities.ClusterDeploymentScore) {
    dm.ClustersScore[score.ClusterId] = score
    dm.AllocatedScore[score.ClusterId] = entities.NewClusterDeploymentScore(score.ClusterId)
}


func (dm *DeploymentMatrix) FindBestTargetForGroups(groups []string) string {
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


// Allocate groups and update scores.
func (dm *DeploymentMatrix) allocateGroups(clusterId string, groupId string,groups []string) {
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
func (dm *DeploymentMatrix) generateGroupId(groups []string) string {
    sort.Strings(groups)
    concatenated := ""
    for _, s := range groups {
        concatenated = concatenated + s
    }
    return concatenated
}

