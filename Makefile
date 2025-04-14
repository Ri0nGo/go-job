APP_MASTER_NAME := go-job-master
APP_NODE_NAME := go-job-node
UPLOAD_JOB_DIR := upload_job
NODE_UPLOAD_JOB_DIR := node_upload_job

.PHONY: run-n
run-n:
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
	@if exist data\$(NODE_UPLOAD_JOB_DIR) (rmdir /s /q data\$(NODE_UPLOAD_JOB_DIR))

else

clean:
	@echo "clean linux dir: $(UPLOAD_JOB_DIR)"
	@rm -rf data/$(UPLOAD_JOB_DIR)
	@rm -rf data/$(NODE_UPLOAD_JOB_DIR)

endif


.PHONY:
docker-build:
	@nowtime=$$(date +%Y%m%d%H%M%S) && \
	rm -rf /cicd/projects/go-job && \
	git clone git@github.com:Ri0nGo/go-job.git /cicd/projects/go-job && \
	cd /cicd/projects/go-job && \
	go env -w GOPROXY=https://goproxy.cn,direct && \
	go env -w GO111MODULE=on && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/node/go-job-node cmd/master/master.go && \
	cp node/Dockerfile ./build/node && \
	cd ./build/node && \
	docker build -t go-job-node:$$nowtime . -f Dockerfile