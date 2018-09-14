//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

package scorer

import (
    pbConductor "github.com/nalej/grpc-conductor-go"
    "github.com/nalej/conductor/internal/entities"
    "google.golang.org/grpc"
    "context"
    "github.com/rs/zerolog/log"
    "time"
)

type SimpleScorer struct {}

func NewSimpleScorer() Scorer {
    return SimpleScorer{}
}

// For a existing set of deployment requirements score potential candidates.
//  params:
//   requirements to be fulfilled
//  return:
//   candidates score
func (s SimpleScorer) ScoreRequirements (requirements *entities.Requirements, musicians []string) (*entities.ClusterScore, error) {
    for _, target := range(musicians) {
        // Set up a connection to the server.
        conn, err := grpc.Dial(target, grpc.WithInsecure())
        if err != nil {
            log.Fatal().Errs("unable to connect: %v", []error{err})
        }
        defer conn.Close()
        c := pbConductor.NewMusicianClient(conn)

        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()

        log.Info().Str("musician", target).Msg("query score")
        req:=pbConductor.ClusterScoreRequest{Requirements:"this are my requirements", RequestId:"1"}
        res, err := c.Score(ctx,&req)

        if err != nil {
            return nil, err
        }

        log.Info().Str("cluster",target).Interface("response", res)

    }
    return nil,nil
}