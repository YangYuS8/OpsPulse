# OpsPulse

OpsPulse is a terminal-first homelab heartbeat console powered by a lightweight Go agent.

OpsPulse is a lightweight local-first DevOps console focused on one job: showing real host telemetry reported by a real agent.

Core principles:

- the Agent runs on real hosts, not inside Docker Compose
- `server`, `web`, and `caddy` can run through Docker Compose
- when no Agent is connected, the UI shows an empty state instead of fake demo data

OpsPulse v0.1 focuses only on real Agent-reported nodes, resource signals, service health, and event visibility.

![OpsPulse screenshot placeholder](docs/screenshot.png)

## Current Scope

- node heartbeat
- CPU / memory / disk / load
- Docker running/exited count
- systemd whitelist services
- SSE event stream
- Caddy behind existing Nginx

## Not Goals

- remote command execution
- Prometheus replacement
- full monitoring platform
- fake demo data
- complex alerting system

## v0.1 Checklist

- [x] terminal-first console layout
- [x] real Agent heartbeat ingestion
- [x] CPU / memory / disk / load telemetry
- [x] Docker summary and container detail pane
- [x] systemd whitelist service pane
- [x] HTTP / TCP checks
- [x] SSE event stream
- [x] event de-noising on state change
- [x] Compose deployment behind existing reverse proxy
- [ ] screenshot capture for `docs/screenshot.png`

## Local Development

### Requirements

- Go 1.23+
- Node.js 22+
- npm 10+

### 1. Start the Server

```bash
cd server
OPS_AGENT_TOKEN=change-me go run ./cmd/server
```

The server listens on `http://localhost:8080` by default.

### 2. Start the Frontend

```bash
cd web
npm install
npm run dev -- --host
```

Open `http://localhost:5173`.

The frontend uses relative paths for `/api` and `/healthz`. In development, Vite proxies those requests to `http://localhost:8080`.

### 3. Start the Host Agent

Create a host-specific config file, for example:

```yaml
node_id: node-local-01
server_url: http://localhost:8080
token: change-me
interval: 15s
docker_enabled: true
service_whitelist:
  - docker
  - ssh
```

Then run:

```bash
cd agent
go run ./cmd/agent --config /path/to/agent.yaml
```

Notes:

- the dashboard only shows nodes after a real host Agent starts reporting
- if no Agent is connected, the node pane shows an empty but honest state

## Command Bar

The bottom `opspulse>` command bar is a frontend-only interaction surface. It does not execute remote commands.

Supported commands:

- `help`
- `nodes`
- `inspect <nodeId|hostname>`
- `events`
- `clear`

## Docker Compose

Compose is only responsible for:

- `server`
- `web`
- `caddy`

It does not start the Agent.

### Start the Stack

```bash
export OPS_AGENT_TOKEN='replace-with-long-random-token'
docker compose up -d --build
```

If your environment has trouble reaching `proxy.golang.org`, you can explicitly set a Go module proxy:

```bash
docker compose build --build-arg GOPROXY=https://goproxy.cn,direct server
docker compose up -d
```

Caddy remains part of the stack, but only as the internal OpsPulse entrypoint. It does not occupy host `80/443`.

Default mapping:

```text
127.0.0.1:8090 -> caddy:80
```

Override it with:

```bash
export OPS_CADDY_PORT=8090
```

Access URL:

- `http://127.0.0.1:8090/`

Routing:

- `/api/*`, `/healthz` -> `server:8080`
- everything else -> `web:3000`

### Show Real Nodes Behind Compose

To see real nodes in the Compose-hosted UI, start the Agent on a host and point it to the Compose entrypoint.

Example config:

```yaml
node_id: node-host-01
server_url: http://127.0.0.1:8090
token: replace-with-long-random-token
interval: 15s
docker_enabled: true
service_whitelist:
  - docker
  - ssh
  - caddy
```

Then run on the host:

```bash
cd agent
go build -o opspulse-agent ./cmd/agent
./opspulse-agent --config /path/to/agent.yaml
```

If `server_url` points to `http://127.0.0.1:8090`, the Agent goes through Caddy. If you prefer to connect directly to the server, you can point it to `http://127.0.0.1:8080`, assuming you explicitly expose that port yourself.

## Existing Nginx Reverse Proxy Example

If the host already runs Nginx, you can proxy a public or internal domain to `127.0.0.1:8090`:

```nginx
server {
    listen 80;
    server_name ops.example.com;

    location / {
        proxy_pass http://127.0.0.1:8090;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

## Agent Configuration

Deployment example config: `agent/agent.example.yaml`

The Agent remains a single binary with YAML configuration. It reports:

- hostname
- uptime
- CPU usage
- memory usage
- disk usage
- load average
- metrics history for recent heartbeats
- Docker running/exited count and container details
- systemd whitelist services
- HTTP/TCP health checks from whitelist-style config

## API Overview

- `POST /api/v1/agents/heartbeat` - Agent heartbeat ingestion (Bearer Token auth)
- `GET /api/v1/overview` - dashboard overview
- `GET /api/v1/nodes` - node list
- `GET /api/v1/nodes/:nodeId` - node detail
- `GET /api/v1/events` - event stream history
- `GET /api/v1/stream` - SSE live stream
- `GET /healthz` - health check

## systemd Agent Example

See: `deploy/systemd/opspulse-agent.service`

## Security Notes

- Agent-to-server traffic uses Bearer Token authentication.
- The frontend command bar is local UI only, not remote execution.
- The dashboard does not expose source IP addresses.
- The project does not fake nodes or pretend the system is connected when no Agent exists.

## Project Note

See: `docs/PROJECT_NOTE.md`

## Possible Future Extensions

This v0.1 intentionally stops early, but the current structure leaves room for:

- Kubernetes cluster status
- GitHub Actions / Gitea Actions visibility
- richer historical trend views
- log browsing
- constrained control actions
