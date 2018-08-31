//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//

package dhttp

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "io"
    "net/http"

    "github.com/daishogroup/derrors"
    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
)

var loggerTestHandler = log.WithField("package", "client")
// TestBody is a struct thar return the HTTP handler.
type TestBody struct {
    Name  string `json:"name, omitempty"`
    Value int    `json:"value, omitempty"`
    Check bool   `json:"check, omitempty"`
}

// IsValid checks if the struct is correct.
func (tb *TestBody) IsValid() bool {
    return tb.Name != "" && tb.Value != 0 && tb.Check != false
}

// NewValidTestBody generates a valid object.
func NewValidTestBody() *TestBody {
    return NewTestBody("valid", 1, true)
}

// NewInvalidTestBody generates a invalid object.
func NewInvalidTestBody() *TestBody {
    return NewTestBody("invalid", 0, false)
}

// NewTestBody is the basic constructor of the object.
func NewTestBody(name string, value int, check bool) *TestBody {
    return &TestBody{Name: name, Value: value, Check: check}
}

// TestHandler is a basic HTTP server.
type TestHandler struct {
    basePath string
    Err      bool
}

// SetRoutes generate all routes of the server.
func (h *TestHandler) SetRoutes(router *mux.Router) {
    loggerTestHandler.Debug("Setting cluster routes")
    router.HandleFunc("/test/post", h.post).Methods("POST")
    router.HandleFunc("/test/put", h.put).Methods("PUT")
    router.HandleFunc("/test/get", h.get).Methods("GET")
    router.HandleFunc("/test/getWithBody", h.getWithBody).Methods("GET")
    router.HandleFunc("/test/delete", h.delete).Methods("DELETE")
    router.HandleFunc("/test/deleteWithBody", h.deleteWithBody).Methods("DELETE")
    router.HandleFunc("/test/upload", h.upload).Methods("PUT")
    router.HandleFunc("/test/getWithHeaders", h.getWithHeaders).Methods("GET")
}

// post is a new user to an existing cluster.
//   params:
//     w The HTTP writer structure.
//     r The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (h *TestHandler) post(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    loggerTestHandler.Debug("post")
    tb := &TestBody{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        loggerTestHandler.Error(err)
        RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(err.Error()))
    } else {
        err := json.Unmarshal(b, tb)
        if err == nil && tb.IsValid() {
            RespondWithJSON(w, http.StatusOK,
                NewSuccessfulOperation(fmt.Sprintf("post: %v", tb)))
        } else {
            RespondWithError(w, http.StatusBadRequest,
                derrors.NewOperationError(fmt.Sprintf("error: %v", tb)))
        }
    }
}

// put is a new user to an existing cluster.
//   params:
//     w The HTTP writer structure.
//     r The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (h *TestHandler) put(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    loggerTestHandler.Debug("put")
    tb := &TestBody{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        loggerTestHandler.Error(err)
        RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(err.Error()))
    } else {
        err = json.Unmarshal(b, tb)
        if err == nil && tb.IsValid() {
            RespondWithJSON(w, http.StatusOK,
                NewSuccessfulOperation(fmt.Sprintf("put: %v", tb)))
        } else {
            RespondWithError(w, http.StatusBadRequest,
                derrors.NewOperationError(fmt.Sprintf("put: %v", tb)))
        }
    }
}

// get is a new user to an existing cluster.
//   params:
//     w The HTTP writer structure.
//     r The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (h *TestHandler) get(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    loggerTestHandler.Debug("get")

    // Check for daisho-test header
    header := r.Header.Get("X-Daisho-Test")
    headerAdd := ""
    if header != "" {
        headerAdd = fmt.Sprintf(" X-Daisho-Test: %s", header)
    }

    if h.Err {
        RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError("error get"))
    } else {
        RespondWithJSON(w, http.StatusOK, NewSuccessfulOperation("get" + headerAdd))
    }

}

// get is a new user to an existing cluster with a bunch of headers.
//   params:
//     w The HTTP writer structure.
//     r The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (h *TestHandler) getWithHeaders(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    loggerTestHandler.Debug("getWithHeaders")

    // Check for daisho-test header

    header := r.Header.Get("X-Daisho-Test")
    headerAdd := ""
    if header != "" {
        headerAdd = fmt.Sprintf(" X-Daisho-Test: %s", header)
    }

    testHeaders := http.Header{}
    testHeaders.Set("TestingHeader", "TestingValue")

    if h.Err {
        RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError("error get"))
    } else {
        RespondWithHeaderJSON(w, http.StatusOK, NewSuccessfulOperation("get" + headerAdd), testHeaders)
    }

}

// getWithBody is a new user to an existing cluster.
//   params:
//     w The HTTP writer structure.
//     r The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (h *TestHandler) getWithBody(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    loggerTestHandler.Debug("getWithBody")
    tb := &TestBody{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        loggerTestHandler.Error(err)
        RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(err.Error()))
    } else {
        err = json.Unmarshal(b, tb)
        if err == nil && tb.IsValid() {
            RespondWithJSON(w, http.StatusOK,
                NewSuccessfulOperation(fmt.Sprintf("getWithBody: %v", tb)))
        } else {
            RespondWithError(w, http.StatusBadRequest,
                derrors.NewOperationError(fmt.Sprintf("getWithBody: %v", tb)))
        }
    }
}

// delete is a new user to an existing cluster.
//   params:
//     w The HTTP writer structure.
//     r The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (h *TestHandler) delete(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    loggerTestHandler.Debug("deleteWithBody")

    if h.Err {
        RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError("error delete"))
    } else {
        RespondWithJSON(w, http.StatusOK, NewSuccessfulOperation("delete"))
    }
}

// deleteWithBody is a new user to an existing cluster.
//   params:
//     w The HTTP writer structure.
//     r The HTTP request structure.
//   returns:
//     A 200 OK if the operation is successful or the error.
func (h *TestHandler) deleteWithBody(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    loggerTestHandler.Debug("deleteWithBody")
    tb := &TestBody{}
    b, err := ioutil.ReadAll(r.Body)
    if err != nil {
        loggerTestHandler.Error(err)
        RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError(err.Error()))
    } else {
        err := json.Unmarshal(b, tb)
        if err == nil && tb.IsValid() {
            RespondWithJSON(w, http.StatusOK,
                NewSuccessfulOperation(fmt.Sprintf("deleteWithBody: %v", tb)))
        } else {
            RespondWithError(w, http.StatusBadRequest,
                derrors.NewOperationError(fmt.Sprintf("deleteWithBody: %v", tb)))
        }
    }
}


func (h *TestHandler) upload(w http.ResponseWriter, r *http.Request) {
    loggerTestHandler.Info("upload")

    defer r.Body.Close()

    tb := &TestBody{}
    reader, err := r.MultipartReader()
    if err != nil {
        RespondWithError(w, http.StatusBadRequest, derrors.NewOperationError("impossible to get multipart reader", err))
    }

    for {
        p, err := reader.NextPart()
        if err == io.EOF {
            return
        }
        if err != nil {
            log.Fatal(err)
        }
        _, err = ioutil.ReadAll(p)
        if err != nil {
            log.Fatal(err)
        }
    }

    RespondWithJSON(w, http.StatusOK,
        NewSuccessfulOperation(fmt.Sprintf("upload: %v", tb)))
}

// NewTestHandler is the basic constructor of the TestHandler.
func NewTestHandler() *TestHandler {
    return &TestHandler{}
}
