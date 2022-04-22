.PHONY: build test clean prepare update docker

GO = CGO_ENABLED=1 GO111MODULE=on go

MICROSERVICES=cmd/device-video-analytics-mqtt
.PHONY: $(MICROSERVICES)

DOCKERS=docker_device_generic_mqtt
.PHONY: $(DOCKERS)

VERSION=$(shell cat ./VERSION)
GIT_SHA=$(shell git rev-parse HEAD)

GOFLAGS=-ldflags "-X github.com/TIBCOSoftware/labs-air/edgexfoundry/device-video-analytics-mqtt.Version=$(VERSION)"

build: $(MICROSERVICES)

cmd/device-video-analytics-mqtt:
	go mod tidy
	$(GO) build $(GOFLAGS) -o $@ ./cmd

test:
	$(GO) test ./... -cover

clean:
	rm -f $(MICROSERVICES)

run:
	cd bin && ./edgex-launch.sh

docker: $(DOCKERS)

docker_device_generic_mqtt:
	docker build --no-cache \
		--label "git_sha=$(GIT_SHA)" \
		-t tibcosoftware/labs-air-edgex-device-video-analytics-mqtt:$(VERSION) \
		.

dockerarm64:
	docker build --no-cache \
		--label "git_sha=$(GIT_SHA)" \
		-t tibcosoftware/labs-air-edgex-device-video-analytics-mqtt-arm64:$(VERSION) \
		.

dockerbuildxarm64:
	docker buildx build \
		--platform linux/arm64 \
		--label "git_sha=$(GIT_SHA)" \
		-t magallardo/docker-device-video-analytics-mqtt-arm64:$(VERSION) \
		.

