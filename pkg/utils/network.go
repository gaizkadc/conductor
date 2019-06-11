/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package utils

import (
    "fmt"
    "github.com/nalej/deployment-manager/pkg/common"
    "net"
)

// Compute the next consecutive IP address
// params:
//  ip
//  inc number of addresses to increment
// return:
//  the next ip address
func NextIP(ip net.IP, inc uint) net.IP {
    i := ip.To4()
    v := uint(i[0])<<24 + uint(i[1])<<16 + uint(i[2])<<8 + uint(i[3])
    v += inc
    v3 := byte(v & 0xFF)
    v2 := byte((v >> 8) & 0xFF)
    v1 := byte((v >> 16) & 0xFF)
    v0 := byte((v >> 24) & 0xFF)
    return net.IPv4(v0, v1, v2, v3)
}


// Return how would be the VSA of a potential entry.
// params:
//  serviceName
//  organizationId
//  appInstanceId
// return:
//  the fqdn
func GetVSAName(serviceName string, organizationId string, appInstanceId string) string {
    value := fmt.Sprintf("%s-%s-%s.service.nalej", common.FormatName(serviceName), organizationId[0:10],
        appInstanceId[0:10])
    return value
}