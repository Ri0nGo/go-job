# go-job

go-job 是一个任务执行平台，区分master和node两部分，master负责任务的增删改查等，node负责执行任务

## 目录结构说明

### 整体结构

```shell
.
├── cmd      # 启动目录
├── config   # 配置文件
├── internal # 内部目录
├── master   # 存放master程序
└── node     # 存放node程序
```

### node 目录结构

```shell
node
├── api                  # api接口
├── pkg                
│   ├── config     # 配置文件解析
│   ├── executor   # 执行器
│   └── job        # 任务管理
├── router               # 路由
└── service              # 服务层
```

### master 目录结构

```shell
├─api              # api 接口
├─database         # 数据库，DAO层
├─pkg 
│  ├─config        # 配置文件
│  ├─ioc           # 依赖注入， wire
│  └─middleware    # 中间件
├─repo             # repo层
├─router           # 路由层
└─service          # service层
```

## 容器运行说明

### docker 启动master节点

```shell
make build-master-image
docker run -d --name go-job-node -v $(pwd)/master.yaml:/app/config/master.yaml -p 8080:8080 <BUILD_DOCKER_IMAGE>
```

### docker 启动 python 环境的 node节点

```shell
make build-node-py-image  PIP_FILE=/your/pip.txt
docker run -d --name go-job-node -v $(pwd)/node.yaml:/app/config/node.yaml -p 8081:8081 <BUILD_DOCKER_IMAGE>
```