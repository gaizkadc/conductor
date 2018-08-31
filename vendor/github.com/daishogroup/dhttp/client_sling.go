//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//

package dhttp

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "io"
    "mime/multipart"
    log "github.com/sirupsen/logrus"
    "github.com/daishogroup/derrors"
    "github.com/daishogroup/dhttp/errors"
    "github.com/dghubble/sling"
    "crypto/tls"
)

var logger = log.WithField("package", "dhttp.client_sling")

// ClientSling wrapper structure for the client using the sling library.
type ClientSling struct {
    // Configuration entry
    config RestConfig
}

// Create a new rest sling using a given configuration element.
//  params:
//   conf Configuration object.
//  return:
//   Instantiated rest sling client.
func NewClientSling(conf RestConfig) Client {
    return ClientSling{conf}
}


// Return a preconfigured client using the configuration type.
//  return:
//   Configured client with the corresonding configuration headers.
func (rest ClientSling) getPreconfiguredClient() *sling.Sling {
    var tr http.Transport
    if rest.config.SkipVerification {
        tr = http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
    }else{
        tr =  http.Transport{}
    }

    httpClient :=  &http.Client{Transport: &tr}
    basicClient := sling.New().Client(httpClient).Base(rest.GetBasePath())

    // Set headers
    for header, value := range(rest.config.Headers) {
        basicClient.Set(header, value)
    }

    switch rest.config.ClientType {
    case RestBasic:
        return basicClient
    case RestBasicHTTPS:
        return basicClient
    case RestBasicAuth:
        return basicClient.SetBasicAuth(*rest.config.User,*rest.config.Password)
    case RestOAuth:
        return basicClient.Set("Authorization",fmt.Sprintf("Bearer %s",*rest.config.Secret))
    case RestBasicHTTPSAuth:
        return basicClient.SetBasicAuth(*rest.config.User, *rest.config.Password)
    case RestApiKeyHTTPS:
        return basicClient
    default:
        logger.Errorf("preconfigured client [%s] has no preconfiguration", rest.config.ClientType)
        return nil
    }
}

// Get method.
//   params:
//     path 	Uri without the base path.
//     output 	The output object.
//   returns:
//     HTTP Response. Error, if there is an internal error.
func (rest ClientSling) Get(path string, output interface{}) Response {
    client := rest.getPreconfiguredClient()
    request, err := client.Get(path).Request()
    logger.Debug(request.RequestURI)
    if err != nil {
        return extractResponse(nil, nil, nil, err)
    }
    failure := new(json.RawMessage)
    request.Close = true
    response, err := client.Do(request, output, failure)
    f := string(*failure)
    return extractResponse(output, &f, response, err)
}

// GetRaw is a request that does not require a output type, it returns a string variable.
//  params:
//      path Uri without the connection address
//  returns:
//      HTTP Response.
func (rest ClientSling) GetRaw(path string) Response {
    client := rest.getPreconfiguredClient()
    request, err := client.Get(path).Request()
    logger.Debug(request.RequestURI)
    if err != nil {
        return extractResponse(nil, nil, nil, err)
    }
    failure := new(json.RawMessage)
    result := new(json.RawMessage)
    request.Close = true
    response, err := client.Do(request, result, nil)
    f := string(*failure)
    return extractResponse(string(*result), &f, response, err)
}

// GetWithBody method get with body.
//   params:
//     path 	Uri without the base path.
//     body     Body to be sent.
//     output 	The output object.
//   returns:
//     HTTP Response. Error, if there is an internal error.
func (rest ClientSling) GetWithBody(path string, body interface{}, output interface{}) Response {
    client := rest.getPreconfiguredClient()
    request, err := client.Get(path).BodyJSON(body).Request()
    logger.Debug(request.RequestURI)
    if err != nil {
        return extractResponse(nil, nil, nil, err)
    }
    failure := new(json.RawMessage)
    request.Close = true
    response, err := client.Do(request, output, failure)
    f := string(*failure)
    return extractResponse(output, &f, response, err)
}

// Post method.
//   params:
//     path 	Uri without the base path.
//     body 	Body object.
//     output 	The output object.
//   returns:
//     HTTP Response. Error, if there is an internal error.
func (rest ClientSling) Post(path string, body interface{}, output interface{}) Response {
    client := rest.getPreconfiguredClient()
    request, err := client.Post(path).BodyJSON(body).Request()
    logger.Debug(request.RequestURI)
    if err != nil {
        return extractResponse(nil, nil, nil, err)
    }
    failure := new(json.RawMessage)
    request.Close = true
    response, err := client.Do(request, output, failure)
    f := string(*failure)
    return extractResponse(output, &f, response, err)
}

// Put method.
//   params:
//     path 	Uri without the base path.
//     body 	Body object.
//     output 	The output object.
//   returns:
//     HTTP Response. Error, if there is an internal error.
func (rest ClientSling) Put(path string, body interface{}, output interface{}) Response {
    client := rest.getPreconfiguredClient()
    request, err := client.Put(path).BodyJSON(body).Request()
    logger.Debug(request.RequestURI)
    if err != nil {
        return extractResponse(nil, nil, nil, err)
    }
    failure := new(json.RawMessage)
    request.Close = true
    response, err := client.Do(request, output, failure)
    f := string(*failure)
    return extractResponse(output, &f, response, err)
}

