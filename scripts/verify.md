# arlox verification

## Automated (recommended)

From the repo root — exports `$(go env GOPATH)/bin` if needed:

```bash
make verify
```

Or run the script directly:

```bash
source scripts/env.sh && ensure_arlox_path && ensure_local_tool_paths
./scripts/verify.sh
```

The script creates temp projects only — nothing in your cwd.

---

## What `make verify` covers

| Step | Checks |
|------|--------|
| 1 | Full create (`demo --backend --web --app --web-ui tailwind`) |
| 1b | Partial workspace + unique Postgres ports |
| 2–3 | Duplicate skip, workspace folder sync, `add --web` |
| 4 | Backend `go build` |
| 5 | Web `npm run build` + Vitest |
| 5c | Web-only create with `--web-ui shadcn` + build |
| 5d | Web-only create with `--web-ui antd` + build |
| 6 | Flutter analyze |
| 7 | Guardrails (skills, `design-system.mdc`, contracts) |
| 8–10 | repair, skills update, upgrade smoke |
| 11 | Brownfield `arlox init` |

Go tests also render all six design system overlays: `go test ./internal/generate/...`

---

## Manual checklist

Run from a throwaway directory (not inside the arlox repo).

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
cd /tmp
rm -rf arlox-demo arlox-demo2
```

### 1. Create all stacks

```bash
arlox create arlox-demo --backend --web --app --web-ui tailwind
```

Expect: all stacks, `arlox-demo.code-workspace`, `contracts/auth.md`, `web/.arlox/design-system.json`.

### 2. No duplicates

```bash
arlox create arlox-demo --backend --web --app --web-ui tailwind
```

Expect: all stacks **skipped (already exists)**.

### 3. Add after partial create

```bash
arlox create arlox-demo2 --backend
cd arlox-demo2
arlox add --web --web-ui shadcn
```

Expect: web generated with `components.json`; backend skipped.

### 4. Backend health (needs Postgres)

```bash
cd /tmp/arlox-demo/backend
make setup
make run   # curl -s localhost:8080/health/live
```

### 5. Web build

```bash
cd /tmp/arlox-demo/web
npm install && npm run build && npm test
```

### 6. Flutter analyze

```bash
cd /tmp/arlox-demo/app
flutter pub get && flutter analyze
```

### 7. Guardrails present

```bash
test -f /tmp/arlox-demo/.cursor/skills/add-feature-fullstack/SKILL.md
test -f /tmp/arlox-demo/contracts/auth.md
test -f /tmp/arlox-demo/web/.arlox/design-system.json
test -f /tmp/arlox-demo/web/.cursor/rules/design-system.mdc
test -f /tmp/arlox-demo/backend/.cursor/rules/karpathy.mdc
```

### 8. Design system variants

```bash
arlox create arlox-antd --web --web-ui antd --out /tmp
cd /tmp/arlox-antd/web && npm install && npm run build
```

### 9. Brownfield init

```bash
mkdir -p /tmp/brownfield/backend
printf 'module example\n\ngo 1.22\n' > /tmp/brownfield/backend/go.mod
arlox init /tmp/brownfield --name brownfield
test -f /tmp/brownfield/Makefile
```

### 10. Skills update

```bash
cd /tmp/arlox-demo
arlox skills update
```

Expect: `learned/` folders untouched.

---

## Full-stack feature (manual / agent)

1. Add `contracts/<feature>.md` at workspace root
2. Backend via `add-feature-backend` → `make test`
3. Web via `add-feature-web` (respect `design-system.json`) → `npm run build`
4. App via `add-feature-mobile` → `flutter analyze`

Use workspace `add-feature-fullstack` skill for strict ordering.
