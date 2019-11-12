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

import "time"

type CacheEntry struct {
	// timestamp when this value was updated
	TimeStamp time.Time
	Value     interface{}
}

// Generic interface for caching status values.
type Cache interface {

	// Put new entries into the cache identified by a unique key.
	//  params:
	//   key Unique key
	//   value Value to be stored
	Put(key string, value interface{})

	// Get an entry value identified by the key.
	//  params:
	//   key Unique key
	//  return:
	//   cache entry
	Get(key string) (*CacheEntry, error)

	// Get all the cached entries
	// return:
	//  map with the cached entries
	GetAll() map[string]CacheEntry
}
