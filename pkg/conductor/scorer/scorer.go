//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

package scorer

import (
    "github.com/nalej/conductor/internal/entities"
)

// Common interface for deployment scorers.
type Scorer interface {

    // For a existing set of deployment requirements score potential candidates.
    //  params:
    //   requirements to be fulfilled
    //   array of musician addresses to be queried
    //  return:
    //   candidates score
    ScoreRequirements (requirements *entities.Requirements, musicians []string) (*entities.ClusterScore, error)
}
