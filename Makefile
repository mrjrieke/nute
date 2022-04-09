GIT ?= git
GO_VARS ?=
GO ?= go
COMMIT := $(shell $(GIT) rev-parse HEAD)
VERSION ?= $(shell $(GIT) describe --tags ${COMMIT} 2> /dev/null || echo "$(COMMIT)")
BUILD_TIME := $(shell LANG=en_US date +"%F_%T_%z")
ROOT := .
LD_FLAGS := -X $(ROOT).Version=$(VERSION) -X $(ROOT).Commit=$(COMMIT) -X $(ROOT).BuildTime=$(BUILD_TIME)
GOBIN ?= ./bin

.PHONY: help clean 
help:
	@echo "Please use \`make <ROOT>' where <ROOT> is one of"
	@echo "  dependencies to go install the dependencies"
	@echo "  nute   to build the main binary for current platform"
	@echo "  test         to run unittests"

depend:
	go mod tidy

clean:
	rm -f bin

helloworldgio: */*.go */*/*.go */*/*/*.go
	$(GO_VARS) $(GO) build -tags "gioboot" -o="$(ROOT)/examples/helloworld/bin/hellogio" -ldflags="$(LD_FLAGS)" $(ROOT)/examples/helloworld/hellogio/main.go
	$(GO_VARS) $(GO) build -tags "g3nboot" -o="$(ROOT)/examples/helloworld/bin/worldg3n" -ldflags="$(LD_FLAGS)" $(ROOT)/examples/helloworld/worldg3n/main.go

helloworldfyne: */*.go */*/*.go */*/*/*.go
	$(GO_VARS) $(GO) build -tags "fyneboot" -o="$(ROOT)/examples/helloworld/bin/hellofyne" -ldflags="$(LD_FLAGS)" $(ROOT)/examples/helloworld/hellofyne/main.go
	$(GO_VARS) $(GO) build -tags "g3nboot" -o="$(ROOT)/examples/helloworld/bin/worldg3n" -ldflags="$(LD_FLAGS)" $(ROOT)/examples/helloworld/worldg3n/main.go

helloworldmobile: */*.go */*/*.go */*/*/*.go
	$(GO_VARS) $(GO) build -o="$(ROOT)/examples/helloworld/bin/hellogio" -ldflags="$(LD_FLAGS)" $(ROOT)/examples/helloworld/hellogio/main.go
	$(GO_VARS) $(GO) build -o="$(ROOT)/examples/helloworld/bin/worldg3n" -ldflags="$(LD_FLAGS)" $(ROOT)/examples/helloworld/worldmobile/main.go

cleanmsdk:
	rm mashupsdk/mashupsdk_grpc.pb.go; rm mashupsdk/mashupsdk.pb.go

mashupsdk: */*.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative mashupsdk/mashupsdk.proto

nute: *.go */*.go */*/*.go
	$(GO_VARS) $(GO) build -o="$(GOBIN)/nute" -ldflags="$(LD_FLAGS)" $(ROOT)/main.go

test: *.go */*.go */*/*.go
	$(GO_VARS) $(GO) test -v .
	$(GO_VARS) $(GO) test -v ./tests
