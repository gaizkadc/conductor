/*
 * Copyright 2018 Nalej
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
 */

package statuscollector

import "github.com/nalej/conductor/internal/entities"

// Interface to be fulfilled by any StatusCollector implementation. In a few words a status collector is a service
// running in the background collecting status information from the cluster where it was deployed. This is done
// repetitively waiting between every check. When the service is queried. It returns the latest observation known.
type  StatusCollector interface {

    // Start the collector
    // return:
    //  Error if any
    Run() error

    // Stop the collector.
    // return:
    //  Error if any
    Finalize(killSignal bool) error

    // Get the current status.
    // return:
    //  Current status of the cluster.
    GetStatus() (*entities.Status, error)

    // Return the status collector name.
    // return:
    //  Name of this collector.
    Name() string

    // Return a description of this status collector.
    // return:
    //  Description of this collector.
    Description() string

}
