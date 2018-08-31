//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Serialization tests

package derrors

import (
    "encoding/json"
    "errors"
    "testing"
)

func TestFromJsonGenericError(t *testing.T) {
    cause := errors.New("error cause")
    msg := "Error message"
    toSend := NewGenericError(msg, cause)

    data, err := json.Marshal(toSend)
    assertEquals(t, nil, err, "expecting no error")
    retrieved, err := FromJSON(data)

    assertEquals(t, nil, err, "message should be deserialized")
    //assertEquals(t, GenericErrorType, retrieved.Type(), "type mismatch")
    assertEquals(t, toSend, retrieved, "structure should match")
}

func TestFromJsonEntityError(t *testing.T){

    entity := "serializableEntity"
    cause := errors.New("another cause")
    msg := "Other message"
    toSend := NewEntityError(entity, msg, cause)

    data, err := json.Marshal(toSend)
    assertEquals(t, nil, err, "expecting no error")
    retrieved, err := FromJSON(data)

    assertEquals(t, nil, err, "message should be deserialized")
    //assertEquals(t, GenericErrorType, retrieved.Type(), "type mismatch")
    assertEquals(t, toSend, retrieved, "structure should match")

}

func TestFromJsonConnectionError(t *testing.T){

    URL := "http://url-that-fails.com"
    cause := errors.New("yet another cause")
    msg := "Yet another message"
    toSend := NewConnectionError(msg, cause).WithParams(URL)

    data, err := json.Marshal(toSend)
    assertEquals(t, nil, err, "expecting no error")
    retrieved, err := FromJSON(data)

    assertEquals(t, nil, err, "message should be deserialized")
    //assertEquals(t, GenericErrorType, retrieved.Type(), "type mismatch")
    assertEquals(t, toSend, retrieved, "structure should match")

}

func TestFromJsonOperationError(t *testing.T){

    param1 := "param1"
    param2 := "param2"
    param3 := "param3"
    cause := errors.New("operation failed")
    msg := "operation failure"
    toSend := NewOperationError(msg, cause).WithParams(param1, param2, param3)

    data, err := json.Marshal(toSend)
    assertEquals(t, nil, err, "expecting no error")
    retrieved, err := FromJSON(data)

    assertEquals(t, nil, err, "message should be deserialized")
    //assertEquals(t, GenericErrorType, retrieved.Type(), "type mismatch")
    assertEquals(t, toSend, retrieved, "structure should match")

}