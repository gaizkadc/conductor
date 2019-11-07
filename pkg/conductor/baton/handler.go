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

// Service in charge of processing deployment gRPC requests.

package baton

import (
	"context"

	"github.com/nalej/conductor/internal/entities"
	pbCommon "github.com/nalej/grpc-common-go"
	pbConductor "github.com/nalej/grpc-conductor-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
)

//Handler struct with a related Manager
type Handler struct {
	c *Manager
}

//NewHandler creates a new Handler with its given Manager
func NewHandler(c *Manager) *Handler {
	return &Handler{c}
}

//Deploy validates and queues a deployment request
func (h *Handler) Deploy(ctx context.Context, request *pbConductor.DeploymentRequest) (*pbCommon.Success, error) {
	log.Debug().Interface("deploymentRequest", request).Msg("Deploy")
	if err := entities.ValidDeploymentRequest(request); err != nil {
		return nil, conversions.ToGRPCError(err)
	}

	// Enqueue request for later processing
	log.Debug().Msgf("enqueue request %s", request.RequestId)
	err := h.c.PushRequest(request)
	if err != nil {
		return nil, err
	}

	return &pbCommon.Success{}, nil
}

func (h *Handler) Undeploy(ctx context.Context, request *pbConductor.UndeployRequest) (*pbCommon.Success, error) {
	log.Debug().Interface("undeployRequest", request).Msg("Undeploy")
	if err := entities.ValidUndeployRequest(request); err != nil {
		return nil, conversions.ToGRPCError(err)
	}

	toUndeploy := entities.UndeployRequest{
		OrganizationId: request.OrganizationId,
		AppInstanceId:  request.AppInstanceId,
	}
	err := h.c.Undeploy(&toUndeploy)
	if err != nil {
		log.Error().Msgf("Unable to undeploy application %s", request.AppInstanceId)
		return nil, err
	}
	log.Debug().Msgf("Application %s undeployed", request.AppInstanceId)
	return &pbCommon.Success{}, nil
}

func (h *Handler) DrainCluster(ctx context.Context, request *pbConductor.DrainClusterRequest) (*pbCommon.Success, error) {
	log.Info().Msg("drain cluster was invoked through the GRPC api but it is not implemented!!!!")
	return &pbCommon.Success{}, nil
}
