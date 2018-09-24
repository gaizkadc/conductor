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
    "sync"
    "time"
)

// Time to wait between checks in the queue.
const CheckSleepTime = 2

type Manager struct {
    // queue for incoming messages
    queue *queue.Queue
    // ScorerMethod
    ScorerMethod scorer.Scorer
    // Mutex for queue operations
    mux sync.RWMutex
}

func NewManager(queue *queue.Queue, scorer scorer.Scorer, port uint32) *Manager {
    // instantiate a server
    return &Manager{queue: queue, ScorerMethod: scorer}
}

// Check iteratively if there is anything to be processed in the queue.
func (c *Manager) Run() {
    for {
        for c.AvailableRequests() {
            c.ProcessDeploymentRequest()
        }
        // log.Debug().Msg("no requests to consume")
        time.Sleep(time.Second*CheckSleepTime)
    }
}

func(c *Manager) ProcessDeploymentRequest(){
    req := c.NextRequest()
    if req == nil {
        log.Error().Msg("the queue was unexpectedly empty")
        return
    }

    scoreRequest := entities.Requirements{RequestID: req.RequestId,
        Disk: req.Disk, CPU: req.Cpu, Memory: req.Memory}

    scoreResult, err := c.ScorerMethod.ScoreRequirements (&scoreRequest)

    if err != nil {
        log.Error().Err(err).Msgf("error scoring request %s",scoreRequest.RequestID)
        return
    }

    log.Info().Msgf("conductor maximum score for %s is for cluster %s among %d possible",
        scoreResult.RequestID, scoreResult.ClusterID, scoreResult.TotalEvaluated)

    // TODO elaborate plan, modify system model accordingly
    // Elaborate deployment plan
}



// Thread-safe method to access queued requests
func(c *Manager) NextRequest() *pbConductor.DeploymentRequest {
    c.mux.Lock()
    toReturn := c.queue.PopFront().(*pbConductor.DeploymentRequest)
    defer c.mux.Unlock()
    return toReturn
}

// Thread-safe function to find whether there are more requests available or not.
func(c *Manager) AvailableRequests() bool {
    c.mux.RLock()
    available := c.queue.Len()!=0
    defer c.mux.RUnlock()
    return available
}

// Push a new request to the que for later processing.
//  params:
//   req entry to be enqueued
func (c *Manager) PushRequest(req *pbConductor.DeploymentRequest) error {
    c.mux.Lock()
    c.queue.PushBack(req)
    defer c.mux.Unlock()
    return nil
}







