# go-job

go-job 是一个任务执行平台，区分master和node两部分，master负责任务的增删改查等，node负责执行任务

## Feature
- [x] 用户管理
- [x] 任务增删改查
- [x] 任务记录增删改查
- [x] 节点增删改查



## 容器运行说明

### docker 启动 master节点

```shell
make build-master-image
docker run -d --name go-job-node -v $(pwd)/master.yaml:/app/config/master.yaml -p 8080:8080 <BUILD_DOCKER_IMAGE>
```

### docker 启动 python 环境的 node节点

```shell
make build-node-py-image  PIP_FILE=/your/pip.txt
docker run -d --name go-job-node -v $(pwd)/node.yaml:/app/config/node.yaml -p 8081:8081 <BUILD_DOCKER_IMAGE>
```
