# Conductor

Conductor is the application super-orchestrator for the Nalej platform. This component
elaborates deployment plans and controls the lifecycle of any application deployed using
Nalej.  

Conductor works in conjunction with musicians deployed on every application cluster. 

## Musician

Musicians score how well they can deploy application fragments. In order to do this, musicians
capture cluster metrics and evaluate the suitability of the deployment in the current cluster. 
This information is collected by Conductor who decides what is the best deployment plan that
satisfies existing application constraints. 


## Getting Started

Check the following entries before deploying Conductor.

### Prerequisites

* CA certificate shared among all clusters
* Storage volume for local database
* Nalej-bus 
* Network manager
* Unified logging
* System model 



### Build and compile

In order to build and compile this repository use the provided Makefile:

```
make all
```

This operation generates the binaries for this repo, download dependencies,
run existing tests and generate ready-to-deploy Kubernetes files.

### Run tests

Tests are executed using Ginkgo. To run all the available tests:

```
make test
```

### Update dependencies

Dependencies are managed using Godep. For an automatic dependencies download use:

```
make dep
```

In order to have all dependencies up-to-date run:

```
dep ensure -update -v
```

## Known issues
* Slow Prometheus startup times may delay Musicians bootstrapping.
* Musicians score deployments using temporal metrics. If a recently deployed musician is 
requested to score a deployment, it may result in inaccurate scores due to the lack of references.

## Contributing
​
Please read [contributing.md](contributing.md) for details on our code of conduct, and the process for submitting pull requests to us.
​
​
## Versioning
​
We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 
​
## Authors
​
See also the list of [contributors](https://github.com/nalej/grpc-utils/contributors) who participated in this project.
​
## License
This project is licensed under the Apache 2.0 License - see the [LICENSE-2.0.txt](LICENSE-2.0.txt) file for details.