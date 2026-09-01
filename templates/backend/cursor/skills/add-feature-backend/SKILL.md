---
name: add-feature-backend
description: >-
  Ship a new Go API resource under modules/<name> with tests green and a
  contracts/<name>.md API contract for web/app. Use when adding a backend
  feature, new REST resource, or modules/<name> domain in this Go service.
---

# Add Feature (Backend)

## User intent

When done, the API exposes the new resource at documented paths, `make test` passes, and **`contracts/<name>.md`** at the workspace root describes endpoints for web and Flutter agents.

## Read first

1. **`.arlox/project.json`** (workspace root) — `stacks.backend.module`, `stacks.backend.dbName`, naming (`kebab`, `snake`).
2. **`contracts/README.md`** — envelope, auth header, contract format.
3. **`contracts/<name>.md`** — create or update before coding if missing.
4. **`.cursor/rules/architecture.mdc`** and **`.cursor/rules/karpathy.mdc`** (if present).
5. **`.cursor/skills/add-feature-backend/learned/README.md`** — prior backend patterns (not API contracts).

**Do not use** for cross-stack work — use workspace `add-feature-fullstack`.

## Success criteria

| Check | Must be true |
|-------|----------------|
| Module wired | `repository` → `service` → `handler` → routes registered |
| Tests | `make test` exits 0 |
| Contract | `contracts/<name>.md` exists and matches implemented endpoints |
| Scope | No drive-by edits outside this feature |

If the feature name or scope is ambiguous, ask before scaffolding.

## Steps

1. **Contract** — write or confirm `contracts/<name>.md` (workspace root). Do not invent endpoints web will need later without documenting them.
2. **Scaffold** `modules/<name>/` — `repository/`, `service/`, `handler/`, `routes/`.
3. **Wire DI** in `cmd/server/main.go` and register routes in `internal/server/router.go`.
4. **Models** — add GORM models under `internal/models/` only if needed.
5. **Verify** — `make test` (and `make lint` if available). Loop until green.
6. **Learn** — run **`.cursor/skills/learn/SKILL.md`** (mandatory). Append to `learned/README.md`; do not put API contracts in `learned/`.

### Layer rules

| Layer | Allowed | Forbidden |
|-------|---------|-----------|
| Repository | GORM, models | HTTP, gin.Context, business rules |
| Service | Business logic | GORM, gin.Context |
| Handler | Bind → service → `response.*` | Business logic, raw `c.JSON` |
| Routes | Register on passed group | New router frameworks |

## Verify

```bash
make test
make lint   # if golangci-lint is installed
```

## Learn

Run **`.cursor/skills/learn/SKILL.md`** after verify passes. Record implementation patterns and gotchas in **`learned/README.md`**. API shape stays in **`contracts/<name>.md`**.

## Anti-patterns

- Business logic in handlers or GORM outside repository
- API contract only in `learned/` instead of `contracts/`
- Skipping `contracts/<name>.md` before implementation
- Refactoring unrelated modules
- Ignoring `.arlox/project.json` when naming DB tables or module paths
