//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Node command line operations.

package main

import (
    "strings"

    "gopkg.in/alecthomas/kingpin.v2"

    "github.com/daishogroup/system-model/client"
    "github.com/daishogroup/system-model/entities"
    "strconv"
)

type NodeCommand struct {
    GlobalCommand

    addNodeNetworkID * string
    addNodeClusterID * string
    addNodeName      * string
    addNodeDesc      * string
    addNodeLabels * string
    addNodePublicIP  * string
    addNodePrivateIP * string
    addNodeInstalled * bool
    addNodeUsername  * string
    addNodePassword  * string
    addNodeSSHKey    * string

    listNodesNetworkID * string
    listNodesClusterID * string

    getNodeNetworkID * string
    getNodeClusterID * string
    getNodeNodeID    * string

    removeNodeNetworkID * string
    removeNodeClusterID * string
    removeNodeNodeID    * string

    updateNodeNetworkID * string
    updateNodeClusterID * string
    updateNodeNodeID    * string
    updateNodeName      * string
    updateNodeDesc      * string
    updateNodeLabels * string
    updateNodePublicIP  * string
    updateNodePrivateIP * string
    updateNodeInstalled * string
    updateNodeUsername  * string
    updateNodePassword  * string
    updateNodeSSHKey    * string
    updateNodeStatus    * string

}

func NewNodeCommand(app * kingpin.Application, global GlobalCommand) * NodeCommand {

    c := &NodeCommand{
        GlobalCommand : global,
    }

    nodeCommand := app.Command("node", "Node commands")

    cmd := nodeCommand.Command("add", "Add a new node").Action(c.addNode)
    c.addNodeNetworkID = cmd.Arg(NetworkID, NetworkIDDesc).Required().String()
    c.addNodeClusterID = cmd.Arg(ClusterID, ClusterIDDesc).Required().String()
    c.addNodeName = cmd.Arg("name", "Node name").Required().String()
    c.addNodeDesc = cmd.Flag("desc", "Node description").Default("").String()
    c.addNodeLabels = cmd.Flag("labels", "Comma separated list of labels").Default("").String()
    c.addNodePublicIP = cmd.Arg("publicIP", "Public IP address").Required().String()
    c.addNodePrivateIP = cmd.Arg("privateIP", "Public IP address").Required().String()
    c.addNodeInstalled = cmd.Flag("installed", "Flag to indicate if the node is deployed").Default("false").Bool()
    c.addNodeUsername = cmd.Arg("username", "Username used to connect through SSH").Required().String()
    c.addNodePassword = cmd.Flag("password", "Password used to connect through SSH").Default("").String()
    c.addNodeSSHKey = cmd.Flag("sshKey", "SSH key used to connect through SSH").Default("").String()

    cmd = nodeCommand.Command("list", "List nodes").Action(c.listNodes)
    c.listNodesNetworkID = cmd.Arg(NetworkID, NetworkIDDesc).Required().String()
    c.listNodesClusterID = cmd.Arg(ClusterID, ClusterIDDesc).Required().String()

    cmd = nodeCommand.Command("get", "Get a node").Action(c.getNode)
    c.getNodeNetworkID = cmd.Arg(NetworkID, NetworkIDDesc).Required().String()
    c.getNodeClusterID = cmd.Arg(ClusterID, ClusterIDDesc).Required().String()
    c.getNodeNodeID = cmd.Arg(NodeID, NodeIDDesc).Required().String()

    cmd = nodeCommand.Command("delete", "Delete a node").Action(c.deleteNode)
    c.removeNodeNetworkID = cmd.Arg(NetworkID, NetworkIDDesc).Required().String()
    c.removeNodeClusterID = cmd.Arg(ClusterID, ClusterIDDesc).Required().String()
    c.removeNodeNodeID = cmd.Arg(NodeID, NodeIDDesc).Required().String()

    cmd = nodeCommand.Command("update", "Update a node").Action(c.updateNode)
    c.updateNodeNetworkID = cmd.Arg(NetworkID, NetworkIDDesc).Required().String()
    c.updateNodeClusterID = cmd.Arg(ClusterID, ClusterIDDesc).Required().String()
    c.updateNodeNodeID = cmd.Arg(NodeID, NodeIDDesc).Required().String()
    c.updateNodeName = cmd.Flag("name", "Node name").Default(NotAssigned).String()
    c.updateNodeDesc = cmd.Flag("desc", "Node description").Default(NotAssigned).String()
    c.updateNodeLabels = cmd.Flag("labels", "Comma separated list of labels").Default(NotAssigned).String()
    c.updateNodePublicIP = cmd.Flag("publicIP", "Public IP address").Default(NotAssigned).String()
    c.updateNodePrivateIP = cmd.Flag("privateIP", "Private IP address").Default(NotAssigned).String()
    c.updateNodeInstalled = cmd.Flag("installed", "Flag to indicate if the node is deployed").Default(NotAssigned).String()
    c.updateNodeUsername = cmd.Flag("username", "Username used to connect through SSH").Default(NotAssigned).String()
    c.updateNodePassword = cmd.Flag("password", "Password used to connect through SSH").Default(NotAssigned).String()
    c.updateNodeSSHKey = cmd.Flag("sshKey", "SSH key used to connect through SSH").Default(NotAssigned).String()
    c.updateNodeStatus = cmd.Flag("status", "Node status").
        Enum(NotAssigned, string(entities.NodeUnchecked), string(entities.NodeReadyToInstall),
            string(entities.NodeInstalling), string(entities.NodeInstalled),
            string(entities.NodeUninstalling), string(entities.NodePrecheckError),
            string(entities.NodeError))

    return c
}

