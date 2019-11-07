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

package provider

import "github.com/nalej/derrors"

// Intermediate struct representing a KV tuple with a
// byte array as key.
type KVTuple struct {
	Key   []byte
	Value []byte
}

type KeyValueProvider interface {

	// Close the database
	// return:
	//  error if any
	Close() derrors.Error

	// Get a value in the bucket with the given key. If the key is not found,
	// the returned value is nil.
	// params:
	//  bucket
	//  key
	// return:
	//  found item, nil if the key was not found
	Get(bucket []byte, key []byte) ([]byte, derrors.Error)

	// Get all the pairs stored in a bucket.
	// params:
	//  bucket
	// return:
	//  array of kv tuples or error if any
	GetAllPairsInBucket(bucket []byte) ([]KVTuple, derrors.Error)

	// Put a new value for a key
	// params:
	//  bucket to be used
	//  key
	//  value
	// return:
	//  error if any
	Put(bucket []byte, key []byte, value []byte) derrors.Error

	// Delete a key from a bucket
	// params:
	//  bucket
	//  key
	// return:
	//  error if any
	Delete(bucket []byte, key []byte) derrors.Error

	// Get all the current buckets
	// return:
	//  Array of byte arrays with the bucket names
	GetBuckets() [][]byte
}
