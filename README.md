# arlox

> **Status: Active Development 🚧**
> `arlox` is currently in active development. APIs, CLI commands, and template architectures are evolving rapidly. Feedback, issues, and contributions are welcome!

Scaffold and orchestrate modern **multi-stack workspaces** (Go backend, React/Vite web, and Flutter mobile app) with unified `.code-workspace` configurations, root `Makefile` automation, and **AI Agent guardrails** (Cursor & Antigravity rules and skills) baked in.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)](go.mod)

---

## 🎯 What is arlox intended for?

Building fullstack products often requires managing separate stacks (API backend, web dashboard, and mobile client) with repetitive setup, conflicting configs, and inconsistent developer workflows. 

`arlox` solves this by giving you:

* **Production-ready Scaffolds:**
  * 🐹 **Backend:** Go + Gin REST API with JWT auth, GORM/PostgreSQL migrations, health checks, and modular architecture.
  * ⚛️ **Web:** Vite + React + TypeScript + Tailwind CSS with Zustand session state, React Query, and auth routing.
  * 📱 **Mobile App:** Flutter app with clean navigation, Dio networking, secure storage, and Material 3 design.
* **Unified Workspace:** Single `.code-workspace` tailored for **Cursor**, **VS Code**, and **Antigravity IDE**.
* **One-Command Orchestration:** Root `Makefile` with `make dev` (starts backend, health-polls until ready, and launches web frontend), `make test`, `make doctor`, and `make status`.
* **AI Agent Guardrails:** Pre-configured `.cursor/rules`, Karpathy-style engineering guidelines, `AGENTS.md`, and self-documenting `.cursor/skills` in each stack so AI coding assistants write clean, consistent code.
* **Modular Stack Management:** Generate all stacks together or start small and incrementally add stacks later with `arlox add`.

---

## Installation

### Option 1: 1-Line Installer (Recommended — Zero Manual PATH Setup)

Installs `arlox`, automatically links it into your system `$PATH`, and persists configuration in your shell rc:

```bash
curl -fsSL https://raw.githubusercontent.com/sureshmopidevi/arlox/main/install.sh | bash
```

### Option 2: Go Install (Direct compile via Go toolchain)

If you already have `~/go/bin` configured in your PATH:

```bash
go install github.com/sureshmopidevi/arlox/cmd/arlox@latest
```

### Option 3: From Source (For contributors / local development)

```bash
git clone https://github.com/sureshmopidevi/arlox.git ~/arlox
cd ~/arlox
./install.sh
```

Then create a project:

```bash
cd ~/Projects
arlox create myapp
```

---

## Quick reference

```bash
arlox                         # help
arlox version                 # print version (also: arlox -v)
arlox doctor                  # check toolchains (go, node, npm, flutter, git) & PATH
arlox upgrade                 # git pull + rebuild + reinstall from ~/arlox
arlox create myapp            # interactive — pick stacks, shows path first
arlox create myapp --backend --web --app
cd myapp && arlox add --app   # add stacks later
cd myapp && arlox repair      # audit & restore missing configs, rules, or deps
arlox uninstall               # remove arlox and legacy binaries from PATH
make verify                   # dev smoke test (temp dir only, auto-deleted)
make help
```

Release version lives in [`internal/version/VERSION`](internal/version/VERSION). **Bump it on every CLI or template change** (see [`.cursor/rules/versioning.mdc`](.cursor/rules/versioning.mdc)).

To update an installed binary after pulling changes:

```bash
arlox upgrade
# or: arlox upgrade --source ~/arlox
# or: ARLOX_HOME=~/arlox arlox upgrade --no-pull
```

---

## What you get

```text
myapp/
  myapp.code-workspace   # open in Cursor, VS Code, or Antigravity IDE
  .cursor/              # fullstack sequential skill
  backend/              # Go + Gin (own git repo)
  web/                  # React + Tailwind (own git repo)
  app/                  # Flutter home page (own git repo)
```

---

## Agent guardrails

Each stack ships `.cursor/rules`, `.cursor/skills/add-feature-*` (with auto `learned/` logs), Karpathy guidelines, and `AGENTS.md`.

```bash
arlox skills update    # refresh templates; keeps your learned/ files
```

---

## Makefile (optional)

| Command | What |
|---------|------|
| `source ./install.sh` | **Use this** — full setup |
| `make verify` | Automated test (temp folder, not your cwd) |
| `make doctor` | Check tools on PATH |
| `make uninstall` | Remove binary |

See [scripts/verify.md](scripts/verify.md) for manual verification.

---

## License

This project is licensed under the [MIT License](LICENSE).
