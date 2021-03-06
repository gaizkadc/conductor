/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package statuscollector

import (
	"context"
	"time"

	"github.com/nalej/conductor/internal/entities"

	"github.com/nalej/derrors"

	"github.com/nalej/grpc-monitoring-go"

	"github.com/rs/zerolog/log"
)

type MetricsAPICollector struct {
	client grpc_monitoring_go.MetricsCollectorClient

	// Milliseconds to sleep between calls.
	sleepDuration time.Duration
	ticker        *time.Ticker

	// Cached status
	// TODO: Evaluate potential ways to have a more efficient provider.
	cached Cache

	// Information needed to create API request
	organizationId string
	clusterId      string
}

const apiTimeout = time.Second * 15

const (
	cpuKey  = "cpu"
	memKey  = "mem"
	diskKey = "disk"
)

func NewMetricsAPICollector(client grpc_monitoring_go.MetricsCollectorClient, organizationId string, clusterId string, sleepTime uint32) StatusCollector {
	c := &MetricsAPICollector{
		client:         client,
		sleepDuration:  time.Duration(sleepTime) * time.Millisecond,
		cached:         NewSimpleCache(),
		organizationId: organizationId,
		clusterId:      clusterId,
	}

	return c
}

// Start the collector
// return:
//  Error if any
func (coll *MetricsAPICollector) Run() error {
	log.Info().Msg("starting metrics api status collector...")

	err := coll.gatherStats()
	if err != nil {
		log.Warn().Err(err).Msg("error initializing cache with stats; continuing gather loop anyway")
	}

	coll.ticker = time.NewTicker(coll.sleepDuration)

	for {
		select {
		case <-coll.ticker.C:
			err = coll.gatherStats()
			if err != nil {
				log.Warn().Err(err).Msg("error collecting status from metrics api. continuing.")
			}
		}
	}

	return nil
}

func (coll *MetricsAPICollector) gatherStats() error {
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()
	req := &grpc_monitoring_go.ClusterSummaryRequest{
		OrganizationId: coll.organizationId,
		ClusterId:      coll.clusterId,
	}

	summary, err := coll.client.GetClusterSummary(ctx, req)
	if err != nil {
		return err
	}

	log.Debug().Interface("stats", summary).Msg("collected cluster summary")

	coll.cached.Put(cpuKey, summary.GetCpuMillicores())
	coll.cached.Put(memKey, summary.GetMemoryBytes())
	coll.cached.Put(diskKey, summary.GetUsableStorageBytes())

	return nil
}

// Stop the collector.
// return:
//  Error if any
func (coll *MetricsAPICollector) Finalize(killSignal bool) error {
	log.Info().Msg("finalize was called")
	if coll.ticker != nil {
		coll.ticker.Stop()
	}
	return nil
}

// Get the current status.
func (coll *MetricsAPICollector) GetStatus() (*entities.Status, error) {
	// Build the status and return it
	cacheContent := coll.cached.GetAll()
	if cacheContent == nil {
		log.Debug().Msg("no status entry found in cache; trying to retrieve now")
		err := coll.gatherStats()
		if err != nil {
			log.Warn().Err(err).Msg("error gathering stats")
			return nil, derrors.NewNotFoundError("no cache entries found and unable to retrieve stats", err)
		}
		cacheContent = coll.cached.GetAll()
	}

	status := &entities.Status{
		Timestamp: time.Now(),
		MemFree:   float64(cacheContent[memKey].Value.(*grpc_monitoring_go.ClusterStat).GetAvailable()),
		CPUIdle:   float64(cacheContent[cpuKey].Value.(*grpc_monitoring_go.ClusterStat).GetAvailable()),
		DiskFree:  float64(cacheContent[diskKey].Value.(*grpc_monitoring_go.ClusterStat).GetAvailable()),
		CPUNum:    float64(cacheContent[cpuKey].Value.(*grpc_monitoring_go.ClusterStat).GetTotal() / 1000), // millicores to cores
	}

	return status, nil
}

// Return the status collector name.
// return:
//  Name of this collector.
func (coll *MetricsAPICollector) Name() string {
	return "Metrics API status collector"
}

// Return a description of this status collector.
// return:
//  Description of this collector.
func (coll *MetricsAPICollector) Description() string {
	return "Status collector that calls out to the metrics-collector API to gather status information"
}
