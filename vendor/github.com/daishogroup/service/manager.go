//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the service manager.

package service

import (
    "os"
    "os/signal"
    "syscall"
)

// Calling this method the manager can launch an service and control the entire life-cycle.
func Launch(srv Instance) error {

    // Set up channel on which to send signal notifications.
    // We must use a buffered channel or risk missing the signal
    // if we're not ready to receive when the signal is sent.
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

    err := srv.Run()
    if err != nil {
        return err
    }

    // Loop work cycle with accept interrupt by system signal
    for {
        select {
        case killSignal := <-interrupt:
            srv.Finalize(killSignal == os.Kill )
            return nil
        }
    }
    // Never happen, but need to complete code.
    return nil
}