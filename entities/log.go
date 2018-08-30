//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//

package entities

type LogEntries struct {
    Entries [] string `json:"entries"`
}

func NewLogEntries(entries [] string) *LogEntries {
    return &LogEntries{Entries: entries}
}
