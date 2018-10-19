/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package scorer

import (
    pbConductor "github.com/nalej/grpc-conductor-go"
    "github.com/nalej/conductor/internal/entities"
    "context"
    "github.com/rs/zerolog/log"
    "time"
    "github.com/nalej/grpc-utils/pkg/tools"
    "github.com/nalej/conductor/pkg/conductor"
    "errors"
    "github.com/google/uuid"
    "github.com/nalej/conductor/pkg/utils"
    "fmt"
)

type SimpleScorer struct {
    musicians *tools.ConnectionsMap
}

func NewSimpleScorer() Scorer {
    return SimpleScorer{musicians: conductor.GetMusicianClients()}
}

// For a existing set of deployment requirements score potential candidates.
//  params:
//   requirements to be fulfilled
//  return:
//   candidates score
func (s SimpleScorer) ScoreRequirements (organizationId string, requirements *entities.Requirements) (*entities.ClusterScore, error) {
    if requirements == nil {
        nil_error := errors.New("impossible to score nil requirements")
        log.Error().Err(nil_error)
        return nil, nil_error
    }
    scores := s.collectScores(organizationId, requirements)

    if scores == nil {
        no_scores := errors.New("no available scores found")
        log.Error().Err(no_scores).Msg("simple scorer could not collect any score")
        return nil, no_scores
    }
    // evaluate scores

    // check what can we say with the returned values
    // Find maxi score
    max := float32(-1.0)
    var candidate *pbConductor.ClusterScoreResponse = nil
    evaluated := 0
    for _, musicianScore := range scores {
        if musicianScore.Score > max {
            candidate = musicianScore
        }
        evaluated = evaluated + 1
    }

    if candidate == nil {
        no_candidate := errors.New("no candidate fulfils the requirements")
        log.Error().Err(no_candidate)
        return nil, no_candidate
    }

    finalScore := entities.ClusterScore{RequestID: candidate.RequestId,
                                        ClusterID: candidate.ClusterId,
                                        Score: candidate.Score,
                                        TotalEvaluated: evaluated}

    log.Debug().Str("component","conductor").Interface("score",finalScore).Msg("final score found")
    return &finalScore,nil
}

// Internal method to query known clusters about requirements scoring.
func (s SimpleScorer) collectScores(organizationId string, requirements *entities.Requirements)[]*pbConductor.ClusterScoreResponse{

    clusters := conductor.UpdateClusterConnections(organizationId)
    if clusters == nil || len(clusters) == 0 {
        log.Error().Msgf("no clusters found for oganization %s", organizationId)
        return nil
    }


    // we expect as many scores as musicians we have
    log.Debug().Msgf("we have %d known clusters",len(clusters))
    collected_scores := make([]*pbConductor.ClusterScoreResponse,0,len(clusters))
    found_scores := 0
    for _, c := range  clusters {
        log.Debug().Interface("musician", c).Msg("conductor query score")

        conn, err := s.musicians.GetConnection(fmt.Sprintf("%s:%d",c,utils.MUSICIAN_PORT))
        if err != nil {
            log.Error().Err(err).Msgf("impossible to get connection for %s",c)
        }

        c := pbConductor.NewMusicianClient(conn)

        res := s.queryMusician(c,requirements)

        if res != nil {
            log.Info().Interface("response",res).Msg("musician responded with score")
            collected_scores = append(collected_scores,res)
            found_scores = found_scores + 1
        } else {
            log.Warn().Msgf("querying musician %s failed, ignore it",c)
        }
    }


    if found_scores==0 {
        log.Debug().Msg("not found scores")
        collected_scores = nil
    }

    log.Debug().Msgf("returned score %v", collected_scores)
    return collected_scores
}

// Private function to query a target musician about the score of a given set of requirements.
func (s SimpleScorer) queryMusician(musicianClient pbConductor.MusicianClient, requirements *entities.Requirements) *pbConductor.ClusterScoreResponse{

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    req:=pbConductor.ClusterScoreRequest{RequestId: uuid.New().String(),
        Disk: requirements.Disk,
        Memory: requirements.Memory,
        Cpu: requirements.CPU}
    res, err := musicianClient.Score(ctx,&req)

    if err != nil {
        log.Error().Err(err).Msg("errors found querying musician")
        return nil
    }

    return res
}