func (cmd * NodeCommand) getLabels(labels string) []string {
    result := make([]string, 0)
    if labels != "" {
        result = append(result, strings.Split(labels, ",") ...)
    }
    return result
}

func (cmd * NodeCommand) addNode(c *kingpin.ParseContext) error {
    nodeRest := client.NewNodeClientRest(cmd.IP.String(), *cmd.Port)
    toAdd := entities.NewAddNodeRequest(
        * cmd.addNodeName, * cmd.addNodeDesc, cmd.getLabels(*cmd.addNodeLabels),
        * cmd.addNodePublicIP, * cmd.addNodePrivateIP, * cmd.addNodeInstalled,
        * cmd.addNodeUsername, * cmd.addNodePassword, * cmd.addNodeSSHKey)
    added, err := nodeRest.Add(* cmd.addNodeNetworkID, * cmd.addNodeClusterID, * toAdd)
    return cmd.printResultOrError(added, err)
}

func (cmd * NodeCommand) listNodes(c *kingpin.ParseContext) error {
    nodeRest := client.NewNodeClientRest(cmd.IP.String(), *cmd.Port)
    nodes, err := nodeRest.List(* cmd.listNodesNetworkID, * cmd.listNodesClusterID)
    return cmd.printResultOrError(nodes, err)
}

func (cmd * NodeCommand) getNode(c *kingpin.ParseContext) error {
    nodeRest := client.NewNodeClientRest(cmd.IP.String(), *cmd.Port)
    node, err := nodeRest.Get(* cmd.getNodeNetworkID, * cmd.getNodeClusterID, * cmd.getNodeNodeID)
    return cmd.printResultOrError(node, err)
}

func (cmd * NodeCommand) deleteNode(c *kingpin.ParseContext) error {
    nodeRest := client.NewNodeClientRest(cmd.IP.String(), *cmd.Port)
    err := nodeRest.Remove(* cmd.removeNodeNetworkID, * cmd.removeNodeClusterID, * cmd.removeNodeNodeID)
    return cmd.printResultOrError(entities.NewSuccessfulOperation("Node deleted"), err)
}

func (cmd * NodeCommand) updateNode(c * kingpin.ParseContext) error {
    nodeRest := client.NewNodeClientRest(cmd.IP.String(), *cmd.Port)
    update := entities.NewUpdateNodeRequest()
    if cmd.updateNodeName != nil && * cmd.updateNodeName != NotAssigned {
        update.WithName(* cmd.updateNodeName)
    }
    if cmd.updateNodeDesc != nil && * cmd.updateNodeDesc != NotAssigned {
        update.WithDescription(* cmd.updateNodeDesc)
    }
    if cmd.updateNodeLabels != nil && * cmd.updateNodeLabels != NotAssigned {
        update.WithLabels(cmd.getLabels(*cmd.updateNodeLabels))
    }
    if cmd.updateNodePublicIP != nil && * cmd.updateNodePublicIP != NotAssigned {
        update.WithPublicIP(* cmd.updateNodePublicIP)
    }
    if cmd.updateNodePrivateIP != nil && * cmd.updateNodePrivateIP != NotAssigned {
        update.WithPrivateIP(* cmd.updateNodePrivateIP)
    }
    if cmd.updateNodeInstalled != nil && * cmd.updateNodeInstalled != NotAssigned {
        b, err := strconv.ParseBool(* cmd.updateNodeInstalled)
        if err != nil {
            kingpin.Fatalf("invalid boolean value")
        }
        update.WithInstalled(b)
    }
    if cmd.updateNodeUsername != nil && * cmd.updateNodeUsername != NotAssigned {
        update.WithUsername(* cmd.updateNodeUsername)
    }
    if cmd.updateNodePassword != nil && * cmd.updateNodePassword != NotAssigned {
        update.WithPassword(* cmd.updateNodePassword)
    }
    if cmd.updateNodeSSHKey != nil && * cmd.updateNodeSSHKey != NotAssigned {
        update.WithSSHKey(* cmd.updateNodeSSHKey)
    }
    if cmd.updateNodeStatus != nil && * cmd.updateNodeStatus != NotAssigned {
        update.WithStatus((entities.NodeStatus)(* cmd.updateNodeStatus))
    }

    updated, err := nodeRest.Update(* cmd.updateNodeNetworkID, * cmd.updateNodeClusterID, * cmd.updateNodeNodeID, * update)
    return cmd.printResultOrError(updated, err)
}