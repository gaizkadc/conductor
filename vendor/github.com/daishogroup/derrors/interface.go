//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// Definition of the error interface.

package derrors

// DaishoError defines the interface for all Daisho-defined errors.
type DaishoError interface {
    // Error returns the string representation of the error.
    Error() string
    // Type returns the ErrorType associated with the current DaishoError.
    Type() ErrorType
    // DebugReport returns a detailed error report including the stack information.
    DebugReport() string
    // StackTrace returns an array with the calling stack that created the error.
    StackTrace() [] StackEntry
}