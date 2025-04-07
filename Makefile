APP_MASTER_NAME := go-job-master
APP_NODE_NAME := go-job-node

.PHONY: run-node
run-node:
	@echo "running node $(APP_NODE_NAME)..."
	go run cmd/node/node.go

.PHONY: wire-gen
wire-gen:
	wire gen master/pkg/ioc/wire.go

.PHONY: run-master
run-master: wire-gen
	@echo "running master $(APP_MASTER_NAME)..."
	go run cmd/master/master.go