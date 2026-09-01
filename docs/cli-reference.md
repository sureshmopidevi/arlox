# CLI reference

Complete command reference for **arlox**. Run `arlox <command> --help` for flag details.

Current version: see [`internal/version/VERSION`](../internal/version/VERSION).

---

## Global

```bash
arlox              # help
arlox version      # print version
arlox -v           # same as version
arlox doctor       # PATH, toolchains, workspace detection
```

---

## `arlox create [name]`

Create a new workspace under the current directory (or `--out`).

```bash
arlox create myapp
arlox create myapp --backend --web --app
arlox create myapp --web --web-ui shadcn --out ~/Projects
arlox create myapp --open
```

| Flag | Description |
|------|-------------|
| `--backend` | Include Go + Gin API |
| `--web` | Include React + Vite admin |
| `--app` | Include Flutter app |
| `--web-ui <id>` | Web design system (see [web-design-systems.md](web-design-systems.md)) |
| `--module` | Go module path (default `github.com/example/<name>-backend`) |
| `--org` | Flutter org (default `com.example`) |
| `--out` | Parent directory for workspace folder |
| `--open` | Launch IDE with `.code-workspace` |

Without stack flags, an interactive multi-select runs. With `--web` and no `--web-ui`, an interactive design-system picker runs (TTY only).

---

## `arlox add`

Add missing stacks to the current workspace.

```bash
cd myapp
arlox add
arlox add --web --web-ui antd
```

Same stack and `--web-ui` flags as create. Must run from workspace root or a stack subdirectory.

---

## `arlox init [path]`

Brownfield workspace orchestration. Default path: current directory.

```bash
arlox init
arlox init . --name myapp
arlox init ~/Projects/legacy --deps
```

| Flag | Default | Description |
|------|---------|-------------|
| `--name` | inferred | Workspace project name |
| `--force` | false | Overwrite modified root skills |
| `--deps` | false | Install stack dependencies |

Requires at least one detected stack (`backend/go.mod`, etc.).

---

## `arlox repair`

Restore missing workspace and stack files; report drift.

```bash
cd myapp
arlox repair
arlox repair --force --deps=false --check-drift
```

| Flag | Default | Description |
|------|---------|-------------|
| `--force` | false | Overwrite locally modified skills/rules |
| `--deps` | true | `go mod tidy`, `npm install`, `flutter pub get` |
| `--check-drift` | true | Compare files to `.origin-manifest.json` |

---

## `arlox skills update`

Refresh `.cursor` templates from embedded arlox version.

```bash
cd myapp
arlox skills update
arlox skills update --force
```

Preserves `learned/` subfolders unless `--force`.

---

## `arlox upgrade`

Rebuild and reinstall arlox from source (default `~/arlox`).

```bash
arlox upgrade
arlox upgrade --no-pull --source ~/arlox
```

---

## `arlox uninstall`

Remove arlox binary from PATH locations.

```bash
arlox uninstall
```

---

## Design system IDs (`--web-ui`)

| ID | Library |
|----|---------|
| `tailwind` | Tailwind utilities only (default in non-interactive mode) |
| `shadcn` | shadcn/ui |
| `antd` | Ant Design |
| `mui` | Material UI |
| `chakra` | Chakra UI |
| `mantine` | Mantine |

Invalid IDs fail with a list of valid values.

---

## Exit codes

Commands exit `1` on failure (validation error, stack generation failure, etc.). `arlox` with no subcommand prints help and exits `0`.
