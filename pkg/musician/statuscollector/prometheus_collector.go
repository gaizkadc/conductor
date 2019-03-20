/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */



package statuscollector

import (
    "github.com/prometheus/client_golang/api/prometheus/v1"
    "github.com/prometheus/common/model"
    "github.com/prometheus/client_golang/api"
    "github.com/rs/zerolog/log"
    "time"
    "context"
    "github.com/nalej/conductor/internal/entities"
    "errors"
)


const (
    // describe the set of queries we use for Prometheus to collecto monitor stats
    // accumulated free memory for all the nodes in the cluster
    PROM_MEM_QUERY="sum(node_memory_MemFree)"
    // accumulated idle time for all the cpus in the cluster
    // PROM_CPU_QUERY="sum(node_cpu{mode=~\"idle\"})"
    PROM_CPU_QUERY="count(count by (nodename,cpu) (node_cpu))"
    // the available disk space is the one of the node with the largest available space
    PROM_DISK_QUERY="max(node_filesystem_free{mountpoint=\"/\"})"
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

// Get the available memory in the cluster
// return:
//  Value with the result of the query, error if any.
func(c *PrometheusClient) GetMemory() (float64,error) {
    value, err := c.runQuery(PROM_MEM_QUERY)
    if err != nil {
        log.Error().Msg("error querying memory")
        return -1, err
    }
    // We expect a vector value
    var vectorValue model.Vector
    vectorValue = (*value).(model.Vector)
    if vectorValue.Len() > 1 {
        log.Error().Msg("mem query returned more than one entry")
        log.Error().Interface("vectorValue", vectorValue).Msg("data returned")
        return -1, errors.New("mem query returned more than one entry")
    }
    return float64(vectorValue[0].Value), nil
}

// Get the cpu load.
// return:
//  Value with the result of the query, error if any.
func(c *PrometheusClient) GetCPU() (float64, error) {
    value, err := c.runQuery(PROM_CPU_QUERY)
    if err != nil {
        log.Error().Msg("error querying CPU")
        return -1, err
    }
    // We expect a vector value
    var vectorValue model.Vector
    vectorValue = (*value).(model.Vector)
    if vectorValue.Len() > 1 {
        log.Error().Msg("cpu query returned more than one entry")
        log.Error().Interface("vectorValue", vectorValue).Msg("data returned")
        return -1, errors.New("cpu query returned more than one entry")
    }

    return float64(vectorValue[0].Value), nil
}

// Get the available disk space.
// return:
//  Value with the result of the query, error if any.
func(c *PrometheusClient) GetDisk() (float64, error) {
    value, err := c.runQuery(PROM_DISK_QUERY)
    if err != nil {
        log.Error().Msg("error querying disk")
        return -1, err
    }
    // We expect a vector value
    var vectorValue model.Vector
    vectorValue = (*value).(model.Vector)
    if vectorValue.Len() > 1 {
        log.Error().Msg("disk query returned more than one entry")
        log.Error().Interface("vectorValue", vectorValue).Msg("data returned")
        return -1, errors.New("disk query returned more than one entry")
    }

    return float64(vectorValue[0].Value), nil
}

func (c *PrometheusClient) runQuery(query string) (*model.Value, error) {
    value, err := c.api.Query(context.Background(), query, time.Now())
    if err != nil {
        log.Error().Msg("err")
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
    cached Cache
}


func NewPrometheusStatusCollector(address string, sleepTime uint32) StatusCollector {
    // Build a client
    client := NewPrometheusClient(address)
    sleepDuration := time.Duration(time.Millisecond) * time.Duration(sleepTime)
    cache := NewSimpleCache()
    return &PrometheusStatusCollector{*client, sleepDuration, cache}
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
            mem, err := coll.client.GetMemory()
            if err != nil {
                log.Error().Msgf("error requesting memory %s",err)
            } else {
                log.Debug().Msgf("memory: %f", mem)
                coll.cached.Put("memory", mem)
            }

            cpu, err := coll.client.GetCPU()
            if err != nil {
                log.Error().Msgf("error requesting cpu %s",err)
            } else {
                log.Debug().Msgf("cpu: %f", cpu)
                coll.cached.Put("cpu", cpu)
            }

            disk, err := coll.client.GetDisk()
            if err != nil {
                log.Error().Msgf("error requesting disk %s",err)
            } else {
                log.Debug().Msgf("disk: %f", disk)
                coll.cached.Put("disk", disk)
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

    mem, err := coll.cached.Get("memory")
    if err != nil {
        log.Error().Msg(err.Error())
        return nil, err
    }
    cpu, err := coll.cached.Get("cpu")
    if err != nil {
        log.Error().Msg(err.Error())
        return nil, err
    }
    disk, err := coll.cached.Get("disk")
    if err != nil {
        log.Error().Msg(err.Error())
        return nil, err
    }

    return &entities.Status{
        Timestamp: time.Now(),
        Mem: mem.Value.(float64),
        CPU: cpu.Value.(float64),
        Disk: disk.Value.(float64)}, nil

    return nil, nil
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


