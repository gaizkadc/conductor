//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//

package dhttp

import (
    "github.com/daishogroup/derrors"
    "github.com/stretchr/testify/suite"
    "io/ioutil"
    "os"
)

const host = "localhost"

//GetTest is the name of the TestGet
const GetTest = "TestGet"

//GetRawTest is the name of the GetRawTest
const GetRawTest = "TestGetRaw"

//GetRawNotFoundTest is the name of the TestGetRawNotFound
const GetRawNotFoundTest = "TestGetRawNotFound"

//GetNotFoundTest is the name of the TestGetNotFound
const GetNotFoundTest = "TestGetNotFound"

//GetWithBodyTest is the name of the TestGetWithBody
const GetWithBodyTest = "TestGetWithBody"

//GetWithBodyNotFoundTest is the name of the TestGetWithBodyNotFound
const GetWithBodyNotFoundTest = "TestGetWithBodyNotFound"

//GetWithBodyNotValidTest is the name of the TestGetWithBodyNotValid
const GetWithBodyNotValidTest = "TestGetWithBodyNotValid"

//PostTest is the name of the TestPost
const PostTest = "TestPost"

//PostNotValidTest is the name of the TestPostNotValid
const PostNotValidTest = "TestPostNotValid"

//PostNotFoundTest is the name of the TestPostNotFound
const PostNotFoundTest = "TestPostNotFound"

//PutTest is the name of the TestPut
const PutTest = "TestPut"

//PutNotValidTest is the name of the TestPutNotValid
const PutNotValidTest = "TestPutNotValid"

//PutNotFoundTest is the name of the TestPutNotFound
const PutNotFoundTest = "TestPutNotFound"

//DeleteTest is the name of the TestDelete
const DeleteTest = "TestDelete"

//DeleteNotFoundTest is the name of the TestDeleteNotFound
const DeleteNotFoundTest = "TestDeleteNotFound"

//DeleteWithBodyTest is the name of the TestDeleteWithBody
const DeleteWithBodyTest = "TestDeleteWithBody"

//DeleteWithBodyNotValidTest is the name of the TestDeleteWithBodyNotValid
const DeleteWithBodyNotValidTest = "TestDeleteWithBodyNotValid"

//DeleteWithBodyNotFoundTest is the name of the TestDeleteWithBodyNotFound
const DeleteWithBodyNotFoundTest = "TestDeleteWithBodyNotFound"

//UploadValidTest is th ename of the TestUploadtest
const UploadValidTest = "TestUpload"

// GetWithHeaders is the name of the TestGetWithHeaders
const GetWithHeaders = "TestGetWithHeaders"

// ClientTestSuite is the generic class to generate the client tests.
type ClientTestSuite struct {
    suite.Suite
    client Client
}

// GetURL creates the URL using the base path.
func (t *ClientTestSuite) GetURL(path string) string {
    return t.client.GetBasePath() + path
}

// TestGet basic get method test.
func (t *ClientTestSuite) TestGet() {
    url := t.GetURL("/test/get")
    result := &SuccessfulOperation{}
    response := t.client.Get(url, result)
    httpResponse, ok := response.Result.(*SuccessfulOperation)
    t.True(ok)
    t.Equal("get X-Daisho-Test: test123", httpResponse.Operation)
    t.NoError(response.Error, "error must be nil")
}

// TestGetNotFound is get method test using a URI not valid.
func (t *ClientTestSuite) TestGetNotFound() {
    url := t.GetURL("/test/getNotFound")
    result := &SuccessfulOperation{}
    response := t.client.Get(url, result)
    t.Error(response.Error, "method must fail")
    t.Equal(response.Error.Type(), derrors.ConnectionErrorType)

}

// TestGetRaw basic get raw method test.
func (t *ClientTestSuite) TestGetRaw() {
    url := t.GetURL("/test/get")
    response := t.client.GetRaw(url)
    t.NoError(response.Error, "error must be nil")
}

// TestGetRawNotFound is get method test using a URI not valid.
func (t *ClientTestSuite) TestGetRawNotFound() {
    url := t.GetURL("/test/getNotFound")
    response := t.client.GetRaw(url)
    t.Error(response.Error, "method must fail")
    t.Equal(response.Error.Type(), derrors.ConnectionErrorType)
}
// TestGetWithBody basic get with body method test.
func (t *ClientTestSuite) TestGetWithBody() {
    url := t.GetURL("/test/getWithBody")
    result := &SuccessfulOperation{}
    response := t.client.GetWithBody(url, NewValidTestBody(), result)
    t.NoError(response.Error, "error must be nil")
}

// TestGetWithBodyNotValid is get method test using a Body not valid.
func (t *ClientTestSuite) TestGetWithBodyNotValid() {
    url := t.GetURL("/test/getWithBody")
    result := &SuccessfulOperation{}
    response := t.client.GetWithBody(url, NewInvalidTestBody(), result)
    t.Error(response.Error, "method must fail")
    t.NotEqual(response.Error.Type(), derrors.ConnectionErrorType)
}

// TestGetWithBodyNotFound is get method test using a URI not valid.
func (t *ClientTestSuite) TestGetWithBodyNotFound() {
    url := t.GetURL("/test/getWithBodyNotFound")
    result := &SuccessfulOperation{}
    response := t.client.GetWithBody(url, NewValidTestBody(), result)
    t.Error(response.Error, "method must fail")
    t.Equal(response.Error.Type(), derrors.ConnectionErrorType)
}

