# arlox

Scaffold multi-stack workspaces with **Cursor rules and skills** baked in.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)](go.mod)

---

## Installation

### Option 1: Go Install (Recommended for external users — no git clone needed)

If you have Go installed (Go 1.22+; generated backend stacks use Go 1.22+):

```bash
go install github.com/sureshmopidevi/arlox/cmd/arlox@latest
```

Make sure `~/go/bin` is in your `PATH` (e.g. `export PATH="$HOME/go/bin:$PATH"` in your `~/.zshrc` or `~/.bashrc`).

### Option 2: From Source (For contributors / local development)

```bash
git clone https://github.com/sureshmopidevi/arlox.git ~/arlox
cd ~/arlox
source ./install.sh
```

This builds, installs to `~/go/bin`, ensures Go bin is in `~/.zshrc`, and makes `arlox` available immediately in your current terminal session.

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
arlox upgrade                 # git pull + rebuild + reinstall from ~/arlox
arlox create myapp            # interactive — pick stacks, shows path first
arlox create myapp --backend --web --app
cd myapp && arlox add --app   # add stacks later
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
