# Main Makefile for dmarc-cat
#
# Copyright 2018 © by Ollivier Robert
#

GO=		go
GOBIN=  ${GOPATH}/bin

BIN=	dmarc-cat

SRCS= analyze.go file.go main.go parse.go types.go utils.go

OPTS=	-ldflags="-s -w" -v

all: ${BIN}

${BIN}: ${SRCS}
	${GO} build -o ${BIN} ${OPTS} .

test:
	${GO} test -v .

lint:
	gometalinter

install: ${BIN}
	${GO} install ${OPTS} .

clean:
	${GO} clean -v

push:
	git push --all
	git push --tags
