# Path to go
GO ?= go

# Binary name
BINARY := nx-vatcheck

# Go stuff
PKG := ./...
LDFLAGS := -s -w

TAG := $(shell git describe --tags `git rev-list --tags --max-count=1`)

.PHONY: image check test

# The release number & build date are stamped into the binary.
build: LDFLAGS += -X 'main.buildTag=$(TAG)'
build: LDFLAGS += -X 'main.buildDate=$(shell date -u '+%Y/%m/%d %H:%M:%S')'
build:
	@echo "Building $(BINARY) statically"
	cd cmd/vatcheck && GOOS=linux CGO_ENABLED=0 $(GO) build -a -ldflags "$(LDFLAGS)" -v -o $(BINARY)

# Create Docker image and make sure code is fmt'ed, checked and tested before we build
image: | check build
	@echo "Building docker image"
	docker build --rm --force-rm=true --tag=registry.nexway.build/$(BINARY):$(TAG) .

# Run several automated source checks to get rid of the most simple issues.
# This helps keeping code review focused on application logic.
# github.com/alecthomas/gometalinter
check:
	@echo "Checking gometalinter output"
	@! gometalinter $(PKG) | \
	  grep -vE '(conf\/yaml\.go)'

# "go test -i" builds dependencies and installs them into GOPATH/pkg,
# but does not run the tests.
test:
	@ $(GO) test -i $(PKG)
	@echo "Running tests"
	@! $(GO) test $(PKG) | grep FAIL

default: build
