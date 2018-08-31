//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Cluster cli operations.

package main

import (
    "strconv"

    "gopkg.in/alecthomas/kingpin.v2"

    "github.com/daishogroup/system-model/client"
    "github.com/daishogroup/system-model/entities"
)

type ClusterCommand struct {
    GlobalCommand

    addClusterNetworkID   *string
    addClusterName        *string
    addClusterDescription *string
    addClusterType        *string
    addClusterLocation    *string
    addAdminMail          *string

    updateNetworkID   *string
    updateClusterID   *string
    updateName        *string
    updateDescription *string
    updateType        *string
    updateLocation    *string
    updateEmail       *string
    updateStatus      *string
    updateDrain       *string
    updateCordon      *string

    listClustersNetworkID *string

    getClusterNetworkID *string
    getClusterClusterID *string

    deleteClusterNetworkID *string
    deleteClusterClusterID *string
}

func NewClusterCommand(app *kingpin.Application, global GlobalCommand) *ClusterCommand {

    c := &ClusterCommand{
        GlobalCommand: global,
    }

    clusterCmd := app.Command("cluster", "Cluster commands")

    cmd := clusterCmd.Command("add", "Add a new cluster").Action(c.addCluster)
    c.addClusterNetworkID = cmd.Arg(NetworkID, NetworkIDDesc).Required().String()
    c.addClusterName = cmd.Arg("name", "Name of the cluster to add").Required().String()
    c.addClusterDescription = cmd.Flag("desc", "Description of the cluster to add").Default("").String()
    c.addClusterType = cmd.Arg("type", "Type of cluster").Required().Enum(
        string(entities.CloudType),
        string(entities.EdgeType),
        string(entities.GatewayType))
    c.addClusterLocation = cmd.Flag("location", "Cluster location").Default("").String()
    c.addAdminMail = cmd.Flag("email", "Cluster Admin Mail").Default("").String()

    cmd = clusterCmd.Command("update", "Update a selected cluster").Action(c.updateCluster)
    c.updateNetworkID = cmd.Arg(NetworkID, NetworkIDDesc).Required().String()
    c.updateClusterID = cmd.Arg(ClusterID, ClusterIDDesc).Required().String()
    c.updateName = cmd.Flag("name", "Name of the cluster to update").Default(NotAssigned).String()
    c.updateDescription = cmd.Flag("desc", "Description of the cluster to update").Default(NotAssigned).String()
    c.updateType = cmd.Flag("type", "Type of the cluster to update").Enum(
        string(entities.GatewayType),
        string(entities.CloudType),
        string(entities.EdgeType),
    )
    c.updateLocation = cmd.Flag("location", "Location of the cluster to update").Default(NotAssigned).String()
    c.updateEmail = cmd.Flag("email", "Email of the cluster to update").Default(NotAssigned).String()
    c.updateStatus = cmd.Flag("status", "Status of the cluster to update").Enum(
        string(entities.ClusterCreated),
        string(entities.ClusterReadyToInstall),
        string(entities.ClusterInstalling),
        string(entities.ClusterInstalled),
        string(entities.ClusterUninstalling),
        string(entities.ClusterError),
    )
    c.updateDrain = cmd.Flag("drain", "Flag drain of the cluster to update").Default(NotAssigned).String()
    c.updateCordon = cmd.Flag("cordon", "Flag cordon of the cluster to update").Default(NotAssigned).String()

    cmd = clusterCmd.Command("list", "List the clusters in a network").Action(c.listClusters)
    c.listClustersNetworkID = cmd.Arg(NetworkID, NetworkIDDesc).Required().String()

    cmd = clusterCmd.Command("get", "Get a cluster").Action(c.getCluster)
    c.getClusterNetworkID = cmd.Arg(NetworkID, NetworkIDDesc).Required().String()
    c.getClusterClusterID = cmd.Arg(ClusterID, ClusterIDDesc).Required().String()

    cmd = clusterCmd.Command("delete", "Delete a cluster").Action(c.deleteCluster)
    c.deleteClusterNetworkID = cmd.Arg(NetworkID, NetworkIDDesc).Required().String()
    c.deleteClusterClusterID = cmd.Arg(ClusterID, ClusterIDDesc).Required().String()
    return c
}

func (cmd *ClusterCommand) addCluster(c *kingpin.ParseContext) error {
    clusterRest := client.NewClusterClientRest(cmd.IP.String(), *cmd.Port)
    toAdd := entities.NewAddClusterRequest(* cmd.addClusterName, * cmd.addClusterDescription,
        (entities.ClusterType)(*cmd.addClusterType), * cmd.addClusterLocation, * cmd.addAdminMail)
    added, err := clusterRest.Add(*cmd.addClusterNetworkID, * toAdd)
    return cmd.printResultOrError(added, err)
}

func (cmd *ClusterCommand) listClusters(c *kingpin.ParseContext) error {
    clusterRest := client.NewClusterClientRest(cmd.IP.String(), *cmd.Port)
    clusters, err := clusterRest.ListByNetwork(* cmd.listClustersNetworkID)
    return cmd.printResultOrError(clusters, err)
}

func (cmd *ClusterCommand) getCluster(c *kingpin.ParseContext) error {
    clusterRest := client.NewClusterClientRest(cmd.IP.String(), *cmd.Port)
    cluster, err := clusterRest.Get(* cmd.getClusterNetworkID, * cmd.getClusterClusterID)
    return cmd.printResultOrError(cluster, err)
}

func (cmd *ClusterCommand) deleteCluster(c *kingpin.ParseContext) error {
    clusterRest := client.NewClusterClientRest(cmd.IP.String(), *cmd.Port)
    err := clusterRest.Delete(*cmd.deleteClusterNetworkID, *cmd.deleteClusterClusterID)
    return cmd.printResultOrError(entities.NewSuccessfulOperation("Cluster deleted"), err)
}

func (cmd *ClusterCommand) updateCluster(c *kingpin.ParseContext) error {
    clusterRest := client.NewClusterClientRest(cmd.IP.String(), *cmd.Port)
    update := entities.NewUpdateClusterRequest()
    if cmd.updateName != nil && * cmd.updateName != NotAssigned {
        update.WithName(* cmd.updateName)
    }
    if cmd.updateDescription != nil && * cmd.updateDescription != NotAssigned {
        update.WithDescription(* cmd.updateDescription)
    }

    if cmd.updateType != nil && * cmd.updateType != "" {
        update.WithType((entities.ClusterType)(* cmd.updateType))
    }
    if cmd.updateLocation != nil && * cmd.updateLocation != NotAssigned {
        update.WithLocation(* cmd.updateLocation)
    }
    if cmd.updateEmail != nil && * cmd.updateEmail != NotAssigned {
        update.WithEmail(* cmd.updateEmail)
    }
    if cmd.updateStatus != nil && * cmd.updateStatus != "" {
        update.WithClusterStatus((entities.ClusterStatus)(* cmd.updateStatus))
    }
    if cmd.updateDrain != nil && * cmd.updateDrain != NotAssigned {
        b, err := strconv.ParseBool(* cmd.updateDrain)
        if err != nil {
            kingpin.Fatalf("invalid boolean value")
        }
        update.WithDrain(b)
    }
    if cmd.updateCordon != nil && * cmd.updateCordon != NotAssigned {
        b, err := strconv.ParseBool(* cmd.updateCordon)
        if err != nil {
            kingpin.Fatalf("invalid boolean value")
        }
        update.WithCordon(b)
    }

    updated, err := clusterRest.Update(* cmd.updateNetworkID, * cmd.updateClusterID, * update)
    return cmd.printResultOrError(updated, err)
}
