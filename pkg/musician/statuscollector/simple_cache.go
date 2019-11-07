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

// Simple cache implementation using a map with string ids.

package statuscollector

import (
	"sync"
	"time"
)

type SimpleCache struct {
	pool map[string]CacheEntry
	mux  sync.Mutex
}

func NewSimpleCache() *SimpleCache {
	return &SimpleCache{pool: make(map[string]CacheEntry, 0)}
}

// Put new entries into the cache identified by a unique key.
//  params:
//   key Unique key
//   value Value to be stored
func (c *SimpleCache) Put(key string, value interface{}) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.pool[key] = CacheEntry{time.Now(), value}
}

// Get an entry value identified by the key.
//  params:
//   key Unique key
//  return:
//   stored interface or error if not found
func (c *SimpleCache) Get(key string) (*CacheEntry, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	res, found := c.pool[key]
	if !found {
		return nil, nil
	}
	return &res, nil
}

// Get all the cached entries
// return:
//  map with the cached entries
func (c *SimpleCache) GetAll() map[string]CacheEntry {
	c.mux.Lock()
	defer c.mux.Unlock()

	if len(c.pool) == 0 {
		return nil
	}
	return c.pool
}
