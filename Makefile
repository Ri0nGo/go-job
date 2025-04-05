APP_NODE_NAME := go-job-node

.PHONY: run-node
run-node:
	@echo "running node $(APP_NODE_NAME)..."
	go run cmd/node/node.go