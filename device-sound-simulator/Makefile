.PHONY: build test clean prepare update docker

GO = CGO_ENABLED=1 GO111MODULE=on go

MICROSERVICES=cmd/device-sound-simulator
.PHONY: $(MICROSERVICES)

DOCKERS=docker_device_sound_simulator
.PHONY: $(DOCKERS)

VERSION=$(shell cat ./VERSION)
GIT_SHA=$(shell git rev-parse HEAD)

GOFLAGS=-ldflags "-X github.com/TIBCOSoftware/labs-air-infra/edgexfoundry/device-sound-simulator.Version=$(VERSION)"

build: $(MICROSERVICES)

cmd/device-sound-simulator:
	go mod tidy
	$(GO) build $(GOFLAGS) -o $@ ./cmd

test:
	$(GO) test ./... -cover

clean:
	rm -f $(MICROSERVICES)

docker: $(DOCKERS)

docker_device_sound_simulator:
	docker build --no-cache \
		--build-arg http_proxy \
        --build-arg https_proxy \
		--label "git_sha=$(GIT_SHA)" \
		-t tibcosoftware/labs-air-edgex-device-sound-simulator:$(VERSION) \
		.

dockerarm64:
	docker build --no-cache \
		--build-arg http_proxy \
        --build-arg https_proxy \
		--label "git_sha=$(GIT_SHA)" \
		-t tibcosoftware/labs-air-edgex-device-sound-simulator-arm64:$(VERSION) \
		.

dockerbuildxarm64:
	docker buildx build \
		--build-arg http_proxy \
        --build-arg https_proxy \
		--platform linux/arm64 \
		--label "git_sha=$(GIT_SHA)" \
		-t magallardo/docker-device-sound-simulator-arm64:$(VERSION) \
		.

