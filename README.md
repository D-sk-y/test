# test

基于 Go 的模块化后端服务框架，内置消息总线、双风格 HTTP API 和 MySQL 持久化支持。

## 项目结构

```
test/
├── cmd/
│   └── main.go              # 入口：初始化配置、数据库、模块、API 网关
├── pkg/
│   ├── api/
│   │   ├── gateway.go       # HTTP 网关，管理 v1/v2 路由，监听 :18080
│   │   ├── v1/
│   │   │   └── handler.go   # RESTful 风格 API（GET/POST/PUT/DELETE）
│   │   └── v2/
│   │       └── handler.go   # RPC-over-HTTP 风格 API（全 POST）
│   ├── config/
│   │   └── config.go        # 读取 build/conf/config.yaml 配置
│   └── mod/
│       ├── module.go        # Module 接口 + ModuleBase 管理器
│       ├── bus.go           # 单例消息总线，模块间解耦通信
│       ├── message.go       # 消息结构体
│       ├── store.go         # 内存任务存储（支持后续切换数据库）
│       ├── system.go        # System 模块：生命周期管理
│       └── worker.go        # Worker 模块：业务处理
├── db/
│   └── sql.go               # sqlx.DB 包装类
├── build/
│   ├── conf/
│   │   └── config.yaml      # 配置文件
│   └── test.exe             # 编译产物（make build 生成）
├── Makefile
├── go.mod
└── go.sum
```

## 快速开始

### 环境要求

- Go >= 1.25
- MySQL（如需数据库功能）

### 配置

编辑 `build/conf/config.yaml`：

```yaml
db:
  driver: mysql
  host: 127.0.0.1
  port: 3306
  user: root
  password: "your_password"
  name: test
```

### 构建与运行

```bash
# 安装依赖
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
```

## 架构设计

### 模块系统

所有业务模块实现 `Module` 接口，由 `ModuleBase` 统一管理生命周期：

```go
type Module interface {
    Initialize(ctx context.Context) error
    DeInitialize() error
}
```

当前内置模块：

| 模块 | 职责 |
|------|------|
| **System** | 系统级控制，持有全局 Context 和 CancelFunc |
| **Worker** | 业务逻辑处理，通过 Bus 接收消息 |
| **Gateway** | HTTP API 网关，管理路由和请求分发 |

### 消息总线

`Bus` 是单例（`GetBus()`），提供发布-订阅模式，模块间通过消息名称通信：

```go
bus := mod.GetBus()
bus.CreateMsg("task:created", task)
```

### API 设计

项目同时提供两种 API 风格用于对比学习，共享同一个 `Store` 实例：

#### v1 — RESTful

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/tasks` | 列出所有任务 |
| POST | `/api/v1/tasks` | 创建任务 |
| GET | `/api/v1/tasks/{id}` | 获取单个任务 |
| PUT | `/api/v1/tasks/{id}` | 更新任务 |
| DELETE | `/api/v1/tasks/{id}` | 删除任务 |

#### v2 — RPC-over-HTTP

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v2/listTasks` | 列出所有任务 |
| POST | `/api/v2/createTask` | 创建任务 |
| POST | `/api/v2/getTask` | 获取单个任务 |
| POST | `/api/v2/updateTask` | 更新任务 |
| POST | `/api/v2/deleteTask` | 删除任务 |

### 请求示例

**v1 创建任务：**

```bash
curl -X POST http://localhost:18080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"name": "学习 Go"}'
```

**v2 创建任务：**

```bash
curl -X POST http://localhost:18080/api/v2/createTask \
  -H "Content-Type: application/json" \
  -d '{"name": "学习 Go"}'
```

响应格式统一：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": "task-001",
    "name": "学习 Go",
    "status": "pending",
    "created_at": "2026-08-11T10:00:00Z",
    "updated_at": "2026-08-11T10:00:00Z"
  }
}
```

## 待扩展

- [ ] Store 切换为 MySQL 持久化（`NewStoreWithDB`）
- [ ] Bus 增加订阅接口（当前仅支持发布）
- [ ] 配置支持环境变量覆盖
- [ ] 优雅关闭时关闭数据库连接

## 技术栈

| 组件 | 库 |
|------|-----|
| HTTP 路由 | `net/http`（Go 1.22+ 增强路由） |
| 数据库 | `jmoiron/sqlx` + `go-sql-driver/mysql` |
| 配置解析 | `go.yaml.in/yaml/v3` |
| 构建工具 | GNU Make |
