//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the configuration of the server

package server

import (
    log "github.com/sirupsen/logrus"
    "strconv"
    "errors"
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
    // Address where the API service will listen requests.
    Port uint16
    // Use in-memory providers
    UseInMemoryProviders bool
    // Use file system providers
    UseFileSystemProvider bool
    FileSystemBasePath string
    // Name of the default admin user
    DefaultAdminUser string
    // Admin default password
    DefaultAdminPassword string
}

// Validate the configuration.
func (conf * Config) Validate() error {
    /*if conf.UseInMemoryProviders && conf.UseFileSystemProvider {
        return errors.New("only one type of provider is allowed")
    }*/
    if conf.UseFileSystemProvider && conf.FileSystemBasePath == ""{
        return errors.New("filesystem path required")
    }
    return nil
}

// Print writes in the log the current API configuration.
func (conf *Config) Print() {
    loggerConfig.Info("HTTP Server Address: [", conf.Port, "]")
    loggerConfig.Info("UseInMemoryProviders: " + strconv.FormatBool(conf.UseInMemoryProviders))
    loggerConfig.Info("UseFileSystemProvider: " + strconv.FormatBool(conf.UseFileSystemProvider) + " basePath: " + conf.FileSystemBasePath)
    loggerConfig.Infof("Default admin user: %s" + conf.DefaultAdminUser)
    loggerConfig.Infof("Default admin password: %s" + conf.DefaultAdminPassword)
}