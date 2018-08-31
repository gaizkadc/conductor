//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains a parser for gravitational log-forwarder messages. By looking into the code available
// on https://github.com/gravitational/logging-app/blob/3402830ee1d1d91f1b94d04cd21fa2237b6689a8/cmd/wstail/tail.go
// it seems that we are receiving a json array of string, where each string corresponds to a json object.
//

package logger

import (
    "encoding/json"
    "fmt"

    "github.com/nalej/conductor/entities"
    "github.com/nalej/conductor/errors"
    "github.com/daishogroup/derrors"
)

type RawLogResponse struct {
    SEntries [] string
}

type GenericLogEntry struct {
    Type    string `json:"type"`
    Payload string `json:"payload"`
}

func (gle *GenericLogEntry) ToString() string {
    // It seems that all are marked as data.
    if gle.Type == "data" {
        return fmt.Sprintf("%s", gle.Payload)
    }
    return fmt.Sprintf("%s - %s", gle.Type, gle.Payload)
}

type LogResponse struct {
    Entries [] GenericLogEntry `json:"entries"`
}

func NewLogReponse(entries [] GenericLogEntry) *LogResponse {
    return &LogResponse{entries}
}

func (lr *LogResponse) ToLogEntries() *entities.LogEntries {
    messages := make([] string, 0)
    for _, entry := range lr.Entries {
        messages = append(messages, entry.ToString())
    }
    return entities.NewLogEntries(messages)
}

func ParseLogResponse(payload string) (*LogResponse, derrors.DaishoError) {
    var msgList [] string
    if err := json.Unmarshal([] byte(payload), &msgList); err != nil {
        return nil, derrors.NewOperationError(errors.UnmarshalError, err)
    }
    entries := make([] GenericLogEntry, 0)
    for _, msg := range msgList {
        logEntry, err := ParseLogEntry(msg)
        if err != nil {
            return nil, err
        }
        entries = append(entries, *logEntry)
    }
    result := NewLogReponse(entries)
    return result, nil
}

func ParseLogEntry(entry string) (*GenericLogEntry, derrors.DaishoError) {
    le := &GenericLogEntry{}
    if err := json.Unmarshal([] byte(entry), &le); err != nil {
        return nil, derrors.NewOperationError(errors.UnmarshalError, err)
    }
    return le, nil
}

func ParseLogEntries(entry string) (*entities.LogEntries, derrors.DaishoError) {
    gen, err := ParseLogResponse(entry)
    if err != nil {
        return nil, err
    }
    return gen.ToLogEntries(), nil
}
