# arlox

> **Status: Active Development 🚧** — CLI behavior and template architectures may evolve.

`arlox` scaffolds backend, web, and app stacks into one workspace. Each selected
stack can use a different technology while retaining stable `backend/`, `web/`,
and `app/` directory names, workspace automation, and technology-specific
Cursor rules and skills.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)](go.mod)

## Supported variants

- **Backend:** Go + Gin, Python + FastAPI, Node.js + Express, Node.js + Fastify,
  Java + Spring Boot
- **Web:** React + Vite, Next.js, Vue + Vite, Svelte + Vite, Angular, Nuxt
- **App:** Flutter, React PWA, native iOS (SwiftUI), native Android
  (Kotlin + Jetpack Compose)

See [Workspace and variants](docs/workspace-guide.md) for prerequisites,
generated structure, test commands, and Cursor assets.

## Installation

```bash
curl -fsSL https://raw.githubusercontent.com/sureshmopidevi/arlox/main/install.sh | bash
```

Or install with Go:

```bash
go install github.com/sureshmopidevi/arlox/cmd/arlox@latest
```

For local development:

```bash
git clone https://github.com/sureshmopidevi/arlox.git ~/arlox
cd ~/arlox
./install.sh
```

## Create a workspace

Interactive creation asks for the project name, stacks, and then one variant
for each selected stack. App selection is `Flutter | React PWA | Native`; choosing
Native opens a second `iOS | Android` prompt.

```bash
arlox create myapp
```

For scripts, pass variants explicitly:

```bash
arlox create myapp \
  --backend python \
  --web nextjs \
  --app react-pwa
```

Bare stack flags remain backward compatible:

```bash
arlox create myapp --backend --web --app
# equivalent variants: backend=go, web=react-vite, app=flutter
```

Variants can also be added later:

```bash
cd myapp
arlox add --web vue
```

See the [CLI reference](docs/cli-reference.md) for every value and flag.

## Generated workspace

```text
myapp/
  myapp.code-workspace
  Makefile
  .cursor/                # cross-stack order and full-stack skills
  backend/                # selected backend variant, own git repository
  web/                    # selected web variant, own git repository
  app/                    # selected app variant, own git repository
```

Each generated stack records its concrete variant in
`.origin-manifest.json` and includes rules for its architecture and tooling,
plus `dev-workflow`, `add-feature-*`, `learn`, and `apply-pending` skills.

```bash
arlox skills update       # preserve learned/ and locally edited files
```

## Other commands

```bash
arlox version
arlox doctor
arlox repair
arlox upgrade
arlox uninstall
make verify
```

See [scripts/verify.md](scripts/verify.md) for verification commands.

## License

[MIT](LICENSE)
