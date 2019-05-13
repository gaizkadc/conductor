/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package queue

import (
    "context"
    "github.com/nalej/conductor/pkg/conductor/baton"
    "github.com/nalej/nalej-bus/pkg/queue/application/ops"
    "github.com/nalej/conductor/internal/entities"
    "github.com/rs/zerolog/log"
    "time"
)

// Timeout between incoming messages
const ApplicationOpsTimeout = time.Minute * 60

// Control incoming requests for the application ops topic

type ApplicationOpsHandler struct {
    // reference baton
    baton *baton.Manager
    // consumer for this queue
    cons *ops.ApplicationOpsConsumer

}

// Instantiate a new object to consume and process entries from the application ops queue
// params:
//  baton to decide how to proceed
//  consumer to get entries from the queue
// return:
//  instance of an application ops queue
func NewApplicationOpsHandler(baton *baton.Manager, cons *ops.ApplicationOpsConsumer) ApplicationOpsHandler {
    return ApplicationOpsHandler{baton: baton, cons: cons}
}


// This operations runs a set of subroutines feeding the corresponding channels for this handler.
func(h ApplicationOpsHandler) Run() {
    go h.consumeDeploymentRequest()
    go h.consumeUndeployRequest()
    go h.waitRequests()
}

// Endless loop waiting for requests
func (h ApplicationOpsHandler) waitRequests() {
    log.Debug().Msg("wait for requests to be received by the application ops queue")
    for {
        ctx, cancel := context.WithTimeout(context.Background(), ApplicationOpsTimeout)
        // in every iteration this loop consumes data and sends it to the corresponding channels
        currentTime := time.Now()
        err := h.cons.Consume(ctx)
        cancel()
        select {
        case <- ctx.Done():
            // the timeout was reached
            log.Debug().Msgf("no message received since %s",currentTime.Format(time.RFC3339))
        default:
            // we received something or an error
            if err != nil {
                log.Error().Err(err).Msg("error consuming data from application ops")
            }
        }
    }
}

func(h ApplicationOpsHandler) consumeDeploymentRequest () {
    log.Debug().Msg("waiting for deployment requests...")
    for {
        received := <- h.cons.Config.ChDeploymentRequest
        log.Debug().Interface("deploymentRequest", received).Msg("<- incoming deployment request")
        err := h.baton.PushRequest(received)
        if err != nil {
            log.Error().Err(err).Msg("failed processing deployment request")
        }
    }
}

func(h ApplicationOpsHandler) consumeUndeployRequest () {
    log.Debug().Msg("waiting for undeploy requests...")
    for {
        received := <- h.cons.Config.ChUndeployRequest
        log.Debug().Interface("undeployRequest", received).Msg("<- incoming undeploy request")
        aux := entities.UndeployRequest{OrganizationId: received.OrganizationId, AppInstanceId: received.AppInstanceId}
        err := h.baton.Undeploy(&aux)
        if err != nil {
            log.Error().Err(err).Msg("failed processing undeploy request")
        }
    }
}