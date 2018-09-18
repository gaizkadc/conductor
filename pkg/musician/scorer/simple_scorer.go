//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

package scorer

import (
    pbConductor "github.com/nalej/grpc-conductor-go"
    "github.com/rs/zerolog/log"
)

type SimpleScorer struct {

}

func NewSimpleScorer() Scorer {
    return &SimpleScorer{}
}

func(s *SimpleScorer) Score(request *pbConductor.ClusterScoreRequest) (*pbConductor.ClusterScoreResponse, error){
    log.Debug().Msg("simple scorer mentioned")
    return &pbConductor.ClusterScoreResponse{RequestId: "cluster score reponse", Score: 0.1}, nil
}