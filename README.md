# Conductor

Conductor is an application super-orchestrator and scheduler for the Nalej platform.

The current solution has two components: conductor and musician. Conductor is the central Nalej scheduler, while
musicians are locally available to a cluster. Musicians must be accessible to conductor in order to ensure a correct
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

Execute a musician pointing the nalej prometheus service.
```bash
./bin/conductor musician -o http://192.168.99.100:31080
{"level":"info","time":1537868707,"message":"running init config"}
{"level":"info","time":1537868707,"message":"no config file was set"}
{"level":"info","time":1537868707,"message":"launching musician..."}
{"level":"info","time":1537868707,"message":"Running server..."}
{"level":"info","time":1537868707,"message":"starting Prometheus status collector..."}
{"level":"info","port":5001,"time":1537868707,"message":"Launching gRPC server"}
```

## Conductor

In order to run conductor a list of available musician addresses must be provided. By default musicians listen
to incoming grpc connections in port 5001.
```bash
./bin/conductor run -m 127.0.0.1:5001
{"level":"info","time":1537872175,"message":"running init config"}
{"level":"info","time":1537872175,"message":"no config file was set"}
{"level":"info","time":1537872175,"message":"launching conductor..."}
{"level":"info","time":1537872175,"message":"Connected to address at 127.0.0.1:5001"}
{"level":"info","address":"127.0.0.1:5001","time":1537872175,"message":"musician address correctly added"}
{"level":"info","time":1537872175,"message":"Running server..."}
{"level":"info","port":5000,"time":1537872175,"message":"Launching gRPC server"}
```

Conductor should be now up and listening to incoming grpc connections in port 5000. For manual testing, use grpc_cli
with the following example:

```bash
grpc_cli call localhost:5000 conductor.Conductor.Deploy "request_id: 'req_001', app_id: {organization_id: 'org_001', app_descriptor_id: 'app_001'}, cpu:0.3, disk:2000, memory:3000"
```



