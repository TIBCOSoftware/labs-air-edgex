.PHONY: build test clean docker run

GO=CGO_ENABLED=0 GO111MODULE=on go

MICROSERVICES=cmd/device-opcua-go
.PHONY: $(MICROSERVICES)

DOCKERS=docker_device_opcua_go
.PHONY: $(DOCKERS)

VERSION=$(shell cat ./VERSION)
GIT_SHA=$(shell git rev-parse HEAD)

GOFLAGS=-ldflags "-X github.com/TIBCOSoftware/labs-air-infra/edgexfoundry/device-opcua-go.Version=$(VERSION)"

build: $(MICROSERVICES)
	$(GO) install -tags=safe

cmd/device-opcua-go:
	$(GO) build $(GOFLAGS) -o $@ ./cmd

test:
	go test ./... -cover

clean:
	rm -f $(MICROSERVICES)

run:
	cd bin && ./edgex-launch.sh

docker: $(DOCKERS)

docker_device_opcua_go:
	docker build --no-cache \
		--label "git_sha=$(GIT_SHA)" \
		-t tibcosoftware/labs-air-edgex-device-opcua-go:$(VERSION) \
		.

dockerarm64:
	docker build --no-cache \
		--label "git_sha=$(GIT_SHA)" \
		-t tibcosoftware/labs-air-edgex-device-opcua-go-arm64:$(VERSION) \
		.

