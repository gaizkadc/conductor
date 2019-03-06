/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package plandesigner

import (
    "github.com/google/uuid"
    "github.com/nalej/conductor/internal/structures"
    pbApplication "github.com/nalej/grpc-application-go"
    pbAuthx "github.com/nalej/grpc-authx-go"
    pbDevice "github.com/nalej/grpc-device-go"
    pbOrganization "github.com/nalej/grpc-organization-go"
    "github.com/nalej/derrors"
    "github.com/nalej/conductor/internal/entities"
    "github.com/nalej/conductor/pkg/utils"
    "context"
    "github.com/rs/zerolog/log"
)

/*
 * This plan designer designs a plan for scenarios where one or several services are indicated to run using
 * replication in all the available clusters of the organization.
 */

type SimpleReplicaPlanDesigner struct {
    // Applications client
    appClient pbApplication.ApplicationsClient
    // Organizations client
    orgClient pbOrganization.OrganizationsClient
    // Connections helper
    connHelper *utils.ConnectionsHelper
    // Authx client
    authxClient pbAuthx.AuthxClient
}

func NewSimpleReplicaPlanDesigner (connHelper *utils.ConnectionsHelper) PlanDesigner {
    connectionsSM := connHelper.GetSystemModelClients()
    appClient := pbApplication.NewApplicationsClient(connectionsSM.GetConnections()[0])
    orgClient := pbOrganization.NewOrganizationsClient(connectionsSM.GetConnections()[0])
    connectionsAuthx := connHelper.GetAuthxClients()
    authxClient := pbAuthx.NewAuthxClient(connectionsAuthx.GetConnections()[0])
    return &SimpleReplicaPlanDesigner{appClient: appClient, orgClient: orgClient, connHelper: connHelper, authxClient:authxClient}
}


func(p *SimpleReplicaPlanDesigner) DesignPlan(app entities.AppInstance,
score entities.DeploymentScore, request entities.DeploymentRequest) (*entities.DeploymentPlan, error) {

    // Build deployment stages for the application
    retrievedDesc,err :=p.appClient.GetAppDescriptor(context.Background(),
        &pbApplication.AppDescriptorId{OrganizationId: app.OrganizationId, AppDescriptorId: app.AppDescriptorId})
    if err!=nil{
        theErr := derrors.NewGenericError("error recovering application instance", err)
        log.Error().Err(theErr).Msg("error recovering application instance")
        return nil, theErr
    }

    // get organization name
    org, err := p.orgClient.GetOrganization(context.Background(),
        &pbOrganization.OrganizationId{OrganizationId: app.OrganizationId})
    if err != nil {
        theErr := derrors.NewGenericError("error when retrieving organization data", err)
        log.Error().Err(err).Msgf("error when retrieving organization data")
        return nil, theErr
    }

    // Get a local representation of the object
    toDeploy := entities.NewAppDescriptorFromGRPC(retrievedDesc)


    planId := uuid.New().String()
    log.Info().Str("planId",planId).Msg("start building the plan")

    // There must be one fragment per service group
    // Each service group with a set of stages following the stages defined in the DeploymentPlan
    // Store the group name and the corresponding deployment order.
    log.Debug().Str("appDescriptor",toDeploy.AppDescriptorId).Msg("analyze group internal dependencies")
    groupsOrder := make(map[string][][]entities.Service)
    for _, g := range toDeploy.Groups {
        log.Debug().Str("appDescriptor",toDeploy.AppDescriptorId).Str("serviceGroupId",g.ServiceGroupId).
            Msg("compute dependency graph for service group")
        dependencyGraph := NewDependencyGraph(g.Services)
        order, err := dependencyGraph.GetDependencyOrderByGroups()
        if err != nil {
            return nil, err
        }
        groupsOrder[g.Name] = order
    }

    // Build deployment matrix
    deploymentMatrix := structures.NewDeploymentMatrix(score)

    // Compute the list of groups to be deployed per cluster
    clustersMap, groupReplicas, err := p.findTargetClusters(toDeploy, deploymentMatrix)
    if err != nil {
        return nil, err
    }

    // Instantiate the number of replicas we need for every group
    groupInstances := make(map[string][]entities.ServiceGroupInstance)
    for serviceGroupId, numReplicas := range groupReplicas {
        log.Debug().Str("serviceGroupId", serviceGroupId).Int("numReplicas", numReplicas).Msg("instantiate service groups")
        instancesReq := pbApplication.AddServiceGroupInstancesRequest{
            AppDescriptorId: app.AppDescriptorId,
            AppInstanceId: app.AppInstanceId,
            OrganizationId: app.OrganizationId,
            ServiceGroupId: serviceGroupId,
            NumInstances: int32(numReplicas),
        }
        serviceInstances, err := p.appClient.AddServiceGroupInstances(context.Background(), &instancesReq)
        if err != nil {
            log.Error().Err(err).Msg("it was impossible to instantiate the service")
            return nil, err
        }
        createdInstances := make([]entities.ServiceGroupInstance,len(serviceInstances.ServiceGroupInstances))
        for i, theInstance := range serviceInstances.ServiceGroupInstances {
            createdInstances[i] = entities.NewServiceGroupInstanceFromGRPC(theInstance)
        }
        groupInstances[serviceInstances.ServiceGroupInstances[0].Name] = createdInstances
    }

    fragments, err := p.buildFragmentsPerCluster(toDeploy,clustersMap, app, groupsOrder, groupInstances, planId, org)

    if err != nil {
        log.Error().Err(err).Msg("impossible to build deployment fragments")
        return nil, err
    }

    // Fill variables
    finalFragments := p.fillVariables(fragments)

    // Now that we have all the fragments, build the deployment plan
    newPlan := entities.DeploymentPlan{
        AppInstanceId: app.AppInstanceId,
        DeploymentId: planId,
        OrganizationId: app.OrganizationId,
        Fragments: finalFragments,
    }

    log.Info().Str("appDescriptorId",app.AppDescriptorId).Str("planId",newPlan.DeploymentId).
        Int("number of fragments",len(newPlan.Fragments)).
        Interface("plan",newPlan).
        Msg("a plan was generated")

    return &newPlan, nil
}


