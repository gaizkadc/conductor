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
    log.Debug().Interface("request",request).Msg("musician simple scorer queried")
    // check
    status, err := s.collector.GetStatus()

    if err != nil {
        log.Error().Err(err)
        return nil, err
    }

    foundScores := make([]*pbConductor.DeploymentScore,0)

    // Compute the combinations of requirements
    sets := s.getCombinations(request.Requirements)
    for _, s := range sets {
        var totalCPU float32 = 0
        var totalMem float32 = 0
        var totalStorage float32 = 0
        // store here the concatenation of service names
        listOfServices := make([]string,0)
        // compute the score for this combination of services
        for _, r := range s {
            listOfServices = append(listOfServices,r.GroupServiceInstanceId)
            totalCPU = totalCPU + (float32(r.Cpu) * float32(r.Replicas))
            totalMem = totalMem + (float32(r.Memory) * float32(r.Replicas))
            totalStorage = totalStorage + (float32(r.Storage) * float32(r.Replicas))

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
        scoreForGroup := &pbConductor.DeploymentScore{
            Score:                 score,
            AppInstanceId:         s[0].AppInstanceId,
            GroupServiceInstances: listOfServices,
        }

        foundScores = append(foundScores, scoreForGroup)
    }

    response := &pbConductor.ClusterScoreResponse{
        ClusterId: os.Getenv(utils.MUSICIAN_CLUSTER_ID),
        RequestId: request.RequestId,
        Score: foundScores,
    }


    log.Debug().Interface("score request",request).Interface("score", foundScores).Msg("returned scores")

    // TODO recover cluster id from a cluster environment variable
    return response, nil
}

// Local function to return all the combinations of requirements to check.
// params:
//  reqs array of requirements
// return:
//  array of arrays with all the permutations.
//  E.G.: [A, B, C] -> [[A], [B], [C], [A, B], [A, C], [B, C], [A, B, C]]
func (s *SimpleScorer) getCombinations(reqs []*pbConductor.Requirement) [][]*pbConductor.Requirement {
    length := uint(len(reqs))

    subsets := make([][]*pbConductor.Requirement,0)

    // Go through all possible combinations of objects
    // from 1 (only first object in subset) to 2^length (all objects in subset)
    for subsetBits := 1; subsetBits < (1 << length); subsetBits++ {
        var subset []*pbConductor.Requirement

        for object := uint(0); object < length; object++ {
            // checks if object is contained in subset
            // by checking if bit 'object' is set in subsetBits
            if (subsetBits>>object)&1 == 1 {
                // add object to subset
                subset = append(subset, reqs[object])
            }
        }
        // add subset to subsets
        subsets = append(subsets, subset)
    }
    return subsets
}
