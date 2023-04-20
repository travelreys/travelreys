DOCKER_REGISTRY := asia-southeast1-docker.pkg.dev/tiinyplanet-379603
DOCKER_REPOSITORY := travelreys-api
DOCKER_IMAGE := travelreys-api
VERSION := $(shell grep 'VERSION' pkg/common/version.go | awk '{ print $$4 }' | tr -d '"')

.PHONY: build protoc protoc clean test

build:
	@echo "Building..."
	go build -o build/server cmd/server/*.go
	go build -o build/coordinator cmd/coordinator/*.go

test:
	go test ./...

clean:
	rm -rf build

docker:
	docker build -t travelreys-api .
	docker tag travelreys-api $(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY)/$(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY)/$(DOCKER_IMAGE):$(VERSION)
