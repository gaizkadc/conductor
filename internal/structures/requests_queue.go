/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package structures

import (
    "github.com/phf/go-queue/queue"
    "sync"
    "github.com/nalej/conductor/internal/entities"
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

    // Clear the queue
    Clear()

    // queue length
    Len() int
}

// Basic queue in memory solution.
type MemoryRequestQueue struct {
    // queue for incoming messages
    queue *queue.Queue
    // Mutex for queue operations
    mux sync.RWMutex
}

func NewMemoryRequestQueue () RequestsQueue {
    toReturn := MemoryRequestQueue{queue: queue.New()}
    toReturn.queue.Init()
    return &toReturn
}

// Thread-safe method to access queued requests
func(q *MemoryRequestQueue) NextRequest() *entities.DeploymentRequest {
    q.mux.Lock()
    defer q.mux.Unlock()
    toReturn := q.queue.PopFront().(*entities.DeploymentRequest)
    return toReturn
}

// Thread-safe function to find whether there are more requests available or not.
func(q *MemoryRequestQueue) AvailableRequests() bool {
    q.mux.RLock()
    defer q.mux.RUnlock()
    available := q.queue.Len()!=0
    return available
}

// Push a new request to the que for later processing.
//  params:
//   req entry to be enqueued
func (q *MemoryRequestQueue) PushRequest(req *entities.DeploymentRequest) error {
    q.mux.Lock()
    defer q.mux.Unlock()
    q.queue.PushBack(req)
    return nil
}

func (q *MemoryRequestQueue) Clear() {
    q.mux.Lock()
    defer q.mux.Unlock()
    q.queue.Init()
}

func (q *MemoryRequestQueue) Len() int{
    q.mux.Lock()
    defer q.mux.Unlock()
    return q.queue.Len()
}