// Delete method.
//   params:
//     path Uri without the connection address.
//     output 	The output object
//   returns:
//     HTTP Response. Error, if there is an internal error.
func (rest ClientSling) Delete(path string, output interface{}) Response {
    client := rest.getPreconfiguredClient()
    request, err := client.Delete(path).Request()
    logger.Debug(request.RequestURI)
    if err != nil {
        return extractResponse(nil, nil, nil, err)
    }
    failure := new(json.RawMessage)
    request.Close = true
    response, err := client.Do(request, output, failure)
    f := string(*failure)
    return extractResponse(output, &f, response, err)
}

// DeleteWithBody method delete with body.
//   params:
//     path Uri without the connection address.
//     body Body to be sent.
//     output 	The output object
//   returns:
//     HTTP Response. Error, if there is an internal error.
func (rest ClientSling) DeleteWithBody(path string, body interface{}, output interface{}) Response {
    client := rest.getPreconfiguredClient()
    request, err := client.Delete(path).BodyJSON(body).Request()
    logger.Debug(request.RequestURI)
    if err != nil {
        return extractResponse(nil, nil, nil, err)
    }
    failure := new(json.RawMessage)
    request.Close = true
    response, err := client.Do(request, output, failure)
    f := string(*failure)
    return extractResponse(output, &f, response, err)
}


// Upload the content stored at the filepath and wait for a given response.
//   params:
//     path         Uri without the connection address.
//     body         If body is provided, use it as a reader, otherwise read the content from filename
//     filepath     Target file to be uploaded.
//     fieldname    The name of the field where the file content will be added.
//     output       The output object.
//  returns:
//     HTTP Response. Error, if there is an internal error.
func (rest ClientSling) Upload(path string, body io.Reader, filename string, fieldname string, output interface{}) Response {
    if body == nil {
        // Open it
        f, err := os.Open(filename)
        if  err != nil {
            errToReturn := derrors.NewOperationError(errors.OpeningFileError,err)
            return extractResponse(nil, nil, nil, errToReturn)
        }

        // Go for it
        fi, err := f.Stat()
        if err != nil {
            errToReturn := derrors.NewOperationError(errors.OpeningFileError,err)
            return extractResponse(nil, nil, nil, errToReturn)
        }
        body = f
        filename = fi.Name()
    } else {
        filename = filepath.Base(filename)
    }

    r, w := io.Pipe()
    mpw := multipart.NewWriter(w)
    // We need to do this asynchronously in a go routine, because we need
    // to read and write concurrently through the pipe
    errChan := make(chan error, 1)
    go func() {
        defer w.Close()
        defer mpw.Close()
        part, err := mpw.CreateFormFile(fieldname, filename)
        if err != nil {
            errChan <- err
            return
        }
        _, err = io.Copy(part, body)
        if err != nil {
            errChan <- err
            return
        }
        errChan <- nil

        return
    }()

    // Get our preconfigured client.
    client := rest.getPreconfiguredClient()

    request, err := client.Put(path).Set("Content-Type", mpw.FormDataContentType()).Body(r).Request()
    if err != nil {
        return extractResponse(nil, nil, nil, err)
    }

    failure := new(json.RawMessage)
    request.Close=true
    response, err := client.Do(request, output,failure)
    fstring := string(*failure)

    // Check for reader error
    err = <-errChan
    if err != nil {
        channelErr := derrors.NewOperationError(errors.IOError, err)
        close(errChan)
        return extractResponse(nil, nil, nil, channelErr)
    }
    close(errChan)

    // Everything was OK, return what we got
    return extractResponse(output, &fstring, response, err)

}



// GetBasePath return the current base path of the client.
//   returns:
//     The base path url.
func (rest ClientSling) GetBasePath() string {
    if rest.config.HTTPS {
        return fmt.Sprintf("https://%s:%d", rest.config.Host, rest.config.Port)
    }
    return fmt.Sprintf("http://%s:%d", rest.config.Host, rest.config.Port)
}

func extractError(failure *string, statusCode int) derrors.DaishoError {
    dError := &derrors.GenericError{}
    err := json.Unmarshal([]byte(*failure), dError)
    if err != nil {
        return derrors.NewConnectionError(fmt.Sprintf("HTTP Error [%d]: %s", statusCode, *failure)).
            WithParams(statusCode)
    }
    return dError.WithParams(statusCode)
}

// extractResponse is a helper to create a Response object.
//   params:
//     result 	Result object.
//     response HTTP Response object.
//     err 		Error description.
//   returns:
//     HTTP Response. Error, if there is an internal error.
func extractResponse(result interface{}, failure *string, response *http.Response, err error) Response {

    if err != nil || response == nil {
        return NewResponse(nil, nil, derrors.AsDaishoError(err, errors.ConnectionError))
    }

    // HTTP status OK
    if response.StatusCode >= 200 && response.StatusCode < 300 {
        // Get the incoming headers and return it.
        return NewResponseWithHeader(result, &response.StatusCode, response.Header, nil)
    }

    //HTTP status FAIL
    dError := extractError(failure, response.StatusCode)
    return NewResponse(nil, &response.StatusCode, dError)

}

