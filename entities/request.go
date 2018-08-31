//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//

package entities

import (
    "fmt"
    "regexp"

    "github.com/nalej/conductor/errors"
    "github.com/daishogroup/derrors"
    "github.com/ghodss/yaml"
    log "github.com/sirupsen/logrus"
    "k8s.io/helm/pkg/strvals"

    "github.com/daishogroup/system-model/entities"
)

var logger = log.WithField("package", "entities")

// Request to deploy a new app instance. While this structure may be replicated in the current version with
// the AddAppInstanceRequest, it allows for future differences. If not, we can consolidate both in future versions.
type DeployAppRequest struct {
    // Application descriptor identifier.
    AppDescriptorId string `json:"appDescriptorId, omitempty"`
    // Application instance name.
    Name string `json:"name, omitempty"`
    // Application description.
    Description string `json:"description, omitempty"`
    // TODO Considering just one service at this point.
    // Instance label.
    Label string `json:"label, omitempty"`
    // Labels for the deployment that will be set on every element associated to the app in k8s.
    // Expect those to be daisho.group/networkId, daisho.group/clusterId, etc.
    // Labels are automatically set by conductor.
    Labels map[string]string `json:"labels"`
    // Arguments for launching the app.
    Arguments string `json:"arguments, omitempty"`
    // Persistence size required by the app
    PersistenceSize string `json:"persistentSize"`
    // Storage type
    StorageType entities.AppStorageType `json:"storageType"`
}

// MaxDeployNameLength is the maximum length an application name can have in Helm/Kubernetes.
const MaxDeployNameLength = 63
// AppNameRegex is the regular expression app names must conform to according to Helm/Kubernetes.
const AppNameRegex = "^[a-z]([-a-z0-9]*[a-z0-9])?$"

func NewDeployAppRequest(
    name string,
    appDescriptorId string,
    description string,
    label string,
    labels map[string]string,
    arguments string,
    persistenceSize string,
    storageType entities.AppStorageType) DeployAppRequest {
    return DeployAppRequest{
        appDescriptorId,
        name,
        description,
        label,
        labels,
        arguments,
        persistenceSize,
        storageType}
}

func (deployRequest *DeployAppRequest) String() string {
    return fmt.Sprintf("%#v", deployRequest)
}

// IsValid checks that the application required parameters are present and the names conforms to the underlying
// HELM/Kubernetes specifications.
func (deployRequest *DeployAppRequest) IsValid() derrors.DaishoError {

    if deployRequest.Name == "" || deployRequest.AppDescriptorId == "" {
        return derrors.NewEntityError(deployRequest, errors.MissingRequiredFields)
    }

    if len(deployRequest.Name) > MaxDeployNameLength {
        return derrors.NewEntityError(deployRequest, errors.CannotDeleteApp)
    }

    nameRegex := regexp.MustCompile(AppNameRegex)
    if !nameRegex.MatchString(deployRequest.Name) {
        return derrors.NewEntityError(deployRequest, errors.InvalidAppName)
    }

    return nil
}

// Transform a deploy request into an application instance request.
//   params:
//     request The add application instance request.
//   returns:
//     An instance with an UUID.
/*
func ToAddAppInstance(request DeployAppRequest) * entities.AddAppInstanceRequest {
    return entities.NewAddAppInstanceRequest(
        request.AppDescriptorId,
        request.Name,
        request.Description,
        request.Label,
        request.Arguments,
        request.PersistenceSize,
        request.StorageType,
        PORT????)
}
*/

type AppStartRequest struct {
    //Command specific parameters here
    Name            string `json:"name, omitempty"`
    Namespace       string `json:"namespace, omitempty"`
    ManifestName    string `json:"manifest, omitempty"`
    ManifestVersion string `json:"version, omitempty"`
    RawValue        []byte `json:"rawvalue, omitempty"`
    Labels        map[string] string `json:"labels" omitempty`

    ValueFiles   []string `json:"-"`
    ChartPath    string   `json:"-"`
    DryRun       bool     `json:"-"`
    DisableHooks bool     `json:"-"`
    Replace      bool     `json:"-"`
    Verify       bool     `json:"-"`
    Keyring      string   `json:"-"`
    Values       []string `json:"-"`
    NameTemplate string   `json:"-"`
    Version      string   `json:"-"`
    Timeout      int64    `json:"-"`
    Wait         bool     `json:"-"`
    RepoURL      string   `json:"-"`
    Devel        bool     `json:"-"`
}

func NewAppStartRequest(name string, manifestName string,
    manifestVersion string, args [] string, labels map[string]string) *AppStartRequest {

    // get the settings if any
    values := map[string]interface{}{}
    marshalledSettings := make([]byte, 0)
    if len(args) != 0 {
        problem := false
        for _, n := range args {
            if err := strvals.ParseInto(n, values); err != nil {
                logger.Errorf("Problem parsing argument %s [%s]", n, err)
                problem = true
            }
        }
        // marshall it to be sent
        if !problem {
            aux, err := yaml.Marshal(values)
            marshalledSettings = aux
            if err != nil {
                logger.Errorf("Problem marshalling %s", values)
            }
        }
    }

    return &AppStartRequest{
        //Initialize default parameters here
        Name:            name,
        Namespace:       "",
        ManifestName:    manifestName,
        ManifestVersion: manifestVersion,
        ChartPath:       "",
        RawValue:        marshalledSettings,
        Labels:          labels,
        DisableHooks:    false,
        Replace:         false,
        Verify:          false,
        Keyring:         "", //defaultKeyring(),
        NameTemplate:    "",
        Version:         "",
        Timeout:         300,
        Wait:            false, //
        RepoURL:         "",
        Devel:           false,
    }
}

type AppStartResponse struct {
    Info         string `json:"info, omitempty"`
    LastDeployed string `json:"updated, omitempty"`
    Status       string `json:"status, omitempty"`
    Resources    string `json:"resources, omitempty"`
    Namespace    string `json:"namespace, omitempty"`

    TestStarted   string `json:"testStarted, omitempty"`
    TestCompleted string `json:"testCompleted, omitempty"`
    TestResult    string `json:"testResult, omitempty"`
    Notes         string `json:"notes, omitempty"`
}

// To be sent when we have to stop an application.
type AppStopRequest struct {
    //Command specific parameters here
    Name         string `json:"-"`
    DryRun       bool   `json:"-"`
    DisableHooks bool   `json:"-"`
    Purge        bool   `json:"-"` //Currently always purge
    Timeout      int64
}

func NewAppStopRequest(instanceName string) *AppStopRequest {
    return &AppStopRequest{
        Name:         instanceName,
        DryRun:       false,
        DisableHooks: true,
        Purge:        true,
        Timeout:      10,
    }
}

type AppStopResponse struct {
    Stopped bool   `json:"stopped, omitempty"`
    Info    string `json:"-"`
}

func NewAppStopResponse() *AppStopResponse {
    return &AppStopResponse{false, ""}
}
