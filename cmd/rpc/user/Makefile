.PHONY: all build run gotool install clean help

NAME=open_im_user
BIN_DIR=../../../bin/

OS:= $(or $(os),linux)
ARCH:=$(or $(arch),amd64)
all: gotool build

ifeq ($(OS),windows)

BINARY_NAME=${NAME}.exe

else

BINARY_NAME=${NAME}

endif

build:
	CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH}; go build -ldflags="-w -s" -o ${BINARY_NAME}

run:
	@go run ./

gotool:
	go fmt ./
	go vet ./

install:build
	mv ${BINARY_NAME} ${BIN_DIR}

clean:
	@if [ -f ${BINARY_NAME} ] ; then rm ${BINARY_NAME} ; fi





