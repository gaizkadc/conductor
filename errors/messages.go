//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Definition of the error messages.

package errors

// InvalidEntity message indicating that the associated entity cannot be validated.
const InvalidEntity = "invalid entity, check mandatory fields"
// MarshalError message to indicate errors with the json.Marshal operation.
const MarshalError = "cannot marshal structure"
// UnmarshalError message to indicate errors with the json.Unmarshal operation.
const UnmarshalError = "cannot unmarshal structure"
// OpFail message to indicate that a complex operation has failed.
const OpFail = "operation failed"
// MissingRESTParameter message to indicate that a required parameter is missing.
const MissingRESTParameter = "missing rest parameter"


// ClustersNotAvailable message to indicate that the scheduling fails as there are not available clusters.
const ClustersNotAvailable = "no available clusters"
const NodeNotAvailable = "node not available"
const ConnectionError = "error connecting with external entity"
const CannotStopApp = "cannot stop application"
const CannotDeleteApp = "cannot delete application"
const CannotUpdateApp = "cannot update application"
const CannotStoreApp = "cannot store application"
const NoApplicationsAvailable = "no application descriptors available"
const CannotDeployApp = "cannot deploy application"
const CannotRetrieveApp = "cannot retrieve application"
const UndeploySuccess = "application was successfully undeployed"

const MaxAppNameLength = "application name cannot exceed 63 characters"

const InvalidAppName = "application name must conform to [a-z]([-a-z0-9]*[a-z0-9])?"

const MissingRequiredFields = "missing application required fields"

// HTTPConnectionError message to indicate that the communication with an external entity has failed.
const HTTPConnectionError = "HTTP connection error"