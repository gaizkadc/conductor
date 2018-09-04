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

// This file contains the service manager.

package service

import (
	"os"
	"os/signal"
	"syscall"
)

// Launch triggers the execution of the instance wrapping it to capture life-cycle signals.
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
			srv.Finalize(killSignal == os.Kill)
			return nil
		}
	}
	// Never happen, but need to complete code.
	return nil
}
