# Workspace guide

How arlox organizes multi-stack projects and how teams work across backend, web, and app.

---

## Workspace structure

```text
myapp/
  myapp.code-workspace     # Open in Cursor, VS Code, or Antigravity
  Makefile                 # make dev, make test, make status, …
  docker-compose.yml       # Postgres 16 (unique host port per project name)
  contracts/               # Shared API docs for cross-stack features
  README.md
  .cursor/                 # Fullstack skill + workspace rules

  backend/                 # Separate git repo — Go + Gin
  web/                     # Separate git repo — React + Vite
  app/                     # Separate git repo — Flutter
```

Each stack has its own `.git`, `.cursor/rules`, `.cursor/skills`, and `AGENTS.md`.

---

## Partial workspaces

Create only what you need:

```bash
arlox create myapp --backend --web
arlox add --app
```

`make status` and the root Makefile skip missing stacks automatically.

---

## API contracts (`contracts/`)

Cross-stack API documentation lives at the **workspace root**, not inside backend skills.

```text
contracts/
  README.md      # Envelope rules, auth header, how to add features
  auth.md        # Seed doc for login / me endpoints
  issues.md      # Example feature contract (you add these)
```

### Full-stack feature order

Use `.cursor/skills/add-feature-fullstack`:

1. **Contract** — write `contracts/<feature>.md`
2. **Backend** — implement API, `make test`
3. **Web** — types, UI, route (respect `web/.arlox/design-system.json`)
4. **App** — Flutter feature slice, `flutter analyze`

Never parallelize the same feature across stacks.

---

## Brownfield: `arlox init`

Adopt arlox orchestration in an **existing** repo that already has stack folders:

```bash
cd ~/Projects/legacy-api   # contains backend/go.mod, etc.
arlox init . --name legacy-api
```

**What init does**

- Detects stacks via marker files (`go.mod`, `package.json`, `pubspec.yaml`)
- Writes workspace files: Makefile, `.code-workspace`, `contracts/`, root `.cursor/`
- Installs cursor skills if missing

**What init never does**

- Overwrite source inside `backend/`, `web/`, or `app/`

Flags:

| Flag | Default | Purpose |
|------|---------|---------|
| `--name` | infer from dir or `.code-workspace` | Project name |
| `--force` | false | Overwrite locally modified root skills |
| `--deps` | false | Run `go mod tidy` / `npm install` / `flutter pub get` |

---

## Database (Postgres)

Host port is **deterministic from project name** (range 5433–7432) so multiple workspaces can run Postgres concurrently.

```bash
make db-up
make backend.setup
```

Port appears in `docker-compose.yml` and `backend/configs/local/app.env`.

---

## Day-to-day commands

From project root:

| Command | Purpose |
|---------|---------|
| `make dev` | Backend + web dev (health poll then Vite) |
| `make test` | All present stacks |
| `make status` | Which stacks exist |
| `make doctor` | Toolchain check per stack |
| `make backend.run` | API only |
| `make web.dev` | Vite dev server |
| `make app.run` | Flutter |

---

## Maintaining the workspace

| Situation | Command |
|-----------|---------|
| Missing Makefile, `.env`, cursor files | `arlox repair` |
| Refresh skills after arlox upgrade | `arlox skills update` |
| Overwrite locally edited skills | `arlox repair --force` or `arlox skills update --force` |
| See scaffold drift vs template | `arlox repair` (drift section) |

Each stack stores `.origin-manifest.json` from create time for drift detection.

---

## AI agent guardrails

| Layer | Location |
|-------|----------|
| Architecture | `<stack>/.cursor/rules/architecture.mdc` |
| Styling (web) | `web/.cursor/rules/design-system.mdc` |
| Feature skills | `add-feature-backend`, `add-feature-web`, `add-feature-mobile` |
| Full stack | `.cursor/skills/add-feature-fullstack` at workspace root |

Web agents must read [Web design systems](web-design-systems.md) before building UI.