// Build the fragments to be sent to every cluster
func (p* SimpleReplicaPlanDesigner) buildFragmentsPerCluster(
    desc entities.AppDescriptor,
    clustersMap map[string][]entities.ServiceGroup,
    app entities.AppInstance,
    groupsOrder map[string][][]entities.Service,
    groupInstances map[string][]entities.ServiceGroupInstance,
    planId string,
    org *pbOrganization.Organization) ([]entities.DeploymentFragment, derrors.Error) {

        toReturn := make([]entities.DeploymentFragment, 0)
    // combine all the groups per cluster into the corresponding fragment
    for cluster, listGroups := range clustersMap {
        // collect stages per group and generate one fragment
        // UUID for this fragment
        fragmentUUID := uuid.New().String()
        stages := make([]entities.DeploymentStage, 0)

        log.Debug().Str("cluster", cluster).Int("numGroupsToDeploy", len(listGroups)).Msg("design plan for cluster")
        for _, g := range listGroups {

            // take one instance from the available list
            availableGroupInstances := groupInstances[g.Name]
            localGroupInstance := availableGroupInstances[0]
            // remove one entry
            groupInstances[g.Name] = availableGroupInstances[1:]

            // create the stages corresponding to this group
            log.Debug().Str("appDescriptor", app.AppDescriptorId).Str("groupName",g.Name).
                Str("cluster",cluster).Interface("sequences", groupsOrder).
                Interface("groupInstanceToDeploy",localGroupInstance).Msg("create stages for deployment sequence")

            for _, sequence := range groupsOrder[g.Name] {
                // this stage must deploy the services following this order
                stage, err := p.buildDeploymentStage(desc,fragmentUUID, localGroupInstance, sequence)
                if err != nil {
                    log.Error().Err(err).Str("fragmentId",fragmentUUID).Msg("impossible to build stage")
                    return nil, derrors.NewGenericError("impossible to build stage", err)
                }
                stages = append(stages, *stage)
            }
        }
        // one fragment per group
        fragment := entities.DeploymentFragment{
            ClusterId: cluster,
            OrganizationId: org.OrganizationId,
            AppInstanceId: app.AppInstanceId,
            AppDescriptorId: app.AppDescriptorId,
            // To be filled in global instances
            //NalejVariables: ,
            FragmentId: fragmentUUID,
            Stages: stages,
            AppName: app.Name,
            DeploymentId: planId,
            OrganizationName: org.Name,
        }
        toReturn = append(toReturn, fragment)
    }
    return toReturn, nil
}


