/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 *
 */

package utils

import "net"

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