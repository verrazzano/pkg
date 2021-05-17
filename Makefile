# Copyright (c) 2021, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

GO ?= GO111MODULE=on go

.PHONY: go-build
go-build:
	$(GO) build ./...

.PHONY: go-test
go-test: go-build
	$(GO) test ./...

.PHONY: go-lint
go-lint: install-linter
	golangci-lint run -E misspell,gofmt,goimports ./...

.PHONY: install-linter
install-linter:
ifeq (, $(shell command -v golangci-lint))
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.38.0
	$(eval LINTER=$(GOPATH)/bin/golangci-lint)
else
	$(eval LINTER=$(shell command -v golangci-lint))
endif