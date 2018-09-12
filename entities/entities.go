//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

package entities

import "time"

// System status representation.
type Status struct {
    Timestamp time.Time `json:"timestamp"`
    Mem float64 `json: "mem"`
    CPU float64 `json: "cpu"`
    Disk float64 `json: "disk"`
}