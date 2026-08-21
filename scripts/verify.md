# vibeit v1 verification

## Automated (recommended)

From the repo root — exports `$(go env GOPATH)/bin` if needed:

```bash
make verify
```

Or run the script directly:

```bash
source scripts/env.sh && ensure_vibeit_path && ensure_local_tool_paths
./scripts/verify.sh
```

## Manual

Run from a throwaway directory (not inside the vibeit repo).

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
# or: eval $(make -C ~/vibeit path)
cd /tmp
rm -rf vibeit-demo vibeit-demo2
```

## 1. Create all stacks

```bash
vibeit create vibeit-demo --backend --web --app
```

Expect: `backend/`, `web/`, `app/`, `vibeit-demo.code-workspace`, `.cursor/skills/add-feature-fullstack`.

## 2. No duplicates

```bash
vibeit create vibeit-demo --backend --web --app
```

Expect: all stacks **skipped (already exists)**.

## 3. Add after partial create

```bash
vibeit create vibeit-demo2 --backend
cd vibeit-demo2
vibeit add --web
```

Expect: web generated; backend skipped; workspace folders include both.

## 4. Backend health (needs local Postgres)

```bash
cd /tmp/vibeit-demo/backend
make setup
make run   # in another terminal: curl -s localhost:8080/health/live
```

## 5. Web build

```bash
cd /tmp/vibeit-demo/web
npm install && npm run build
```

## 6. Flutter analyze

```bash
cd /tmp/vibeit-demo/app
flutter pub get && flutter analyze
```

## 7. Guardrails present

```bash
test -f /tmp/vibeit-demo/.cursor/skills/add-feature-fullstack/SKILL.md
test -f /tmp/vibeit-demo/backend/.cursor/rules/karpathy.mdc
test -f /tmp/vibeit-demo/web/.cursor/rules/tailwind.mdc
test -f /tmp/vibeit-demo/app/.cursor/skills/add-feature/learned/README.md
```

## 8. NO_COLOR

```bash
NO_COLOR=1 vibeit create vibeit-nocolor --backend
```

Expect: same layout, no ANSI escapes.

## 9–11. Skills (manual / agent)

- Add a feature via backend `add-feature` skill → `learned/<feature>.md` appears
- Full-stack feature → sequential backend → web → app
- `vibeit skills update` leaves `learned/` untouched
