//
// Copyright (C) 2018 Nalej Group - All Rights Reserved
//
// Set of common operations for grpc functions.

package tools

import (
    "github.com/rs/zerolog/log"
    "google.golang.org/grpc/test/bufconn"
    "google.golang.org/grpc"
    "time"
    "context"
    "net"
)

const bufSize = 1024*1024



func LaunchServer(server *grpc.Server, listener *bufconn.Listener) {
    go func() {
        if err := server.Serve(listener); err != nil {
            log.Fatal().Errs("failed to listen: %v", []error{err})
        }
    }()
}


func GetConn (listener bufconn.Listener) (*grpc.ClientConn, error){
    ctx := context.Background()
    conn, err := grpc.DialContext(ctx, "bufnet",
        grpc.WithDialer(func(string, time.Duration)(net.Conn, error){
            return listener.Dial()
        }), grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    return conn, nil
}
