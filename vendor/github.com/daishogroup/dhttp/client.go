//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//

package dhttp

import (
    "net/http"
    "io"
    "github.com/daishogroup/derrors"
)

// Client is the abstraction of the REST client.
type Client interface {
    // Get method.
    //   params:
    //     path 	Uri without the base path.
    //     output 	The output object.
    //   returns:
    //     HTTP Response. Error, if there is an internal error.
    Get(path string, output interface{}) Response

    // GetWithBody method get with body.
    //   params:
    //     path 	Uri without the base path.
    //     body     Body to be sent.
    //     output 	The output object.
    //   returns:
    //     HTTP Response. Error, if there is an internal error.
    GetWithBody(path string, body interface{}, output interface{}) Response

    // GetRaw is a request that does not require a output type, it returns a string variable.
    //  params:
    //      path Uri without the connection address
    //  returns:
    //      HTTP Response.
    GetRaw(path string) Response

    // Post method.
    //   params:
    //     path 	Uri without the base path.
    //     body 	Body object.
    //     output 	The output object.
    //   returns:
    //     HTTP Response. Error, if there is an internal error.
    Post(path string, body interface{}, output interface{}) Response

    // Put method.
    //   params:
    //     path 	Uri without the base path.
    //     body 	Body object.
    //     output 	The output object.
    //   returns:
    //     HTTP Response. Error, if there is an internal error.
    Put(path string, body interface{}, output interface{}) Response

    // Delete method.
    //   params:
    //     path Uri without the connection address.
    //     output 	The output object
    //   returns:
    //     HTTP Response. Error, if there is an internal error.
    Delete(path string, output interface{}) Response

    // DeleteWithBody method delete with body.
    //   params:
    //     path Uri without the connection address.
    //     body Body to be sent.
    //     output 	The output object
    //   returns:
    //     HTTP Response. Error, if there is an internal error.
    DeleteWithBody(path string, body interface{}, output interface{}) Response

    // Upload the content stored at the filepath and wait for a given response.
    //   params:
    //     path         Uri without the connection address.
    //     filepath     Target file to be uploaded.
    //     fieldname    The name of the field where the file content will be added.
    //     output       The output object.
    //  returns:
    //     HTTP Response. Error, if there is an internal error.
    Upload(path string, body io.Reader, filepath string, fieldname string, output interface{}) Response

    // GetBasePath return the current base path of the client.
    //   returns:
    //     The base path url.
    GetBasePath() string
}

// Response is a object with the HTTP response.
type Response struct {
    // The result object.
    Result interface{}
    // The http status code.
    Status *int
    // The error description.
    Error derrors.DaishoError
    // Headers
    Headers http.Header
}

// Ok checks if the response was successful and did not contain errors.
func (r *Response) Ok() bool {
    return r.Error == nil
}

// NewResponse is the constructor for the HTTP Response object.
//   params:
//     result 	The result object.
//     status 	The HTTP status code.
//     error 	The error description
//   returns:
//     A HTTP response object.
func NewResponse(result interface{}, status *int, error derrors.DaishoError) Response {
    return Response{Result: result, Status: status, Error: error, Headers: nil}
}

// NewResponse is the constructor for the HTTP Response object.
//   params:
//     result 	The result object.
//     status 	The HTTP status code.
//     header   The HTTP header.
//     error 	The error description
//   returns:
//     A HTTP response object.
func NewResponseWithHeader(result interface{}, status *int, header http.Header, error derrors.DaishoError) Response {
    return Response{Result: result, Status: status, Error: error, Headers: header}
}


// Add a new header to the response.
//  params:
//   key        The header key.
//   content    The content of the header.
func (r *Response) WithHeader(key string, content string) {
    r.Headers.Set(key,content)
}