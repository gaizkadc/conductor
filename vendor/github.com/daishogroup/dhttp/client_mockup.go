//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//

package dhttp

import (
    "io/ioutil"
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/dhttp/errors"
    "io"
)

// ClientMockup simulates responses of a REST API.
type ClientMockup struct {
    // Internal index with the next GET Method.
    getIndex int
    // Slice with the response of each GET call.
    getF [] func(path string) Response
    // Internal index with the next GET body Method.
    getWithBodyIndex int
    // Slice with the response of each GET body call.
    getWithBodyF [] func(path string, body interface{}) Response
    // Internal index with the next POST Method.
    postIndex int
    // Slice with the response of each POST call.
    postF [] func(path string, body interface{}) Response
    // Internal index with the next PUT Method.
    putIndex int
    // Slice with the response of each PUT call.
    putF [] func(path string, body interface{}) Response
    // Internal index with the next DELETE Method.
    deleteIndex int
    // Slice with the response of each DELETE call.
    deleteF [] func(path string) Response
    // Internal index with the next DELETE with bodyMethod.
    deleteWithBodyIndex int
    // Slice with the response of each DELETE with body call.
    deleteWithBodyF [] func(path string, body interface{}) Response

}

// Get method.
//   params:
//     path 	Uri without the base path.
//     output 	The output object.
//   returns:
//     HTTP Response. Error, if there is an internal error.
func (client *ClientMockup) Get(path string, output interface{}) Response {
    if client.getIndex >= len(client.getF) {
        client.getIndex = 0
    }
    selectedGet := client.getF[client.getIndex]
    client.getIndex ++
    return selectedGet(path)
}

// GetWithBody method get with body.
//   params:
//     path 	Uri without the base path.
//     body     Body to be sent.
//     output 	The output object.
//   returns:
//     HTTP Response. Error, if there is an internal error.
func (client *ClientMockup) GetWithBody(path string, body interface{}, output interface{}) Response {
    if client.getWithBodyIndex >= len(client.getWithBodyF) {
        client.getIndex = 0
    }
    selectedGet := client.getWithBodyF[client.getWithBodyIndex]
    client.getWithBodyIndex ++
    return selectedGet(path, body)
}

// GetRaw is a request that does not require a output type, it returns a string variable.
//  params:
//      path Uri without the connection address
//  returns:
//      HTTP Response.
func (client *ClientMockup) GetRaw(path string) Response {
    if client.getIndex >= len(client.getF) {
        client.getIndex = 0
    }
    selectedGet := client.getF[client.getIndex]
    client.getIndex ++
    return selectedGet(path)
}

// Post method.
//   params:
//     path 	Uri without the base path.
//     body 	Body object.
//     output 	The output object.
//   returns:
//     HTTP Response. Error, if there is an internal error.
func (client *ClientMockup) Post(path string, body interface{}, output interface{}) Response {
    if client.postIndex >= len(client.postF) {
        client.postIndex = 0
    }
    selectedPost := client.postF[client.postIndex]
    client.postIndex ++
    return selectedPost(path, body)
}

// Put method.
//   params:
//     path 	Uri without the base path.
//     body 	Body object.
//     output 	The output object.
//   returns:
//     HTTP Response. Error, if there is an internal error.
func (client *ClientMockup) Put(path string, body interface{}, output interface{}) Response {
    if client.putIndex >= len(client.putF) {
        client.putIndex = 0
    }
    selectedPut := client.putF[client.putIndex]
    client.putIndex ++
    return selectedPut(path, body)
}

// Delete method.
//   params:
//     path Uri without the connection address.
//     output 	The output object.
//   returns:
//     HTTP Response. Error, if there is an internal error.
func (client *ClientMockup) Delete(path string, output interface{}) Response {
    if client.deleteIndex >= len(client.deleteF) {
        client.deleteIndex = 0
    }
    selectedDelete := client.deleteF[client.deleteIndex]
    client.deleteIndex ++
    return selectedDelete(path)
}

