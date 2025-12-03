# Kratos Project Template

这是一个基于 [Kratos](https://github.com/go-kratos/kratos) 框架的项目模板，提供了完整的项目结构和基础设施代码，方便快速启动新项目。

## 项目结构

```
.
├── api/                    # API 定义（protobuf）
│   └── demo/v1/           # Demo API 定义
├── cmd/                    # 应用入口
│   └── app/               # 主程序
├── configs/               # 配置文件
├── internal/              # 内部代码
│   ├── conf/             # 配置定义
│   ├── global/           # 全局变量
│   ├── server/           # 服务器初始化
│   └── service/          # 业务服务
├── provider/             # 基础设施提供者
│   ├── cache/            # Redis 缓存
│   ├── db/               # 数据库（MySQL/PostgreSQL/SQLite）
│   ├── logger/           # 日志
│   └── storage/          # 对象存储（MinIO/S3/OSS/COS）
└── third_party/          # 第三方 proto 文件
```

## 功能特性

- ✅ **HTTP/gRPC 双协议支持**：同时支持 HTTP RESTful API 和 gRPC
- ✅ **多数据库支持**：MySQL、PostgreSQL、SQLite
- ✅ **Redis 缓存**：集成 Redis 客户端
- ✅ **对象存储**：支持 MinIO、S3、OSS、COS（可扩展）
- ✅ **结构化日志**：基于 zap 的日志系统
- ✅ **健康检查**：内置健康检查端点
- ✅ **CORS 支持**：跨域资源共享配置
- ✅ **配置管理**：支持环境变量覆盖

## 快速开始

### 1. 安装依赖工具

```bash
make init
```

### 2. 生成代码

```bash
make all
```

### 3. 配置数据库和 Redis

编辑 `configs/config.yaml` 或使用环境变量：

```bash
export DB_USER=root
export DB_PASSWORD=password
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=demo

export REDIS_HOST=localhost
export REDIS_PORT=6379
```

### 4. 运行项目

```bash
make build
./bin/app -conf configs
```

或者直接运行：

```bash
go run cmd/app/main.go cmd/app/gen.go -conf configs
```

## 开发指南

### 添加新的 API

1. 在 `api/demo/v1/` 目录下定义新的 proto 文件
2. 运行 `make api` 生成代码
3. 在 `internal/service/` 中实现服务逻辑
4. 在 `internal/server/` 中注册服务

### 添加新的数据库模型

1. 在 `internal/models/` 中定义模型
2. 在 `internal/global/global.go` 中添加自动迁移

### 使用对象存储

在配置文件中启用对象存储：

```yaml
data:
  object_storage:
    enabled: true
    provider: minio
    endpoint: localhost:9000
    access_key_id: minioadmin
    secret_access_key: minioadmin
    bucket_name: demo
```

然后在代码中使用：

```go
import "kratos-project-template/provider/storage"

storage := storage.Get()
err := storage.PutObject(ctx, "key", data)
```

## API 端点

- `GET /demo/hello?name=World` - Hello 接口
- `GET /demo/health` - 健康检查

## 环境变量

所有配置项都支持通过环境变量覆盖，使用点号分隔的路径：

- `SERVER_HTTP_ADDR` - HTTP 服务地址
- `SERVER_GRPC_ADDR` - gRPC 服务地址
- `DATA_DATABASE_SOURCE` - 数据库连接字符串
- `DATA_REDIS_ADDR` - Redis 地址
- `LOG_LEVEL` - 日志级别
- `LOG_FORMAT` - 日志格式（json/console）

## 构建和部署

### 本地构建

```bash
make build
```

### Docker 构建

```bash
docker build -t kratos-project-template .
```

### 生产环境建议

1. 使用环境变量管理敏感配置
2. 启用 JSON 格式日志
3. 配置适当的日志级别
4. 启用对象存储（如需要）
5. 配置数据库连接池大小

## 许可证

MIT

## 贡献

欢迎提交 Issue 和 Pull Request！

