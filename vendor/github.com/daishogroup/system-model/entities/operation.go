//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the specification of the operation entities.

package entities

// Successful operation is a object to send by endpoints that return void.
type SuccessfulOperation struct {
    //Name of the method.
    Operation string `json:"operation"`
}

// Basic constructor.
func NewSuccessfulOperation(operation string) SuccessfulOperation{
    return SuccessfulOperation{operation}
}

