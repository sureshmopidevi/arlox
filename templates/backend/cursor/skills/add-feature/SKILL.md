# Skill: Add Feature

Checklist for adding a new domain module.

## Step 1 — Read examples first
Read all files in `.cursor/skills/add-feature/learned/` before writing a single line of code. They show real patterns already established in this project.

## Step 2 — Create module files
```
modules/<name>/repository/repository.go
modules/<name>/service/service.go
modules/<name>/handler/handler.go
modules/<name>/routes/routes.go
```

Follow layer rules (see `.cursor/rules/architecture.mdc`):
- Repository: GORM only
- Service: logic only
- Handler: bind + `response.*` only
- Routes: register on the group passed in

## Step 3 — Wire in main.go
Add to `cmd/server/main.go`:
```go
<name>Repo := <name>Repo.New(db)
<name>Svc  := <name>Service.New(<name>Repo)
<name>H    := <name>Handler.New(<name>Svc)
```
Pass `<name>H` into `server.NewRouter(...)`.

## Step 4 — Register in router.go
Add parameter to `NewRouter` signature and call:
```go
<name>Routes.Register(api, h)
```

## Step 5 — Verify
```bash
make test
make lint
```
Both must pass before declaring done.

## Step 6 — Record learning
Append `.cursor/skills/add-feature/learned/<name>.md` with:
- What endpoints were added
- Any non-obvious patterns used
- Gotchas encountered

Update `.cursor/skills/add-feature/learned/README.md` index.
