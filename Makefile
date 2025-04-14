############# DEBUG #############
APP_MASTER_NAME := go-job-master
APP_NODE_NAME := go-job-node
UPLOAD_JOB_DIR := upload_job
NODE_UPLOAD_JOB_DIR := node_upload_job


#############  CICD  #############
VERSION ?= $(shell date +%Y%m%d%H%M%S)
BUILD_DIR := /cicd/projects/go-job

# cicd master
MASTER_IMAGE_NAME := go-job-master
MASTER_IMAGE := $(MASTER_IMAGE_NAME):$(VERSION)

# cicd node
NODE_IMAGE_NAME := go-job-node
NODE_IMAGE := $(NODE_IMAGE_NAME):$(VERSION)


.PHONY: run-n
run-n:
	@echo "running node $(APP_NODE_NAME)..."
	go run cmd/node/node.go

.PHONY:
build-node:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/node/go-job-node cmd/master/master.go

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

build-master-image:
	@echo ">>> 开始构建 go-job master 镜像: $(MASTER_IMAGE)"
	@rm -rf $(BUILD_DIR)
	@git clone git@github.com:Ri0nGo/go-job.git $(BUILD_DIR)
	@cd $(BUILD_DIR) && \
		go env -w GOPROXY=https://goproxy.cn,direct && \
		go env -w GO111MODULE=on && \
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/master/go-job-master cmd/master/master.go && \
		cp master/Dockerfile ./build/master && \
		cd ./build/master && \
		docker build -t $(MASTER_IMAGE) . -f Dockerfile
	@echo ">>> 镜像构建完成: $(MASTER_IMAGE)"

.PHONY:



# 必须传递pip文件
ifndef PIP_FILE
$(error PIP_FILE is not set. Please run with: make docker-build-node PIP_FILE=/your/pip.txt)
endif

build-node-py-image:
	@echo ">>> 开始构建 go-job node 镜像: $(NODE_IMAGE)"
	@rm -rf $(BUILD_DIR)
	@git clone git@github.com:Ri0nGo/go-job.git $(BUILD_DIR)
	@cd $(BUILD_DIR) && \
		go env -w GOPROXY=https://goproxy.cn,direct && \
		go env -w GO111MODULE=on && \
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/node/go-job-node cmd/node/node.go && \
		cp node/Dockerfile ./build/node && \
		cp $(PIP_FILE) ./build/node && \
		cd ./build/node && \
		docker build -t $(NODE_IMAGE) . -f Dockerfile
	@echo ">>> 镜像构建完成: $(NODE_IMAGE)"