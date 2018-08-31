//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package config

import (
    "testing"

    "github.com/daishogroup/system-model/entities"
    "github.com/daishogroup/system-model/provider/configstorage"
    "github.com/stretchr/testify/suite"
)

type ManagerHelper struct {
    configProvider *configstorage.MockupConfigProvider
    configMgr Manager
}

func NewManagerHelper() ManagerHelper {
    var configProvider = configstorage.NewMockupConfigProvider()
    var configMgr = NewManager(configProvider)
    return ManagerHelper{configProvider, configMgr}
}

type TestHelper struct {
    suite.Suite
    manager ManagerHelper
}

func (helper *TestHelper) SetupSuite() {
    managerHelper := NewManagerHelper()
    helper.manager = managerHelper
}

// The SetupTest method is called before every test on the suite.
func (helper *TestHelper) SetupTest() {
    helper.manager.configProvider.Clear()
}

func TestHandlerSuite(t *testing.T) {
    suite.Run(t, new(TestHelper))
}

func (helper *TestHelper) TestSetGet() {
    newConfig := entities.NewConfig("1h")
    err := helper.manager.configMgr.SetConfig(*newConfig)
    helper.Nil(err, "config should be added")
    retrieved, err := helper.manager.configMgr.GetConfig()
    helper.Nil(err, "config should be retrieved")
    helper.Equal(newConfig.LogRetention, retrieved.LogRetention, "log retention should match")
}