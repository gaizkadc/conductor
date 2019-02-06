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
    "github.com/nalej/grpc-utils/pkg/tools"
    "github.com/rs/zerolog/log"
    "time"
)

const MusicianQueryTimeout = time.Minute

type SimpleScorer struct {
    connHelper *utils.ConnectionsHelper
    musicians *tools.ConnectionsMap
}

func NewSimpleScorer(connHelper *utils.ConnectionsHelper) Scorer {
    return SimpleScorer{musicians: connHelper.GetClusterClients(), connHelper: connHelper}
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
    collectedScores := make([]*pbConductor.ClusterScoreResponse,0,len(s.connHelper.ClusterReference))
    found_scores := 0

    for clusterId, clusterHost := range s.connHelper.ClusterReference {

        log.Debug().Msgf("conductor query musician cluster %s at %s", clusterId, clusterHost)

        conn, err := s.musicians.GetConnection(fmt.Sprintf("%s:%d",clusterHost,utils.APP_CLUSTER_API_PORT))
        if err != nil {
            log.Error().Err(err).Msgf("impossible to get connection for %s",clusterHost)
        }

        c := pbAppClusterApi.NewMusicianClient(conn)

        res := s.queryMusician(c,requirements)

        if res != nil {
            log.Info().Interface("response",res).Msg("musician responded with score")
            collectedScores = append(collectedScores,res)
            found_scores = found_scores + 1
        } else {
            log.Warn().Msgf("querying musician %s failed, ignore it",c)
        }
    }

    if found_scores==0 {
        log.Debug().Msg("not found scores")
        collectedScores = nil
    }

    log.Debug().Msgf("returned score %v", collectedScores)
    return collectedScores
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