//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Definition of a port entity for applications

package entities

// ApplicationPort structure that defines a port exposed by an application.
type ApplicationPort struct {
    // Name of the port to disambiguate multiple endpoints.
    Name string `json:"name"`
    // Protocol of the open port: TCP or UDP.
    Protocol string `json:"protocol"`
    // Port is the abstracted Service port, which can be any port other pods use to access the Service.
    Port int32 `json:"port"`
    // TargetPort is the port the container accepts traffic on; by default contains the same value as
    // Port but can be a string literal name.
    TargetPort string `json:"targetPort"`
    // NodePort contains the port exposed publicly on the cluster IP address.
    NodePort int32 `json:"nodePort"`
}

// NewApplicationPort creates a new application port description.
func NewApplicationPort(name string, protocol string, port int32, targetPort string, nodePort int32) * ApplicationPort {
    return &ApplicationPort{name, protocol, port, targetPort, nodePort}
}