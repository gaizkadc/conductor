#
# Copyright (C) 2018 Nalej Group - All Rights Reserved
#
# Makefile for Nalej projects. It provides build, test, and package targets.
#

# Name of the target applications to be built
TARGETS=conductor

# Name of the components that will be build and packaged as Daisho core packages.
COMPONENTS=conductor

NAME=conductor

# Target directory to store binaries and results
TARGET=conductor

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

# Build information
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Use ldflags to pass commit and branch information
# TODO: Integrate this into the compilation process
LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

COVERAGE_FILE=$(TARGET)/coverage.out

.PHONY: all
all: dep build test

.PHONY: dep
dep:
	$(info >>> Updating dependencies...)
	dep ensure

.PHONY: test test-race test-coverage
test:
	$(info >>> Launching tests...)
	$(GOTEST) ./...

test-race:
	$(info >>> Launching tests... (Race detector enabled))
	$(GOTEST) -race ./...

test-coverage:
    $(info >>> Launching tests... (Coverage enabled))
    $(GOTEST) -coverprofile=$(COVERAGE_FILE) -covermode=atomic  ./...

.PHONY: build-all build build-linux
build-all: build build-linux

build:
	$(info >>> Building ...)
	for app in $(TARGETS); do \
            $(GOBUILD) -o $(TARGET)/"$$app" ./cmd/"$$app" ; \
	done

# Cross compilation to obtain a linux binary
build-linux:
	$(info >>> Bulding for Linux...)
	for app in $(TARGETS); do \
    	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(TARGET)/bin/linux_amd64/"$$app" ./cmd/"$$app" ; \
	done

# Package all images and components
.PHONY: package package-create-dir create-package
package: build-linux package-create-dir create-package

package-create-dir:
	mkdir -p $(TARGET)/images
	mkdir -p $(TARGET)/packages

create-package:
	$(info >>> Packaging ...)
	for component in $(COMPONENTS); do \
        docker build -t apiserver:5000/nalej/"$$component" -f components/"$$component"/Dockerfile $(TARGET)/ ; \
        mkdir $(TARGET)/images/"$$component" ; \
        docker save apiserver:5000/nalej/"$$component" > $(TARGET)/images/"$$component"/image.tar ; \
        cp components/"$$component"/component.yaml $(TARGET)/images/"$$component"/. ; \
        cd $(TARGET)/images/"$$component"/ && tar cvzf core-"$$component".tar.gz * && cd - ; \
        mv $(TARGET)/images/"$$component"/core-"$$component".tar.gz $(TARGET)/packages ; \
    done

# Check the codestyle using gometalinter
.PHONY: checkstyle
checkstyle:
	gometalinter --disable-all --enable=golint --enable=vet --enable=errcheck --enable=goconst --vendor ./...

.PHONY: clean
clean:
	$(info >>> Cleaning project...)
	$(GOCLEAN)
	rm -Rf $(TARGET)
