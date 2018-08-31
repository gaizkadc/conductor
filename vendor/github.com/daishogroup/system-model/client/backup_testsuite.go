//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the backup Client TestSuite

package client

import (
    "github.com/stretchr/testify/suite"
)

func TestExport(suite *suite.Suite, backup Backup) {
    n, err := backup.Export("all")
    suite.Nil(err, "error must be Nil")
    suite.NotNil(n, "network must not be Nil")
}

func TestImport(suite *suite.Suite, backup Backup) {

    // create data for Import
    backup.(*BackupMockup).InitMockup()
    backupData, err := backup.Export("all")

    // clean up created data
    backup.(*BackupMockup).ClearMockup()


    err = backup.Import("all", backupData)
    suite.Nil(err, "error must be  Nil")

    importData, _ := backup.Export("all")
    // Ignore error as equal test will show failure
    suite.Equal(len(backupData.Clusters), len(importData.Clusters), "Backup data and resotre data must be same")
}

