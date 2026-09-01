# arlox

Scaffold and orchestrate **multi-stack workspaces** — Go backend, React/Vite web, and Flutter app — with a unified IDE workspace, root Makefile, Docker Postgres, and AI agent guardrails built in.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go)](go.mod)

> **Active development.** CLI commands and templates evolve quickly. Feedback and contributions welcome.

---

## Table of contents

- [What you get](#what-you-get)
- [Prerequisites](#prerequisites)
- [Install arlox](#install-arlox)
- [Quick start — your first project](#quick-start--your-first-project)
- [CLI reference](#cli-reference)
- [Generated project layout](#generated-project-layout)
- [Day-to-day workflow](#day-to-day-workflow)
- [Database (Postgres)](#database-postgres)
- [AI agent guardrails](#ai-agent-guardrails)
- [Maintaining arlox & workspaces](#maintaining-arlox--workspaces)
- [Contributing & verification](#contributing--verification)
- [License](#license)

---

## What you get

| Layer | What arlox provides |
|-------|---------------------|
| **Backend** | Go + Gin REST API, JWT auth, GORM/PostgreSQL, health checks, modular layout |
| **Web** | React 19 + Vite + TypeScript + Tailwind, Zustand, React Query, Vitest |
| **App** | Flutter + Riverpod + go_router, Dio, secure storage |
| **Workspace** | `.code-workspace` for Cursor, VS Code, Antigravity IDE |
| **Orchestration** | Root `Makefile` — `make dev`, `make test`, `make doctor`, `make status` |
| **AI guardrails** | `.cursor/rules`, `.cursor/skills`, Karpathy guidelines, `AGENTS.md` per stack |

Start with all stacks or add them later (`arlox add`). Partial workspaces work — missing stacks are skipped automatically.

---

## Prerequisites

Install these **before** generating stacks that need them:

| Tool | Required for | Check |
|------|----------------|-------|
| [Go](https://go.dev/dl/) 1.26+ | arlox CLI, backend stack | `go version` |
| [Git](https://git-scm.com/) | all stacks | `git version` |
| [Node.js](https://nodejs.org/) + npm | web stack | `node -v && npm -v` |
| [Flutter](https://flutter.dev/) | app stack | `flutter doctor` |
| [Docker](https://www.docker.com/) (optional) | local Postgres via Compose | `docker compose version` |

Run **`arlox doctor`** anytime to see what is on your PATH and whether GOPATH/bin is configured.

---

## Install arlox

### Option 1 — One-line installer (recommended)

Installs the binary, links it into your PATH, and updates your shell rc:

```bash
curl -fsSL https://raw.githubusercontent.com/sureshmopidevi/arlox/main/install.sh | bash
```

### Option 2 — Go install

Requires `~/go/bin` on your PATH:

```bash
go install github.com/sureshmopidevi/arlox/cmd/arlox@latest
```

### Option 3 — From source (contributors)

```bash
git clone https://github.com/sureshmopidevi/arlox.git ~/arlox
cd ~/arlox
./install.sh
```

Confirm:

```bash
arlox version
```

---

## Quick start — your first project

### 1. Create a workspace

Interactive (pick stacks in the terminal):

```bash
cd ~/Projects
arlox create myapp
```

Or specify stacks upfront:

```bash
arlox create myapp --backend --web --app
```

Useful flags:

| Flag | Purpose |
|------|---------|
| `--module github.com/you/myapp-backend` | Go module path (default: `github.com/example/<name>-backend`) |
| `--org com.yourcompany` | Flutter org (default: `com.example`) |
| `--out ~/Projects` | Parent directory for the new folder |
| `--open` | Open `.code-workspace` in Cursor, VS Code, or Antigravity after create |

During create you will see **live terminal output** for `go mod tidy`, `npm install`, `flutter create`, etc.

### 2. Open in your editor

```bash
cursor --classic myapp/myapp.code-workspace   # Cursor IDE
code myapp/myapp.code-workspace               # VS Code
agy-ide myapp/myapp.code-workspace            # Antigravity
```

### 3. Start Postgres (backend projects)

From the project root:

```bash
cd myapp
make db-up          # starts Postgres via docker-compose (host port 5433)
make backend.setup  # env file + migrations prep
```

### 4. Run the stack

```bash
make dev            # backend (background) → health poll → web dev server
```

In another terminal, for mobile:

```bash
make app.run
```

### 5. Verify everything

```bash
make status         # which stacks exist
make doctor         # toolchain + stack health
make test           # backend + web (Vitest) + app tests
```

---

## CLI reference

### Core commands

```bash
arlox                         # help
arlox version                 # print version (also: arlox -v)
arlox doctor                  # toolchains, PATH, workspace detection
```

### Create & extend projects

```bash
arlox create myapp                    # interactive stack picker
arlox create myapp --backend --web    # partial workspace
cd myapp && arlox add --app            # add stacks later
```

### Maintain workspaces

```bash
cd myapp && arlox repair              # restore missing configs, rules, skills, deps
cd myapp && arlox repair --force      # overwrite locally modified skills/rules
cd myapp && arlox skills update       # refresh .cursor templates (keeps learned/)
```

`repair` also reports **scaffold drift** — files that differ from the original arlox template (see `.origin-manifest.json` per stack).

### Maintain arlox itself

```bash
arlox upgrade                         # git pull + rebuild + reinstall from ~/arlox
arlox upgrade --source ~/arlox        # explicit source path
arlox upgrade --no-pull               # rebuild without git pull
arlox uninstall                       # remove binary from PATH
```

After upgrade you will see step-by-step success lines (`build complete`, `installed to …`, `reinstalled` or `upgraded X → Y`).

---

## Generated project layout

```text
myapp/
  myapp.code-workspace      # open this in your IDE
  docker-compose.yml        # Postgres 16 (host port 5433)
  Makefile                  # orchestrates all stacks
  README.md
  .cursor/                  # fullstack skill + rules

  backend/                  # own git repo — Go + Gin + JWT + GORM
  web/                      # own git repo — React + Vite + Tailwind + Vitest
  app/                      # own git repo — Flutter + Riverpod
```

Each stack includes its own `.cursor/rules`, `.cursor/skills`, and `AGENTS.md`.

---

## Day-to-day workflow

Run these from the **project root** (inside `myapp/`):

| Command | What it does |
|---------|----------------|
| `make help` | List all targets |
| `make status` | Show which stacks exist |
| `make dev` | Start backend, wait for `/health/live`, launch web dev server |
| `make test` | Run tests for all present stacks |
| `make build` | Build backend + web |
| `make doctor` | Check toolchains per stack |
| `make db-up` | Start Postgres (Docker) |
| `make db-down` | Stop Postgres container |

Per-stack shortcuts: `make backend.run`, `make web.dev`, `make app.run`, `make web.test`, etc.

Add a stack you skipped at create time:

```bash
arlox add --web
```

---

## Database (Postgres)

Generated projects include **`docker-compose.yml`** with Postgres 16 on **host port 5433** (avoids clashing with a local Postgres on 5432).

Backend config (`backend/configs/local/app.env`):

```text
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=<project-name>
```

Without Docker, the backend Makefile falls back to Homebrew Postgres + `createdb` on macOS.

---

## AI agent guardrails

arlox is designed for **AI-assisted development** in Cursor and Antigravity:

- **Per-stack rules** — architecture, Tailwind, Makefile conventions, Karpathy guidelines
- **Skills** — `add-feature-backend`, `add-feature-web`, `add-feature-mobile`, `add-feature-fullstack`
- **Fullstack order** — backend → web → app, one stack at a time, with API contract docs in `learned/`

Refresh templates after upgrading arlox:

```bash
cd myapp && arlox skills update
```

Locally modified skill files are preserved unless you pass `--force`.

---

## Maintaining arlox & workspaces

| Situation | Command |
|-----------|---------|
| Updated arlox source repo | `arlox upgrade` |
| Missing Makefile, `.env`, or `.cursor` files | `cd myapp && arlox repair` |
| New arlox version, want latest skills | `arlox skills update` |
| Check what changed vs original scaffold | `arlox repair` (drift section) |
| Tooling problems | `arlox doctor` |

---

## Contributing & verification

For arlox developers (from the repo root):

```bash
git clone https://github.com/sureshmopidevi/arlox.git ~/arlox
cd ~/arlox
./install.sh

make test      # go test ./...
make verify    # full smoke test in a temp dir (create, build, test, repair, upgrade)
make doctor    # local toolchain check
```

Release version: [`internal/version/VERSION`](internal/version/VERSION) (currently **0.10.0**). Bump on every CLI or template change — see [`.cursor/rules/versioning.mdc`](.cursor/rules/versioning.mdc) and [`AGENTS.md`](AGENTS.md).

Manual verification checklist: [`scripts/verify.md`](scripts/verify.md).

---

## License

MIT — see [LICENSE](LICENSE).
