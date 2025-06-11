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

# server
SERVER_IP=xx.xx.xx.xx

# git
COMMIT_FILE=last_commit.txt
CUR_DIR := $(CURDIR)


.PHONY: run-n
run-n:
	@echo "running node $(APP_NODE_NAME)..."
	go run cmd/node/node.go

.PHONY: build-node
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
	go run cmd/master/master.go  -c ./config/master-prod.yaml

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
			cp master/Dockerfile . && \
			docker build -t $(MASTER_IMAGE) . -f Dockerfile
	@echo ">>> 镜像构建完成: $(MASTER_IMAGE)"

.PHONY:

build-node-py-image:
	@if [ -z "$(PIP_FILE)" ]; then \
		echo "Error: PIP_FILE is not set. Please run with: make build-node-py-image PIP_FILE=/your/pip.txt"; \
		exit 1; \
	fi
	@echo ">>> 开始构建 go-job node 镜像: $(NODE_IMAGE)"
	@rm -rf $(BUILD_DIR)
	@git clone git@github.com:Ri0nGo/go-job.git $(BUILD_DIR)
	@cd $(BUILD_DIR) && \
		if [ -f $(CUR_DIR)/$(COMMIT_FILE) ]; then \
					echo ">>> 本次构建的 Git 提交记录如下：" && \
			git log --date=format:'%Y-%m-%d %H:%M:%S' \
					--pretty=format:"%h|%an|%ad|%s" `cat $(CUR_DIR)/$(COMMIT_FILE)`..HEAD | \
			while IFS="|" read -r hash author date msg; do \
					printf "%-10s %-15s %-20s %s\n" "$$hash" "$$author" "$$date" "$$msg"; \
			done; \
		else \
			echo ">>> 未找到提交记录文件 $(COMMIT_FILE)，显示最近 3 条提交：" && \
			git rev-parse HEAD > $(CUR_DIR)/$(COMMIT_FILE); \
			git log -n 4 --date=format:'%Y-%m-%d %H:%M:%S' \
					--pretty=format:"%h|%an|%ad|%s" | \
			while IFS="|" read -r hash author date msg; do \
					printf "%-10s %-15s %-20s %s\n" "$$hash" "$$author" "$$date" "$$msg"; \
			done; \
		fi && \
		cp node/Dockerfile . && \
		docker build -t $(NODE_IMAGE) . -f Dockerfile
	@echo ">>> 镜像构建完成: $(NODE_IMAGE)"

	# 此处的镜像仓库可以自己修改
	@echo ">>> 开始推送镜像到阿里云镜像仓库: $(REMOTE_MASTER_IMAGE)"
	@docker tag $(MASTER_IMAGE) $(REMOTE_MASTER_IMAGE)
	@docker push $(REMOTE_MASTER_IMAGE)
	@echo ">>> 镜像推送完成: $(REMOTE_MASTER_IMAGE)"

	# 添加秘钥自动部署 （可选）#
	@echo ">>> 开始更新阿里云容器镜像"
	@ssh root@$(SERVER_IP) 'set -e; cd /root/shell/images-shell && make dpi-go-job-master version=$(VERSION) GJC=/root/shell/init/go-job/compose.yaml'
	@echo ">>> 阿里云容器镜像更新完成"


.PHONY:
set-env:
	export workspace="D:\MyProjects\go-job"


#############################################
#		仅供参考			                    #
#		此处为服务器拉取镜像仓库脚本             #
#############################################

GO_JOB_ADMIN_IMAGE_NAME := xxx
GO_JOB_MASTER_IMAGE_NAME := xxx
GO_JOB_NODE_PY_IMAGE_NAME := xxx

REGISTRY := xxx


.PHONY: help dpi-go-job-node dpi-go-job-master

### 显示所有可用命令
help:
	@echo "可用命令如下"
	@grep -E '^[a-zA-Z0-9_-]+:' Makefile | awk -F: '{ \
			target=$$1; \
			printf "  %-25s\t拉取镜像; eg. make %s version=img_version\n", target, target; \
	}'


### 拉取 go-job node 镜像; eg. make dpi-go-job-master version=img_version
dpi-go-job-node:
	@if [ -z "$(version)" ]; then \
			echo "Error: tag is not set. Please run with: make build-node-py-image version=img_version"; \
			exit 1; \
	fi
	@IMAGE=$(REGISTRY)/go-job-node:$(version); \
	echo ">>> 拉取镜像: $$IMAGE"; \
	docker pull $$IMAGE

	@echo ">>> 开始更新 go-job-node tag: $(version)"
	@sed -i 's#\(.*rion/go-job-node:\)[^"]*#\1$(version)#g' $$GJC
	@docker compose  -f $$GJC  up -d --force-recreate $(GO_JOB_NODE_PY_IMAGE_NAME)
	@echo ">>> 更新镜像完, 服务名: $(GO_JOB_NODE_PY_IMAGE_NAME)"

### 拉取go-job master 镜像; eg. make dpi-go-job-master version=img_version
dpi-go-job-master:
	@if [ -z "$(version)" ]; then \
			echo "Error: version is not set. Please run with: make dpi-go-job-master version=img_version"; \
			exit 1; \
	fi
	@IMAGE=$(REGISTRY)/go-job-master:$(version); \
	echo ">>> 拉取镜像: $$IMAGE"; \
	docker pull $$IMAGE

	@echo ">>> 开始更新 go-job-master tag: $(version)"
	@sed -i 's#\(.*rion/go-job-master:\)[^"]*#\1$(version)#g' $$GJC
	@docker compose  -f $$GJC  up -d --force-recreate $(GO_JOB_MASTER_IMAGE_NAME)
	@echo ">>> 更新镜像完, 服务名: $(GO_JOB_MASTER_IMAGE_NAME)"


### 拉取go-job admin 镜像; eg. make dpi-go-job-admin version=img_version
dpi-go-job-admin:
	@if [ -z "$(version)" ]; then \
			echo "Error: version is not set. Please run with: make dpi-go-job-admin version=img_version"; \
			exit 1; \
	fi
	@IMAGE=$(REGISTRY)/go-job-admin:$(version); \
	echo ">>> 拉取镜像: $$IMAGE"; \
	docker pull $$IMAGE

	@echo ">>> 开始更新 go-job-admin tag: $(version)"
	@sed -i 's#\(.*rion/go-job-admin:\)[^"]*#\1$(version)#g' $$GJC
	@docker compose  -f $$GJC  up -d --force-recreate $(GO_JOB_ADMIN_IMAGE_NAME)
	@echo ">>> 更新镜像完, 服务名: $(GO_JOB_ADMIN_IMAGE_NAME)"


