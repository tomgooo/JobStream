# JobStream - 职位聚合 & 实时推送平台

> 一个为「找工作」场景量身定制的 Go 微服务项目：聚合多渠道职位、根据用户订阅偏好做智能匹配，并通过长连接/推送服务实时通知用户。

> ⚠️ 当前处于早期开发阶段（v0），功能与架构仍在迭代中。

---

## 🎯 项目目标

- 聚合多个渠道的招聘信息（招聘网站 / 内推 / 自建数据源等）
- 支持用户配置订阅条件（城市、薪资、技术栈等）
- 定期/实时匹配符合条件的职位，并推送给在线用户
- 练习并展示一套现代后端工程能力：
  - Go 语言 + 微服务架构
  - tRPC-Go / RPC 通信 / TCP 长连接
  - MySQL + Redis
  - Docker /（后续）Kubernetes 部署
  - CI/CD、监控、日志等工程实践

---

## 🏗 技术栈（规划）

> ✅ 已完成 / ⏳ 进行中 / 📝 计划中

- **语言**
  - ✅ Go 1.x

- **Web & 微服务框架**
  - ⏳ go-zero（REST API & 部分 RPC 服务）
  - 📝 tRPC-Go（内部高性能 RPC & 推送服务）
  - 📝 自定义 TCP/WebSocket 长连接（基于 tRPC-Go/tnet）

- **数据存储**
  - ✅ MySQL（主业务数据：用户、职位、订阅等）
  - ✅ Redis（缓存、限流、会话/在线用户状态等）

- **基础设施 / 运维**
  - ✅ Docker / Docker Compose（本地一键启动）
  - 📝 Kubernetes（Kind/Minikube 本地集群）
  - 📝 Prometheus + Grafana（监控 & 可视化）

- **工程实践**
  - ⏳ GitHub Actions（CI：测试 & 构建）
  - 📝 代码规范与简单的目录约定
  - 📝 单元测试 & 简单压测脚本

---

## 🔧 快速开始（开发环境）

### 1. 环境要求

- Go 1.x（建议 1.21+）
- Git
- Docker Desktop

### 2. 克隆项目

git clone https://github.com/<your-name>/jobstream.git
cd jobstream

### 3. 启动基础依赖（MySQL & Redis）

项目使用 Docker 启动 MySQL 和 Redis，无需本机单独安装数据库。

```
# 在项目根目录
docker compose up -d
# 或旧版本 Docker：docker-compose up -d
```

启动完成后包含：

- MySQL：监听 `localhost:3306`
- Redis：监听 `localhost:6379`

> 初始数据库、用户和密码可以在 `docker/mysql/init.sql` 与 `docker-compose.yml` 中查看和修改。



### 4. 运行后端服务（最小可用版本）

> v0 阶段先提供一个简单的 HTTP 服务，暴露 `/healthz` 等基础接口，确认环境无误。

```
go run ./cmd/jobstream
```

访问：

- 健康检查：http://localhost:8080/healthz

## 📁 项目结构（v0 规划）

> 随着项目迭代，目录会不断完善，这里先给出一个目标结构。

```
jobstream/
  ├── cmd/
  │   └── jobstream/          # 主入口（v0 单体 / 网关）
  ├── internal/
  │   ├── user/               # 用户模块（注册/登录/鉴权等）
  │   ├── job/                # 职位模块（职位信息、搜索等）
  │   ├── subscription/       # 订阅模块（用户偏好、匹配逻辑）
  │   └── common/             # 公共工具（配置、日志、中间件等）
  ├── services/               # 后续拆出的微服务（api-gateway, user-service 等）
  ├── configs/                # 配置文件（YAML/JSON）
  ├── docker/
  │   ├── mysql/
  │   │   └── init.sql        # MySQL 初始化脚本
  │   └── docker-compose.yml  # Docker Compose 配置（也可放根目录）
  ├── docs/
  │   ├── arch-overview.md    # 架构概览（TODO）
  │   ├── service-design.md   # 服务设计（TODO）
  │   └── dev-guide.md        # 开发者指南（TODO）
  ├── go.mod
  ├── go.sum
  └── README.md
```

------

## 🗺 开发路线图（Roadmap）

> 这是项目规划的里程碑，方便自己和面试官快速理解项目演进过程。

- **Milestone 1：最小可用版本（单体 / 少量服务）**

  - 用户注册 / 登录

  - 职位录入 / 查询

  - 订阅条件配置 + 简单匹配接口

    > - “正在进行中：已完成基础 HTTP 服务 + MySQL/Redis 环境”。

- **Milestone 2：微服务化 + RPC**

  - 拆分 api-gateway / user-service / job-service / matching-service
  - 引入 go-zero RPC

- **Milestone 3：推送服务 & tRPC-Go**

  - tRPC-Go push-service
  - 长连接 / 自定义协议推送新职位

- **Milestone 4：DevOps & 监控**

  - Docker / Docker Compose 完善
  - GitHub Actions CI
  - Prometheus + Grafana 监控面板

- **Milestone 5：K8s 部署 & 性能压测**

  - 在 Kind/Minikube 跑一套精简版集群
  - 简单 QPS / 延迟测试与优化

------

## 💼 面向面试官的亮点（预留）

> TODO：当主要功能完成后，这里会总结本项目的技术亮点、架构思路和踩坑经验，方便在面试中快速讲解。

- 微服务架构与服务拆分思路
- tRPC-Go 在推送场景中的应用
- 高并发 / 高可用相关设计（缓存、限流、降级等）
- CI/CD + 监控链路