// Return a map with the list of groups to be deployed per cluster.
// return:
//  map with the list of groups to be deployed per cluster clusterId -> [group0, group1...]
//  map with the number of replicas per group id
//  error if any
func (p* SimpleReplicaPlanDesigner) findTargetClusters(
    desc entities.AppDescriptor,
    deploymentMatrix *structures.DeploymentMatrix) (map[string][]entities.ServiceGroup, map[string]int, derrors.Error) {

    resultClusters := make(map[string][]entities.ServiceGroup,0)
    resultReplicas := make(map[string]int,0)

    for _, g := range desc.Groups {
        targets, err := deploymentMatrix.FindBestTargetsForReplication(g)
        if err != nil {
            log.Error().Err(err).Msg("impossible to find best targets for replication")
            return nil, nil, err
        }
        // Add the number of replicas we need for this group
        resultReplicas[g.ServiceGroupId] = len(targets)

        // Add the group per cluster
        for _, t := range targets {
            current, found := resultClusters[t]
            if !found {
                resultClusters[t] = []entities.ServiceGroup{g}
            } else {
                resultClusters[t] = append(current, g)
            }
        }
    }
    return resultClusters, resultReplicas, nil
}


func (p*SimpleReplicaPlanDesigner) existsService (sequence []entities.Service, serviceName string ) *entities.Service {

    for i:=0;i<len(sequence);i++ {
        if sequence[i].Name == serviceName {
            return &sequence[i]
        }
    }
    return nil
}

// For a given sequence of services, it generates the corresponding deployment stage. This includes the
// instantiation of new services in a service group instance.
func(p *SimpleReplicaPlanDesigner) buildDeploymentStage(desc entities.AppDescriptor, fragmentUUID string, group entities.ServiceGroupInstance,
    sequence []entities.Service) (*entities.DeploymentStage, error) {

    serviceNames := make(map[string]*pbApplication.ServiceInstance, 0)    // variable to store the service names
    // fill a map with the available service instances
    for _, serv := range group.ServiceInstances {
        serviceNames[serv.Name] = serv.ToGRPC()
    }

    deviceSecurityRules := make ([]entities.DeviceGroupSecurityRuleInstance, 0)
    publicSecurityRules := make ([]entities.PublicSecurityRuleInstance, 0)

    // NP-852. Send the new information regarding security rule instances (DeviceGroupRules and PublicRules)
    for _, rule := range desc.Rules {
        service, exists := serviceNames[rule.TargetServiceName]
        if exists {
            if rule.Access == entities.Public {
                publicSecurityRules = append(publicSecurityRules, *entities.NewPublicSercurityRuleInstance(*service, rule))
            } else if rule.Access == entities.DeviceGroup{
                sgJwtSecrets := make ([]string, 0)
                for _, sg := range rule.DeviceGroupIds {
                    secret, err := p.authxClient.GetDeviceGroupSecret(context.Background(), &pbDevice.DeviceGroupId{
                        OrganizationId: rule.OrganizationId,
                        DeviceGroupId:  sg,
                    })
                    if err != nil {
                        log.Error().Err(err).Str("organization_id", rule.OrganizationId).
                            Str("device_group_id", sg).Msg("error getting the Jwt secret")
                    }else{
                        sgJwtSecrets = append(sgJwtSecrets, secret.Secret)
                    }
                }
                deviceSecurityRules = append(deviceSecurityRules, *entities.NewDeviceGroupSecurityRuleInstance(*service, rule, sgJwtSecrets))
            }
        }
    }

    ds := entities.DeploymentStage{
        StageId: uuid.New().String(),
        FragmentId: fragmentUUID,
        Services: group.ServiceInstances,
        PublicRules:publicSecurityRules,
        DeviceGroupRules:deviceSecurityRules,
    }
    return &ds,nil
}



// Fill the fragments with the corresponding variables per group. This has to be done after the generation of the fragments
// to correctly fill the entries with the corresponding values.
func (p *SimpleReplicaPlanDesigner) fillVariables(fragmentsToDeploy []entities.DeploymentFragment) []entities.DeploymentFragment {
    toChange := make(map[int]map[string]string,0)
    allServices := make(map[string]string,0)
    for fragmentIndex, f := range fragmentsToDeploy {
        // Create the service entries we need for this fragment
        variables := make(map[string]string,0)
        for _, stage := range f.Stages {
            for _, serv := range stage.Services {
                key, value := GetDeploymentVariableForService(serv)
                variables[key] = value
                // TODO: this might be modified depending on decissions about how to declare variables
                allServices[key] = value
            }
        }
        toChange[fragmentIndex] = variables
    }

    for fragmentIndex, variables := range toChange {
        fragmentsToDeploy[fragmentIndex].NalejVariables = variables
        // This fragment must know all the variables for the sake of completeness
        for key, variable := range allServices {
            _, found := fragmentsToDeploy[fragmentIndex].NalejVariables[key]
            if !found {
                // we couldn't find it, add the first entry we know
                fragmentsToDeploy[fragmentIndex].NalejVariables[key] = variable
            }
        }
    }

    return fragmentsToDeploy
}
