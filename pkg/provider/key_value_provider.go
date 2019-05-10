/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package provider

import "github.com/nalej/derrors"

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

}
