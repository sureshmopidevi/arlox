# vibeit

Scaffold multi-stack workspaces with **Cursor rules and skills** baked in.

## Install (one command)

```bash
cd ~/vibeit
source ./install.sh
```

That builds, installs to `~/go/bin`, adds Go bin to `~/.zshrc`, and makes `vibeit` work **immediately** in your current terminal.

Then create a project:

```bash
cd ~/Projects
vibeit create myapp
```

## Quick reference

```bash
vibeit create myapp              # interactive — pick stacks, shows path first
vibeit create myapp --backend --web --app
cd myapp && vibeit add --app     # add stacks later
make verify                      # dev smoke test (temp dir only, auto-deleted)
make help
```

## What you get

```text
myapp/
  myapp.code-workspace
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
