//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Application related commands.

package main

import (
    "fmt"
    "strings"

    "gopkg.in/alecthomas/kingpin.v2"

    "github.com/daishogroup/system-model/client"
    "github.com/daishogroup/system-model/entities"
)


// Make '--images' a cumulative flag
type imageList []string

func (i *imageList) Set(value string) error {
    // Verify the image format, it MUST contain both repo and tag
    items := strings.Split(value, ":")
    if len(items) != 2 || len(items[0]) == 0 || len(items[1]) == 0 {
        return fmt.Errorf("Wrong format of image '%s' - it must contain both repository and tag", value)
    }
    *i = append(*i, value)
    return nil
}

func (i *imageList) String() string {
    return ""
}

func (i *imageList) IsCumulative() bool {
    return true
}

func ImageList(s kingpin.Settings) (target *imageList) {
    target = new(imageList)
    s.SetValue((*imageList)(target))
    return
}

type ApplicationCommand struct {
    GlobalCommand

    listDescNetworkID * string

    addDescNetworkID      * string
    addDescName           * string
    addDescDescription    * string
    addDescServiceName    * string
    addDescServiceVersion * string
    addDescLabel          * string
    addDescPort           * int
    addDescImages         * imageList

    getDescNetworkID    * string
    getDescDescriptorID * string

    deleteDescNetworkID  * string
    deleteDescDescriptorID * string

    addInstNetworkID      * string
    addInstDescriptorID   * string
    addInstName           * string
    addInstDescription    * string
    addInstLabel          * string
    addInstArguments      * string
    addInstPersistentSize * string
    addInstStorageType    * string

    listInstNetworkID * string

    getInstNetworkID  * string
    getInstDeployedID * string


    updateInstNetworkID      * string
    updateInstDeployedID     * string
    updateInstDescriptorID   * string
    updateInstClusterID      * string
    updateInstDescription    * string
    updateInstStatus         * string
    updateInstClusterAddress * string

    deleteInstNetworkID  * string
    deleteInstDeployedID * string

}

