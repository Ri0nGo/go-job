# go-job

go-job 是一个任务执行平台，包含master和node两部分，master负责用户端的数据展示和操作，node负责执行任务和数据回传


## Feature
- [x] 支持用户管理
- [x] 支持任务增删改查
- [x] 仅支持运行python脚本
- [x] 支持秒级定时任务
- [x] 支持任务记录查询
- [x] 支持节点增删改查
- [x] 支持节点健康检测
- [x] 支持节点依赖包安装和查询
- [x] 新增首页数据看板
- [x] 支持绑定邮箱和修改用户密码
- [x] 支持多用户
- [x] 支持配置QQ和163的SMTP服务器，并实现高可用



## 容器运行说明

### 1. docker 启动 master节点

```shell
make build-master-image
docker run -d --name go-job-node \
  -v $(pwd)/master.yaml:/app/config/master.yaml \
  -v /data/go-job/data/:/app/data \
  -p 8080:8080 \
  <BUILD_DOCKER_IMAGE>
```

### docker 启动 python 环境的 node节点

```shell
make build-node-py-image  PIP_FILE=/your/pip.txt
docker run -d --name go-job-node \
  -v $(pwd)/node.yaml:/app/config/node.yaml \
  -v /data/go-job-node/data:/app/data \
  -p 8080:8080 \
  <BUILD_DOCKER_IMAGE>
```

### 2. 使用 docker-compose 启动节点

`docker compose up -d`

## 前端地址

👉 [Ri0nGo/go-job-admin](https://github.com/Ri0nGo/go-job-admin)
