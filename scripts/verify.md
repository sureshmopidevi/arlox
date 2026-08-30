# arlox v1 verification

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

## Manual

Run from a throwaway directory (not inside the arlox repo).

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
# or: eval $(make -C ~/arlox path)
cd /tmp
rm -rf arlox-demo arlox-demo2
```

## 1. Create all stacks

```bash
arlox create arlox-demo --backend --web --app
```

Expect: `backend/`, `web/`, `app/`, `arlox-demo.code-workspace`, `.cursor/skills/add-feature-fullstack`.

## 2. No duplicates

```bash
arlox create arlox-demo --backend --web --app
```

Expect: all stacks **skipped (already exists)**.

## 3. Add after partial create

```bash
arlox create arlox-demo2 --backend
cd arlox-demo2
arlox add --web
```

Expect: web generated; backend skipped; workspace folders include both.

## 4. Backend health (needs local Postgres)

```bash
cd /tmp/arlox-demo/backend
make setup
make run   # in another terminal: curl -s localhost:8080/health/live
```

## 5. Web build

```bash
cd /tmp/arlox-demo/web
npm install && npm run build
```

## 6. Flutter analyze

```bash
cd /tmp/arlox-demo/app
flutter pub get && flutter analyze
```

## 7. Guardrails present

```bash
test -f /tmp/arlox-demo/.cursor/skills/add-feature-fullstack/SKILL.md
test -f /tmp/arlox-demo/backend/.cursor/rules/karpathy.mdc
test -f /tmp/arlox-demo/web/.cursor/rules/tailwind.mdc
test -f /tmp/arlox-demo/app/.cursor/skills/add-feature-mobile/learned/README.md
```

## 8. NO_COLOR

```bash
NO_COLOR=1 arlox create arlox-nocolor --backend
```

Expect: same layout, no ANSI escapes.

## 9–11. Skills (manual / agent)

- Add a feature via backend `add-feature-backend` skill → `learned/<feature>.md` appears
- Full-stack feature → sequential backend → web → app
- `arlox skills update` leaves `learned/` untouched