// TestPost basic post method test.
func (t *ClientTestSuite) TestPost() {
    url := t.GetURL("/test/post")
    result := &SuccessfulOperation{}
    response := t.client.Post(url, NewValidTestBody(), result)
    t.NoError(response.Error, "error must be nil")
}

// TestPostNotValid is post method test using a Body not valid.
func (t *ClientTestSuite) TestPostNotValid() {
    url := t.GetURL("/test/post")
    result := &SuccessfulOperation{}
    response := t.client.Post(url, NewInvalidTestBody(), result)
    t.Error(response.Error, "method must fail")
    t.NotEqual(response.Error.Type(), derrors.ConnectionErrorType)
}

// TestPostNotFound is post method test using a URI not valid.
func (t *ClientTestSuite) TestPostNotFound() {
    url := t.GetURL("/test/postNotFound")
    result := &SuccessfulOperation{}
    response := t.client.Post(url, NewValidTestBody(), result)
    t.Error(response.Error, "method must fail")
    t.Equal(response.Error.Type(), derrors.ConnectionErrorType)
}
// TestPut basic put method test.
func (t *ClientTestSuite) TestPut() {
    url := t.GetURL("/test/put")
    result := &SuccessfulOperation{}
    response := t.client.Put(url, NewValidTestBody(), result)

    t.NoError(response.Error, "error must be nil")
}

// TestPutNotValid is put method test using a Body not valid.
func (t *ClientTestSuite) TestPutNotValid() {
    url := t.GetURL("/test/put")
    result := &SuccessfulOperation{}
    response := t.client.Put(url, NewInvalidTestBody(), result)
    t.Error(response.Error, "method must fail")
    t.NotEqual(response.Error.Type(), derrors.ConnectionErrorType)
}

// TestPutNotFound is put method test using a URI not valid.
func (t *ClientTestSuite) TestPutNotFound() {
    url := t.GetURL("/test/putNotFound")
    result := &SuccessfulOperation{}
    response := t.client.Put(url, NewValidTestBody(), result)
    t.Error(response.Error, "method must fail")
    t.Equal(response.Error.Type(), derrors.ConnectionErrorType)
}

// TestDelete basic delete method test.
func (t *ClientTestSuite) TestDelete() {
    url := t.GetURL("/test/delete")
    result := &SuccessfulOperation{}
    response := t.client.Delete(url, result)
    t.NoError(response.Error, "error must be nil")
}

// TestDeleteNotFound is delete method test using a URI not valid.
func (t *ClientTestSuite) TestDeleteNotFound() {
    url := t.GetURL("/test/deleteNotFound")
    result := &SuccessfulOperation{}
    response := t.client.Delete(url, result)
    t.Error(response.Error, "method must fail")
    t.Equal(response.Error.Type(), derrors.ConnectionErrorType)
}

// TestDeleteWithBody basic delete with body method test.
func (t *ClientTestSuite) TestDeleteWithBody() {
    url := t.GetURL("/test/deleteWithBody")
    result := &SuccessfulOperation{}
    response := t.client.DeleteWithBody(url, NewValidTestBody(), result)

    t.NoError(response.Error, "error must be nil")
}
// TestDeleteWithBodyNotValid is delete method test using a Body not valid.
func (t *ClientTestSuite) TestDeleteWithBodyNotValid() {
    url := t.GetURL("/test/deleteWithBody")
    result := &SuccessfulOperation{}
    response := t.client.DeleteWithBody(url, NewInvalidTestBody(), result)
    t.Error(response.Error, "method must fail")
    t.NotEqual(response.Error.Type(), derrors.ConnectionErrorType)
}

// TestDeleteWithBodyNotFound is delete method test using a URI not valid.
func (t *ClientTestSuite) TestDeleteWithBodyNotFound() {
    url := t.GetURL("/test/deleteWithBodyNotFound")
    result := &SuccessfulOperation{}
    response := t.client.DeleteWithBody(url, NewValidTestBody(), result)
    t.Error(response.Error, "method must fail")
    t.Equal(response.Error.Type(), derrors.ConnectionErrorType)
}

func (t *ClientTestSuite) TestUpload() {
    url := t.GetURL("/test/upload")
    result := &SuccessfulOperation{}

    file, err := ioutil.TempFile(os.TempDir(), "prefix")
    t.NoError(err, "we have to create a file to test this feature!!!")
    // remove it when finished
    defer os.Remove(file.Name())

    // Put something inside
    for i:=0;i<5;i++ {
        _, err = file.WriteString("this is a testing file\n")
        t.NoError(err, "we have to write some data in the file for testing!!!")
    }

    response := t.client.Upload(url, nil, file.Name(), "file", result)
    t.NoError(response.Error, "no error was expected")

}

func (t *ClientTestSuite) TestGetWithHeaders() {
    url := t.GetURL("/test/getWithHeaders")
    result := &SuccessfulOperation{}
    response := t.client.Get(url, result)
    httpResponse, ok := response.Result.(*SuccessfulOperation)
    t.True(ok)
    t.Equal("get X-Daisho-Test: test123", httpResponse.Operation)
    t.Equal("TestingValue", response.Headers.Get("TestingHeader"))
    t.NoError(response.Error, "error must be nil")
}