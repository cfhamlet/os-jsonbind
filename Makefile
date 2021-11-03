PROJECT     := github.com/cfhamlet/os-jsonbind
GOPATH      := $(shell go env GOPATH)
GOIMPORTS   := $(GOPATH)/bin/goimports
GOFUMPT     := $(GOPATH)/bin/gofumpt
GOLINT      := $(GOPATH)/bin/golangci-lint
SRC         := $(shell find . -type f -name '*.go' -print)
PKG         := ./...
TESTS       := .
TESTFLAGS   :=
GOFLAGS     :=
LDFLAGS     :=

.PHONY: all
all: test

.PHONY: test
test: TESTFLAGS += -race -v
test: test-lint
test: test-unit
test: test-bench
test: test-cover

.PHONY: test-lint
test-lint:$(GOLINT)
	@echo
	@echo "==> Running lint test <=="
	GO111MODULE=on $(GOLINT) run $(TESTS)

.PHONY: test-unit
test-unit:
	@echo
	@echo "==> Running unit tests <=="
	GO111MODULE=on go test -v $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS)

.PHONY: test-bench
test-bench:
	@echo
	@echo "==> Running benchmark tests <=="
	GO111MODULE=on go test $(GOFLAGS) -bench $(TESTS) $(PKG)

.PHONY: test-cover
test-cover:
	@echo
	@echo "==> Running unit tests with coverage <=="
	@./scripts/coverage.sh

.PHONY: format
format: $(GOIMPORTS) $(GOFUMPT)
	@echo
	@echo "==> Formatting <=="
	GO111MODULE=on go list -f '{{.Dir}}' ./... | xargs $(GOIMPORTS) -w 
	GO111MODULE=on go list -f '{{.Dir}}' ./... | xargs $(GOFUMPT) -w 

$(GOFUMPT):
	@echo 
	@echo "==> Installing gofumpt <=="
	(cd /; GO111MODULE=on go get -u mvdan.cc/gofumpt)

$(GOIMPORTS):
	@echo
	@echo "==> Installing goimports <=="
	(cd /; GO111MODULE=on go get -u golang.org/x/tools/cmd/goimports)

$(GOLINT):
	@echo
	@echo "==> Installing golangci-lint <=="
	(cd /; curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin latest)

.PHONY cloc:
cloc:
	@ cloc .