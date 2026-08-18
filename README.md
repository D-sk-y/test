# test

基于 Go 的模块化后端服务框架，内置消息总线、双风格 HTTP API 和 MySQL 持久化支持。

## 项目简介

`test` 是一个用于学习与演示的 Go 后端服务项目，采用模块化架构设计：

- 通过 `Module` 接口统一管理各业务模块的生命周期（初始化 / 反初始化）
- 内置单例消息总线（`Bus`），实现模块间基于发布-订阅模式的解耦通信
- 同时提供 RESTful（v1）与 RPC-over-HTTP（v2）两种 API 风格，便于对比学习
- 数据层基于 `sqlx` 直接操作 MySQL，任务数据持久化存储

## 功能特性

- 模块化生命周期管理：`ModuleBase` 统一注册、初始化、反初始化所有模块
- 消息总线：单例 `Bus`，支持发布消息与订阅消息（`Subscribe`）
- 双风格 HTTP API：
  - v1 RESTful：`GET / POST /api/v1/tasks`、`PUT /api/v1/tasks/{id}`
  - v2 RPC-over-HTTP：`POST /api/v2/hello`（其余路由暂被注释）
- MySQL 持久化：任务（Task）的增、查、改、删直接操作数据库
- 优雅退出：监听 `SIGINT` / `SIGTERM` 信号，取消上下文并反初始化所有模块
- 定时状态输出：运行期间每 5 秒打印一次应用状态

## 技术栈

| 组件 | 说明 |
|------|------|
| 语言 | Go 1.25.1 |
| HTTP 路由 | `net/http`（Go 1.22+ 增强路由，方法 + 路径模式） |
| 数据库 | `jmoiron/sqlx` + `go-sql-driver/mysql` |
| 配置解析 | `go.yaml.in/yaml/v3` |
| 构建工具 | GNU Make |

## 目录结构

```
test/
├── cmd/
│   └── main.go              # 程序入口：加载配置、连接数据库、初始化模块、等待退出信号
├── pkg/
│   ├── api/
│   │   ├── gateway.go       # HTTP 网关模块，注册 v1/v2 路由，监听 :18080
│   │   ├── v1/
│   │   │   └── handler.go   # RESTful 风格 API（GET/POST/PUT）
│   │   └── v2/
│   │       └── handler.go   # RPC-over-HTTP 风格 API（POST /api/v2/hello）
│   ├── config/
│   │   └── config.go        # 读取可执行文件同目录 conf/config.yaml 配置
│   └── mod/
│       ├── module.go        # Module 接口 + ModuleBase 模块管理器
│       ├── bus.go           # 单例消息总线，支持发布与订阅
│       ├── message.go       # 消息结构体
│       ├── system.go        # System 模块：持有全局 Context 与 CancelFunc
│       └── worker.go        # Worker 模块：订阅消息并消费
├── entity/
│   └── task.go              # Task 实体定义（db / json 标签）
├── db/
│   └── sql.go               # DB 封装：基于 sqlx 的任务增删查改
├── build/
│   ├── conf/
│   │   └── config.yaml      # 数据库等配置文件
│   └── test.exe             # 编译产物（make build 生成）
├── Makefile                 # 构建、运行、测试等命令
├── go.mod / go.sum          # Go 模块依赖
├── TODO.md                  # 待办事项
└── test.go.bak              # 历史备份文件（已注释）
```

## 快速开始

### 环境要求

- Go >= 1.25
- MySQL（任务数据持久化依赖）
- GNU Make（可选，也可直接使用 `go` 命令）

### 配置

编辑 `build/conf/config.yaml`，配置数据库连接信息：

```yaml
db:
  driver: mysql
  host: 127.0.0.1
  port: 3306
  user: root
  password: 123456
  name: go
```

> 注意：程序从可执行文件同目录下的 `conf/config.yaml` 读取配置，即 `build/conf/config.yaml`。

### 数据库准备

程序启动时会连接 MySQL 并操作 `tasks` 表，需提前创建数据库与表：

```sql
CREATE DATABASE IF NOT EXISTS go;
USE go;
CREATE TABLE IF NOT EXISTS tasks (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  status VARCHAR(64) NOT NULL DEFAULT 'pending',
  created_at VARCHAR(64),
  updated_at VARCHAR(64)
);
```

### 构建与运行

```bash
# 整理依赖
make tidy

# 仅构建
make build

# 构建并运行
make run

# 格式化代码
make fmt

# 静态检查
make vet

# 运行测试
make test

# 清理构建产物
make clean
```

## 架构设计

### 模块系统

所有业务模块实现 `Module` 接口，由 `ModuleBase` 统一管理生命周期：

```go
type Module interface {
    Initialize(c context.Context) error
    DeInitialize() error
}
```

当前内置模块：

| 模块 | 职责 |
|------|------|
| **System** | 系统级控制，持有全局 Context 和 CancelFunc |
| **Worker** | 业务处理，订阅 Bus 消息并消费 |
| **ApiGateway** | HTTP API 网关，管理路由和请求分发 |

### 消息总线

`Bus` 是单例（`GetBus()`），提供发布-订阅模式，模块间通过消息名称通信：

```go
bus := mod.GetBus()
bus.CreateMsg("Worker", "hello")   // 发布消息
bus.Subscribe("Worker", ch)        // 订阅消息
```

### API 设计

项目同时提供两种 API 风格用于对比学习，共享同一个 `db.DB` 实例：

#### v1 — RESTful

| 方法 | 路径 | 说明 | 状态 |
|------|------|------|------|
| GET | `/api/v1/tasks` | 列出所有任务 | 已启用 |
| POST | `/api/v1/tasks` | 创建任务 | 已启用 |
| GET | `/api/v1/tasks/{id}` | 获取单个任务 | 已注释 |
| PUT | `/api/v1/tasks/{id}` | 更新任务 | 已启用 |
| DELETE | `/api/v1/tasks/{id}` | 删除任务 | 已注释 |

#### v2 — RPC-over-HTTP

| 方法 | 路径 | 说明 | 状态 |
|------|------|------|------|
| POST | `/api/v2/hello` | 发送消息到 Worker | 已启用 |
| POST | `/api/v2/listTasks` | 列出所有任务 | 已注释 |
| POST | `/api/v2/createTask` | 创建任务 | 已注释 |
| POST | `/api/v2/getTask` | 获取单个任务 | 已注释 |
| POST | `/api/v2/updateTask` | 更新任务 | 已注释 |
| POST | `/api/v2/deleteTask` | 删除任务 | 已注释 |

### 请求示例

**v1 创建任务：**

```bash
curl -X POST http://localhost:18080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"name": "学习 Go", "status": "pending"}'
```

**v1 列出任务：**

```bash
curl http://localhost:18080/api/v1/tasks
```

**v1 更新任务：**

```bash
curl -X PUT http://localhost:18080/api/v1/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "学习 Go", "status": "done"}'
```

**v2 发送消息：**

```bash
curl -X POST http://localhost:18080/api/v2/hello \
  -H "Content-Type: application/json" \
  -d '{"msg": "hello worker"}'
```

v1 接口直接返回 JSON 数据（任务对象或数组），无统一包装结构。

## 待扩展

- [ ] 启用 v1 的 `GET /api/v1/tasks/{id}` 与 `DELETE /api/v1/tasks/{id}` 路由
- [ ] 启用 v2 的 `listTasks` / `createTask` / `getTask` / `updateTask` / `deleteTask` 路由
- [ ] 配置支持环境变量覆盖
- [ ] 优雅关闭时关闭数据库连接
