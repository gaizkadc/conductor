/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package scorer

import (
    pbConductor "github.com/nalej/grpc-conductor-go"
    "github.com/rs/zerolog/log"
    "github.com/nalej/conductor/pkg/musician/statuscollector"
    "os"
    "github.com/nalej/conductor/pkg/utils"
)

type SimpleScorer struct {
    collector statuscollector.StatusCollector
}

func NewSimpleScorer(collector statuscollector.StatusCollector) Scorer {
    return &SimpleScorer{collector}
}

// TODO collect cluster id
func(s *SimpleScorer) Score(request *pbConductor.ClusterScoreRequest) (*pbConductor.ClusterScoreResponse, error){
    log.Debug().Msg("musician simple scorer queried")
    // check
    status, err := s.collector.GetStatus()

    if err != nil {
        log.Error().Err(err)
        return nil, err
    }

    log.Debug().Interface("status",status).Msg("musician found status")
    // compute score based on requested and available
    dCPU := (1-float32(status.CPU)) - request.Cpu
    dMem := (float32(status.Mem) - request.Memory) / float32(status.Mem)
    dDisk := (float32(status.Disk) - request.Disk) / float32(status.Disk)

    var score float32
    if dCPU * dMem * dDisk < 0 {
        score = -1
    }

    score = dCPU + dMem + dDisk

    //log.Debug().Str("component", "musician").Msgf("(%f-%f) + (%f-%f) +(%f-%f)",
    //    float32(status.CPU),request.Cpu,float32(status.Mem),request.Memory,float32(status.Disk),request.Disk)
    log.Debug().Str("component", "musician").Msgf("%f + %f + %f = %f",dCPU, dMem, dDisk, score)

    // TODO recover cluster id from a cluster environment variable
    return &pbConductor.ClusterScoreResponse{RequestId: request.RequestId, Score: score,
        ClusterId: os.Getenv(utils.MUSICIAN_CLUSTER_ID)}, nil
}