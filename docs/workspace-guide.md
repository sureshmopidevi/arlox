# Workspace and variants

## Variant catalog and prerequisites

All variants require Git. Install only the toolchains needed by the variants
you generate.

- **Backend**
  - `go` — Go 1.22+
  - `python` — Python 3.11+ and pip
  - `node-express`, `node-fastify` — Node.js and npm
  - `java` — Java and Maven
- **Web**
  - `react-vite`, `nextjs`, `vue`, `svelte`, `angular`, `nuxt` — Node.js and npm
- **App**
  - `flutter` — Flutter SDK
  - `react-pwa` — Node.js and npm
  - `ios` — macOS with Swift and Xcode command-line tools
  - `android` — JDK, Android SDK, and Gradle

## Generated structure

Stack directory names do not depend on the selected technology:

```text
project/
  project.code-workspace
  Makefile
  README.md
  .cursor/
    rules/fullstack-order.mdc
    skills/add-feature-fullstack/SKILL.md
  backend/                  # present when selected
    .origin-manifest.json   # stack and concrete variant
    .cursor/
  web/                      # present when selected
    .origin-manifest.json
    .cursor/
  app/                      # present when selected
    .origin-manifest.json
    .cursor/
```

Each stack is initialized as its own Git repository. `arlox add` recognizes
existing variants through their marker files and manifest.

## Verify a generated variant

- Backend:
  - `go`: `go test ./...`
  - `python`: `pytest`
  - `node-express`: `npm test`
  - `node-fastify`: `npm test`
  - `java`: `mvn test`
- Web:
  - `react-vite`: `npm run build`
  - `nextjs`: `npm test`
  - `vue`: `npm test`
  - `svelte`: `npm test`
  - `angular`: `ng test --watch=false`
  - `nuxt`: `npm test`
- App:
  - `flutter`: `flutter test`
  - `react-pwa`: `npm test`
  - `ios`: `xcodebuild test`
  - `android`: `gradle test`

The root workspace exposes stable orchestration targets:
`make backend.test`, `make web.build`, and `make app.verify`.

## Cursor rules and skills

Every variant includes universal `karpathy.mdc` and learning guidance plus
technology-specific architecture, tooling, and workflow instructions.

- Go uses Make; Python uses uv-oriented workflows.
- Node, React PWA, and most web variants use npm; React/Vite, Next.js, Vue,
  Svelte, and Nuxt include framework/Tailwind guidance.
- Java uses Maven; Angular has Angular CLI rules.
- iOS uses Xcode/Swift rules; Android uses Gradle rules.
- Flutter provides Flutter-specific architecture and analysis workflows.

Each stack also includes `dev-workflow`, `learn`, `apply-pending`, and an
appropriate feature skill: `add-feature-backend`, `add-feature-web`,
`add-feature-mobile`, `add-feature-pwa`, `add-feature-ios`, or
`add-feature-android`. Learned entries remain local when skills are updated.
