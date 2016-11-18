# Path to go
GO ?= go

# Binary name
BINARY := vatcheck-svc

# Go stuff
PKG := ./...
LDFLAGS := -s -w

USERNAME := 'gopherskatowice'
TAG := $(shell git describe --tags `git rev-list --tags --max-count=1`)

.PHONY: image test

# The release number & build date are stamped into the binary.
build: LDFLAGS += -X 'main.buildTag=$(TAG)'
build: LDFLAGS += -X 'main.buildDate=$(shell date -u '+%Y/%m/%d %H:%M:%S')'
build:
	@echo "Building $(BINARY) statically"
	cd cmd/vatcheck && GOOS=linux CGO_ENABLED=0 $(GO) build -tags netgo -ldflags "$(LDFLAGS)" -v -o $(BINARY)

check:
	@echo "Checking gometalinter output"
	@ gometalinter --vendor --exclude "dot imports|unkeyed fields" --disable gotype --deadline 30s ./...

# Create Docker image and make sure code is fmt'ed, checked and tested before we build
image: | build
	@echo "Building docker image"
	docker build --rm --force-rm=true --tag=$(USERNAME)/$(BINARY):$(TAG) .

# "go test -i" builds dependencies and installs them into GOPATH/pkg,
# but does not run the tests.
test:
	@ $(GO) test -i `glide novendor`
	@echo "Running tests"
	@! $(GO) test -race `glide novendor` | grep FAIL

default: build
