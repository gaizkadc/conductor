/*
 * Copyright 2018 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package scorer

import (
    pbConductor "github.com/nalej/grpc-conductor-go"
    "github.com/nalej/conductor/internal/entities"
    "context"
    "github.com/rs/zerolog/log"
    "time"
    "github.com/nalej/conductor/tools"
    "github.com/nalej/conductor/pkg/conductor"
    "errors"
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
func (s SimpleScorer) ScoreRequirements (requirements *entities.Requirements) (*entities.ClusterScore, error) {
    if requirements == nil {
        nil_error := errors.New("impossible to score nil requements")
        log.Error().Err(nil_error)
        return nil, nil_error
    }
    scores := s.collectScores(requirements)

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
func (s SimpleScorer) collectScores(requirements *entities.Requirements) []*pbConductor.ClusterScoreResponse{
    // we expect as many scores as musicians we have
    musicians := s.musicians.GetConnections()
    collected_scores := make([]*pbConductor.ClusterScoreResponse,0,len(musicians))
    found_scores := 0
    for _, conn := range  musicians {
        log.Debug().Interface("musician", conn.Target()).Msg("conductor query score")

        c := pbConductor.NewMusicianClient(conn)

        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()

        req:=pbConductor.ClusterScoreRequest{RequestId: requirements.RequestID,
                                            Disk: requirements.Disk,
                                            Memory: requirements.Memory,
                                            Cpu: requirements.CPU}
        res, err := c.Score(ctx,&req)

        if err != nil {
            log.Error().Err(err).Msg("errors found querying musician")
        } else {
            if res==nil{
                log.Error().Err(errors.New("musician returned nil response"))
            } else {
                log.Info().Interface("response",res).Msg("musician responded with score")
                collected_scores = append(collected_scores,res)
                found_scores = found_scores + 1
            }
        }
    }


    if found_scores==0 {
        log.Debug().Msg("not found scores")
        collected_scores = nil
    }

    log.Debug().Msgf("returned score %v", collected_scores)
    return collected_scores
}