//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package entities

import (
    "fmt"
    "time"
)

// Cluster contains a set of nodes that make the computing system where the user can deploy apps on.
type Config struct {
    // LogRetention configures the retention policy for the loggging subsystem.
    LogRetention string `json:"logRetention"`
}

func NewConfig(logRetention string) * Config {
    return &Config{logRetention}
}

func (c *Config) Valid() bool {
    // Validate log retention
    _, err := time.ParseDuration(c.LogRetention)
    return err == nil
}

func (c *Config) String() string {
    return fmt.Sprintf("%#v", c)
}

