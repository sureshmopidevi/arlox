---
name: add-feature-mobile
description: >-
  Ship a Flutter feature slice with GoRouter route, Riverpod wiring, and API
  calls matching contracts/<feature>.md; flutter analyze clean. Use when adding
  a mobile screen, app feature, or Flutter resource that consumes the backend API.
---

# Add Feature (Mobile)

## User intent

When done, the user can navigate to a new screen, see data from the backend, and the feature follows the project's clean-architecture slice under `lib/features/<name>/`.

## Read first

1. **`.arlox/project.json`** (workspace root) — `stacks.app.package` (snake_case pubspec name), `stacks.app.org`.
2. **`contracts/<feature>.md`** (workspace root, preferred) or `../backend/.cursor/skills/add-feature-backend/learned/<feature>.md` as fallback.
3. **`.cursor/rules/architecture.mdc`**.
4. **`.cursor/skills/add-feature-mobile/learned/README.md`**.

**Do not use** alone for multi-stack work — prefer workspace `add-feature-fullstack`.

## Success criteria

| Check | Must be true |
|-------|----------------|
| Slice | `data/` → `domain/` → `presentation/` under `lib/features/<name>/` |
| Route | `GoRoute` registered; protected routes use `authRedirect` |
| Contract | API paths match `contracts/<feature>.md` |
| Analyze | `flutter analyze` reports no issues |
| Learn | Entry appended via `learn/SKILL.md` |

Confirm feature name if ambiguous before scaffolding.

## Steps

1. **Scaffold** feature slice — datasource, DTO, repository, entity, use case, providers, screen.
2. **Datasource** — inject `DioClient`; map errors with `DioClient.mapError`.
3. **Repository** — domain interface returns entities; impl maps DTO → entity.
4. **Providers** — feature-scoped Riverpod: dio → datasource → repository → use case → UI state.
5. **Screen** — `ConsumerWidget`; handle `AsyncValue` loading/error/data.
6. **Route** — add `GoRoute` in `lib/core/router/app_router.dart`.
7. **Verify** — `flutter analyze`.
8. **Learn** — run **`.cursor/skills/learn/SKILL.md`** (mandatory).

## Verify

```bash
flutter analyze
```

Run `dart run custom_lint` if the project uses it.

## Learn

Run **`.cursor/skills/learn/SKILL.md`** after verify passes. Append to **`learned/README.md`**.

## Anti-patterns

- `StatefulWidget` + manual plumbing instead of Riverpod
- Calling Dio directly from widgets
- Skipping GoRouter registration
- Inventing backend endpoints
- Cross-feature imports (shared code belongs in `core/`)
