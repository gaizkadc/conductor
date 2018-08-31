//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//

package entities

import (
    "fmt"

    "github.com/stretchr/testify/assert"
    "testing"
)

const (
    networkTestID = "n1"
    networkTestName = "Test network"
    networkTestDescription = "Network description"
    networkTestAdminName = "The Admin"
    networkTestAdminPhone = "1234 1234 1234"
    networkTestAdminEmail = "admin@admins.com"
    networkTestEdgenetID = "deadbeef00123456"
)

func TestNewNetworkWithId(t *testing.T) {
    network := NewNetworkWithID(
        networkTestID,
        networkTestName,
        networkTestDescription,
        networkTestAdminName,
        networkTestAdminPhone,
        networkTestAdminEmail)
    assert.NotNil(t, network, "Network should be defined")
    expected := `&entities.Network{ID:"n1", EdgenetID:"",` +
        ` Name:"Test network",` +
        ` Description:"Network description",` +
        ` AdminName:"The Admin",` +
        ` AdminPhone:"1234 1234 1234",` +
        ` AdminEmail:"admin@admins.com",` +
        ` Operator:(*entities.User)(nil)}`
    assert.Equal(t, expected, network.String())
}

func TestToNetwork(t *testing.T) {
    request := NewAddNetworkRequest(
        networkTestName,
        networkTestDescription,
        networkTestAdminName,
        networkTestAdminPhone,
        networkTestAdminEmail)

    // Fill in other request values
    request.EdgenetID = networkTestEdgenetID

    network := ToNetwork(*request)

    // Pull some trick to make this testable, as ID is generated
    // Maybe we should make a GenerateUUID mockup.
    id := network.ID
    expected := fmt.Sprintf(
        `&entities.Network{ID:"%s", EdgenetID:"%s",` +
        ` Name:"%s",` +
        ` Description:"%s",` +
        ` AdminName:"%s",` +
        ` AdminPhone:"%s",` +
        ` AdminEmail:"%s",` +
        ` Operator:(*entities.User)(nil)}`,
        id, networkTestEdgenetID, networkTestName,
        networkTestDescription, networkTestAdminName,
        networkTestAdminPhone, networkTestAdminEmail)
    assert.Equal(t, expected, network.String())
}
