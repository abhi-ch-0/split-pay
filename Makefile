PROTO_DIR=proto
GO_OUT_DIR=generated
PROTO_FILES=$(wildcard $(PROTO_DIR)/*.proto)

SQL_DIR=sql
DATABASE=split_pay_db

.PHONY: all proto clean build run

all: proto build

create_database:
	createdb $(DATABASE)
	psql $(DATABASE)

execute_sql:
	@echo "Enter your PostgreSQL username:"
	@read USERNAME; \
	echo "Executing all SQL files in $(SQL_DIR)"; \
	for file in $(SQL_DIR)/*.sql; do \
		echo "Executing $$file"; \
		psql -U $$USERNAME -d $(DATABASE) -f "$$file"; \
	done

proto:
	protoc -I=$(PROTO_DIR) --go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative --go-grpc_out=$(GO_OUT_DIR) --go-grpc_opt=paths=source_relative $(PROTO_FILES)

clean:
	rm -f $(GO_OUT_DIR)/*.pb.go

build:
	go build -o split-pay-app .

run: build
	./split-pay-app