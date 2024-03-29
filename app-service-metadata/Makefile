.PHONY: build test clean docker

GO=CGO_ENABLED=1 GO111MODULE=on go

# VERSION file is not needed for local development, In the CI/CD pipeline, a temporary VERSION file is written
# if you need a specific version, just override below
# TODO: If your service is not being upstreamed to Edgex Foundry, you need to determine the best approach for
#       setting your service's version for non-development builds.
APPVERSION=$(shell cat ./VERSION)

# This pulls the version of the SDK from the go.mod file. If the SDK is the only required module,
# it must first remove the word 'required' so the offset of $2 is the same if there are multiple required modules
SDKVERSION=$(shell cat ./go.mod | grep 'github.com/edgexfoundry/app-functions-sdk-go/v2 v' | sed 's/require//g' | awk '{print $$2}')

MICROSERVICE=app-service-metadata
GOFLAGS=-ldflags "-X github.com/edgexfoundry/app-functions-sdk-go/v2/internal.SDKVersion=$(SDKVERSION) -X github.com/edgexfoundry/app-functions-sdk-go/v2/internal.ApplicationVersion=$(APPVERSION)"

# TODO: uncomment and remove default once files are in a Github repository or
#       remove totally including usage below
GIT_SHA=$(shell git rev-parse HEAD)


build:
	go mod tidy
	$(GO) build $(GOFLAGS) -o $(MICROSERVICE)

docker:
	docker build --no-cache \
	    --build-arg http_proxy \
	    --build-arg https_proxy \
		-f Dockerfile \
		--label "git_sha=$(GIT_SHA)" \
		-t tibcosoftware/labs-air-edgex-app-service-metadata:$(APPVERSION) \
		.

dockerarm64:
	docker build --no-cache \
	    --build-arg http_proxy \
	    --build-arg https_proxy \
		-f Dockerfile \
		--label "git_sha=$(GIT_SHA)" \
		-t tibcosoftware/labs-air-edgex-app-service-metadata-arm64:$(APPVERSION) \
		.

# The test-attribution-txt.sh scripts are required for upstreaming to EdgeX Foundry.
# TODO: Remove bin folder and reference to script below if NOT upstreaming to EdgeX Foundry.
test:
	go mod tidy
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) vet ./...
	gofmt -l .
	[ "`gofmt -l .`" = "" ]
	./bin/test-attribution-txt.sh

clean:
	rm -f $(MICROSERVICE)


