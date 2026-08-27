# vibeit

Scaffold multi-stack workspaces with **Cursor rules and skills** baked in.

## Installation

### Option 1: Go Install (Recommended for external users — no git clone needed)

If you have Go installed (see `go.mod` for the required version; generated backend stacks use Go 1.22+):

```bash
go install github.com/sureshmopidevi/vibeit/cmd/vibeit@latest
```

Make sure `~/go/bin` is in your `PATH` (e.g. `export PATH="$HOME/go/bin:$PATH"` in your `~/.zshrc` or `~/.bashrc`).

### Option 2: From Source (For contributors / local development)

```bash
git clone https://github.com/sureshmopidevi/vibeit.git ~/vibeit
cd ~/vibeit
source ./install.sh
```

This builds, installs to `~/go/bin`, ensures Go bin is in `~/.zshrc`, and makes `vibeit` available immediately in your current terminal session.

Then create a project:

```bash
cd ~/Projects
vibeit create myapp
```

## Quick reference

```bash
vibeit                         # help
vibeit version                 # print version (also: vibeit -v)
vibeit upgrade                 # git pull + rebuild + reinstall from ~/vibeit
vibeit create myapp            # interactive — pick stacks, shows path first
vibeit create myapp --backend --web --app
cd myapp && vibeit add --app   # add stacks later
make verify                    # dev smoke test (temp dir only, auto-deleted)
make help
```

Release version lives in [`internal/version/VERSION`](internal/version/VERSION). **Bump it on every CLI or template change** (see [`.cursor/rules/versioning.mdc`](.cursor/rules/versioning.mdc)).

To update an installed binary after pulling changes:

```bash
vibeit upgrade
# or: vibeit upgrade --source ~/vibeit
# or: VIBEIT_HOME=~/vibeit vibeit upgrade --no-pull
```

## What you get

```text
myapp/
  myapp.code-workspace   # open in Cursor, VS Code, or Antigravity IDE
  .cursor/              # fullstack sequential skill
  backend/              # Go + Gin (own git repo)
  web/                  # React + Tailwind (own git repo)
  app/                  # Flutter home page (own git repo)
```

## Agent guardrails

Each stack ships `.cursor/rules`, `.cursor/skills/add-feature` (with auto `learned/` logs), Karpathy guidelines, and `AGENTS.md`.

```bash
vibeit skills update    # refresh templates; keeps your learned/ files
```

## Makefile (optional)

| Command | What |
|---------|------|
| `source ./install.sh` | **Use this** — full setup |
| `make verify` | Automated test (temp folder, not your cwd) |
| `make doctor` | Check tools on PATH |
| `make uninstall` | Remove binary |

See [scripts/verify.md](scripts/verify.md) for manual verification.

## License

MIT
