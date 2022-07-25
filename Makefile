.PHONY: all build run dev install clean docker-build docker-push

VERSION="2.0.5"
BIN_FILE = "min-gateway"

all: install build

install:
	@go mod tidy

run: build
	@./bin/${BIN_FILE}

dev:
	@air

build:
	@CGO_ENABLED=0 GOOS=linux go build -a -o ./bin/${BIN_FILE}

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

docker-build: build
	@docker build -t docker.in-mes.com/min/min-gateway:${VERSION} -f ./docker/Dockerfile .

docker-push:
	@docker push docker.in-mes.com/min/min-gateway:${VERSION}