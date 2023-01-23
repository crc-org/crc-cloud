VERSION ?= 0.0.1
CONTAINER_MANAGER ?= podman
# Image URL to use all building/pushing image targets
IMG ?= quay.io/crcont/crc-cloud:${VERSION}

# Go and compilation related variables
GOPATH ?= $(shell go env GOPATH)
BUILD_DIR ?= out
SOURCE_DIRS = cmd pkg test
# https://golang.org/cmd/link/
LDFLAGS := $(VERSION_VARIABLES) -extldflags='-static' ${GO_EXTRA_LDFLAGS}
GCFLAGS := all=-N -l

# Add default target
.PHONY: default
default: install

# Create and update the vendor directory
.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: check
check: build test lint

# Start of the actual build targets

.PHONY: install
install: $(SOURCES)
	go install -ldflags="$(LDFLAGS)" $(GO_EXTRA_BUILDFLAGS) ./cmd

$(BUILD_DIR)/crc-cloud: $(SOURCES)
	GOOS=linux GOARCH=amd64 go build -gcflags="$(GCFLAGS)" -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/crc-cloud $(GO_EXTRA_BUILDFLAGS) ./cmd


 
.PHONY: build 
build: clean $(BUILD_DIR)/crc-cloud

.PHONY: test
test:
	go test -race --tags build -v -ldflags="$(VERSION_VARIABLES)" ./pkg/... ./cmd/...

.PHONY: clean ## Remove all build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f $(GOPATH)/bin/crc-cloud

.PHONY: fmt
fmt:
	@gofmt -l -w $(SOURCE_DIRS)

$(GOPATH)/bin/golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2

# Run golangci-lint against code
.PHONY: lint
lint: $(GOPATH)/bin/golangci-lint
	$(GOPATH)/bin/golangci-lint run

# Build the container image
.PHONY: oci-build
oci-build: 
	${CONTAINER_MANAGER} build -t ${IMG} -f oci/Containerfile .

# Push the container image
.PHONY: oci-push
container-push:
	${CONTAINER_MANAGER} push ${IMG}
	
