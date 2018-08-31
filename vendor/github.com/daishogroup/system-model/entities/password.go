//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
// User password representation to be used.
//

package entities

import (
    "golang.org/x/crypto/bcrypt"
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/system-model/errors"
)

// Data structure to store and deal with passwords
type Password struct{
    // user identifier
    UserID string `json: "userId, omitempty"`
    //Hash representation
    Hash *[]byte `json: "hash, omitempty"`
}

// Create a new password.
//  params:
//    userID   User identifier string.
//    password Password string.
//  return:
//    New password entity.
//    Error if any.
func NewPassword(userID string, password *string) (*Password, derrors.DaishoError){
    // Generate "hash" to store from user password
    if password!=nil {
        hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
        if err != nil {
            return nil, derrors.NewGenericError(errors.ErrorPassword,err)
        }
        return &Password{userID,&hash}, nil
    }

    return &Password{userID, nil}, nil
}

// Compare this password with an incoming one.
//  params:
//    password String password to be compared with.
//  return:
//    True if both entries match. False if any of them is nil.
func(pass *Password) CompareWith(password string) bool {
    if pass.Hash == nil {
        // By now we simply return false if the object password is nil
        return false
    }
    if err := bcrypt.CompareHashAndPassword(*pass.Hash, []byte(password)); err != nil {
        return false
    }
    return true
}