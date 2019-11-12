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

package queue

import (
	"context"
	"github.com/nalej/conductor/pkg/conductor/baton"
	"github.com/nalej/nalej-bus/pkg/queue/infrastructure/ops"
	"github.com/rs/zerolog/log"
	"time"
)

// Timeout between incoming messages
const InfrastructureOpsTimeout = time.Minute * 60

// Control incoming requests for the application ops topic

type InfrastructureOpsHandler struct {
	// reference baton
	baton *baton.Manager
	// consumer for this queue
	cons *ops.InfrastructureOpsConsumer
}

// Instantiate a new object to consume and process entries from the infrastructure ops queue
// params:
//  baton to decide how to proceed
//  consumer to get entries from the queue
// return:
//  instance of an infrastructure ops queue
func NewInfrastructureOpsHandler(baton *baton.Manager, cons *ops.InfrastructureOpsConsumer) InfrastructureOpsHandler {
	return InfrastructureOpsHandler{baton: baton, cons: cons}
}

// This operations runs a set of subroutines feeding the corresponding channels for this handler.
func (h InfrastructureOpsHandler) Run() {
	go h.consumeDrainRequest()
	go h.waitRequests()
}

// Endless loop waiting for requests
func (h InfrastructureOpsHandler) waitRequests() {
	log.Debug().Msg("wait for requests to be received by the infrastructure ops queue")
	for {
		ctx, cancel := context.WithTimeout(context.Background(), InfrastructureOpsTimeout)
		// in every iteration this loop consumes data and sends it to the corresponding channels
		currentTime := time.Now()
		err := h.cons.Consume(ctx)
		cancel()
		select {
		case <-ctx.Done():
			// the timeout was reached
			log.Debug().Msgf("no message received since %s", currentTime.Format(time.RFC3339))
		default:
			// we received something or an error
			if err != nil {
				log.Error().Err(err).Msg("error consuming data from application ops")
			}
		}
	}
}

func (h InfrastructureOpsHandler) consumeDrainRequest() {
	log.Debug().Msg("waiting for drain requests...")
	for {
		received := <-h.cons.Config.ChDrainRequest
		log.Debug().Interface("drain", received).Msg("<- incoming drain request")
		h.baton.DrainCluster(received)
	}
}
