/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package plandesigner

import (
    "github.com/nalej/conductor/internal/entities"
    "fmt"
    "github.com/nalej/derrors"
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
    // Nalej service suffix
    NalejServiceSuffix = "service.nalej"
    // Timeout for GRPC operations
    PlanDesignerGRPCTimeout = 5 * time.Second
)


// Generate the tuple key and value for a nalej service to be represented.
// params:
//  serv service instance to be processed
// return:
//  variable name, variable value
func GetDeploymentVariableForService(serv entities.ServiceInstance) (string, string) {
    key := fmt.Sprintf(NalejVariablePrefix,strings.ToUpper(serv.Name))
    value := fmt.Sprintf("%s-%s-%s.%s", formatName(serv.Name), serv.OrganizationId[0:10],
        serv.AppInstanceId[0:10], NalejServiceSuffix)
    return key,value
}


// Format a string removing white spaces and going lowercase
func formatName(name string) string {
    aux := strings.ToLower(name)
    // replace any space
    aux = strings.Replace(aux, " ", "", -1)
    return aux
}