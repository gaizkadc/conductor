/*
 * Copyright 2019 Nalej
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
 *
 */

package utils

import (
	"fmt"
	"net"
	"strings"
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
	value := fmt.Sprintf("%s-%s-%s", formatName(serviceName), organizationId[0:10],
		appInstanceId[0:10])
	return value
}


// Format a string removing white spaces and going lowercase
func formatName(name string) string {
	aux := strings.ToLower(name)
	// replace any space
	aux = strings.Replace(aux, " ", "", -1)
	return aux
}