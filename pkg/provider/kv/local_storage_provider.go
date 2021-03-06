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

package kv

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/nalej/conductor/pkg/provider"
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

// Implementation of a local provider key/value solution using bolt.
// The current implementation is a generic wrapper

type LocalDB struct {
	db *bolt.DB
	mu sync.Mutex
}

// Create a new local db using the given file path
// params:
//  filePath where the database is physically stored
// return:
//  database instance or error if any
func NewLocalDB(filePath string) (provider.KeyValueProvider, derrors.Error) {
	db, err := bolt.Open(filePath, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		log.Error().Err(err).Msg("impossible create local database")
		return nil, derrors.NewInternalError("impossible to create local database", err)
	}
	return &LocalDB{db: db}, nil
}

func (ldb *LocalDB) Close() derrors.Error {
	err := ldb.db.Close()
	if err != nil {
		return derrors.NewInternalError("impossible to close database", err)
	}
	return nil
}

func (ldb *LocalDB) Get(bucket []byte, key []byte) ([]byte, derrors.Error) {
	ldb.mu.Lock()
	defer ldb.mu.Unlock()
	var result []byte
	err := ldb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return errors.New(fmt.Sprintf("bucket %s not found", bucket))
		}
		result = b.Get(key)
		return nil
	})
	if err != nil {
		e := derrors.NewInternalError(fmt.Sprintf("error getting value from key "+
			"%s in bucket %s", key, bucket))
		return nil, e
	}
	return result, nil
}

func (ldb *LocalDB) Put(bucket []byte, key []byte, value []byte) derrors.Error {
	ldb.mu.Lock()
	defer ldb.mu.Unlock()
	err := ldb.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		err = b.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return derrors.NewInternalError("error setting db value", err)
	}
	return nil
}

func (ldb *LocalDB) Delete(bucket []byte, key []byte) derrors.Error {
	ldb.mu.Lock()
	defer ldb.mu.Unlock()
	err := ldb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return errors.New(fmt.Sprintf("bucket %s not found", bucket))
		}
		// delete this key from the bucket
		return b.Delete(key)
	})
	if err != nil {
		return derrors.NewInternalError(fmt.Sprintf("impossible to delete key %s from bucket %s", key, bucket), err)
	}
	return nil
}

func (ldb *LocalDB) GetBuckets() [][]byte {
	ldb.mu.Lock()
	defer ldb.mu.Unlock()
	listBuckets := make([][]byte, 0)
	ldb.db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(bucketName []byte, b *bolt.Bucket) error {
			listBuckets = append(listBuckets, bucketName)
			return nil
		})
		return nil
	})
	return listBuckets
}

func (ldb *LocalDB) GetAllPairsInBucket(bucket []byte) ([]provider.KVTuple, derrors.Error) {
	ldb.mu.Lock()
	defer ldb.mu.Unlock()
	result := make([]provider.KVTuple, 0)

	err := ldb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return errors.New(fmt.Sprintf("bucket %s not found", bucket))
		}
		b.ForEach(func(k, v []byte) error {
			result = append(result, provider.KVTuple{k, v})
			return nil
		})

		return nil
	})
	if err != nil {
		e := derrors.NewInternalError(fmt.Sprintf("error getting pairs from bucket %s", bucket))
		return nil, e
	}
	return result, nil
}
