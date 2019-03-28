/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package entities

import "time"

// System Status representation.
type Status struct {
    Timestamp time.Time `json:"timestamp,omitempty"`
    MemFree   float64   `json: "mem_free,omitempty"`
    CPUNum    float64   `json: "cpu_num,omitempty"`
    CPUIdle   float64   `json: "cpu_idle,omitempty"`
    DiskFree  float64   `json: "disk_free,omitempty"`
}
