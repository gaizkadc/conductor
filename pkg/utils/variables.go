/*
 * Copyright (C) 2018 Nalej Group - All Rights Reserved
 *
 */

package utils

// This a common list of global variables



const (
    // Environment variable to define the system model address
    IT_SYSTEM_MODEL = "IT_SYSTEM_MODEL"
    // Environment variable to define the networking manager address
    IT_NETWORKING_MANAGER = "IT_NETWORKING_MANAGER"
    // ID for a musician cluster
    MUSICIAN_CLUSTER_ID = "CLUSTER_ID"
)


// Standard conductor port
var CONDUCTOR_PORT uint32 = 5000
// Standard system model port
var SYSTEM_MODEL_PORT uint32 = 8800
// App cluster api port
var APP_CLUSTER_API_PORT uint32 =443
// Relevant ports for the system
var MUSICIAN_PORT uint32 = 5100
// Standard deployment manager port
var DEPLOYMENT_MANAGER_PORT uint32 = 443
// Networking service port
var NETWORKING_SERVICE_PORT uint32 = 8000
// Standard authx port
var AUTHX_PORT uint32 = 8010
// Standard unifiedLogging port
var UNIFIED_LOGGING_PORT = 8323
