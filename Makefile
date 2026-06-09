.PHONY: build update-deps compile install test testacc fmt lint clean tools docs

# Set the default goal to 'build' so that running 'make' will build the provider
.DEFAULT_GOAL := build

# Version: single source of truth for the provider version.
# Override from CLI: make build VERSION=0.2.0
VERSION ?= 0.1.5

# Disable VCS stamping so builds work in repositories without git metadata available to Go.
GOFLAGS += -buildvcs=false
export GOFLAGS

# Provider registry address (must match main.go ServeOpts.Address)
REGISTRY = registry.terraform.io/HPE/opsramp

BINARY = terraform-provider-opsramp
TESTACC_RUN ?= TestAcc
TESTACC_TIMEOUT ?= 45m
TESTACC_PATH ?= ./internal/resources/...

# Detect OS and architecture via Go so it works on both platforms
GOOS   := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

ifeq ($(GOOS),windows)
  PLUGIN_ARCH=windows_amd64
  EXE = .exe
  PLUGIN_BASE = $(APPDATA)/terraform.d/plugins
else
  PLUGIN_ARCH=linux_amd64
  EXE =
  PLUGIN_BASE = $(HOME)/.terraform.d/plugins
endif

PLUGIN_DIR = $(PLUGIN_BASE)/$(REGISTRY)/$(VERSION)/$(GOOS)_$(GOARCH)

# Update dependencies to latest versions, tidy go.mod, and download all modules
update-deps:
	go get -u ./...
	go mod tidy
	go mod download all

compile:
	go build -ldflags="-X 'github.com/HPE/terraform-provider-opsramp.Version=$(VERSION)'" \
		-o bin/$(BINARY)$(EXE) .

# Build = deps + compile
build: update-deps compile

# Install the provider to the local Terraform plugins directory
install: build

ifeq ($(GOOS),windows)
	@if not exist "$(subst /,\,$(PLUGIN_DIR))" mkdir "$(subst /,\,$(PLUGIN_DIR))"
	copy "bin\$(BINARY)$(EXE)" "$(subst /,\,$(PLUGIN_DIR))\$(BINARY)$(EXE)"
else
	mkdir -p "$(PLUGIN_DIR)"
	cp "bin/$(BINARY)$(EXE)" "$(PLUGIN_DIR)/$(BINARY)$(EXE)"
endif

# Run all unit tests
test:
	go test -v ./...

# Run acceptance tests only (requires TF_ACC=1 and optional API env vars)
testacc: export TF_ACC = 1
testacc: 
	go test -count=1 -v -run "$(TESTACC_RUN)" -timeout "$(TESTACC_TIMEOUT)" $(TESTACC_PATH)	

# Install project tooling used for provider docs generation.
# Uses go install @latest so go.mod is never modified.
tools:
	go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest

# Generate provider documentation from schema descriptions
# --provider-name overrides the default (folder-derived) name so schema lookups resolve correctly.
# Calls the binary installed by 'make tools' to avoid touching go.mod.
docs: tools
	tfplugindocs generate --provider-name opsramp

# Format all Go and Terraform code
fmt:
	go fmt ./...
	terraform fmt -recursive ./examples/

# Lint the code (requires golangci-lint)
lint:
	golangci-lint run ./...

# Clean build artifacts
clean:
ifeq ($(GOOS),windows)
	@if exist bin rmdir /s /q bin
else
	rm -rf bin/
endif

