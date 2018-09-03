//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//

//!!!!! This is a temporal folder until we decide how to proceed with entities

package entities

import (
    "encoding/json"
    "fmt"
    "strconv"
    "math"
)

// A SampleValue is a representation of a value for a given sample at a given
// time.
type SampleValue float64

// MarshalJSON implements json.Marshaler.
func (v SampleValue) MarshalJSON() ([]byte, error) {
    return json.Marshal(v.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (v *SampleValue) UnmarshalJSON(b []byte) error {
    if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
        return fmt.Errorf("sample value must be a quoted string")
    }
    f, err := strconv.ParseFloat(string(b[1:len(b)-1]), 64)
    if err != nil {
        return err
    }
    *v = SampleValue(f)
    return nil
}

// Equal returns true if the value of v and o is equal or if both are NaN. Note
// that v==o is false if both are NaN. If you want the conventional float
// behavior, use == to compare two SampleValues.
func (v SampleValue) Equal(o SampleValue) bool {
    if v == o {
        return true
    }
    return math.IsNaN(float64(v)) && math.IsNaN(float64(o))
}

func (v SampleValue) String() string {
    return strconv.FormatFloat(float64(v), 'f', -1, 64)
}

// SamplePair pairs a SampleValue with a Timestamp.
type SamplePair struct {
    Timestamp Time
    Value     SampleValue
}

// MarshalJSON implements json.Marshaler.
func (s SamplePair) MarshalJSON() ([]byte, error) {
    t, err := json.Marshal(s.Timestamp)
    if err != nil {
        return nil, err
    }
    v, err := json.Marshal(s.Value)
    if err != nil {
        return nil, err
    }
    return []byte(fmt.Sprintf("[%s,%s]", t, v)), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SamplePair) UnmarshalJSON(b []byte) error {
    v := [...]json.Unmarshaler{&s.Timestamp, &s.Value}
    return json.Unmarshal(b, &v)
}

func (s SamplePair) String() string {
    return fmt.Sprintf("%s @[%s]", s.Value, s.Timestamp)
}


type PrometheusMemoryStatus struct {
    Data struct {
        Result []struct {
            Metric struct {
                Name                   string `json:"__name__, omitempty"`
                ControllerRevisionHash string `json:"controller_revision_hash, omitempty"`
                Daemon                 string `json:"daemon, omitempty"`
                Grafanak8Sapp          string `json:"grafanak8sapp, omitempty"`
                Instance               string `json:"instance, omitempty"`
                Job                    string `json:"job, omitempty"`
                K8SApp                 string `json:"k8s_app, omitempty"`
                KubernetesNamespace    string `json:"kubernetes_namespace, omitempty"`
                KubernetesPodName      string `json:"kubernetes_pod_name, omitempty"`
                Nodename               string `json:"nodename, omitempty"`
                PodTemplateGeneration  string `json:"pod_template_generation, omitempty"`
            } `json:"metric"`
            //Value SamplePair `json:"value"`
            Value []interface{} `json:"value"`
        } `json:"result"`
        ResultType string `json:"resultType"`
    } `json:"data"`
    Status string `json:"status"`
}