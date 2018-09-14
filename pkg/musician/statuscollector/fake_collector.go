/*
 * Copyright 2018 Nalej
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

package statuscollector

import "github.com/nalej/conductor/internal/entities"

type FakeCollector struct {

}

func NewFakeCollector() StatusCollector {
    return &FakeCollector{}
}

func(c *FakeCollector) Run() error {
    return nil
}

func (c *FakeCollector) Finalize(killSignal bool) error {
    return nil
}

func (c *FakeCollector) GetStatus() (*entities.Status, error) {
    toReturn := entities.Status{CPU:0.1, Mem: 0.2, Disk: 0.3}
    return &toReturn, nil
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
