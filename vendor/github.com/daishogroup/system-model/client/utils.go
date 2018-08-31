//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//

package client

import (
    "net/url"
    "strconv"
)

const DefaultPort = 80

func ParseHostPort(basePath string) (string, int) {
    uri, err := url.Parse(basePath)
    if err != nil {
        panic(err)
    }
    port := uri.Port()
    if port == "" {
        return uri.Hostname(), DefaultPort
    }
    parsedPort, err := strconv.ParseInt(port, 10, 32)
    if err != nil {
        panic(err)
    }
    return uri.Hostname(), int(parsedPort)

}