func NewApplicationCommand(app * kingpin.Application, global GlobalCommand) * ApplicationCommand {

    c := &ApplicationCommand{
        GlobalCommand : global,
    }
    appCommand := app.Command("application", "Application commands")

    descCommand := appCommand.Command("descriptor", "Application descriptor commands")

    cmd := descCommand.Command("list", "List all descriptors").Action(c.listDescriptors)
    c.listDescNetworkID = cmd.Arg(NetworkID, NetworkIDDesc).HintAction(nil).Required().String()

    cmd = descCommand.Command("add", "Add a new descriptor").Action(c.addDescriptor)
    c.addDescNetworkID = cmd.Arg(NetworkID, NetworkIDDesc).Required().String()
    c.addDescName = cmd.Arg("name", "Descriptor name").Required().String()
    c.addDescDescription = cmd.Flag("description", "Descriptor description").Default("").String()
    c.addDescServiceName = cmd.Arg("serviceName", "Name of the service as in the manifest").
        Required().String()
    c.addDescServiceVersion = cmd.Arg("serviceVersion", "Version of the service as in the manifest").
        Required().String()
    c.addDescLabel = cmd.Arg("label", "Descriptor label").Required().String()
    c.addDescPort = cmd.Arg("appPort", "Application exposed port").Required().Int()
    c.addDescImages  = ImageList(cmd.Flag("images", "Required images"))


    cmd = descCommand.Command("get", "Get a descriptor").Action(c.getDescriptor)
    c.getDescNetworkID = cmd.Arg(NetworkID, NetworkIDDesc).HintAction(nil).Required().String()
    c.getDescDescriptorID = cmd.Arg(DescriptorID, DescriptorIDDesc).
        HintAction(nil).Required().String()

    cmd = descCommand.Command("delete", "Delete a descriptor").Action(c.deleteDescriptor)
    c.deleteDescNetworkID = cmd.Arg(NetworkID, NetworkIDDesc).HintAction(nil).Required().String()
    c.deleteDescDescriptorID = cmd.Arg(DescriptorID, DescriptorIDDesc).HintAction(nil).Required().String()

    instCommand := appCommand.Command("instance", "Application instance commands")

    cmdi := instCommand.Command("add", "Add a new instance").Action(c.addInstance)
    c.addInstNetworkID = cmdi.Arg(NetworkID, NetworkIDDesc).Required().String()
    c.addInstDescriptorID = cmdi.Arg(DescriptorID, DescriptorIDDesc).Required().String()
    c.addInstName = cmdi.Arg("name", "Instance name").Required().String()
    c.addInstDescription = cmdi.Flag("description", "Instance description").Default("").String()
    c.addInstLabel = cmdi.Arg("label", "Instance label").Required().String()
    c.addInstArguments = cmdi.Flag("arguments", "Instance arguments").Default("").String()
    c.addInstPersistentSize = cmdi.Flag("persistenceSize", "Persistence size required").String()
    c.addInstStorageType = cmdi.Flag("storageType", "Storage Type").
        Enum(string(entities.AppStorageDefault),
        string(entities.AppStoragePersistent),
        string(entities.AppStorageNetPersistent))

    cmdi = instCommand.Command("list", "List instances").Action(c.listInstances)
    c.listInstNetworkID = cmdi.Arg(NetworkID, NetworkIDDesc).HintAction(nil).Required().String()

    cmdi = instCommand.Command("get", "Get instance").Action(c.getInstance)
    c.getInstNetworkID = cmdi.Arg(NetworkID, NetworkIDDesc).HintAction(nil).Required().String()
    c.getInstDeployedID = cmdi.Arg(DeployedID, DeployedIDDesc).HintAction(nil).Required().String()

    cmdi = instCommand.Command("update", "Update an instance").Action(c.updateInstance)
    c.updateInstNetworkID = cmdi.Arg(NetworkID, NetworkIDDesc).Required().String()
    c.updateInstDeployedID = cmdi.Arg(DeployedID, DeployedIDDesc).Required().String()
    c.updateInstClusterID = cmdi.Flag(ClusterID, ClusterIDDesc).Default(NotAssigned).String()
    c.updateInstDescription = cmdi.Flag("description", "Instance description").Default(NotAssigned).String()
    c.updateInstStatus = cmdi.Flag("status", "Instance status").
        Enum(string(entities.AppInstInit),
            string(entities.AppInstReady),
            string(entities.AppInstNotReady),
            string(entities.AppInstError))
    c.updateInstClusterAddress = cmdi.Flag("clusterAddress", "Cluster address").Default(NotAssigned).String()

    cmdi = instCommand.Command("delete", "Delete an instance").Action(c.deleteInstance)
    c.deleteInstNetworkID = cmdi.Arg(NetworkID, NetworkIDDesc).HintAction(nil).Required().String()
    c.deleteInstDeployedID = cmdi.Arg(DeployedID, DeployedIDDesc).HintAction(nil).Required().String()

    return c
}

func (cmd * ApplicationCommand) addDescriptor(c *kingpin.ParseContext) error {
    appRest := client.NewApplicationClientRest(cmd.IP.String(), *cmd.Port)
    toAdd := entities.NewAddAppDescriptorRequest(
        * cmd.addDescName, * cmd.addDescDescription, * cmd.addDescServiceName,
        * cmd.addDescServiceVersion, * cmd.addDescLabel, * cmd.addDescPort,
        * cmd.addDescImages)
    descriptors, err := appRest.AddApplicationDescriptor(* cmd.addDescNetworkID, * toAdd)
    return cmd.printResultOrError(descriptors, err)
}

func (cmd * ApplicationCommand) listDescriptors(c *kingpin.ParseContext) error {
    appRest := client.NewApplicationClientRest(cmd.IP.String(), *cmd.Port)
    descriptors, err := appRest.ListDescriptors(* cmd.listDescNetworkID)
    return cmd.printResultOrError(descriptors, err)
}

