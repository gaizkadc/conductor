/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */


package statuscollector

import "time"

type CacheEntry struct {
    // timestamp when this value was updated
    TimeStamp time.Time
    Value interface{}
}

// Generic interface for caching status values.
type Cache interface {

    // Put new entries into the cache identified by a unique key.
    //  params:
    //   key Unique key
    //   value Value to be stored
    Put (key string, value interface{})

    // Get an entry value identified by the key.
    //  params:
    //   key Unique key
    //  return:
    //   cache entry
    Get (key string) (*CacheEntry, error)
}
