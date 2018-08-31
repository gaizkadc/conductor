//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//

package dhttp

import (
    "testing"
    "net/http"
    "github.com/daishogroup/derrors"
    "github.com/stretchr/testify/suite"
)

type ClientMockupTestSuite struct {
    ClientTestSuite
}

func (suite *ClientMockupTestSuite) SetupSuite() {
    suite.client = NewClientMockup()
}

func (suite *ClientMockupTestSuite) BeforeTest(suiteName, testName string) {
    mockup := suite.client.(*ClientMockup)
    mockup.Reset()

    switch testName {
    case GetTest:
        mockup.AddGet(func(path string) Response {
            resp := NewSuccessfulOperation("get X-Daisho-Test: test123")
            status := 200
            return NewResponse(resp, &status, nil)
        })
    case GetNotFoundTest:
        mockup.AddGet(func(path string) Response {
            status := 404
            return NewResponse(nil, &status, derrors.NewConnectionError("404"))
        })
    case GetRawTest:
        mockup.AddGet(func(path string) Response {
            resp := NewSuccessfulOperation("get")
            status := 200

            return NewResponse(resp, &status, nil)
        })
    case GetRawNotFoundTest:
        mockup.AddGet(func(path string) Response {
            status := 404
            return NewResponse(nil, &status, derrors.NewConnectionError("404"))
        })
    case GetWithBodyTest:
        mockup.AddGetWithBody(func(path string, body interface{}) Response {
            resp := NewSuccessfulOperation("get")
            status := 200
            return NewResponse(resp, &status, nil)
        })
    case GetWithBodyNotFoundTest:
        mockup.AddGetWithBody(func(path string, body interface{}) Response {
            status := 404
            return NewResponse(nil, &status, derrors.NewConnectionError("404"))
        })
    case GetWithBodyNotValidTest:
        mockup.AddGetWithBody(func(path string, body interface{}) Response {
            status := 500
            return NewResponse(nil, &status, derrors.NewOperationError("invalid"))
        })
    case PostTest:
        mockup.AddPost(func(path string, body interface{}) Response {
            resp := NewSuccessfulOperation("get")
            status := 200

            return NewResponse(resp, &status, nil)
        })
    case PostNotValidTest:
        mockup.AddPost(func(path string, body interface{}) Response {
            status := 500
            return NewResponse(nil, &status, derrors.NewOperationError("invalid"))
        })
    case PostNotFoundTest:
        mockup.AddPost(func(path string, body interface{}) Response {
            status := 404
            return NewResponse(nil, &status, derrors.NewConnectionError("404"))
        })
    case PutTest:
        mockup.AddPut(func(path string, body interface{}) Response {
            resp := NewSuccessfulOperation("get")
            status := 200

            return NewResponse(resp, &status, nil)
        })
    case PutNotValidTest:
        mockup.AddPut(func(path string, body interface{}) Response {
            status := 500
            return NewResponse(nil, &status, derrors.NewOperationError("invalid"))
        })
    case PutNotFoundTest:
        mockup.AddPut(func(path string, body interface{}) Response {
            status := 404
            return NewResponse(nil, &status, derrors.NewConnectionError("404"))
        })
    case DeleteTest:
        mockup.AddDelete(func(path string) Response {
            resp := NewSuccessfulOperation("get")
            status := 200

            return NewResponse(resp, &status, nil)
        })
    case DeleteNotFoundTest:
        mockup.AddDelete(func(path string) Response {
            status := 404
            return NewResponse(nil, &status, derrors.NewConnectionError("404"))
        })
    case DeleteWithBodyTest:
        mockup.AddDeleteWithBody(func(path string, body interface{}) Response {
            resp := NewSuccessfulOperation("get")
            status := 200

            return NewResponse(resp, &status, nil)
        })
    case DeleteWithBodyNotValidTest:
        mockup.AddDeleteWithBody(func(path string, body interface{}) Response {
            status := 500
            return NewResponse(nil, &status, derrors.NewOperationError("invalid"))
        })
    case DeleteWithBodyNotFoundTest:
        mockup.AddDeleteWithBody(func(path string, body interface{}) Response {
            status := 404
            return NewResponse(nil, &status, derrors.NewConnectionError("404"))
        })
    case UploadValidTest:
        mockup.AddPut(func(path string, body interface{}) Response {
            resp := NewSuccessfulOperation("upload")
            status := 200
            return NewResponse(resp, &status, nil)
        })
    case GetWithHeaders:
        mockup.AddGet(func(path string) Response {
            resp := NewSuccessfulOperation("get X-Daisho-Test: test123")
            status := 200
            testHeaders := http.Header{}
            testHeaders.Set("TestingHeader", "TestingValue")
            return NewResponseWithHeader(resp, &status, testHeaders, nil)
        })
    }
}

func TestClientMockup(t *testing.T) {
    suite.Run(t, new(ClientMockupTestSuite))
}
