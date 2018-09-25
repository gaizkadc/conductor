#!/bin/sh

grpc_cli call localhost:5000 conductor.Conductor.Deploy "request_id: 'req_001', app_id: {organization_id: 'org_001', app_descriptor_id: 'app_001'}, cpu:0.3, disk:2000, memory:3000"