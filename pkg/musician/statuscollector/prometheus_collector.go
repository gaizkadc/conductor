/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */



package statuscollector

import (
    "github.com/nalej/derrors"
    "github.com/prometheus/client_golang/api/prometheus/v1"
    "github.com/prometheus/common/model"
    "github.com/prometheus/client_golang/api"
    "github.com/rs/zerolog/log"
    "time"
    "context"
    "github.com/nalej/conductor/internal/entities"
)


const (
    // describe the set of queries we use for Prometheus to collecto monitor stats
    // accumulated free memory for all the nodes in the cluster
    PROM_MEM_FREE_QUERY="sum(node_memory_MemFree)"
    PROM_MEM_FREE_NAME="mem_free"
    // accumulated idle time for all the cpus in the cluster
    PROM_CPU_NUM_QUERY="count(count by (nodename,cpu) (node_cpu))"
    PROM_CPU_NUM_NAME="cpu_num"
    // the available disk space is the one of the node with the largest available space
    PROM_DISK_FREE_QUERY="max(node_filesystem_free{mountpoint=\"/\"})"
    PROM_DISK_FREE_NAME="disk_free"
    // Sum of total idle CPUNum time in the cluster
    PROM_CPU_IDLE_QUERY="sum(node_cpu {mode=\"idle\"}) by (mode)"
    PROM_CPU_IDLE_NAME="cpu_idle"
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
        log.Error().Msgf("error creating Prometheus client %s",err)
        return nil
    }
    api := v1.NewAPI(client)
    return &PrometheusClient{api,address}
}


// Internal function to exec an existing query
// params:
//  query to be executed
// returns:
//  floating value returned from Prometheus
//  error if any
func (c *PrometheusClient) execQuery(query string) (float64, error) {
    value, err := c.runQuery(query)
    if err != nil {
        log.Error().Err(err).Str("query", query).Msg("error querying prometheus")
        return -1, err
    }

    // We expect a vector value
    var vectorValue model.Vector
    vectorValue = (*value).(model.Vector)

    if vectorValue.Len() > 1 {
        log.Error().Str("query", query).Msg("query returned more than one entry")
        log.Error().Interface("vectorValue", vectorValue).Msg("data returned")
        return -1, derrors.NewInternalError("query returned more than one entry")
    }
    toReturn := float64(vectorValue[0].Value)
    return toReturn, nil
}


func (c *PrometheusClient) runQuery(query string) (*model.Value, error) {
    value, err := c.api.Query(context.Background(), query, time.Now())
    if err != nil {
        log.Error().Err(err).Msg("error querying prometheus")
        return nil, err
    }

    return &value, nil
}


// Main structure containing necessary items to support a status collector for Prometheus.
type  PrometheusStatusCollector struct {
    client PrometheusClient
    // Milliseconds to sleep between calls.
    sleepDuration time.Duration
    // Map of queries to be sent to Prometheus
    prometheusQueries map[string]string
    // Cached status
    // TODO: Evaluate potential ways to have a more efficient provider.
    cached Cache
}


func NewPrometheusStatusCollector(address string, sleepTime uint32) StatusCollector {
    // Build a client
    client := NewPrometheusClient(address)
    sleepDuration := time.Duration(time.Millisecond) * time.Duration(sleepTime)
    cache := NewSimpleCache()
    prometheusQueries := map[string]string{
       PROM_MEM_FREE_NAME: PROM_MEM_FREE_QUERY,
       PROM_CPU_NUM_NAME: PROM_CPU_NUM_QUERY,
       PROM_DISK_FREE_NAME : PROM_DISK_FREE_QUERY,
       PROM_CPU_IDLE_NAME: PROM_CPU_IDLE_QUERY,
    }
    return &PrometheusStatusCollector{client: *client, sleepDuration: sleepDuration, cached: cache,
        prometheusQueries: prometheusQueries}
}

// Start the collector
// return:
//  Error if any
func(coll *PrometheusStatusCollector) Run() error {
    log.Info().Msg("starting Prometheus status collector...")

    sleep := time.Tick(coll.sleepDuration)
    for {
        select {
        case <-sleep:
            for queryName, query := range coll.prometheusQueries {
                value, err := coll.client.execQuery(query)
                if err == nil {
                    // log.Debug().Str("query",queryName).Float64("value", value).
                    //    Msgf("%s -> %f",queryName, value)
                    coll.cached.Put(queryName, float64(value))
                } else {
                    log.Error().Err(err).Msgf("error when querying %s", queryName)
                }
            }
        }
    }

    return nil
}

// Stop the collector.
// return:
//  Error if any
func(coll *PrometheusStatusCollector) Finalize(killSignal bool) error {
    log.Info().Msg("finalize was called")
    return nil
}

// Get the current status.
func(coll *PrometheusStatusCollector) GetStatus() (*entities.Status, error) {
    // Build the status and return it
    cacheContent := coll.cached.GetAll()
    if cacheContent == nil {
        return nil, derrors.NewNotFoundError("not found cache entries")
    }

    return &entities.Status {
        Timestamp: time.Now(),
        MemFree:   cacheContent[PROM_MEM_FREE_NAME].Value.(float64),
        CPUIdle:   cacheContent[PROM_CPU_IDLE_NAME].Value.(float64),
        DiskFree:  cacheContent[PROM_DISK_FREE_NAME].Value.(float64),
        CPUNum:    cacheContent[PROM_CPU_NUM_NAME].Value.(float64),
    }, nil

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


