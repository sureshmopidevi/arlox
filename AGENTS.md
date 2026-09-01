# arlox — Agent Guide

CLI and template generator for multi-stack workspaces (backend, web, app).

## Before making changes

1. Read `.cursor/rules/karpathy.mdc` (if present) and `.cursor/rules/versioning.mdc`.
2. Keep changes surgical — templates affect every new project.

## Versioning (required)

Bump **`internal/version/VERSION`** when you change:

- `templates/**`
- `internal/cli`, `internal/generate`, `internal/workspace`, `internal/ui`
- `docs/**` (when user-facing behavior or CLI docs change materially)

Run `make verify` after bumping. See `.cursor/rules/versioning.mdc` for semver rules.

## Key paths

| Path | Purpose |
|------|---------|
| `internal/version/VERSION` | Release version (single source of truth) |
| `templates/` | Embedded scaffold (backend, web, app, workspace, design-systems) |
| `templates/design-systems/` | Per-library web UI overlays (tailwind, shadcn, antd, …) |
| `internal/generate/` | Template render + post-create setup |
| `docs/` | User documentation (design systems, CLI, workspace guide) |
| `scripts/verify.sh` | Smoke test (uses `./bin/arlox`) |

## Verify

```bash
make test
make verify
```
