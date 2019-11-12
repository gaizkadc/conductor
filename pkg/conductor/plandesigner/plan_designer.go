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
	// Suffix to add to the outbound FQDNs
	OutboundSuffix = "-OUT-"
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
	value := fmt.Sprintf("%s.%s", utils.GetVSAName(serviceName, organizationId, appInstanceId), NalejServiceSuffix)
	return key, value
}

func GetDeploymentVariableForOutbound(serviceName string, outboundName string, appInstanceId string, organizationId string) (string, string) {
	if outboundName == "" {
		return "", ""
	}
	key := fmt.Sprintf(NalejVariableOutboundPrefix, strings.ToUpper(outboundName))
	vsaNameCompound := utils.GetVSAName(serviceName, organizationId, appInstanceId) + OutboundSuffix + outboundName
	value := fmt.Sprintf("%s.%s", vsaNameCompound, NalejServiceSuffix)
	return key, value
}

// Format a string removing white spaces and going lowercase
func formatName(name string) string {
	aux := strings.ToLower(name)
	// replace any space
	aux = strings.Replace(aux, " ", "", -1)
	return aux
}
