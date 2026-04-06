# Workspaces Platform

A self-hosted, Kubernetes-native developer workspace platform — similar to [Gitpod](https://gitpod.io) and [Coder](https://coder.com) — built from scratch as a solo portfolio project.

Each workspace is a fully isolated, browser-accessible development environment running inside a Kubernetes namespace, with a terminal (ttyd), optional code-server, and per-workspace subdomain routing via Cloudflare Tunnel.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        User (Browser)                        │
└────────────────────────┬────────────────────────────────────┘
                         │
              ┌──────────▼──────────┐
              │   Cloudflare Tunnel  │  *.yourdomain.com
              └──────────┬──────────┘
                         │ wildcard → ingress-nginx
              ┌──────────▼──────────┐
              │    Frontend (Next.js) │
              └──────────┬──────────┘
                         │ REST API
              ┌──────────▼──────────┐
              │   Backend (Go/Echo)  │
              └──────┬──────┬───────┘
                     │      │
           Redis Pub/Sub   PostgreSQL
                     │
              ┌──────▼──────────────┐
              │  Operator (Go)       │
              │  └─ Subscriber       │  watches Redis channel
              │  └─ Worker Queue     │  async job processing
              │  └─ K8s Provisioner  │  creates namespaces/pods
              └─────────────────────┘
                         │
              ┌──────────▼──────────┐
              │   Kubernetes Cluster │
              │  ├─ ws-{id}-{user}/  │  per-workspace namespace
              │  │  ├─ Deployment    │  ttyd + optional code-server
              │  │  ├─ PVC           │  persistent storage
              │  │  ├─ Service       │  ClusterIP
              │  │  ├─ Ingress       │  WebSocket-enabled
              │  │  └─ RBAC          │  per-tenant isolation
              │  └─ platform-infra/  │
              │     ├─ Elasticsearch │  log aggregation
              │     └─ Fluent Bit    │  log collection (DaemonSet)
              └─────────────────────┘
```

---

## Components

### `workspaces/` — Backend API (Go)

REST API built with [Echo](https://echo.labstack.com/). Handles authentication, workspace lifecycle, template management, and log streaming.

- Workspace CRUD with quota enforcement per user
- Template system with variable substitution and file storage in MinIO
- Real-time log streaming via **SSE** (Server-Sent Events) from Elasticsearch
- Redis cache-aside pattern for workspace and template data
- Publishes workspace lifecycle events to Redis Pub/Sub channel

**Stack:** Go, Echo, sqlx, PostgreSQL, Redis, MinIO, Elasticsearch, Cloudflare API

### `operator/` — Kubernetes Operator (Go)

Standalone service that subscribes to Redis Pub/Sub and provisions/deprovisions Kubernetes resources asynchronously.

- Subscribes to `workspace:events` Redis channel
- Internal job queue (buffered channel) for async processing
- Provisions per-workspace: Namespace, ResourceQuota, RBAC, Deployment, PVC, Service, Ingress
- Renders YAML manifests from MinIO templates using `text/template`
- Applies manifests via Kubernetes dynamic client + RESTMapper

**Stack:** Go, client-go (dynamic client), Redis, MinIO, PostgreSQL

### `frontend-workspaces/` — Frontend (Next.js)

Web UI for managing workspaces, templates, and logs.

- Workspace dashboard with status tracking
- Terminal access via embedded ttyd iframe (token-cached auth)
- Real-time log viewer with SSE-based streaming
- Multi-step template creation with variable/addon management
- Dark industrial aesthetic (IBM Plex Sans/Mono, shadcn/ui)

**Stack:** Next.js 14, TypeScript, TanStack Query, shadcn/ui, Tailwind CSS

---

## Event Flow

```
User creates workspace (POST /workspaces)
    │
    ├─ INSERT workspaces (status=pending)
    ├─ redis.Publish("workspace:events", payload)
    │
    └─► Operator Subscriber receives event
            │
            └─► jobQueue ← WorkspaceJob
                    │
                    └─► handleCreate()
                            ├─ CreateNamespace
                            ├─ CreateResourceQuota
                            ├─ SetupRBAC
                            ├─ ExecuteDeployment (render YAML → apply to K8s)
                            └─ UpdateWorkspaceStatus (running)
```

---

## Infrastructure

Built and runs on a self-hosted homelab:

| Component  | Spec                                              |
| ---------- | ------------------------------------------------- |
| Hypervisor | Proxmox VE on Xeon E5-2673 v3, 62GB RAM           |
| Kubernetes | kubeadm cluster on LXC containers                 |
| CNI        | Flannel                                           |
| Ingress    | ingress-nginx with WebSocket annotations          |
| Storage    | local-path-provisioner                            |
| Tunnel     | Cloudflare Tunnel (wildcard `*.wfdnstore.online`) |
| Registry   | GitLab Container Registry (self-hosted)           |
| CI/CD      | GitLab Runner → build → push → kubectl deploy     |
| Logs       | Fluent Bit → Elasticsearch → SSE stream           |

---

## Getting Started

### Prerequisites

- Kubernetes cluster (kubeadm or any distribution)
- PostgreSQL
- Redis
- MinIO
- Cloudflare Tunnel (for external access)
- `kubectl` configured with cluster access

### Run locally

```bash
# Backend
cd workspaces
go run ./cmd/main.go

# Operator
cd operator
go run ./cmd/main.go

# Frontend
cd frontend-workspaces
npm install
npm run dev
```

---

## Key Design Decisions

**Why Redis Pub/Sub instead of direct job queue in the API?**
Backend and operator are separate services — decoupling them via Redis means the API never blocks on K8s provisioning, and the operator can be restarted independently without losing events (with a persistent queue layer added later).

**Why dynamic client + RESTMapper instead of typed client-go?**
Workspace templates are arbitrary YAML stored in MinIO — the operator doesn't know the resource types at compile time. Dynamic client allows applying any valid Kubernetes manifest at runtime.

**Why ttyd as a sidecar instead of a separate service?**
Co-locating ttyd in the workspace pod means terminal access shares the pod's filesystem and network namespace — the terminal is actually inside the workspace, not a separate container connecting to it.

---

## License

MIT
