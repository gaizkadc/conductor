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

func(s *SimpleScorer) Score(request *pbConductor.ClusterScoreRequest) (*pbConductor.ClusterScoreResponse, error){
    log.Debug().Msg("musician simple scorer queried")
    // check
    status, err := s.collector.GetStatus()

    if err != nil {
        log.Error().Err(err)
        return nil, err
    }

    log.Debug().Interface("status",status).Msg("musician found status")
    // TODO check the coherence of this data type.
    var totalCPU float32 = 0
    var totalMem float32 = 0
    var totalStorage float32 = 0

    for _, req := range request.Requirements {
        totalCPU = totalCPU + float32(req.Cpu)
        totalMem = totalMem + float32(req.Memory)
        totalStorage = totalStorage + float32(req.Storage)
    }



    // compute score based on requested and available
    dCPU := (1-float32(status.CPU)) - totalCPU
    dMem := (float32(status.Mem) - totalMem) / float32(status.Mem)
    dDisk := (float32(status.Disk) - totalStorage) / float32(status.Disk)

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