// DeleteWithBody method delete with body.
//   params:
//     path Uri without the connection address.
//     body Body to be sent.
//     output 	The output object
//   returns:
//     HTTP Response. Error, if there is an internal error.
func (client *ClientMockup) DeleteWithBody(path string, body interface{}, output interface{}) Response {
    if client.deleteWithBodyIndex >= len(client.deleteWithBodyF) {
        client.deleteIndex = 0
    }
    selectedDelete := client.deleteWithBodyF[client.deleteWithBodyIndex]
    client.deleteWithBodyIndex ++
    return selectedDelete(path, body)
}

// Upload the content stored at the filepath and wait for a given response.
//   params:
//     path         Uri without the connection address.
//     filepath     Target file to be uploaded.
//     fieldname    The name of the field where the file content will be added.
//     output       The output object.
//  returns:
//     HTTP Response. Error, if there is an internal error.
func (client *ClientMockup) Upload(path string, reader io.Reader, filepath string, fieldname string, output interface{}) Response {
    if client.putIndex >= len(client.putF) {
        client.putIndex = 0
    }
    selectedUpload := client.putF[client.putIndex]
    client.putIndex ++
    // Read file content and use it as body
    body, err := ioutil.ReadFile(filepath)
    
    if err != nil {
        return extractResponse(nil, nil, nil, derrors.NewOperationError(errors.IOError, err))
    }
    return selectedUpload(path, string(body))
}

// GetBasePath return the base URL for the requests.
func (client *ClientMockup) GetBasePath() string {
    return "http://localhost:8080"
}

// Reset all the lists and parameters.
func (client *ClientMockup) Reset() {
    client.getIndex = 0
    client.postIndex = 0
    client.deleteIndex = 0
    client.putIndex = 0
    client.deleteWithBodyIndex = 0
    client.getWithBodyIndex = 0

    client.getF = [] func(path string) Response{}
    client.postF = [] func(path string, body interface{}) Response{}
    client.deleteF = [] func(path string) Response{}
    client.getWithBodyF = [] func(path string, body interface{}) Response{}
    client.putF = [] func(path string, body interface{}) Response{}
    client.deleteWithBodyF = [] func(path string, body interface{}) Response{}
}

// AddGet adds GET method.
//   params:
//     f 	GET method function.
func (client *ClientMockup) AddGet(f func(path string) Response) {
    client.getF = append(client.getF, f)
}

// AddGetWithBody adds GET method with body.
//   params:
//     f 	GET method function.
func (client *ClientMockup) AddGetWithBody(f func(path string, body interface{}) Response) {
    client.getWithBodyF = append(client.getWithBodyF, f)
}

// AddPost adds POST method.
//   params:
//     f 	POST method function.
func (client *ClientMockup) AddPost(f func(path string, body interface{}) Response) {
    client.postF = append(client.postF, f)
}

// AddPut adds PUT method.
//   params:
//     f 	POST method function.
func (client *ClientMockup) AddPut(f func(path string, body interface{}) Response) {
    client.putF = append(client.putF, f)
}

// AddDelete adds DELETE method.
//   params:
//     f 	DELETE method function.
func (client *ClientMockup) AddDelete(f func(path string) Response) {
    client.deleteF = append(client.deleteF, f)
}

// AddDeleteWithBody adds DELETE with body method.
//   params:
//     f 	DELETE with body method function.
func (client *ClientMockup) AddDeleteWithBody(f func(path string, body interface{}) Response) {
    client.deleteWithBodyF = append(client.deleteWithBodyF, f)
}

// NewClientMockup is the default builder.
func NewClientMockup() *ClientMockup {
    return &ClientMockup{getIndex: 0, getF: [] func(path string) Response{},
        postIndex: 0, postF: [] func(path string, body interface{}) Response{},
        deleteIndex: 0, deleteF: [] func(path string) Response{},
        deleteWithBodyF: [] func(path string, body interface{}) Response{}}
}
