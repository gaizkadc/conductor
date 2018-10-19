/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package entities

import "time"

// System status representation.
type Status struct {
    Timestamp time.Time `json:"timestamp,omitempty"`
    Mem float64 `json: "mem,omitempty"`
    CPU float64 `json: "cpu,omitempty"`
    Disk float64 `json: "disk,omitempty"`
}
