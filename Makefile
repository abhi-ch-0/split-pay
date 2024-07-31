# Directories
PROTO_DIR=protobuf
GO_OUT_DIR=generated

# Proto files
PROTO_FILES=$(wildcard $(PROTO_DIR)/*.proto)

# Targets
.PHONY: all proto clean build run

# Default target: generate proto files and build the project
all: proto build

# Generate Go files from proto files
proto:
	protoc -I=$(PROTO_DIR) --go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative --go-grpc_out=$(GO_OUT_DIR) --go-grpc_opt=paths=source_relative $(PROTO_FILES)

# Clean generated files
clean:
	rm -f $(GO_OUT_DIR)/*.pb.go

# Build the Go project
build:
	go build -o split-pay-app .

# Run the Go project
run: build
	./split-pay-app