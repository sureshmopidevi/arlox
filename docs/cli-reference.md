# CLI reference

## `arlox create [name]`

Creates a workspace. Without a name, a TTY prompt asks for one.

```bash
arlox create shop --backend go --web react-vite --app flutter
arlox create api --backend python
arlox create portal --web nextjs --open
```

Stack flags accept an optional variant value:

- `--backend [go|python|node-express|node-fastify|java]`
- `--web [react-vite|nextjs|vue|svelte|angular|nuxt]`
- `--app [flutter|react-pwa|ios|android]`

A bare flag selects the legacy default, preserving existing scripts:

- `--backend` → `go`
- `--web` → `react-vite`
- `--app` → `flutter`

Other create flags:

- `--module <path>` — Go module path; defaults to
  `github.com/example/<name>-backend`
- `--org <identifier>` — Flutter organization; defaults to `com.example`
- `--out <directory>` — workspace parent; defaults to the current directory
- `--open` — open the generated `.code-workspace`

### Interactive sequence

When no stack flags are supplied, prompts run in this order:

1. Project name, if omitted.
2. Select `backend`, `web`, and/or `app`.
3. Select one backend variant, if backend was selected.
4. Select one web variant, if web was selected.
5. Select `Flutter`, `React PWA`, or `Native`, if app was selected.
6. If Native was selected, select `iOS` or `Android`.

The Native choice is only a prompt branch; the stored variant is `ios` or
`android`.

## `arlox add`

Adds missing stacks to the detected workspace. It accepts the same stack and
variant flags as `create`.

```bash
arlox add --backend node-fastify
arlox add --web nuxt --app android
arlox add --app                 # Flutter default
```

Run it from the workspace root or any generated stack directory. Existing
stacks are skipped.

## Maintenance commands

```bash
arlox skills update             # preserve learned/ and local edits
arlox skills update --force
arlox repair
arlox repair --force
arlox doctor
arlox upgrade [--source <path>] [--no-pull]
arlox uninstall [--yes]
arlox version
```
