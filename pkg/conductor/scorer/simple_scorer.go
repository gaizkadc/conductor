//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

package scorer

import (
    pbConductor "github.com/nalej/grpc-conductor-go"
    "github.com/nalej/conductor/internal/entities"
    "context"
    "github.com/rs/zerolog/log"
    "time"
    "github.com/nalej/conductor/tools"
    "github.com/nalej/conductor/pkg/conductor"
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
    for _, conn := range s.musicians.GetConnections() {
        log.Info().Interface("musician", conn.Target()).Msg("conductor query score")

        c := pbConductor.NewMusicianClient(conn)

        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()

        req:=pbConductor.ClusterScoreRequest{Requirements:"this are my requirements"}
        res, err := c.Score(ctx,&req)

        if err != nil {
            return nil, err
        }

        log.Info().Str("cluster",conn.Target()).Interface("response", res)
    }

    to_return := entities.ClusterScore{Score: 10.0}
    return &to_return, nil

}