/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 */

// Main components to be fulfilled by any scorer.


package scorer

import (
    pbConductor "github.com/nalej/grpc-conductor-go"
)

type Scorer interface {

    // For a given score request return a scoring response.
    //  params:
    //   request to be processed.
    //  return:
    //   score response or error if any.
    Score(request *pbConductor.ClusterScoreRequest) (*pbConductor.ClusterScoreResponse, error)
}
