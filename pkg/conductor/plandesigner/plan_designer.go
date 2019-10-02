/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package plandesigner

import (
	"fmt"
	"github.com/nalej/conductor/internal/entities"
	"github.com/nalej/conductor/pkg/utils"
	"strings"
	"time"
)

// Basic interface to be follow by any plan designer.

type PlanDesigner interface {
	// For any set of requirements an a given score elaborate a deployment plan. If a list of group ids is set
	// only these groups are chosen otherwise, all the groups from the descriptor are used.
	//  params:
	//   app application instance
	//   score obtained by musicians
	//   request deployment request for this plan
	//   groupIds list of group ids to select from the descriptor
	//  return:
	//   A collection of deployment plans each one designed to run in a different cluster.
	DesignPlan(app entities.AppInstance,
		score entities.DeploymentScore, request entities.DeploymentRequest, groupIds []string,
		deployedGroups map[string][]string) (*entities.DeploymentPlan, error)
}

const (
	// Definition of the prefix to be used when defining any Nalej service variable
	NalejVariablePrefix = "NALEJ_SERV_%s"
	// Definition of the prefix to be used when defining any Nalej outbound variable
	NalejVariableOutboundPrefix = "NALEJ_OUTBOUND_%s"
	// Nalej service suffix
	NalejServiceSuffix = "service.nalej"
	// Nalej service suffix
	NalejOutboundSuffix = "outbound.nalej"
	// Timeout for GRPC operations
	PlanDesignerGRPCTimeout = 5 * time.Second
)

// Generate the tuple key and value for a nalej service to be represented.
// params:
//  serviceName
//  appInstanceId
//  organizationId
// return:
//  variable name, variable value
func GetDeploymentVariableForService(serviceName string, appInstanceId string, organizationId string) (string, string) {
	key := fmt.Sprintf(NalejVariablePrefix, strings.ToUpper(serviceName))
	value := fmt.Sprintf("%s-%s-%s.%s", formatName(serviceName), organizationId[0:10],
		appInstanceId[0:10], NalejServiceSuffix)
	return key, value
}

func GetDeploymentVariableForOutbound(serviceName string, appInstanceId string, organizationId string) (string, string) {
	key := fmt.Sprintf(NalejVariableOutboundPrefix, strings.ToUpper(serviceName))
	value := fmt.Sprintf("%s.%s", utils.GetVSAName(serviceName, organizationId, appInstanceId), NalejOutboundSuffix)
	return key, value
}

// Format a string removing white spaces and going lowercase
func formatName(name string) string {
	aux := strings.ToLower(name)
	// replace any space
	aux = strings.Replace(aux, " ", "", -1)
	return aux
}
