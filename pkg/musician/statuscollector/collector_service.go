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

// Service implementation for the status collector.

package statuscollector

import "github.com/nalej/service"

type Service struct {
	StatusCollector StatusCollector
}

//Name get the name of the service.
func (s *Service) Name() string {
	return "Collector service"
}

//Description get a short description of the service proposal.
func (s *Service) Description() string {
	return "Status collector service to collect the current status of the cluster."
}

// Run is the start method is called when the application is initialized.
// This method call is expected to return, so a new go routine should be launched if necessary.
//   returns:
//     An error if the instance cannot be executed.
func (s *Service) Run() error {
	go s.StatusCollector.Run()
	return nil
}

// Finalize is called when the application is shutting down.
// The Wrapper assumes that this method will return fairly quickly.
//   params:
//     killSignal It is true when the process is killed by the system.
func (s *Service) Finalize(killSignal bool) {
	println("Status collector service finished!!")
}

// Main is the service entry point.
func main() {
	srv := &Service{}
	service.Launch(srv)
}
