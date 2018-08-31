//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//

package service

import (
    "testing"
    "errors"
)

type TestService struct {
    value error
}

//Name get the name of the service.
func (s *TestService) Name() string{
    return "TestService"
}

//Description get a short description of the service proposal.
func (s *TestService) Description() string{
    return "Test desrciption"
}

// Run is the start method is called when the application is initialized.
// This method call is expected to return, so a new go routine should be launched if necessary.
//   returns:
//     An error if the instance cannot be executed.
func (s *TestService) Run() error{
    return s.value
}

// Finalize is called when the application is shutting down.
// The Wrapper assumes that this method will return fairly quickly.
//   params:
//     killSignal It is true when the process is killed by the system.
func (s *TestService) Finalize(killSignal bool){

}


func TestLaunchWithError(t *testing.T) {
    service := &TestService{errors.New("test error")}
    result:=Launch(service)
    if result==nil {
        t.Fail()
    }
}

