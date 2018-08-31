//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// Entities defining access control.

package entities

// Role a user can have
type RoleType string

// Installator user
// This role is intended to be used only during installation periods that require authentication.
const InstallerType RoleType = "installer"

// Cloud user
// Requests coming from an slave cluster
const SlaveClusterType RoleType = "slavecluster"

// Application monitor
// Requests coming from an application monitor
const AppMonitorType RoleType = "appmonitor"

// Appmgr connector
const AppManagerType RoleType = "appmgr"

// Global administrator
const GlobalAdmin RoleType = "globaladmin"
// Local administrator
const LocalAdmin RoleType = "localadmin"
// Service technician
const ServiceTech RoleType = "servicetech"
// Developer type
const DeveloperType RoleType = "developer"
// Technical operator
const OperatorType RoleType = "operator"


// Overwrite the string casting operator.
//func (t RoleType) String () string {
//    return fmt.Sprintf("%s",t)
//}

// Return the array with the available roles in the system.
// return:
//  Array of available user roles
func AvailableRoles() []RoleType {
    roles := []RoleType{ InstallerType, SlaveClusterType, AppManagerType, AppMonitorType, GlobalAdmin, LocalAdmin, ServiceTech,
    DeveloperType, OperatorType}
    return roles
}

// Return the array with the available roles in the system.
// return:
//  Array of available user roles
func AvailableRolesString() []string {
    roles := []string{string(InstallerType),
        string(SlaveClusterType),
        string(AppManagerType),
        string(AppMonitorType),
        string(GlobalAdmin),
        string(LocalAdmin),
        string(ServiceTech),
        string(DeveloperType),
        string(OperatorType)}
    return roles
}


// Return an array of available system roles.
//  return:
//   Array of available system roles.
func AvailableSystemRoles() []RoleType {
    roles := []RoleType{InstallerType, SlaveClusterType, AppMonitorType, AppManagerType}
    return roles
}

// Return the array with the available roles in the system.
// return:
//  Array of available user roles
func AvailableUserRoles() []RoleType {
    roles := []RoleType{ GlobalAdmin, LocalAdmin, ServiceTech, DeveloperType, OperatorType}
    return roles
}

// ValidRoleType checks the type enum to determine if the string belongs to the enumeration.
//   params:
//     roleType The type to be checked
//   returns:
//     Whether it is contained in the enum.
func ValidRoleType(roleType RoleType) bool {
    for _, r := range AvailableRoles() {
        if roleType == r {
            return true
        }
    }
    return false
}

// Indicate whether the user corresponds to an internal platform user or not
// return:
//  True if this is an internal user.
func (r RoleType) IsInternalUser() bool {
    switch r {
    case AppManagerType:
        return true
    case AppMonitorType:
        return true
    case SlaveClusterType:
        return true
    default:
        return false
    }
}

// Determine if the role corresponds to a temporal user.
func (r RoleType) IsTemporalUser() bool {
    switch r {
    case InstallerType:
        return true
    default:
        return false
    }
}

// Information about users privilege access.
type UserAccess struct {
    UserID  string `json:"userId, omitempty"`
    Roles    []RoleType `json:"roles, omitempty"`
}

// Instantiate a new user access entry.
//  params:
//     userID   The user ID.
//     role     The role access privilege
//  return:
//     A new user access entry.
func NewUserAccess(userID string, role []RoleType) *UserAccess {
    return &UserAccess{UserID: userID, Roles: role}
}

// Determine if a user access is valid or not.
func (u *UserAccess) IsValid() bool {
    if u.UserID == "" {
        return false
    }
    return false
}

// Entity to provide a new request to add user accesses
type AddUserAccessRequest struct {
    Roles [] RoleType
}

// Create a new access request.
//  params:
//    roles Role to be added.
//  return:
//    A new request.
func NewAddUserAccessRequest (roles []RoleType) *AddUserAccessRequest {
    return &AddUserAccessRequest{Roles: roles}
}

// Reduced information structure for summarizing.
type UserAccessReducedInfo struct {
    UserID  string `json:"userId, omitempty"`
    Roles    [] RoleType `json:"roles, omitempty"`
}

func NewUserAccessReducedInfo(userID string, role []RoleType)  * UserAccessReducedInfo {
    return &UserAccessReducedInfo{UserID: userID, Roles: role}
}