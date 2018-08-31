//
// Copyright (C) 2018 Daisho Group - All Rights Reserved
//

package dhttp

import (
    "net"
    "github.com/daishogroup/derrors"
    "net/http"
    "encoding/json"
    "github.com/daishogroup/dhttp/errors"
    "fmt"
    "time"
)

// GetAvailablePort obtains a free HTTP port.
func GetAvailablePort() (int, error) {
    listener, err := net.Listen("tcp", ":0")
    defer listener.Close()
    if err != nil {
        return -1, err
    }
    return listener.Addr().(*net.TCPAddr).Port, nil
}

// GetURL generates a valid URL.
func GetURL(host string, port int, path string) string {
    return fmt.Sprintf("http://%s:%d%s", host, port, path)
}

// RespondWithError sends an error response.
//   params:
//     w The response writer.
//     code The HTTP response code.
//     error The error message to be sent as JSON response.
func RespondWithError(w http.ResponseWriter, code int, error derrors.DaishoError) {
    RespondWithJSON(w, code, error)
}

// RespondWithJSON sends a JSON as a respond.
//   params:
//     w The response writer.
//     code The HTTP response code.
//     payload The JSON payload.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, err := json.Marshal(payload)
    if err != nil {
        RespondWithError(w, http.StatusInternalServerError, derrors.NewOperationError(errors.MarshalError, err))
    } else {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(code)
        w.Write(response)
    }
}


// It responds with a JSON body and setting a bunch of headers.
//   params:
//     w The response writer.
//     code The HTTP response code.
//     payload The JSON payload.
//     headers The response headers.
func RespondWithHeaderJSON(w http.ResponseWriter, code int, payload interface{}, headers http.Header) {
    response, err := json.Marshal(payload)
    if err != nil {
        RespondWithError(w, http.StatusInternalServerError, derrors.NewOperationError(errors.MarshalError, err))
    } else {
        w.Header().Set("Content-Type", "application/json")
        if headers != nil {
            for name, v := range headers {
                for _, content := range v {
                    w.Header().Set(name, content)
                }
            }
        }
        w.WriteHeader(code)
        w.Write(response)
    }
}

// Check iteratively if a given URL is available or not. This only means that the server is up, running and answering
// requests.
// params:
//  host URL host
//  port the port number
//  url the url
//  tries number of attempt to run
//  seconds to wait between calls
// return:
//  true if the operation was successful
func WaitURLAvailable(host string, port int, tries int, url string, seconds int) bool {
    var i = 0
    var exit = false
    for !exit && i < tries {
        response, _ := http.Get(fmt.Sprintf("http://%s:%d/%s",host,port,url))
        if response != nil {
            exit = true
        } else {
            time.Sleep(time.Duration(seconds) * time.Second)
        }
        i++
    }
    if !exit {
        panic(derrors.NewConnectionError("Server is not launched"))
        return false
    }
    return true
}
