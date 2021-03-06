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

package statuscollector

import "github.com/nalej/conductor/internal/entities"

type FakeCollector struct {
	// Single observed status
	Status *entities.Status
}

func NewFakeCollector() StatusCollector {
	return &FakeCollector{Status: nil}
}

func (c *FakeCollector) SetStatus(status entities.Status) {
	c.Status = &status
}

func (c *FakeCollector) Run() error {
	return nil
}

func (c *FakeCollector) Finalize(killSignal bool) error {
	return nil
}

func (c *FakeCollector) GetStatus() (*entities.Status, error) {
	if c.Status == nil {
		// No status was set, return the basic one.
		return &entities.Status{CPUNum: 0.1, MemFree: 0.2, DiskFree: 0.3}, nil
	}

	return c.Status, nil
}

// Return the status collector name.
// return:
//  Name of this collector.
func (c *FakeCollector) Name() string {
	return "fake collector"
}

// Return a description of this status collector.
// return:
//  Description of this collector.
func (c *FakeCollector) Description() string {
	return "a fake collector"
}
