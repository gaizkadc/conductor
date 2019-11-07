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

package structures

import (
	"github.com/nalej/conductor/internal/entities"
	"sync"
)

// Interface for a queue storing deployment requests
type RequestsQueue interface {

	// Obtain next deployment request
	//  returns:
	//   next deployment request, nil if nothing is ready
	NextRequest() *entities.DeploymentRequest

	// Check if there are more available requests.
	AvailableRequests() bool

	// Push a request into the queue.
	//  params:
	//   req the requirement to be pushed into.
	//  returns:
	//   error if any
	PushRequest(req *entities.DeploymentRequest) error

	// Remove the entry with the indicated appInstanceId.
	// params:
	//  appInstanceId identifier of the instance to be removed
	// returns:
	//  true if removed
	Remove(appInstanceId string) bool

	// Clear the queue
	Clear()

	// queue length
	Len() int
}

// Basic queue in memory solution.
type MemoryRequestQueue struct {
	// queue for incoming messages
	queue []*entities.DeploymentRequest
	// Mutex for queue operations
	mux sync.RWMutex
}

func NewMemoryRequestQueue() RequestsQueue {
	toReturn := MemoryRequestQueue{queue: make([]*entities.DeploymentRequest, 0)}
	return &toReturn
}

// Thread-safe method to access queued requests
func (q *MemoryRequestQueue) NextRequest() *entities.DeploymentRequest {
	q.mux.Lock()
	defer q.mux.Unlock()
	if len(q.queue) == 0 {
		return nil
	}

	toReturn := q.queue[0]
	if len(q.queue) == 1 {
		q.queue = nil
	} else {
		q.queue = q.queue[1:]
	}

	return toReturn
}

// Thread-safe function to find whether there are more requests available or not.
func (q *MemoryRequestQueue) AvailableRequests() bool {
	q.mux.RLock()
	defer q.mux.RUnlock()
	available := len(q.queue) != 0
	return available
}

// Push a new request to the que for later processing.
//  params:
//   req entry to be enqueued
func (q *MemoryRequestQueue) PushRequest(req *entities.DeploymentRequest) error {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.queue = append(q.queue, req)
	return nil
}

func (q *MemoryRequestQueue) Clear() {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.queue = nil
}

func (q *MemoryRequestQueue) Len() int {
	q.mux.RLock()
	defer q.mux.RUnlock()
	return len(q.queue)
}

// Remove the entry with the indicated appInstanceId.
// params:
//  appInstanceId identifier of the instance to be removed
// returns:
//  true if removed
func (q *MemoryRequestQueue) Remove(appInstanceId string) bool {
	q.mux.Lock()
	defer q.mux.Unlock()
	targetIndex := -1
	for i, r := range q.queue {
		if r.AppInstanceId == appInstanceId {
			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		return false
	}

	q.queue = append(q.queue[:targetIndex], q.queue[targetIndex+1:]...)
	return true
}
