# system-model

System Model Service

## Building the code

To build the repository and execute the test use:

```
bazel build ... -- -//vendor/...
bazel test ... -- -//vendor/...
```

## Contributing to the repository

This repository uses Dep to manage the dependencies.

```
dep ensure
dep status
```

To build the code of a specific target use.

```
bazel build //:system-model   
```

To launch the System Model Service.

```
bazel-bin/system-model api      
```

To show the System Model help.

```
bazel-bin/system-model help api       
```

## Adding new dependencies to the repo

To add a new dependency, use:

```
$ dep ensure -add github.com/foo/bar
```

To add new code, add the go files as required and use gazelle to generate the build files:

```
$ bazel run //:gazelle
```

## Using the System Model Go Client

Creating a instance of the Network client:
```
networkClient = client.NewNetworkRest("http://localhost:8800")
```

Creating a Network:
```
addedNetwork, err := networkClient.Add(entities.AddNetworkRequest{
    Name:        "Network3",
    Description: "Description of Network3",
})
 ```
Creating a instance of the Cluster client:
```
clusterClient = client.NewClusterRest("http://localhost:8800")
```
Adding a new cluster:

```
addedCluster,err :=clusterClient.Add("1", entities.AddClusterRequest {
    Name:        "Cluster4",
    Description: "Description Cluster 4",
    Type:        entities.GatewayType,
    Location:    "Madrid",
})
```

## Using the cli

A command line client is available to connect to exposed system model services.

```
$ ./bazel-bin/system_model_cli
usage: system-model-cli [<flags>] <command> [<args> ...]

Command line tool for Daisho System Model

Flags:
  --help          Show context-sensitive help (also try --help-long and --help-man).
  --ip=127.0.0.1  IP address of Daisho System Model
  --port=8800     Port of Daisho System Model
  --debug         Print detailed responses

Commands:
  help [<command>...]
    Show help.

  network add [<flags>] <name>
    Add a new network

  network list
    List all networks

  network get <networkId>
    Get a network

  network delete <networkId>
    Delete a network

  cluster add [<flags>] <networkId> <name> <type>
    Add a new cluster

  cluster update [<flags>] <networkId> <clusterId>
    Update a selected cluster

  cluster list <networkId>
    List the clusters in a network

  cluster get <networkId> <clusterId>
    Get a cluster

  cluster delete <networkId> <clusterId>
    Delete a cluster

  node add [<flags>] <networkId> <clusterId> <name> <publicIP> <privateIP> <username>
    Add a new node

  node list <networkId> <clusterId>
    List nodes

  node get <networkId> <clusterId> <nodeId>
    Get a node

  node delete <networkId> <clusterId> <nodeId>
    Delete a node

  node update [<flags>] <networkId> <clusterId> <nodeId>
    Update a node

  application descriptor list <networkId>
    List all descriptors

  application descriptor add [<flags>] <networkId> <name> <serviceName> <serviceVersion> <label> <appPort>
    Add a new descriptor

  application descriptor get <networkId> <descriptorId>
    Get a descriptor

  application instance add [<flags>] <networkId> <descriptorId> <name> <persistenceSize> <storageType>
    Add a new instance

  application instance list <networkId>
    List instances

  application instance get <networkId> <deployedId>
    Get instance

  application instance update [<flags>] <networkId> <deployedId>
    Update an instance

  application instance delete <networkId> <deployedId>
    Delete an instance

  dump export
    Export all system model entities

  info reduced
    Get the essential info of the system model

  info summary
    Get the number of stored elements

  info reduced-by-network <networkId>
    Get the essential info of the system model by network

  user add <id> <name> <phone> <email>
    Add a new user

  user get <id>
    Get an existing user

  user delete <id>
    Delete an existing user

  user update [<flags>] <id>
    Update an existing user

  user list
    List users

  access add <userId> <role>
    Add new access

  access get <userId>
    Get user access entry

  access delete <userId>
    Delete user access entry

  access list
    List all the users with their privileges
```

## Using persistence providers

To use a filesystem persistence provider use:

```
./bazel-out/local-fastbuild/bin/system-model api --filesystem-persistence --filesystem-basePath=/tmp/newPath/
```

## Default admin user

A default admin user will be generated when running the system model for the first time. 
The name of this user can be modified when running the service.
```
system-model api --filesystem-persistence --default-admin-user=myadmin
```

## Caution

Gazelle has problems with the library golang.org/x/sys/unix (Logrus requires this library), 
because gazelle try to include the CGO files in the target. Put cgo flag to False and the build works.

 
