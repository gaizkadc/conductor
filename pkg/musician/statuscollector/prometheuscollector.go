//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//


package statuscollector

import (
    "github.com/prometheus/client_golang/api/prometheus/v1"
    "github.com/prometheus/common/model"
    "github.com/prometheus/client_golang/api"
    "fmt"
    "time"
    "context"
)


// Simple client to query Prometheus HTTP API.
type PrometheusClient struct {
    // The api endpoint
    api v1.API
    // The target address
    address string
}

func NewPrometheusClient(address string) *PrometheusClient {
    // create a client using the address
    client, err := api.NewClient(api.Config{Address: address})
    if err != nil {
        fmt.Printf("error creating Prometheus client %s",err)
        return nil
    }
    api := v1.NewAPI(client)
    return &PrometheusClient{api,address}
}

// Get the available memory in the cluster
// return:
//
func(c *PrometheusClient) GetMemory() (*model.Value,error) {

    value, err := c.api.Query(context.Background(),"node_memory_MemFree",time.Now())
    if err != nil {
        return nil, err
    }

    return &value, nil
}


// Main structure containing necessary items to support a status collector for Prometheus.
type  PrometheusStatusCollector struct {
    client PrometheusClient
    // Milliseconds to sleep between calls.
    sleepDuration time.Duration
    // Cached status
    // TODO: Evaluate potential ways to have a more efficient storage.
    cached map[string]model.Value
}


func NewPrometheusStatusCollector(address string, sleepTime int) *PrometheusStatusCollector {
    // Build a client
    client := NewPrometheusClient(address)
    sleepDuration := time.Duration(time.Millisecond) * time.Duration(sleepTime)
    return &PrometheusStatusCollector{*client, sleepDuration, make(map[string]model.Value)}
}

// Start the collector
// return:
//  Error if any
func(coll *PrometheusStatusCollector) Run() error {
    fmt.Println("Starting Prometheus status collector...")
    for {
        fmt.Println("Get memory status...")
        mem, err := coll.client.GetMemory()
        if err != nil {
            fmt.Printf("Error requesting memory %s",err)
        } else {
            fmt.Printf("%s\n", *mem)
            coll.cached["memory"] = *mem
        }
        time.Sleep(coll.sleepDuration)
    }
    return nil
}

// Stop the collector.
// return:
//  Error if any
func(coll *PrometheusStatusCollector) Finalize(killSignal bool) error {
    fmt.Println("Finalize was called")
    return nil
}

// Get the current status.
func(coll *PrometheusStatusCollector) GetStatus() string {
    return coll.cached["mem"].String()
}

// Return the status collector name.
// return:
//  Name of this collector.
func(coll *PrometheusStatusCollector) Name() string {
    return "Prometheus status collector"
}

// Return a description of this status collector.
// return:
//  Description of this collector.
func(coll *PrometheusStatusCollector) Description() string {
    return "Status collector based on Prometheus"
}


