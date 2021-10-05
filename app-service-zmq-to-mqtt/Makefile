.PHONY: build test clean docker

GO=CGO_ENABLED=1 GO111MODULE=on go

MICROSERVICES=cmd/app-service-zmq-to-mqtt
.PHONY: $(MICROSERVICES)

VERSION=$(shell cat ./VERSION)
GOFLAGS=-ldflags "-X github.com/TIBCOSoftware/labs-air-infra/edgexfoundry/app-service-zmq-to-mqtt.Version=$(VERSION)"
GIT_SHA=$(shell git rev-parse HEAD)

build: $(MICROSERVICES)
	$(GO) install -tags=safe

$(MICROSERVICES):
	$(GO) build $(GOFLAGS) -o $@ ./cmd

docker:
	docker build --no-cache \
		--label "git_sha=$(GIT_SHA)" \
		-t tibcosoftware/labs-air-edgex-app-service-zmq-to-mqtt:$(VERSION) \
		.

dockerarm64:
	docker build --no-cache \
		--label "git_sha=$(GIT_SHA)" \
		-t tibcosoftware/labs-air-edgex-app-service-zmq-to-mqtt-arm64:$(VERSION) \
		.


dockerbuildxarm64:
	docker buildx build \
		--platform linux/arm64 \
		--label "git_sha=$(GIT_SHA)" \
		-t tibcosoftware/docker-app-service-zmq-to-mqtt-arm64:$(VERSION) --push \
		.

test:
	$(GO) vet ./...
	gofmt -l .
	$(GO) test -coverprofile=coverage.out ./...

clean:
	rm -f $(MICROSERVICES)
