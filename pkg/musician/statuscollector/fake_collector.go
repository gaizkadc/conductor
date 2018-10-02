/*
 * Copyright (C) 2018 Nalej Group -All Rights Reserved
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

func(c *FakeCollector) SetStatus(status entities.Status) {
    c.Status = &status
}

func(c *FakeCollector) Run() error {
    return nil
}

func (c *FakeCollector) Finalize(killSignal bool) error {
    return nil
}

func (c *FakeCollector) GetStatus() (*entities.Status, error) {
    if c.Status == nil {
        // No status was set, return the basic one.
        return &entities.Status{CPU:0.1, Mem: 0.2, Disk: 0.3}, nil
    }

    return c.Status,nil
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
