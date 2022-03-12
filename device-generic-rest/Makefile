.PHONY: build test clean prepare update docker

GO = CGO_ENABLED=1 GO111MODULE=on go

MICROSERVICES=cmd/device-generic-rest
.PHONY: $(MICROSERVICES)

DOCKERS=docker_device_generic_rest
.PHONY: $(DOCKERS)

VERSION=$(shell cat ./VERSION)
GIT_SHA=$(shell git rev-parse HEAD)

GOFLAGS=-ldflags "-X github.com/TIBCOSoftware/labs-air/edgexfoundry/device-generic-rest.Version=$(VERSION)"

build: $(MICROSERVICES)

cmd/device-generic-rest:
	go mod tidy
	$(GO) build $(GOFLAGS) -o $@ ./cmd

test:
	$(GO) test ./... -coverprofile=coverage.out

clean:
	rm -f $(MICROSERVICES)

docker: $(DOCKERS)

docker_device_generic_rest:
	docker build --no-cache \
        --build-arg http_proxy \
        --build-arg https_proxy \
		--label "git_sha=$(GIT_SHA)" \
		-t tibcosoftware/labs-air-edgex-device-generic-rest:$(VERSION) \
		.

dockerarm64:
	docker build --no-cache \
		--build-arg http_proxy \
        --build-arg https_proxy \
		--label "git_sha=$(GIT_SHA)" \
		-t tibcosoftware/labs-air-edgex-device-generic-rest-arm64:$(VERSION) \
		.