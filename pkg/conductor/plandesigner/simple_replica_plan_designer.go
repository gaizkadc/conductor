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
    clustersMap, err := p.findTargetClusters(toDeploy, deploymentMatrix)
    if err != nil {
        return nil, err
    }

    fragments, err := p.buildFragmentsPerCluster(toDeploy,clustersMap, app, groupsOrder, planId, org)

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
    planId string,
    org *pbOrganization.Organization) ([]entities.DeploymentFragment, derrors.Error) {

        toReturn := make([]entities.DeploymentFragment, 0)
    // combine all the groups per cluster into the corresponding fragment
    for cluster, listGroups := range clustersMap {
        // collect stages per group and generate one fragment
        // UUID for this fragment
        fragmentUUID := uuid.New().String()
        stages := make([]entities.DeploymentStage, 0)

        for _, g := range listGroups {
            // Add new ServiceGroupInstance
            newServiceGroupRequest := pbApplication.AddServiceGroupInstanceRequest{
                OrganizationId: app.OrganizationId,
                AppDescriptorId: app.AppDescriptorId,
                AppInstanceId: app.AppInstanceId,
                ServiceGroupId: g.ServiceGroupId,
            }
            groupInstance, err := p.appClient.AddServiceGroupInstance(context.Background(),&newServiceGroupRequest)
            if err != nil {
                log.Error().Err(err).Msg("error creating new service group instance")
                return nil, derrors.NewGenericError("impossible to instantiate service group instance", err)
            }
            localGroupInstance := entities.NewServiceGroupInstanceFromGRPC(groupInstance)

            // create the stages corresponding to this group
            log.Debug().Str("appDescriptor", app.AppDescriptorId).Str("groupName",g.Name).
                Interface("sequences", groupsOrder).Msg("create stages for deployment sequence")

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
func (p* SimpleReplicaPlanDesigner) findTargetClusters(
    desc entities.AppDescriptor,
    deploymentMatrix *structures.DeploymentMatrix) (map[string][]entities.ServiceGroup,derrors.Error) {

    result := make(map[string][]entities.ServiceGroup,0)

    for _, g := range desc.Groups {
        targets, err := deploymentMatrix.FindBestTargetsForReplication(g)
        if err != nil {
            log.Error().Err(err).Msg("impossible to find best targets for replication")
            return nil, err
        }
        // Add the group per cluster
        for _, t := range targets {
            current, found := result[t]
            if !found {
                result[t] = []entities.ServiceGroup{g}
            } else {
                result[t] = append(current, g)
            }
        }
    }
    return result, nil
}



// This local function returns a fragment for a given list of services and its dependency graph
func (p *SimpleReplicaPlanDesigner) buildFragments(
    desc entities.AppDescriptor,
    app entities.AppInstance,
    group entities.ServiceGroup,
    groupsOrder [][]entities.Service,
    targetCluster string,
    planId string,
    org *pbOrganization.Organization,
    ) ([]entities.DeploymentFragment, derrors.Error) {

    fragments := make([]entities.DeploymentFragment,0)
    // collect stages per group and generate one fragment
    // UUID for this fragment
    fragmentUUID := uuid.New().String()

    // Add new ServiceGroupInstance
    newServiceGroupRequest := pbApplication.AddServiceGroupInstanceRequest{
        OrganizationId: app.OrganizationId,
        AppDescriptorId: app.AppDescriptorId,
        AppInstanceId: app.AppInstanceId,
        ServiceGroupId: group.ServiceGroupId,
    }
    groupInstance, err := p.appClient.AddServiceGroupInstance(context.Background(),&newServiceGroupRequest)
    if err != nil {
        log.Error().Err(err).Msg("error creating new service group instance")
        return nil, derrors.NewGenericError("impossible to instantiate service group instance", err)
    }
    localGroupInstance := entities.NewServiceGroupInstanceFromGRPC(groupInstance)

    // create the stages corresponding to this group
    log.Debug().Str("appDescriptor", app.AppDescriptorId).Str("groupName",group.Name).
        Interface("sequences", groupsOrder).Msg("create stages for deployment sequence")
    stages := make([]entities.DeploymentStage, 0)
    for _, sequence := range groupsOrder {
        // this stage must deploy the services following this order
        stage, err := p.buildDeploymentStage(desc,fragmentUUID, localGroupInstance, sequence)
        if err != nil {
            log.Error().Err(err).Str("fragmentId",fragmentUUID).Msg("impossible to build stage")
            return nil, derrors.NewGenericError("impossible to build stage", err)
        }

        stages = append(stages, *stage)
    }

    fragment := entities.DeploymentFragment{
        AppDescriptorId: app.AppDescriptorId,
        OrganizationId:         app.OrganizationId,
        OrganizationName:       org.Name,
        AppInstanceId:          app.AppInstanceId,
        AppName:                app.Name,
        FragmentId:             fragmentUUID,
        DeploymentId:           planId,
        Stages:                 stages,
        ClusterId:              targetCluster,
    }
    fragments = append(fragments, fragment)

    return fragments, nil
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
    instances := make([]entities.ServiceInstance,len(sequence))
    for i, s := range sequence {
        // Instantiate this service
        request := pbApplication.AddServiceInstanceRequest{
            OrganizationId: group.OrganizationId,
            AppDescriptorId: group.AppDescriptorId,
            AppInstanceId: group.AppInstanceId,
            ServiceGroupId: group.ServiceGroupId,
            ServiceGroupInstanceId: group.ServiceGroupInstanceId,
            ServiceId: s.ServiceId,
        }
        instance, err := p.appClient.AddServiceInstance(context.Background(), &request)
        if err != nil {
            log.Error().Err(err).Msg("error when adding a new service instance")
            return nil, err
        }
        instances[i] = entities.NewServiceInstanceFromGRPC(instance)
        serviceNames[s.Name] = instance
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
    log.Debug().Int("services", len(instances)).Int("public rules", len(publicSecurityRules)).
        Int("device group rules", len(deviceSecurityRules)).Msg("deployment stage has been defined")
    ds := entities.DeploymentStage{
        StageId: uuid.New().String(),
        FragmentId: fragmentUUID,
        Services: instances,
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