func (cmd * ApplicationCommand) getDescriptor(c *kingpin.ParseContext) error {
    appRest := client.NewApplicationClientRest(cmd.IP.String(), *cmd.Port)
    descriptor, err := appRest.GetDescriptor(* cmd.getDescNetworkID, * cmd.getDescDescriptorID)
    return cmd.printResultOrError(descriptor, err)
}

func (cmd * ApplicationCommand) deleteDescriptor(c *kingpin.ParseContext) error {
    appRest := client.NewApplicationClientRest(cmd.IP.String(), *cmd.Port)
    err := appRest.DeleteDescriptor(* cmd.deleteDescNetworkID, * cmd.deleteDescDescriptorID)
    return cmd.printResultOrError(entities.NewSuccessfulOperation("Descriptor deleted"), err)
}

func (cmd * ApplicationCommand) addInstance(c *kingpin.ParseContext) error {
    appRest := client.NewApplicationClientRest(cmd.IP.String(), *cmd.Port)
    toAdd := entities.NewAddAppInstanceRequest(
        * cmd.addInstDescriptorID, * cmd.addInstName, * cmd.addInstDescription,
        * cmd.addInstLabel, * cmd.addInstArguments, * cmd.addInstPersistentSize,
        (entities.AppStorageType)(* cmd.addInstStorageType))
    instance, err := appRest.AddApplicationInstance(* cmd.addInstNetworkID, * toAdd)
    return cmd.printResultOrError(instance, err)
}

func (cmd * ApplicationCommand) listInstances(c *kingpin.ParseContext) error {
    appRest := client.NewApplicationClientRest(cmd.IP.String(), *cmd.Port)
    instances, err := appRest.ListInstances(* cmd.listInstNetworkID)
    return cmd.printResultOrError(instances, err)
}

func (cmd * ApplicationCommand) getInstance(c *kingpin.ParseContext) error {
    appRest := client.NewApplicationClientRest(cmd.IP.String(), *cmd.Port)
    instance, err := appRest.GetInstance(* cmd.getInstNetworkID, * cmd.getInstDeployedID)
    return cmd.printResultOrError(instance, err)
}

func (cmd * ApplicationCommand) updateInstance(c *kingpin.ParseContext) error {
    appRest := client.NewApplicationClientRest(cmd.IP.String(), *cmd.Port)

    toUpdate := entities.NewUpdateAppInstRequest()
    if cmd.updateInstClusterID != nil && * cmd.updateInstClusterID != NotAssigned{
        toUpdate = toUpdate.WithClusterID(* cmd.updateInstClusterID)
    }
    if cmd.updateInstDescription != nil && * cmd.updateInstDescription != NotAssigned{
        toUpdate = toUpdate.WithDescription(* cmd.updateInstDescription)
    }
    if cmd.updateInstStatus != nil && * cmd.updateInstStatus != ""{
        toUpdate = toUpdate.WithStatus((entities.AppStatus)(* cmd.updateInstStatus))
    }
    if cmd.updateInstClusterAddress != nil && * cmd.updateInstClusterAddress != NotAssigned{
        toUpdate = toUpdate.WithClusterAddress(* cmd.updateInstClusterAddress)
    }

    updated, err := appRest.UpdateInstance(* cmd.updateInstNetworkID, * cmd.updateInstDeployedID, * toUpdate)
    return cmd.printResultOrError(updated, err)
}

func (cmd * ApplicationCommand) deleteInstance(c *kingpin.ParseContext) error {
    appRest := client.NewApplicationClientRest(cmd.IP.String(), *cmd.Port)
    err := appRest.DeleteInstance(* cmd.deleteInstNetworkID, * cmd.deleteInstDeployedID)
    return cmd.printResultOrError(entities.NewSuccessfulOperation("Instance deleted"), err)
}
