//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//
// Simple cache implementation using a map with string ids.

package statuscollector

import (
    "time"
    "sync"
)

type SimpleCache struct {
    pool map[string] CacheEntry
    sync.RWMutex
}

func NewSimpleCache() *SimpleCache {
    return &SimpleCache{pool:make(map[string] CacheEntry,0)}
}


// Put new entries into the cache identified by a unique key.
//  params:
//   key Unique key
//   value Value to be stored
func(c *SimpleCache) Put (key string, value interface{}) {
    c.Lock()
    c.pool[key] = CacheEntry{time.Now(), value}
    c.Unlock()
}

// Get an entry value identified by the key.
//  params:
//   key Unique key
//  return:
//   stored interface or error if not found
func(c *SimpleCache) Get (key string) (*CacheEntry, error) {
    c.RLock()
    res, found := c.pool[key]
    c.RUnlock()
    if !found {
        return nil, nil
    }
    return &res, nil

}