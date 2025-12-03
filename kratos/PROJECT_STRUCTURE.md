# 项目结构说明

本项目是从 `biya-explorer` 项目中剥离业务代码后创建的项目模板。

## 已剥离的业务代码

以下内容已从原项目中移除，仅保留框架代码：

1. **业务模型** (`internal/models/`)
   - Block, Transaction, Account, Validator 等区块链相关模型
   - 压缩、序列化等业务工具函数

2. **业务服务** (`internal/service/`)
   - 区块链同步服务 (`sync.go`)
   - 归档服务 (`archive.go`)
   - 验证器服务 (`validator.go`)
   - 存储迁移服务 (`storage_migration.go`)
   - 事件总线 (`eventbus.go`)
   - 指标收集 (`metrics.go`)

3. **业务配置**
   - Chain 配置（区块链网络配置）
   - Sync 配置（同步服务配置）
   - Node 配置（节点配置）

4. **业务数据**
   - Injective 特定数据文件
   - OFAC 数据文件
   - 数据库迁移脚本

5. **业务 API**
   - 复杂的区块链查询 API
   - 流式 API（StreamLatestBlocks 等）

## 保留的框架代码

以下基础设施代码已保留并通用化：

1. **服务器框架**
   - HTTP/gRPC 服务器初始化
   - CORS 配置
   - 中间件支持

2. **数据访问层**
   - 数据库连接管理（MySQL/PostgreSQL/SQLite）
   - Redis 缓存客户端
   - 对象存储抽象（MinIO/S3/OSS/COS）

3. **配置管理**
   - Protobuf 配置定义
   - 环境变量支持
   - YAML 配置文件

4. **日志系统**
   - Zap 日志集成
   - 结构化日志
   - 日志级别配置

5. **工具和构建**
   - Makefile 构建脚本
   - Docker 支持
   - 代码生成工具

## 新增的示例代码

为了演示如何使用模板，添加了简单的示例：

1. **Demo Service** (`internal/service/demo.go`)
   - Hello API 示例
   - 健康检查实现

2. **简化的 Proto 定义** (`api/demo/v1/demo.proto`)
   - 基础的 API 定义
   - 健康检查消息

## 使用建议

1. **重命名项目**
   - 修改 `go.mod` 中的模块名
   - 修改所有导入路径
   - 更新配置文件中的服务名

2. **添加业务代码**
   - 在 `internal/models/` 中定义数据模型
   - 在 `internal/service/` 中实现业务逻辑
   - 在 `api/` 中定义 API 接口

3. **配置数据库**
   - 根据需求选择数据库类型
   - 配置连接字符串
   - 添加数据库迁移

4. **扩展功能**
   - 添加认证/授权中间件
   - 集成监控和追踪
   - 添加单元测试和集成测试

## 与原项目的差异

| 特性       | 原项目           | 模板项目         |
| ---------- | ---------------- | ---------------- |
| 业务模型   | 区块链相关       | 无（需自行添加） |
| 业务服务   | 区块链同步、查询 | Demo 服务示例    |
| 配置项     | 包含 Chain、Sync | 仅基础配置       |
| API 复杂度 | 复杂查询、流式   | 简单示例         |
| 依赖       | Injective SDK    | 仅框架依赖       |

## 下一步

1. 根据实际需求添加业务代码
2. 配置开发和生产环境
3. 添加测试覆盖
4. 集成 CI/CD 流程

