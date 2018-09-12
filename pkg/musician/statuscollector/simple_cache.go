//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//
// Simple cache implementation using a map with string ids.

package statuscollector

import "time"

type SimpleCache struct {
    pool map[string] CacheEntry
}

func NewSimpleCache() *SimpleCache {
    return &SimpleCache{make(map[string] CacheEntry,0)}
}


// Put new entries into the cache identified by a unique key.
//  params:
//   key Unique key
//   value Value to be stored
func(c *SimpleCache) Put (key string, value interface{}) {
    c.pool[key] = CacheEntry{time.Now(), value}
}

// Get an entry value identified by the key.
//  params:
//   key Unique key
//  return:
//   stored interface or error if not found
func(c *SimpleCache) Get (key string) (*CacheEntry, error) {
    res, found := c.pool[key]
    if !found {
        return nil, nil
    }
    return &res, nil

}