//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//

package logger

import (
    log "github.com/sirupsen/logrus"
    "github.com/nalej/conductor/entities"
    "github.com/daishogroup/dhttp"
    "github.com/daishogroup/derrors"
    "fmt"
    "strings"
    "net/url"
    "strconv"
)

var logger = log.WithField("package", "logger_client")

const LogsEndpoint = "/v1/log?query=%s"

type Client interface {
    Logs(pods [] string) (*entities.LogEntries, derrors.DaishoError)
}

type MockupClient struct {
    logs [] string
}

func NewMockupClient(logs [] string) Client {
    return &MockupClient{logs: logs}
}

func (client *MockupClient) Logs(pods [] string) (*entities.LogEntries, derrors.DaishoError) {
    return entities.NewLogEntries(client.logs), nil
}

type RestClient struct {
    client dhttp.Client
}

// Constructor of an REST client to connect with an appmgr instance
// params:
//   basePath host address
// return:
//   instance of the rest client
func NewRestClient(basePath string) Client {
    logger.Debugf("Create client pointing at %s", basePath)
    u, err := url.Parse(basePath)
    if err != nil {
        logger.Errorf("Not valid URL [%s]",basePath)
    }
    // TODO check this is a valid port sequence
    port, _ := strconv.Atoi(u.Port())
    conf := dhttp.NewRestBasicConfig(u.Host,port)
    rest := dhttp.NewClientSling(conf)
    return &RestClient{rest}
}

// Logs return the logs from the logger aggregator in Kubernetes.
// params:
//   pods set of kubernetes pods
// return:
//   the log entries
func (client *RestClient) Logs(pods [] string) (*entities.LogEntries, derrors.DaishoError) {
    // In this code we have to extend the sling rest client with a delete with body operation
    logger.Debug("Called get Logs: ", pods)
    // Create a stop request
    if client.client == nil {
        logger.Error("impossible to get logs. The client is null.")
        return nil, derrors.NewGenericError("impossible get logs. The client is null")
    }

    podVariables := make([] string, 0)
    for _, pod := range pods {
        podVariables = append(podVariables, fmt.Sprintf("pod:%s", pod))
    }

    targetUrl := fmt.Sprintf(LogsEndpoint, strings.Join(podVariables, " and "))

    logger.Debug("Send request to logger service: " , targetUrl)

    response := client.client.GetRaw(targetUrl)

    if !response.Ok() {
        return nil, derrors.NewOperationError("error aggregating logs",response.Error)
    }

    result := response.Result.(string)

    logger.Debug(result)
    if result == "null" {
        return entities.NewLogEntries([]string{}), nil
    }

    return ParseLogEntries(result)
}
