APP_MASTER_NAME := go-job-master
APP_NODE_NAME := go-job-node
UPLOAD_JOB_DIR := upload_job

.PHONY: run-n
run-node:
	@echo "running node $(APP_NODE_NAME)..."
	go run cmd/node/node.go

.PHONY: wire
wire:
	wire gen master/pkg/ioc/wire.go

.PHONY: run-wire-m
run-wire-m: wire-gen
	@echo "wire than running $(APP_MASTER_NAME)..."
	go run cmd/master/master.go

.PHONY: run-m
run-m:
	@echo "run master"
	go run cmd/master/master.go

.PHONY: clean

ifeq ($(OS),Windows_NT)

clean:
	@echo "clean windows dir: $(UPLOAD_JOB_DIR)"
	@if exist data\$(UPLOAD_JOB_DIR) (rmdir /s /q data\$(UPLOAD_JOB_DIR))

else

clean:
	@echo "clean linux dir: $(UPLOAD_JOB_DIR)"
	@rm -rf data/$(UPLOAD_JOB_DIR)

endif