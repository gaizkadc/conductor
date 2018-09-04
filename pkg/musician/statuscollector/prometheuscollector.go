//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

package statuscollector

import (
    "net/http"
    "time"
    "fmt"
    "github.com/nalej/conductor/entities"
    "encoding/json"
)


// Simple client to query Prometheus HTTP API.
type PrometheusHTTPClient struct {
    // The http client
    client http.Client
    // The target address
    address string
}

func NewPrometheusHTTPClient(address string) *PrometheusHTTPClient {
    client := http.Client{Timeout: time.Second * 3}
    return &PrometheusHTTPClient{client, address}
}

// Get the available memory in the cluster
func(c *PrometheusHTTPClient) GetMemory() (*entities.PrometheusMemoryStatus,error) {
    target := fmt.Sprintf("%s/api/v1/query?query=node_memory_MemFree",c.address)
    req, err := http.NewRequest(http.MethodGet,target,nil)
    if err != nil {
        fmt.Printf("There was an error creating request for %s",target)
        return nil,err
    }
    res, err := c.client.Do(req)
    if err != nil {
        fmt.Printf("Error saying %s",err)
        return nil, err
    }

    defer res.Body.Close()

    var ent entities.PrometheusMemoryStatus
    if err := json.NewDecoder(res.Body).Decode(&ent); err != nil {
        return nil, err
    }

    return &ent, nil
}


type  PrometheusStatusCollector struct {
    client PrometheusHTTPClient
    // Milliseconds to sleep between calls.
    sleepDuration time.Duration
    // Cached status
    cached map[string]interface{}
}


func NewPrometheusStatusCollector(address string, sleepTime int) *PrometheusStatusCollector {
    // Build a client
    client := NewPrometheusHTTPClient(address)
    sleepDuration := time.Duration(time.Millisecond) * time.Duration(sleepTime)
    return &PrometheusStatusCollector{*client, sleepDuration, make(map[string]interface{})}
}

// Start the collector
// return:
//  Error if any
func(coll *PrometheusStatusCollector) Run() error {
    fmt.Println("Starting status collector...")
    for {
        fmt.Println("Get memory status...")
        mem, err := coll.client.GetMemory()
        if err != nil {
            fmt.Printf("Error requesting memory %s",err)
        } else {
            fmt.Printf("Current mem status: %v\n",*mem)
            coll.cached["memory"] = mem
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
    return coll.GetStatus()
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


