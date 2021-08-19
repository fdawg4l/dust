SHELL=/bin/bash

GO ?= go
GOVERSION ?= go1.11
BIN ?= bin
IGNORE := $(shell mkdir -p $(BIN))

# used to get the short commit id
GITREF := $(shell git describe --always --long --dirty)
GOBUILD := CGO_ENABLED=0 $(GO) build $(RACE) -ldflags "-X github.com/fdawg4l/dust/pkg/build.GitCommitID=$(GITREF)"
.PHONY: all tools clean test check distro \
	goversion goimports gopath govet gofmt golint

.DEFAULT_GOAL := all

ifeq ($(ENABLE_RACE_DETECTOR),true)
	RACE := -race
else
	RACE :=
endif

# utility targets
goversion:
	@echo Checking go version...
	@( $(GO) version | grep -q $(GOVERSION) ) || ( echo "Please install $(GOVERSION) (found: $$($(GO) version))" && exit 1 )

dist: dustd
all: test dist

tools: $(GOIMPORTS)

goimports: $(GOIMPORTS)
$(GOIMPORTS):
	@echo Building $(GOIMPORTS)...
	@$(GO) get -u golang.org/x/tools/cmd/goimports

govet: goversion $(GOIMPORTS)
	@echo Checking go vet...
	@$(GO) tool vet -all -shadow -lostcancel -tests $$(find . -mindepth 1 -maxdepth 1 -type d -not -name vendor -not -name api)
	@echo Checking go imports...
	@$(GOIMPORTS) -local semifreddo -d $$(find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./api/*") 2>&1

dustd := $(BIN)/dustd
dustd : $(dustd)
$(dustd):
	@echo Building server ${dustd}
	$(GOBUILD) -o ${dustd} cmd/dustd/main.go

clean:
	@echo Removing build output directory $(CURDIR)/bin
	@rm -rf $(CURDIR)/bin

test: dist
	$(GO) test -v ./...

docker: dist
	@docker build -t dustd .
	@docker save dustd > dustd-$(GITREF).tar
