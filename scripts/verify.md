# arlox verification

## Automated

From the repository root:

```bash
make verify
```

Or:

```bash
source scripts/env.sh && ensure_arlox_path && ensure_local_tool_paths
./scripts/verify.sh
```

The smoke test skips variants whose required toolchain is unavailable.

## Manual setup

Use a throwaway directory outside this repository:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
cd /tmp
rm -rf arlox-demo-*
```

For each generated stack, confirm:

1. Its marker file exists.
2. `.origin-manifest.json` contains the requested `stack` and `variant`.
3. Technology-specific `.cursor/rules` and `.cursor/skills` exist.
4. The variant test command succeeds.

## Backward-compatible defaults

```bash
arlox create arlox-demo-defaults --backend --web --app
```

Expect `go`, `react-vite`, and `flutter` manifests in `backend/`, `web/`, and
`app/`. Re-running the command must skip existing stacks.

## Explicit variants

Generate one variant at a time with:

```bash
arlox create arlox-demo-<variant> --backend <variant>
arlox create arlox-demo-<variant> --web <variant>
arlox create arlox-demo-<variant> --app <variant>
```

Run the matching check when its prerequisites are installed:

- Backend:
  - `go`: `go test ./...`
  - `python`: `pytest`
  - `node-express`, `node-fastify`: `npm test`
  - `java`: `mvn test`
- Web:
  - `react-vite`: `npm run build`
  - `nextjs`, `vue`, `svelte`, `nuxt`: `npm test`
  - `angular`: `ng test --watch=false`
- App:
  - `flutter`: `flutter test`
  - `react-pwa`: `npm test`
  - `ios`: `xcodebuild test` (macOS/Xcode only)
  - `android`: `gradle test` (requires JDK, Gradle, and Android SDK)

Toolchains are listed in [the workspace guide](../docs/workspace-guide.md).

## Interactive flow

Run `arlox create` in a TTY and verify prompts are sequential:

1. Project name.
2. Stack multi-select.
3. One backend choice, if selected.
4. One web choice, if selected.
5. App choice: Flutter, React PWA, or Native.
6. Native platform: iOS or Android, only after choosing Native.

## Add and guardrails

```bash
arlox create arlox-demo-add --backend python
cd arlox-demo-add
arlox add --web vue
```

Expect both folders in the workspace file and no backend replacement. Confirm
the root full-stack skill and each variant's architecture, tooling,
`dev-workflow`, `add-feature-*`, `learn`, and `apply-pending` assets.

Finally, add a learned entry and run `arlox skills update`; the learned file
must remain unchanged.

## Output guardrail

```bash
NO_COLOR=1 arlox create arlox-demo-nocolor --backend
```

Expect the normal Go backend layout with no ANSI escapes.
