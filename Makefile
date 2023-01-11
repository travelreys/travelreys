.PHONY: build protoc protoc clean test

build:
	@echo "Building..."
	go build -o build/server cmd/server/*.go

test:
	go test ./...

clean:
	rm -rf build
