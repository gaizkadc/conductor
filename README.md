# Conductor

Conductor is an application super-orchestrator for the Nalej platform.

The current solution has two components: conductor and musician. Conductor is the central Nalej scheduler, while
musicians are in-cluster services. Musicians must be accessible to conductor in order to ensure a correct
behavior.

This explanation assumes a minikube environment is up and running.

## Musician

Musicians use Prometheus available cluster monitoring data. For that reason, prior to deploy any musician in the
cluster the corresponding monitoring infrastructure must be deployed. A monitoring infrastructure is defined in
components/monitoring/component.yaml. Deploy the corresponding services by running

```bash
kubectl create namespace nalej
kubect create -f components/monitoring/component.yaml
```
Check the port and ip address of the service with
```bash
minikube service list
|-------------|----------------------|-----------------------------|
|  NAMESPACE  |         NAME         |             URL             |
|-------------|----------------------|-----------------------------|
| default     | kubernetes           | No node port                |
| kube-system | kube-dns             | No node port                |
| kube-system | kubernetes-dashboard | http://192.168.99.100:30000 |
| nalej       | prometheus           | http://192.168.99.100:31080 |
|-------------|----------------------|-----------------------------|
```

The cluster identifier must be available in an environment variable called CLUSTER_ID. If the variable is not set
an error will be displayed and the musician will not start. If the onboard conductor demo is used, it will display
a plausible clusterid to be used in the deployment process.

```bash
./bin/conductor demo
...
{"level":"info","time":1540915189,"message":"The output instance works with id: bd3e31c6-4d48-48d0-9206-49c3592716b6"}
...
```
Now use the id for musician to identify its clusterid.
```bash
export CLUSTER_ID="9c9ccc95-f7c3-436b-8d43-953640ba6724"
```


Execute a musician pointing the Nalej prometheus service.
```bash
./bin/conductor musician -o http://192.168.99.100:31080 --consoleLogging
2018-10-23T15:51:14+02:00 |INFO| launching musician...
2018-10-23T15:51:14+02:00 |INFO| Running server...
2018-10-23T15:51:14+02:00 |INFO| starting Prometheus status collector...
2018-10-23T15:51:14+02:00 |INFO| Launching gRPC server port=5100
```

## Conductor

In order to run conductor, a system model service must be up and running.
```bash
./bin/conductor run --debug --consoleLogging -s localhost:8800
2018-10-23T15:50:48+02:00 |INFO| launching conductor...
2018-10-23T15:50:48+02:00 |INFO| gRPC port port=5000
2018-10-23T15:50:48+02:00 |INFO| System Model URL=localhost:8800
2018-10-23T15:50:48+02:00 |DEBUG| add new connection address=localhost:8800
2018-10-23T15:50:48+02:00 |INFO| Connected to address at localhost:8800
2018-10-23T15:50:48+02:00 |DEBUG| connection successfully added address=localhost:8800
2018-10-23T15:50:48+02:00 |INFO| Running server...
2018-10-23T15:50:48+02:00 |INFO| Launching gRPC server port=5000
```

Conductor should be now up and listening to incoming grpc connections in port 5000. For manual testing, use grpc_cli
with the following example:

```bash
grpc_cli call localhost:5000 conductor.Conductor.Deploy "request_id: 'req_001', app_id: {organization_id: 'org_001', app_descriptor_id: 'app_001'}, cpu:0.3, disk:2000, memory:3000"
```

# Integration tests
Many tests run by conductor require of other Nalej environment components to be up and running. The following table
summarizes the set of expected testing variables.

| Variable  | Example Value | Description |
| ------------- | ------------- |------------- |
| RUN_INTEGRATION_TEST  | true | Run integration tests |
| IT_SYSTEM_MODEL | localhost:8800 | Address of an available system model server |
| CLUSTER_ID | 28602103-1462-43cf-bb38-44e880fa1933 | Cluster id where musician is running |


