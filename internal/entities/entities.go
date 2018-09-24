/*
 * Copyright 2018 Nalej
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
 */

package entities

import "time"

// System status representation.
type Status struct {
    Timestamp time.Time `json:"timestamp"`
    Mem float64 `json: "mem"`
    CPU float64 `json: "cpu"`
    Disk float64 `json: "disk"`
}


type Requirements struct {
    RequestID string
    CPU float32
    Memory float32
    Disk float32
}

// Representation of the score for a potential deployment candidate.
type ClusterScore struct {
    RequestID string
    ClusterID string
    Score float32
    TotalEvaluated int
}
