/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package entities


type Requirements struct {
    CPU float32 `json:"cpu,omitempty"`
    Memory float32 `json:"mem,omitempty"`
    Disk float32 `json:"disk,omitempty"`
}
