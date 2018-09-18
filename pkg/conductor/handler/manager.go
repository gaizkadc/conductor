/*
 * Copyright 2018 Nalej
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
 */

package handler

import (
    "github.com/phf/go-queue/queue"
    "github.com/rs/zerolog/log"
    pbConductor "github.com/nalej/grpc-conductor-go"
    "github.com/nalej/conductor/pkg/conductor/scorer"
    "github.com/nalej/conductor/internal/entities"
)

type Manager struct {
    // Queue for incoming messages
    Queue *queue.Queue
    // ScorerMethod
    ScorerMethod scorer.Scorer
}

func NewManager(queue *queue.Queue, scorer scorer.Scorer, port uint32) *Manager {
    // instantiate a server
    return &Manager{queue, scorer}
}



func(c *Manager) ProcessDeploymentRequest(request *pbConductor.DeploymentRequest) (*pbConductor.DeploymentResponse, error) {
    log.Debug().Msgf("manager queue [%p] contains: -->%s<--\n", &c.Queue,c.Queue)
    // Empty queue process it.
    if c.Queue.Len()==0{
        log.Debug().Str("request_id",request.RequestId).Msg("empty queue process request")

        req:= entities.Requirements{RequestID: request.RequestId,
                                    Disk: request.Disk, CPU: request.Cpu, Memory: request.Memory}

        returned,_ := c.ScorerMethod.ScoreRequirements (&req)
        log.Debug().Msgf("Returned %v",returned)
    } else {
        log.Debug().Str("request_id", request.RequestId).Msg("deployment request send to the queue")
        c.Queue.PushBack(request)
    }
    log.Debug().Msgf("manager queue contains after leaving: -->%s<--\n", c.Queue.String())
    response := pbConductor.DeploymentResponse{RequestId: "this is a response"}
    return &response, nil
}








