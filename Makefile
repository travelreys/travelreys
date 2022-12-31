.PHONY: build protoc clean test

build:
	@echo "Building..."
	go build cmd/server -o build/server

protoc:
	@echo "Generating Go files"
	cd proto && protoc --go_out=. --go-grpc_out=. \
		--go-grpc_opt=paths=source_relative --go_opt=paths=source_relative *.proto
	mv proto/trips.pb.go pkg/trips/collab.pb.go
	mv proto/trips_grpc.pb.go pkg/trips/collab_grpc.pb.go

test:
	go test ./...

clean:
	rm -rf build
