//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
//

package configstorage

import (
    "io/ioutil"
    "os"
    "testing"

    "github.com/daishogroup/system-model/entities"
    "github.com/stretchr/testify/assert"
)

func TestConfigProvider(t *testing.T) {

    dir, _ := ioutil.TempDir("", "testConfig")
    defer os.RemoveAll(dir)

    provider := NewFileSystemProvider(dir)
    config := entities.NewConfig("1d")
    err := provider.Store(*config)
    assert.Nil(t, err, "config should be stored")

    retrieved, err := provider.Get()
    assert.Nil(t, err, "config should be retrieved")
    assert.EqualValues(t, config, retrieved, "config should match")

}
