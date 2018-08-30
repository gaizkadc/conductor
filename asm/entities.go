//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//

package asm

//Get the pods created by an application
type AppPodsRequest struct {
    //Command specific parameters here
    Name      string      `json:"-" uriParameter:"name" required:"true"`
}

func NewAppPodsRequest() *AppPodsRequest {
    return &AppPodsRequest{
        Name : "",
    }
}

func (request *AppPodsRequest) IsValid() bool {
    return request.Name != ""
}

type AppPodsResponse struct {
    Info           string              `json:"info, omitempty"`
    Pods          []string             `json:"pods"`
}

type DaishoApplicationInfo struct {
    Name          string `json:"name, omitempty"`
    App           string `json:"application, omitempty"`
    Version       string `json:"version, omitempty"`
    UpdatedTime   string `json:"updated, omitempty"`
    Status        string `json:"status, omitempty"`
    Revision      int32  `json:"revision, omitempty"`
    Namespace     string `json:"namespace, omitempty"`

    Labels        map[string] string    `json:"labels" omitempty`
    Services      []ServicesInformation `json:"services" omitempty`
}

type ServicePort struct {
    Name   string            `json:"name, omitempty"`
    Protocol string          `json:"protocol, omitempty"`
    Port     int32           `json:"port, omitempty"`
    TargetPort string        `json:"targetPort, omitempty"`
    NodePort    int32        `json:"nodePort, omitempty"`
}

type ServicesInformation struct {
    Namespace   string             `json:"namespace, omitempty"`
    App         string             `json:"app, omitempty"`
    Name        string             `json:"name, omitempty"`

    Type        string             `json:"type, omitempty"`
    ClusterIp   string             `json:"clusterIp, omitempty"`
    ExternalIps []string           `json:"externalIps, omitempty"`
    Ports       []ServicePort      `json:"ports, omitempty"`
}

// ListReleasesResponse is a list of releases.
type AppListResponse struct {
    Info string `json:"info, omitempty"`
    // Count is the expected total number of releases to be returned.
    Count int64 `json:"count"`
    // Next is the name of the next release. If this is other than an empty
    // string, it means there are more results.
    Next string `json:"next"`
    // Total is the total number of queryable applications.
    Total int64 `json:"total"`
    // Releases is the list of found release objects.
    Applications []DaishoApplicationInfo `json:"applications"`
}

func NewAppListResponse(info string, count int64, next string, total int64, apps []DaishoApplicationInfo) * AppListResponse {
    return &AppListResponse{info, count, next, total, apps}
}

// List applications
type AppListRequest struct {
    Limit      int       `json:"limit"`
}

func NewAppListRequest () *AppListRequest {
    request := & AppListRequest {
        Limit:      100,
    }
    return request
}
