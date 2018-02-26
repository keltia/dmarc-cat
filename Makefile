# Main Makefile for dmarc-cat
#
# Copyright 2018 © by Ollivier Robert
#

GOBIN=   ${GOPATH}/bin

BIN=	dmarc-cat

SRCS= analyze.go file.go main.go parse.go types.go utils.go

OPTS=	-ldflags="-s -w" -v

all: ${BIN}

${BIN}: ${SRCS}
	go build -o ${BIN} ${OPTS} .

test:
	go test -v .

lint:
	gometalinter

install: ${BIN}
	go install ${OPTS} .

clean:
	go clean -v

push:
	git push --all
	git push --tags
