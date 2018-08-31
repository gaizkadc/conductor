# dhttp
DHTTP is an utility library to simplify the utilization of http client. It comes with
a set of pre-configured http clients:
* Basic: common HTTP
* Basic auth: includes the authorization header with user and password
* OAuth2: includes the authentication header with the user secret.
* BasicHTTPS: common HTTPS

## Security
DHTTP simplifies the utilization of secure connections using https. The current version
sets HTTPS connections in the following configuration flavors:
* OAuth2
* BasicHTTPS

## Example
```go
// Set the preset configuration.
conf := NewRestOAuthConfig("172.28.128.4", 30000, "amazingsecret")
// Instantiate the client
client := NewClientSling(conf)
//...
// Use the client
response := client.Get("/api/v0/network/list")
```

Or, using a URL:
```go
// Set the preset configuration.
conf, err := NewRestURLConfig("https://user:pass@172.28.128.4:9999")
if err != nil {
    return err
}
// Instantiate the client
client := NewClientSling(conf)
```
