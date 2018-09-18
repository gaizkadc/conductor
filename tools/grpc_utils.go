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
    "sync"
)

const BufSize = 1024*1024

var (
    // Testing listener to be used in testing environments. It uses a singleton pattern to be instantiated
    // only once if required.
    testListener *bufconn.Listener
    once sync.Once
)



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

func GetDefaultListener() *bufconn.Listener {
    once.Do(func(){
        testListener = bufconn.Listen(BufSize)
    })
    return testListener
}
