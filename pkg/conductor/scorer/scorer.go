/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package scorer

import (
    "github.com/nalej/conductor/internal/entities"
)

// Common interface for deployment scorers.
type Scorer interface {

    // For a existing set of deployment requirements score potential candidates.
    //  params:
    //   organizationId where the request is taking place
    //   requirements to be fulfilled
    //  return:
    //   candidates score
    ScoreRequirements (organizationId string, requirements *entities.Requirements) (*entities.ClustersScore, error)
}
