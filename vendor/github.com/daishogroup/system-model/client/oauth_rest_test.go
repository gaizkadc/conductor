//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//
// This file contains the User Client Rest Test

package client

import (
    "testing"

    "github.com/daishogroup/system-model/entities"
    "github.com/stretchr/testify/suite"
    "github.com/daishogroup/system-model/errors"
    //"github.com/daishogroup/derrors"
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/dhttp"
)

type OAuthRestTestSuite struct {
    suite.Suite
    oauth OAuth
    rest *dhttp.ClientMockup
}

func (suite *OAuthRestTestSuite) SetupSuite() {
    suite.rest = dhttp.NewClientMockup()
    suite.oauth = &OAuthRest{suite.rest}
}

func (suite *OAuthRestTestSuite) SetupTest() {
    suite.rest.Reset()

}

func (suite *OAuthRestTestSuite) TestSetSecrets() {

    suite.rest.AddPost(func(path string, body interface{}) dhttp.Response {
        statusCode := 200
        return dhttp.NewResponse(entities.NewSuccessfulOperation(errors.OAuthEntrySet), &statusCode, nil)
    })

    entrySet:=entities.NewOAuthSecrets(OAuthTestUserID)
    entrySet.AddEntry("app1",entities.NewOAuthEntry("clientID1","secret1"))

    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 200
        return dhttp.NewResponse(&entrySet, &statusCode, nil)
    })

    TestSetSecrets(&suite.Suite, suite.oauth)
}

func (suite *OAuthRestTestSuite) TestGetSecret() {
    entrySet:=entities.NewOAuthSecrets(OAuthTestUserID)
    entrySet.AddEntry("app1",entities.NewOAuthEntry("clientID1","secret1"))

    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 200
        return dhttp.NewResponse(&entrySet, &statusCode, nil)
    })

    TestGetSecrets(&suite.Suite, suite.oauth)
}

func (suite *OAuthRestTestSuite) TestDeleteSecrets() {
    entrySet:=entities.NewOAuthSecrets(OAuthTestUserID)
    entrySet.AddEntry("app1",entities.NewOAuthEntry("clientID1","secret1"))

    suite.rest.AddDelete(func(path string) dhttp.Response {
        statusCode := 200
        return dhttp.NewResponse(nil, &statusCode, nil)
    })

    suite.rest.AddGet(func(path string) dhttp.Response {
        statusCode := 500
        return dhttp.NewResponse(nil, &statusCode, derrors.NewOperationError(errors.UserDoesNotExist))
    })

    TestDeleteSecrets(&suite.Suite, suite.oauth)
}



// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestOAuthRest(t *testing.T) {
    suite.Run(t, new(OAuthRestTestSuite))
}