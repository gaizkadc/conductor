# Conductor

Conductor is an application super-orchestrator and scheduler for the Daisho platform. 

## Scheduling policies

Conductor currently supports the following scheduling policies:

### Round-robin
Every application will be deployed on a previously specified type of cluster (cloud, gateway or edge). The cluster 
selected for deployment will be chosen on a [round-robin policy](https://en.wikipedia.org/wiki/Round-robin_scheduling). The clusters will be iteratively selected until an accesible network is found.

## About testing and integration
Conductor is a complex piece of software that involves several components microservices and pieces of software. In order to simplify the testing of this solution, we identify two different scenarios:

###Testing
Common approach using unitary testing to ensure the correctness of data structure, endpoints and scheduling algorithms.

###Workflow integration
More complex testing that involves other external pieces of software and services. In this case, we can expect elements such as nodes, clusters, infrastructure deployments, etc. The workflow integration is intended to offer a more "realistic" testing scenario that may require human supervision and will be closer to the design of experimental 
scenarios.

Refer to the conductor-integration section for more information.


### Command line client

Use the help provided by the CLI:

```
conductor juan$ ./bazel-bin/conductor-cli --help
usage: conductor-cli [<flags>] <command> [<args> ...]

Command line tool for Daisho conductor.

Flags:
  --help          Show context-sensitive help (also try --help-long and --help-man).
  --ip=127.0.0.1  Ip address of conductor service.
  --port=9000     Port number of conductor service.

Commands:
  help [<command>...]
    Show help.

  orchestrator deploy <network.id> <app.descriptor> <app.label> <app.name> [<app.arguments>] [<app.description>] [<app.persistenceSize>] [<app.storage>]
    Deploy an application on a given network

  orchestrator undeploy <network.id> <instance.id>
    Undeploy and application on a given network`
```
### Conductor-integration

Conductor integration makes possible to test conductor in a local environment deployed using colonized. Certain elements must be provided in advance. Take a look at the executable help to have a better idea of what elements are required. Briefly conductor-integration requires,

* Factory compiled targets available
* Colony compiled an ready to deploy elements
* InfluxDB application package from the appdevkit ready
* ASM application manager ready

```
conductor juan$ ./bazel-bin/conductor-integration start --help
usage: conductor-integration start [<flags>]

Launch the service API

Flags:
  --help   Show context-sensitive help (also try --help-long and --help-man).
  --factoryTargetPath="~/daisho/factory/target"
           Path were the daisho factory targets are available.
  --asmPath="~/daisho/appmgr/bazel-bin/asmcli"
           ASM path folder
  --colonyPath="~/daisho/colony/colonize"
           Colony folder were the binary is allocated
  --appdevkitPath="~/daisho/appdevkit"
           Folder were we can find appdevkit
  --debug  Activate debug logging
```

The execution process is automatic and requires human intervention to define whether the integration was successful or not.
