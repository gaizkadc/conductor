package client

//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// This file contains the Backup Client Mockup tests

import (
    "github.com/stretchr/testify/suite"
    "testing"
)

type BackupMockupTestSuite struct {
    suite.Suite
    backup Backup
}

func (suite *BackupMockupTestSuite) SetupTest() {
    suite.backup = NewBackupMockup()
}

func (suite *BackupMockupTestSuite) TestExport() {
    suite.backup.(*BackupMockup).InitMockup()
    TestExport(&suite.Suite,suite.backup)
}

func (suite *BackupMockupTestSuite) TestImport() {
    TestImport(&suite.Suite,suite.backup)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestBackupMockup(t *testing.T) {
    suite.Run(t, new(BackupMockupTestSuite))
}

