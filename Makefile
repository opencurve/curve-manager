.PHONY: init build debug test clean image

tag?= "curve-manager:unknown"

# go env
GOPROXY     := "https://goproxy.cn,direct"
GOOS        := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
GOARCH      := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))
CGO_LDFLAGS := "-static"
CC          := musl-gcc

GOENV := GO111MODULE=on
GOENV += GOPROXY=$(GOPROXY)
GOENV += CC=$(CC)
GOENV += CGO_ENABLED=1 CGO_LDFLAGS=$(CGO_LDFLAGS)
GOENV += GOOS=$(GOOS) GOARCH=$(GOARCH)

# go
GO := go

# output
OUTPUT := bin/pigeon

# build flags
LDFLAGS := -s -w
LDFLAGS += -extldflags "-static -fpic"

BUILD_FLAGS := -a
BUILD_FLAGS += -trimpath
BUILD_FLAGS += -ldflags '$(LDFLAGS)'
BUILD_FLAGS += $(EXTRA_FLAGS)

# debug flags
GCFLAGS := "all=-N -l"

DEBUG_FLAGS := -gcflags=$(GCFLAGS)

# packages
PACKAGES := $(PWD)/cmd/pigeon/main.go

# test flags
TEST_FLAGS := -v
TEST_FLAGS += -p 3
TEST_FLAGS += -cover
TEST_FLAGS += $(DEBUG_FLAGS)

# test env
TEST_ENV := ROOT=$(PWD)

# curve-go-rpc path
RPC_PATH = external/curve-go-rpc

init:
	$(GO) mod init

build:
	$(GOENV) $(GO) build -o $(OUTPUT) $(BUILD_FLAGS) $(PACKAGES)

debug:
	$(GOENV) $(GO) build -o $(OUTPUT) $(DEBUG_FLAGS) $(PACKAGES)

test:
	$(TEST_ENV) $(GO) test $(TEST_FLAGS) ./...

clean:
	rm -rf bin/* docker/pigeon docker/website

image:
	bash image.sh $(tag)
	