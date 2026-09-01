# arlox

Scaffold and orchestrate **multi-stack workspaces** — Go backend, React/Vite web, and Flutter app — with a unified IDE workspace, root Makefile, Docker Postgres, optional web design systems, and AI agent guardrails built in.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go)](go.mod)

> **Active development.** CLI commands and templates evolve quickly. Feedback and contributions welcome.

**Documentation:** [docs/](docs/README.md) — [design systems](docs/web-design-systems.md) · [workspace guide](docs/workspace-guide.md) · [CLI reference](docs/cli-reference.md)

---

## Table of contents

- [What you get](#what-you-get)
- [Prerequisites](#prerequisites)
- [Install arlox](#install-arlox)
- [Quick start — your first project](#quick-start--your-first-project)
- [Web design systems](#web-design-systems)
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
| **Backend** | Go + Gin REST API, JWT auth module, GORM/PostgreSQL, health checks, modular layout |
| **Web** | React 19 + Vite + TypeScript — **your choice of design system** (Tailwind, shadcn/ui, Ant Design, MUI, Chakra, Mantine) |
| **App** | Flutter + Riverpod + go_router, Dio, secure storage, login guards |
| **Workspace** | `.code-workspace` for Cursor, VS Code, Antigravity IDE |
| **Contracts** | `contracts/` — shared API docs for full-stack features |
| **Orchestration** | Root `Makefile` — `make dev`, `make test`, `make doctor`, `make status` |
| **Brownfield** | `arlox init` — add workspace files to existing repos without touching stack source |
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

Interactive (pick stacks, then web design system if web is selected):

```bash
cd ~/Projects
arlox create myapp
```

Or specify everything upfront:

```bash
arlox create myapp --backend --web --app --web-ui shadcn
```

Useful flags:

| Flag | Purpose |
|------|---------|
| `--module github.com/you/myapp-backend` | Go module path (default: `github.com/example/<name>-backend`) |
| `--org com.yourcompany` | Flutter org (default: `com.example`) |
| `--out ~/Projects` | Parent directory for the new folder |
| `--open` | Open `.code-workspace` in Cursor, VS Code, or Antigravity after create |
| `--web-ui <id>` | Web UI library when `--web` is included ([details](#web-design-systems)) |

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
make db-up          # starts Postgres via docker-compose (project-specific host port)
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

## Web design systems

When the **web** stack is included, arlox asks which UI library to scaffold (interactive) or accepts `--web-ui` (scripts/CI).

| `--web-ui` | Library | Notes |
|------------|---------|-------|
| `tailwind` | Tailwind only | Default in non-interactive mode — utilities, no component library |
| `shadcn` | shadcn/ui | Tailwind + CSS variables + `components.json` |
| `antd` | Ant Design | Enterprise admin components |
| `mui` | Material UI | Material Design |
| `chakra` | Chakra UI | Accessible component library |
| `mantine` | Mantine | Full-featured React components |

The choice is saved to **`web/.arlox/design-system.json`**. Agents read this file and `.cursor/rules/design-system.mdc` when adding features — they do not mix libraries.

```bash
arlox create myapp --web --web-ui shadcn
arlox add --web --web-ui antd
```

Full guide: **[docs/web-design-systems.md](docs/web-design-systems.md)**

---

## CLI reference

Summary below. Full flag list: **[docs/cli-reference.md](docs/cli-reference.md)**

### Core commands

```bash
arlox                         # help
arlox version                 # print version (also: arlox -v)
arlox doctor                  # toolchains, PATH, workspace detection
```

### Create & extend projects

```bash
arlox create myapp                    # interactive stack + design system pickers
arlox create myapp --backend --web --web-ui shadcn
cd myapp && arlox add --app           # add stacks later
arlox init                            # brownfield: workspace files on existing stacks
arlox init . --name myapp
```

`init` detects `backend/`, `web/`, and `app/` by marker files and writes orchestration only (Makefile, `.code-workspace`, `contracts/`, root `.cursor/`). It never overwrites stack source. See **[docs/workspace-guide.md](docs/workspace-guide.md)**.

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
  docker-compose.yml        # Postgres 16 (host port unique per project name)
  Makefile                  # orchestrates all stacks
  README.md
  contracts/                # shared API contracts (auth.md seed + your features)
  .cursor/                  # fullstack skill + rules

  backend/                  # own git repo — Go + Gin
  web/                      # own git repo — React + chosen design system
    .arlox/design-system.json
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
arlox add --web --web-ui shadcn
```

---

## Database (Postgres)

Generated projects include **`docker-compose.yml`** with Postgres 16. The **host port is derived from your project name** (range 5433–7432), so `demo` and `demo2` never fight for the same port. It is written to both `docker-compose.yml` and `backend/configs/local/app.env`.

Example (your port will differ by project name):

```text
DB_HOST=localhost
DB_PORT=6127
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=myapp
```

Without Docker, the backend Makefile falls back to Homebrew Postgres + `createdb` on macOS.

---

## AI agent guardrails

arlox is designed for **AI-assisted development** in Cursor and Antigravity:

- **Per-stack rules** — architecture, design system (web), Makefile conventions, Karpathy guidelines
- **Skills** — `add-feature-backend`, `add-feature-web`, `add-feature-mobile`, `add-feature-fullstack`
- **Contracts first** — `contracts/<feature>.md` before implementing cross-stack features
- **Fullstack order** — backend → web → app, one stack at a time

Web features must use the design system from `web/.arlox/design-system.json`. Refresh templates after upgrading arlox:

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
| Existing repo, need workspace orchestration | `arlox init` |
| Tooling problems | `arlox doctor` |

---

## Contributing & verification

For arlox developers (from the repo root):

```bash
git clone https://github.com/sureshmopidevi/arlox.git ~/arlox
cd ~/arlox
./install.sh

make test      # go test ./...
make verify    # full smoke test in a temp dir (create, build, design systems, repair, upgrade)
make doctor    # local toolchain check
```

Release version: [`internal/version/VERSION`](internal/version/VERSION) (currently **0.15.0**). Bump on every CLI or template change — see [`.cursor/rules/versioning.mdc`](.cursor/rules/versioning.mdc) and [`AGENTS.md`](AGENTS.md).

- **Documentation:** [docs/README.md](docs/README.md)
- **Manual verification:** [scripts/verify.md](scripts/verify.md)

---

## License

MIT — see [LICENSE](LICENSE).
