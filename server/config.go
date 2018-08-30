//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Set of configuration parameters for conductor.

package server

import (
    log "github.com/sirupsen/logrus"
)

// Logger for the Config struct.
var loggerConfig = log.WithField(
    "package", "api",
).WithField(
    "struct", "Config",
).WithField(
    "file", "config.go",
)

// Config struct for the API service.
type Config struct {
    // Port where the API service will listen requests.
    Port uint16
    // Address where the system model API is listening
    SystemModelAddress string
    // Address where the system model API is listening
    LoggerAddress string
}

// Write in the log the current API configuration.
func (conf *Config) Print() {
    loggerConfig.Info("BUDO HTTP Server port: [", conf.Port, "]")
    loggerConfig.Info("System model listening on: [", conf.SystemModelAddress, "]")
    loggerConfig.Info("Logger listening on: [", conf.LoggerAddress, "]")
}
