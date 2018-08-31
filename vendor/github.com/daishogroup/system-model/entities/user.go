//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//

package entities

import (
    "fmt"
    "time"
    "strings"
)

// JSONTime is an internal workaround to serialize time struct in a human readable way
type JSONTime struct{
    time.Time
}

func (t *JSONTime) MarshalJSON() ([]byte, error) {
    //do your serializing here
    stamp := fmt.Sprintf("\"%s\"", t.Format(time.RFC3339))
    return []byte(stamp), nil
}
func (t *JSONTime) UnmarshalJSON(value []byte) (err error){
    // unserialize here
    s := strings.Trim(string(value), "\"")
    if s == "null" {
        *t = JSONTime{time.Time{}}
        return err
    }
    parsedTime, err := time.Parse(time.RFC3339, s)
    if err != nil {
        return err
    }
    *t = JSONTime{parsedTime}

    return nil
}

func (t *JSONTime) String() string {
    return fmt.Sprintf("\"%s\"", t.Format(time.RFC3339))
}


// A user of the system. Users are associated with a given network.
type User struct {
    // User identifier.
    ID string `json:"userId, omitempty"`
    // User name.
    Name string `json:"name, omitempty"`
    // User phone number.
    Phone string `json:"phone, omitempty"`
    // User email.
    Email string `json:"email, omitempty"`
    // Date the account was created.
    CreationTime JSONTime `json:"creationTime, omitempty"`
    // Date this account will expire.
    ExpirationTime JSONTime `json:"expirationTime, omitempty"`
}


// Generate a new User.
//   params:
//     name The user name.
//     phone The user phone.
//     email The user email.
//     creationTime The time the user is created.
//     expirationTime When the user account expires.
//   returns:
//     A new user.
func NewUser(name string, phone string, email string, creationTime time.Time, expirationTime time.Time) * User {
    uuid := GenerateUUID(UserPrefix)
    return NewUserWithID(uuid, name, phone, email, creationTime, expirationTime)
}

// Generate a new User.
//   params:
//     id The user identifier.
//     name The user name.
//     phone The user phone.
//     email The user email.
//     password The user password.
//   returns:
//     A new user.
func NewUserWithID(id string, name string, phone string, email string,
    creationTime time.Time, expirationTime time.Time) * User {
    u := & User {id, name, phone, email, JSONTime{creationTime}, JSONTime{expirationTime}}
    return u
}

func (u * User) String() string {
    return fmt.Sprintf("{ID: %s, Name: %s, Phone: %s, Email: %s, CreationTime: %s, ExpirationTime: %s}",
        u.ID, u.Name, u.Phone, u.Email, u.CreationTime.String(), u.ExpirationTime.String())
}

// Merge a existing user with the new values of an update request.
//   params:
//     update Request to update user fields.
//   return:
//     New generated user.
func (u* User) Merge(update UpdateUserRequest) *User {
    if update.Name != nil {
        u.Name = *update.Name
    }

    if update.Phone != nil {
        u.Phone = *update.Phone
    }

    if update.Email != nil {
        u.Email = *update.Email
    }

    return u
}


// Reduced user information for summary operations.
type UserReducedInfo struct{
    // User identifier.
    ID string `json:"userId, omitempty"`
    // User email.
    Email string `json:"email, omitempty"`
}

// Generated a reduced user info.
//   params:
//     id The user identifier.
//     email The user email.
//   returns:
//     New user reduced info object.
func NewUserReducedInfo(id string, email string) * UserReducedInfo {
    u := &UserReducedInfo{id, email}
    return u
}

// REST request to be fulfilled in order to add a new user to the system model.
type AddUserRequest struct {
    // User identifier.
    ID string `json:"userId, omitempty"`
    // User name.
    Name string `json:"name, omitempty"`
    // User phone number.
    Phone string `json:"phone, omitempty"`
    // User email.
    Email string `json:"email, omitempty"`
    // Date the account was created.
    CreationTime JSONTime `json:"creationTime, omitempty"`
    // Date this account will expire.
    ExpirationTime JSONTime `json:"expirationTime, omitempty"`

}

// Check if the request to add a new user is correctly formed.
// return:
//    true if the request is correct, otherwise false
func(request *AddUserRequest) IsValid() bool {
    return request.ID != "" && request.Name != "" && request.Email != ""
}

// Generate a new AddUserRequest.
//   params:
//     id The user identifier.
//     name The user name.
//     phone The user phone.
//     email The user email.
//     password The user password.
//   returns:
//     A new user request.
func NewAddUserRequest(id string, name string, phone string, email string, creationTime time.Time,
    expirationTime time.Time) * AddUserRequest {
    u := & AddUserRequest{id, name, phone, email,
    JSONTime{creationTime}, JSONTime{expirationTime}}
    return u
}

// Return a well-formatted string.
// return:
//     Well formatted string.
func (request * AddUserRequest) String() string {
    return fmt.Sprintf("%#v", request)
}

// Structure to update an existing user data.
type UpdateUserRequest struct {
    // User name.
    Name *string `json:"name, omitempty"`
    // User phone number.
    Phone *string `json:"phone, omitempty"`
    // User email.
    Email *string `json:"email, omitempty"`
}

func NewUpdateUserRequest() *UpdateUserRequest{
    return &UpdateUserRequest{nil, nil, nil}
}

func(r *UpdateUserRequest) IsValid() bool {
    // TODO decide what makes a request to be correct
    return true
}

func(r *UpdateUserRequest) WithName(name string) *UpdateUserRequest{
    r.Name = &name
    return r
}

func(r *UpdateUserRequest) WithPhone(phone string) *UpdateUserRequest{
    r.Phone = &phone
    return r
}

func(r *UpdateUserRequest) WithEmail(email string) *UpdateUserRequest{
    r.Email = &email
    return r
}


// Structure to combine users with their roles.
type UserExtended struct {
    // User identifier.
    ID string `json:"userId, omitempty"`
    // User name.
    Name string `json:"name, omitempty"`
    // User phone number.
    Phone string `json:"phone, omitempty"`
    // User email.
    Email string `json:"email, omitempty"`
    // Date the account was created.
    CreationTime JSONTime `json:"creationTime, omitempty"`
    // Date this account will expire.
    ExpirationTime JSONTime `json:"expirationTime, omitempty"`
    // User Roles
    Roles [] RoleType `json:"roles, omitempty"`
}

func NewUserExtended(id string, name string, phone string, email string, creationTime time.Time,
    expirationTime time.Time, roles []RoleType) UserExtended{
    return UserExtended{ID: id, Name: name, Phone: phone, Email: email,
        CreationTime: JSONTime{creationTime}, ExpirationTime: JSONTime{expirationTime}, Roles: roles}
}

