//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//

package main

import (
    "time"
    "github.com/daishogroup/service"
)

// Service is an example application.
type Service struct {
}

//Name get the name of the service.
func (s *Service) Name() string {
    return "Example Service"
}

//Description get a short description of the service proposal.
func (s *Service) Description() string {
    return "Example description"
}

// Run is the start method is called when the application is initialized.
// This method call is expected to return, so a new go routine should be launched if necessary.
//   returns:
//     An error if the instance cannot be executed.
func (s *Service) Run() error {
    go s.doRun()
    return nil
}

func (s *Service) doRun() {
    for {
        time.Sleep(time.Second)
        println("Hello world!!")
    }
}

// Finalize is called when the application is shutting down.
// The Wrapper assumes that this method will return fairly quickly.
//   params:
//     killSignal It is true when the process is killed by the system.
func (s *Service) Finalize(killSignal bool) {
    println("GOODBYE!!")
}

// Main is the service entry point.
func main() {
    srv := &Service{}
    service.Launch(srv)
}
