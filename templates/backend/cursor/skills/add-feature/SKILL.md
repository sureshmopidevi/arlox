---
name: add-feature
description: >-
  Scaffolds a new Go domain module (repository, service, handler, routes),
  wires it into main/router, verifies with make test, and records an API
  contract in learned/. Use when adding a backend feature, new API resource,
  or new modules/<name> domain in this Go service.
---

# Add Feature (Backend)

Surgical checklist for a new domain module. Touch only what this feature needs.

## When to use

- New API resource / domain module under `modules/`
- Extending the Gin API with CRUD or domain endpoints

**Do not use** for cross-stack work — use workspace `add-feature-fullstack`.

## Before coding

1. Read `.cursor/rules/architecture.mdc` and `.cursor/rules/karpathy.mdc` (if present).
2. Read `.cursor/skills/add-feature/learned/` for prior patterns.
3. State assumptions + success criteria:
   - Endpoints to add
   - Auth required? yes/no
   - Verify: `make test` (+ `make lint` if available)

If the feature name or scope is ambiguous, ask before scaffolding.

## Steps

### 1. Scaffold module

```text
modules/<name>/
  repository/repository.go
  service/service.go
  handler/handler.go
  routes/routes.go
```

Layer rules:

| Layer | Allowed | Forbidden |
|-------|---------|-----------|
| Repository | GORM, models | HTTP, gin.Context, business rules |
| Service | Business logic | GORM, gin.Context |
| Handler | Bind → service → `response.*` | Business logic, raw `c.JSON` |
| Routes | Register on passed group | New router frameworks |

### 2. Wire DI in `cmd/server/main.go`

```go
<name>Repo := <name>Repo.New(db)
<name>Svc  := <name>Service.New(<name>Repo /*, cfg if needed */)
<name>H    := <name>Handler.New(<name>Svc)
```

Pass `<name>H` into `server.NewRouter(...)`.

### 3. Register routes in `internal/server/router.go`

```go
<name>Routes.Register(api, <name>H)
```

Prefer `/api/v1/<resource>` paths consistent with existing modules (`auth`, `ping`).

### 4. Models / migrations

- Add GORM models under `internal/models/` only if needed.
- Keep seed data minimal; only what the feature requires for local demo.

### 5. Verify

```bash
make test
make lint   # if golangci-lint is installed
```

Loop until green. Do not declare done with failing tests.

### 6. Record API contract (required for web/app)

Write `.cursor/skills/add-feature/learned/<name>.md`:

```markdown
# <name> — API contract

## Endpoints
| Method | Path | Auth | Request | Response |
|--------|------|------|---------|----------|
| GET | /api/v1/<name> | yes | — | `{ data: [ ... ] }` |
| POST | /api/v1/<name> | yes | `{ ... }` | `{ data: { ... } }` |

## Errors
| Status | When |
|--------|------|
| 400 | validation |
| 401 | unauthorized |
| 404 | not found |

## Notes for web/app
- Envelope `{ data: T }` — clients unwrap automatically
- Header: `Authorization: Bearer <token>`
```

Update `learned/README.md` index with a one-line entry.

## Anti-patterns

- Business logic in handlers
- Calling GORM from handlers/services incorrectly (GORM only in repository)
- Skipping the learned/ contract doc
- Refactoring unrelated modules
- Adding a second ORM/router/DI framework

## Done checklist

- [ ] Module files created and wired
- [ ] Routes registered
- [ ] `make test` passes
- [ ] `learned/<name>.md` contract written
- [ ] No drive-by edits outside this feature
