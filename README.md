# OpsPulse

OpsPulse 是一个轻量级、可实际部署的 DevOps Dashboard + Agent 系统。

它由 3 个部分组成：

- `agent`：运行在 Linux 节点上的单二进制 Go Agent，主动向 Server 上报状态。
- `server`：Go + SQLite 服务端，负责鉴权、存储、REST API 和实时事件流。
- `web`：基于 SvelteKit + Tailwind CSS 的 Dashboard，提供赛博朋克风格控制中心界面。

## 功能特性

- Agent 单二进制运行，支持 YAML 配置文件。
- Agent 主动连接 Server，不要求节点暴露入站端口。
- 采集主机名、运行时长、CPU、内存、磁盘、负载、Docker 容器统计和白名单 systemd 服务状态。
- Server 提供 Agent 心跳接收、SQLite 存储、节点查询、事件时间线和 SSE 实时推送。
- Dashboard 提供全局状态总览、节点卡片、详情面板、服务状态列表、事件时间线和离线节点识别。
- 默认只读，不实现远程命令执行。

## 项目结构

```text
.
├── agent
├── server
├── web
├── deploy
│   ├── caddy
│   └── systemd
└── docker-compose.yml
```

## 本地开发

### 环境要求

- Go 1.23+
- Node.js 22+
- npm 10+

### 1. 启动 Server

```bash
cd server
OPS_AGENT_TOKEN=change-me go run ./cmd/server
```

服务默认监听 `:8080`，SQLite 文件默认写入 `./opspulse.db`。

### 2. 启动 Dashboard

```bash
cd web
npm install
PUBLIC_API_BASE=http://localhost:8080 npm run dev -- --host
```

默认访问地址为 `http://localhost:5173`。

### 3. 启动 Agent

先复制示例配置：

```bash
mkdir -p /etc/opspulse
cp agent/agent.example.yaml /etc/opspulse/agent.yaml
```

按实际环境修改以下字段：

- `node_id`
- `server_url`
- `token`
- `service_whitelist`

然后运行：

```bash
cd agent
go run ./cmd/agent --config ./agent.example.yaml
```

## Docker Compose 部署

### 1. 设置环境变量

```bash
export OPS_AGENT_TOKEN='replace-with-long-random-token'
export OPS_DOMAIN='ops.example.com'
```

### 2. 启动服务

```bash
docker compose up -d --build
```

### 3. 部署 Agent

将 `agent` 编译为单文件：

```bash
cd agent
go build -o opspulse-agent ./cmd/agent
```

把二进制放到节点的 `/usr/local/bin/opspulse-agent`，配置文件放到 `/etc/opspulse/agent.yaml`。

systemd 服务文件示例见 `deploy/systemd/opspulse-agent.service`。

## API 概览

- `POST /api/v1/agents/heartbeat`：Agent 心跳上报（Bearer Token 鉴权）
- `GET /api/v1/overview`：Dashboard 总览数据
- `GET /api/v1/nodes`：节点列表
- `GET /api/v1/nodes/:nodeId`：节点详情
- `GET /api/v1/events`：事件时间线
- `GET /api/v1/stream`：SSE 实时流

## 安全说明

- Agent 与 Server 通过 Bearer Token 鉴权。
- Dashboard 只展示逻辑节点标识与主机名，不返回真实来源 IP。
- 第一版为只读监控面板，不实现远程命令执行。
- 建议通过 Caddy 提供 HTTPS，并使用强随机 Token。

## 扩展方向

当前代码结构已经为后续扩展预留了清晰边界，可继续增加：

- Kubernetes 集群健康状态
- GitHub Actions / Gitea Actions 流水线概览
- 日志查看与检索
- 受限控制操作（例如重启白名单服务）
