PROTO_DIR=protobuf
GO_OUT_DIR=generated

PROTO_FILES=$(wildcard $(PROTO_DIR)/*.proto)

.PHONY: all proto clean build run

all: proto build

proto:
	protoc -I=$(PROTO_DIR) --go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative --go-grpc_out=$(GO_OUT_DIR) --go-grpc_opt=paths=source_relative $(PROTO_FILES)

clean:
	rm -f $(GO_OUT_DIR)/*.pb.go

build:
	go build -o split-pay-app .

run: build
	./split-pay-app