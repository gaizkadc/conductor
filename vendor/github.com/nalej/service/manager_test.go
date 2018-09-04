/*
 * Copyright 2018 Daisho
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
 */

package service

import (
	"errors"
	"testing"
)

type TestService struct {
	value error
}

//Name get the name of the service.
func (s *TestService) Name() string {
	return "TestService"
}

//Description get a short description of the service proposal.
func (s *TestService) Description() string {
	return "Test description"
}

// Run is the start method is called when the application is initialized.
// This method call is expected to return, so a new go routine should be launched if necessary.
//   returns:
//     An error if the instance cannot be executed.
func (s *TestService) Run() error {
	return s.value
}

// Finalize is called when the application is shutting down.
// The Wrapper assumes that this method will return fairly quickly.
//   params:
//     killSignal It is true when the process is killed by the system.
func (s *TestService) Finalize(killSignal bool) {

}

func TestLaunchWithError(t *testing.T) {
	service := &TestService{errors.New("test error")}
	result := Launch(service)
	if result == nil {
		t.Fail()
	}
}
