/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package queue

import (
    "github.com/nalej/conductor/pkg/conductor/baton"
    "github.com/nalej/nalej-bus/pkg/queue/application/ops"
    "github.com/nalej/conductor/internal/entities"
    "github.com/rs/zerolog/log"
)

// Control incoming requests for the application ops topic

type ApplicationOpsHandler struct {
    // reference baton
    baton *baton.Manager
    // configuration for this queue
    conf *ops.ConfigApplicationOpsConsumer

}

// Instantiate a new object to consume and process entries from the application ops queue
// params:
//  baton to decide how to proceed
//  consumer to get entries from the queue
// return:
//  instance of an application ops queue
func NewApplicationOpsHandler(baton *baton.Manager, conf *ops.ConfigApplicationOpsConsumer) ApplicationOpsHandler {
    return ApplicationOpsHandler{baton: baton, conf: conf}
}


// This operations runs a set of subroutines feeding the corresponding channels for this handler.
func(h ApplicationOpsHandler) Run() {
    go h.consumeDeploymentRequest()
    go h.consumeUndeployRequest()
}

func(h ApplicationOpsHandler) consumeDeploymentRequest () {
    log.Debug().Msg("waiting for deployment requests...")
    for {
        received := <- h.conf.ChDeploymentRequest
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
        received := <- h.conf.ChUndeployRequest
        log.Debug().Interface("undeployRequest", received).Msg("<- incoming undeploy request")
        aux := entities.UndeployRequest{OrganizationId: received.OrganizationId, AppInstanceId: received.AppInstanceId}
        err := h.baton.Undeploy(&aux)
        if err != nil {
            log.Error().Err(err).Msg("failed processing undeploy request")
        }
    }
}