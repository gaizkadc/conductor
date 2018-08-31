# derrors - Daisho Errors

This repository contains the definition of the error for Daisho components.

## General overview

The main purpose of this repository is to improve error reporting facilitating the communication of error states to the
users and allowing deeper reporting of the errors for the developers at the same time.

The DaishoError interfaces defines a set of basic methods that makes a DaishoError compatible with the GolangError but
provides extra functions to track the error origin.

## Building and testing

To update the files, run:

```
'bazel run //:gazelle
```

To build the project, execute:

```
bazel build ...
```

To pass the tests,

```
bazel test ...
```

## How to use it

Defining a new error automatically extracts the StackTrace

```
return derrors.NewEntityError(descriptor, errors.NetworkDoesNotExists, err)
```

Will print the following message when calling Error().

```
[Operation] network does not exists
```

And a detailed one when calling DebugReport().

```
[Operation] network does not exists
Parameters:
P0: []interface {}{"n1b59e008-a9f2-4a25-866a-ace0cabc38b2asdf"}

StackTrace:
ST0: github.com/daishogroup/system-model/server/cluster.(*Manager).AddCluster - server/cluster/manager.go:49
ST1: github.com/daishogroup/system-model/server/cluster.(*Handler).addCluster - server/cluster/handler.go:73
ST2: github.com/daishogroup/system-model/server/cluster.(*Handler).(github.com/daishogroup/system-model/server/cluster.addCluster)-fm - server/cluster/handler.go:44
ST3: net/http.HandlerFunc.ServeHTTP - /private/var/tmp/_bazel_daniel/8d1ee8965f258357de92e48c44b3aaeb/external/go1_8_3_darwin_amd64/src/net/http/server.go:1943
ST4: github.com/gorilla/mux.(*Router).ServeHTTP - vendor/github.com/gorilla/mux/mux.go:151
ST5: net/http.serverHandler.ServeHTTP - /private/var/tmp/_bazel_daniel/8d1ee8965f258357de92e48c44b3aaeb/external/go1_8_3_darwin_amd64/src/net/http/server.go:2569
ST6: net/http.(*conn).serve - /private/var/tmp/_bazel_daniel/8d1ee8965f258357de92e48c44b3aaeb/external/go1_8_3_darwin_amd64/src/net/http/server.go:1826
ST7: runtime.goexit - /private/var/tmp/_bazel_daniel/8d1ee8965f258357de92e48c44b3aaeb/external/go1_8_3_darwin_amd64/src/runtime/asm_amd64.s:2198
```

## But wait, why not call it errors?

We have intentionally avoided the errors package name to avoid conflicts with the golang error package.

## What about internationalization?

The current version does not provide internationalization capabilities for the output messages. However, given that this
repository contains a set of predefined messages, integrating that support in the future should be easy.