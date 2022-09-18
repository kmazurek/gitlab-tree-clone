BINARY_NAME := "gltc"

ARCH := $(shell uname -m)
ifeq ($(ARCH), arm64)
 TAGS += dynamic
endif

all: build
install:
	@go build -tags=${TAGS} -o ${BINARY_NAME} cmd/gitlab_tree_clone.go && cp ${BINARY_NAME} ${GOPATH}/bin
build:
	@go build -tags=${TAGS} -o ${BINARY_NAME} cmd/gitlab_tree_clone.go