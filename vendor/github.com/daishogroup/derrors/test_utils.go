//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Test utils to avoid external dependencies.

package derrors

import (
    "reflect"
    "testing"
)

// AssertEquals utility function for the tests to avoid external dependencies.
func assertEquals(t *testing.T, expected interface {}, current interface{}, message string) {
    if !reflect.DeepEqual(expected, current) {
        t.Errorf("%s\nExpected: %s, Current: %s", message, expected, current)
    }
}

// AssertTrue utility function for the tests to avoid external dependencies.
func assertTrue(t *testing.T, condition bool, message string) {
    if !condition {
        t.Error(message)
    }
}