//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the service interface.

package service

// Instance is the interface that the application must implement this interface to receive event in the different stages
// of the process life-cycle.
type Instance interface {
    //Name get the name of the service.
    Name() string

    //Description get a short description of the service proposal.
    Description() string

    // Run is the start method is called when the application is initialized.
    // This method call is expected to return, so a new go routine should be launched if necessary.
    //   returns:
    //     An error if the instance cannot be executed.
    Run() error

    // Finalize is called when the application is shutting down.
    // The Wrapper assumes that this method will return fairly quickly.
    //   params:
    //     killSignal It is true when the process is killed by the system.
    Finalize(killSignal bool)

}