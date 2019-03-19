/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package scorer

import (
    "context"
    "errors"
    "fmt"
    "github.com/google/uuid"
    "github.com/nalej/conductor/internal/entities"
    "github.com/nalej/conductor/pkg/utils"
    pbAppClusterApi "github.com/nalej/grpc-app-cluster-api-go"
    pbConductor "github.com/nalej/grpc-conductor-go"
    pbInfrastructure "github.com/nalej/grpc-infrastructure-go"
    "github.com/nalej/grpc-utils/pkg/tools"
    "github.com/rs/zerolog/log"
    "time"
)

const MusicianQueryTimeout = time.Minute

type SimpleScorer struct {
    connHelper *utils.ConnectionsHelper
    musicians *tools.ConnectionsMap
    // Infrastructure client
    clusterClient pbInfrastructure.ClustersClient
}

func NewSimpleScorer(connHelper *utils.ConnectionsHelper) Scorer {
    // initialize clients
    pool := connHelper.GetSystemModelClients()
    if pool!=nil && len(pool.GetConnections())==0{
        log.Panic().Msg("system model clients were not started")
        return nil
    }
    conn := pool.GetConnections()[0]
    // Create associated clients
    clusterClient := pbInfrastructure.NewClustersClient(conn)

    return SimpleScorer{musicians: connHelper.GetClusterClients(), connHelper: connHelper, clusterClient: clusterClient}
}

// For a existing set of deployment requirements score potential candidates.
//  params:
//   requirements to be fulfilled
//  return:
//   candidates score
func (s SimpleScorer) ScoreRequirements (organizationId string, requirements *entities.Requirements) (*entities.DeploymentScore, error) {
    if requirements == nil {
        nil_error := errors.New("impossible to score nil requirements")
        log.Error().Err(nil_error)
        return nil, nil_error
    }
    scores := s.collectScores(organizationId, requirements)

    if scores == nil {
        noScores := errors.New("no available scores found")
        log.Error().Err(noScores).Msg("simple scorer could not collect any score")
        return nil, noScores
    }

    clusterScores := entities.NewClustersScore()

    for _, s := range scores {
        // Create a set of scores for different combinations of service groups
        collectedScores := entities.NewClusterDeploymentScore(s.ClusterId)
        for _, x := range s.Score {
            collectedScores.AddScore(x.GroupServiceInstances,x.Score)
        }
        clusterScores.AddClusterScore(collectedScores)
    }

    log.Debug().Str("component","conductor").Interface("score",clusterScores).Msg("final found scores")
    return &clusterScores,nil
}

// Internal method to query known clusters about requirements scoring.
func (s SimpleScorer) collectScores(organizationId string, requirements *entities.Requirements)[]*pbConductor.ClusterScoreResponse{

    err := s.connHelper.UpdateClusterConnections(organizationId)
    if err != nil {
        log.Error().Err(err).Msgf("error updating connections for organization %s", organizationId)
        return nil
    }
    if len(s.connHelper.ClusterReference) == 0 {
        log.Error().Msgf("no clusters found for organization %s", organizationId)
        return nil
    }


    // we expect as many scores as musicians we have
    log.Debug().Msgf("we have %d known clusters",len(s.connHelper.ClusterReference))
    collectedScores := make([]*pbConductor.ClusterScoreResponse,0,0)

    found_scores := 0

    for clusterId, clusterHost := range s.connHelper.ClusterReference {

        // Check what requests can be sent to this cluster
        requestsToSend := s.findRequirementsCluster(organizationId, clusterId, requirements)
        if requestsToSend != nil {
            // there is something to send

            log.Debug().Msgf("conductor query musician cluster %s at %s", clusterId, clusterHost)

            conn, err := s.musicians.GetConnection(fmt.Sprintf("%s:%d",clusterHost,utils.APP_CLUSTER_API_PORT))
            if err != nil {
                log.Error().Err(err).Msgf("impossible to get connection for %s",clusterHost)
            }

            c := pbAppClusterApi.NewMusicianClient(conn)

            res := s.queryMusician(c,requestsToSend)

            if res != nil {
                log.Error().Err(err).Msg("impossible to query musician to obtain requirements score. Ignore it.")
            } else {
                log.Info().Interface("response",res).Msg("musician responded with score")
                collectedScores = append(collectedScores,res)
                found_scores = found_scores + 1
            }
        }
    }

    if found_scores==0 {
        log.Debug().Msg("not found scores")
        collectedScores = nil
    }

    log.Debug().Msgf("returned score %v", collectedScores)
    return collectedScores
}

// Private function to decide what requirements can be sent to a cluster in order to ask the musician. This decision is
// done based on the cluster deployment selector tags. The function returns a requirements entry or nil if nothing to send.
func (s SimpleScorer) findRequirementsCluster(organizationId string, clusterId string,requirements *entities.Requirements) *entities.Requirements {
    cluster, err := s.clusterClient.GetCluster(context.Background(),&pbInfrastructure.ClusterId{OrganizationId: organizationId,ClusterId: clusterId})
    if err != nil {
        log.Error().Err(err).Msg("impossible to return cluster information when checking requirements")
        return nil
    }
    filteredRequirements := entities.NewRequirements()
    for _, req := range requirements.List {
        if req.DeploymentSelectors == nil {
            // no specs, add it
            filteredRequirements.AddRequirement(req)
        } else if cluster.Labels != nil {
            // there as specs, and the cluster has labels. Check it
            allMatch := true
            for k,v := range req.DeploymentSelectors {
                clusterValue, found := cluster.Labels[k]
                if !found || clusterValue != v {
                    allMatch = false
                    break
                }
            }
            log.Debug().Interface("group selectors",req.DeploymentSelectors).
                Interface("cluster labels", cluster.Labels).Bool("match",allMatch).
                Msg("comparing cluster labels")
            if allMatch {
                // add it to the list of requirements
                filteredRequirements.AddRequirement(req)
            }
        }
    }

    if len(filteredRequirements.List) == 0 {
        // no requirements for this cluster
        log.Debug().Str("clusterId",cluster.ClusterId).Msg("no requirements matching the cluster")
        return nil
    }

    return &filteredRequirements
}

// Private function to query a target musician about the score of a given set of requirements.
func (s SimpleScorer) queryMusician(musicianClient pbAppClusterApi.MusicianClient, requirements *entities.Requirements) *pbConductor.ClusterScoreResponse{

    ctx, cancel := context.WithTimeout(context.Background(), MusicianQueryTimeout)
    defer cancel()

    req:=pbConductor.ClusterScoreRequest{
        RequestId: uuid.New().String(),
        Requirements: requirements.ToGRPC(),
    }
    res, err := musicianClient.Score(ctx,&req)

    if err != nil {
        log.Error().Err(err).Msg("errors found querying musician")
        return nil
    }

    return res
}