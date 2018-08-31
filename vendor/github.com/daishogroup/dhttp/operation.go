//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//

package dhttp

// SuccessfulOperation is a object to send by endpoints that return void.
type SuccessfulOperation struct {
    //Name of the method.
    Operation string `json:"operation"`
}

// NewSuccessfulOperation is the basic constructor.
func NewSuccessfulOperation(operation string) *SuccessfulOperation{
    return &SuccessfulOperation{operation}
}
