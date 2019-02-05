/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package plandesigner

import (
    "github.com/nalej/conductor/internal/entities"
    "fmt"
    "strings"
)

// Basic interface to be follow by any plan designer.

type PlanDesigner interface {

    // For any set of requirements an a given score elaborate a deployment plan.
    //  params:
    //   app application instance
    //   services group of services to deploy
    //   score obtained by musicians
    //   request deployment request for this plan
    //  return:
    //   A collection of deployment plans each one designed to run in a different cluster.
    DesignPlan(app *entities.AppInstance,
        score *entities.DeploymentScore, request *entities.DeploymentRequest) (*entities.DeploymentPlan, error)
}

const (
    // Definition of the prefix to be used when defining any Nalej service variable
    NalejVariablePrefix = "NALEJ_SERV_%s"
    // Nalej service suffix
    NalejServiceSuffix = "service.nalej"
)

// Generate the set of Nalej variables for a deployment.
// params:
//  organizationName    name of the organization
//  appInstanceId       application instance
//  desc                deployment descriptor
// return:
//  map with variables and values
func GetDeploymentNalejVariables(organizationName string, appInstanceId string, desc entities.AppDescriptor) map[string]string{
    variables := make(map[string]string,0)
    for _, g := range desc.Groups {
        for _,s := range g.Services {
            value := fmt.Sprintf("%s-%s-%s.%s", formatName(s.Name), formatName(organizationName), appInstanceId[0:5],
                NalejServiceSuffix)
            name := fmt.Sprintf(NalejVariablePrefix,strings.ToUpper(s.ServiceId))
            variables[name]=value
        }
    }
    return variables
}


// Format a string removing white spaces and going lowercase
func formatName(name string) string {
    aux := strings.ToLower(name)
    // replace any space
    aux = strings.Replace(aux, " ", "", -1)
    return aux
}