//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Definition of the supported types of error.

package derrors

// ErrorType defines a new type for creating a enum of error types.
type ErrorType string
// GenericErrorType to be used with general errors.
const GenericErrorType ErrorType = "GenericError"
// ConnectionErrorType is associated with connectivity errors.
const ConnectionErrorType ErrorType = "Connection"
// EntityErrorType is associated with entity related errors including validation, association, etc.
const EntityErrorType ErrorType = "Entity"
// OperationErrorType is associated with failures in external operations.
const OperationErrorType ErrorType = "Operation"
// ProviderErrorType is associated with provider related errors including invalid operations, provider failures, etc.
const ProviderErrorType ErrorType = "Provider"
// OrchestrationErrorType is associated with orchestration related errors including orchestration failures, preconditions, etc.
const OrchestrationErrorType ErrorType = "Orchestration"

// ValidErrorType checks the type enum to determine if the string belongs to the enumeration.
//   params:
//     errorType The type to be checked
//   returns:
//     Whether it is contained in the enum.
func ValidErrorType(errorType ErrorType) bool {
    switch errorType {
    case "" : return false
    case GenericErrorType : return true
    case ConnectionErrorType : return true
    case EntityErrorType : return true
    case OperationErrorType : return true
    case ProviderErrorType : return true
    case OrchestrationErrorType : return true
    default: return false
    }
}