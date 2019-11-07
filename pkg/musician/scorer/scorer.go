/*
 * Copyright 2019 Nalej
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
 *
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
