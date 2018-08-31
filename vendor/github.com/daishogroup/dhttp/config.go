//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Set of utilities to configure a client.
//

package dhttp

import (
    "fmt"
    "net/url"
    "strconv"

    "github.com/daishogroup/derrors"
)

// ApiKeyHeader is the name of the header that contains the service API key.
const ApiKeyHeader = "x-api-key"

// Available preconfigured clients
type RestType int

// Basic HTTP requests
const RestBasic RestType = 0

// HTTP request with BasicAuth headers
const RestBasicAuth RestType = 1

// HTTP request with OAuth authentication
const RestOAuth RestType = 2

// HTTPS client with no additional authentication
const RestBasicHTTPS RestType = 3

// HTTPS client with basic authentication
const RestBasicHTTPSAuth RestType = 4

// HTTPS client with api key
const RestApiKeyHTTPS RestType = 5

// Default ports
const defaultHTTPPort = 80
const defaultHTTPSPort = 443

// Structure used to define the configuration of a rest client in its different flavours.
type RestConfig struct {
    // Type of client configuration we want to use
    ClientType RestType
    // Target host address
    Host string
    // Target host port
    Port int
    // Security enabled
    HTTPS bool
    // Skip cerficate verification
    SkipVerification bool
    // User account
    User *string
    // Password account
    Password *string
    // OAuth secret
    Secret *string
    // Extra headers to be set
    Headers map[string]string
}

// Get a basic client configuration item.
//  params:
//   host Target machine host.
//   port Port number
//  return:
//   Filled configuration item.
func NewRestBasicConfig(host string, port int) RestConfig {
    return RestConfig{
        ClientType:       RestBasic,
        Host:             host,
        Port:             port,
        HTTPS:            false,
        SkipVerification: true,
        User:             nil,
        Password:         nil,
        Secret:           nil,
    }
}

// Get a basic authentication client configuration item.
//  params:
//   host Target machine host.
//   port Port number
//  return:
//   Filled configuration item.
func NewRestBasicAuthConfig(host string, port int, user string, password string) RestConfig {
    return RestConfig{
        ClientType:       RestBasicAuth,
        Host:             host,
        Port:             port,
        HTTPS:            false,
        SkipVerification: true,
        User:             &user,
        Password:         &password,
        Secret:           nil,
    }
}

// Get a basic authentication client configuration item.
//  params:
//   host Target machine host.
//   port Port number
//  return:
//   Filled configuration item.
func NewRestOAuthConfig(host string, port int, secret string) RestConfig {
    return RestConfig{
        ClientType:       RestOAuth,
        Host:             host,
        Port:             port,
        HTTPS:            true,
        SkipVerification: true,
        User:             nil,
        Password:         nil,
        Secret:           &secret,
    }
}

// Get a basic client configuration item using HTTPS connection.
//  params:
//   host Target machine host.
//   port Port number
//  return:
//   Filled configuration item.
func NewRestBasicHTTPS(host string, port int) RestConfig {
    return RestConfig{
        ClientType:       RestBasicHTTPS,
        Host:             host,
        Port:             port,
        HTTPS:            true,
        SkipVerification: true,
        User:             nil,
        Password:         nil,
        Secret:           nil,
    }
}

// Get a basic client configuration item using HTTPS connection and basic authorization.
//  params:
//   host Target machine host.
//   port Port number
//   user Username
//   password Password
//  return:
//   Filled configuration item.
func NewRestBasicAuthorizationHTTPS(host string, port int, user string, password string) RestConfig {
    return RestConfig{
        ClientType:       RestBasicHTTPSAuth,
        Host:             host,
        Port:             port,
        HTTPS:            true,
        SkipVerification: true,
        User:             &user,
        Password:         &password,
        Secret:           nil,
    }
}

// NewRestApiKeyHTTPS create a client with a HTTPS basic configuration using an API key.
func NewRestApiKeyHTTPS(host string, port int, apiKey string) RestConfig {
    return RestConfig{
        ClientType:       RestApiKeyHTTPS,
        Host:             host,
        Port:             port,
        HTTPS:            true,
        SkipVerification: true,
        User:             nil,
        Password:         nil,
        Secret:           nil,
        Headers:          map[string]string{ApiKeyHeader: apiKey},
    }
}

// Get a client configuration by parsing a URL
//  params:
//   url URL string. Path is being ignored.
//  return:
//   Filled configuration item.
func NewRestURLConfig(urlString string) (RestConfig, derrors.DaishoError) {
    config := RestConfig{}

    // Parse the URL
    parsedUrl, err := url.Parse(urlString)
    if err != nil {
        return config, derrors.NewConnectionError(fmt.Sprintf("Cannot parse url %s", urlString), err)
    }

    // Set type
    switch parsedUrl.Scheme {
    case "http":
        config.ClientType = RestBasic
        config.Port = defaultHTTPPort
    case "https":
        config.ClientType = RestBasicHTTPS
        config.Port = defaultHTTPSPort
        config.HTTPS = true
        config.SkipVerification = true
    default:
        return config, derrors.NewConnectionError(fmt.Sprintf("Invalid scheme: %s", parsedUrl.Scheme))
    }

    // Check if username/password is provided
    if parsedUrl.User != nil {
        config.ClientType = RestBasicAuth
        username := parsedUrl.User.Username()
        config.User = &username

        password, set := parsedUrl.User.Password()
        if set {
            config.Password = &password
        }
    }

    // Set hostname and port
    config.Host = parsedUrl.Hostname()
    // If no port or unparsable port - we've already set default above
    if parsedUrl.Port() != "" {
        port, err := strconv.Atoi(parsedUrl.Port())
        if err == nil {
            config.Port = port
        }
    }

    return config, nil
}
