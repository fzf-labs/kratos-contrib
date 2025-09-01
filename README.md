# Kratos-Contrib

`kratos-contrib` 是一个为 [Go-Kratos](https://github.com/go-kratos/kratos) 框架提供的扩展工具集。该项目提供了一系列实用的中间件、引导工具和整合功能，帮助开发者更高效地构建微服务应用。

## 功能特性

- **引导工具**：简化应用的配置和初始化流程
- **中间件集成**：
  - 日志记录（支持 zap、zerolog）
  - 流量控制
  - 指标收集
  - 上下文传递
- **服务注册与发现**：支持 Consul、Etcd、Nacos 等注册中心
- **链路追踪**：集成 OpenTelemetry 实现分布式追踪
- **HTTP & gRPC 支持**：同时支持 HTTP 和 gRPC 协议的服务构建

## 安装

确保您已安装 Go 1.23 或更高版本，然后执行：

```bash
go get github.com/fzf-labs/kratos-contrib
```

## 快速开始

下面是一个基本的使用示例：

```go
package main

import (
	"github.com/fzf-labs/kratos-contrib/bootstrap"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

func main() {
	// 创建服务
	service := bootstrap.NewService(
		"my-service",
		"v1.0.0",
		"A sample service",
	)
	
	// 初始化应用
	cfg, logger, reg, dis := bootstrap.Bootstrap(service)
	
	// 创建 Kratos 应用实例
	app := kratos.New(
		kratos.Name(service.Name),
		kratos.Version(service.Version),
		kratos.Metadata(service.Metadata),
		kratos.Logger(logger),
		kratos.Registrar(reg),
	)
	
	// 启动应用
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## 配置说明

配置文件示例（`configs/config.yaml`）：

```yaml
server:
  http:
    addr: :8000
    timeout: 1s
  grpc:
    addr: :9000
    timeout: 1s
registry:
  type: consul
  consul:
    address: 127.0.0.1:8500
    scheme: http
logger:
  level: info
  encoding: json
  outputPaths:
    - stdout
    - logs/app.log
trace:
  endpoint: http://localhost:14268/api/traces
```

## 中间件使用

### 日志中间件

```go
import (
	"github.com/fzf-labs/kratos-contrib/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware"
)

// 使用日志中间件
middlewares := []middleware.Middleware{
	logging.Server(),
}
```

### 流量限制中间件

```go
import (
	"github.com/fzf-labs/kratos-contrib/middleware/limiter"
	"github.com/go-kratos/kratos/v2/middleware"
)

// 使用限流中间件
middlewares := []middleware.Middleware{
	limiter.NewLimiter(100), // 限制为每秒100个请求
}
```

## 贡献指南

欢迎贡献代码或提交问题！请先 fork 本仓库，然后提交 pull request。

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开一个 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件 # 测试更新
