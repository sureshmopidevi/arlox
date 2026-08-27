---
name: add-feature
description: >-
  Scaffolds a Flutter vertical feature slice (data/domain/presentation),
  wires GoRouter, verifies with flutter analyze, and records learnings.
  Use when adding a mobile/app feature, new screen, or Flutter resource
  that consumes the backend API.
---

# Add Feature (App)

Scaffold a vertical feature slice under `lib/features/<name>/`.

## When to use

- New Flutter screen / feature
- Mobile UI that calls the backend API
- "Add X to the app"

**Do not use** alone for multi-stack work — prefer workspace `add-feature-fullstack`.

## Before coding

1. Read `.cursor/rules/architecture.mdc`.
2. Read `.cursor/skills/add-feature/learned/README.md`.
3. If backend exists, read:
   - `../backend/.cursor/skills/add-feature/learned/<feature>.md`
4. Confirm feature name if ambiguous.
5. Success criteria: route opens, API wired, `flutter analyze` clean.

## Steps

### 1. Scaffold

```text
lib/features/<name>/
  data/
    datasources/<name>_remote_datasource.dart
    models/<name>_dto.dart
    repositories/<name>_repository_impl.dart
  domain/
    entities/<name>.dart
    repositories/<name>_repository.dart
    usecases/get_<name>.dart    # or create_/update_/delete_
  presentation/
    screens/<name>_screen.dart
    providers/<name>_providers.dart
```

### 2. Datasource

- Inject `DioClient` via constructor (providers wire it — do not read providers inside datasource).
- Map `DioException` with `DioClient.mapError` → `ApiException`.
- Paths must match the backend contract.

### 3. Repository

- Domain interface returns entities (not DTOs).
- Impl maps DTO → entity.

### 4. Use case

One public `call` method per use case class. No multi-purpose god use cases.

### 5. Providers

Feature-scoped Riverpod providers in `presentation/providers/`. Wire:

`dioClient` → datasource → repository → use case → presentation state.

### 6. Screen

- `ConsumerWidget` / `ConsumerStatefulWidget` only.
- Handle `AsyncValue`: loading / error (`EmptyStateView` if present) / data.
- Material widgets only — no custom design-token systems unless the project already has one.
- Use `AppConstants.appName` / existing constants — do not hardcode product titles unnecessarily.

### 7. Route

Add a `GoRoute` in `lib/core/router/app_router.dart`.

### 8. Verify

```bash
flutter analyze
```

Fix all issues. Run `dart run custom_lint` if the project uses it.

### 9. Learn

Append to `.cursor/skills/add-feature/learned/README.md`:

```markdown
## YYYY-MM-DD — add <name>
**What**: scaffolded <name> (datasource, repo, use case, screen, route).
**API**: <endpoints consumed>.
**Gotchas**: <notes>.
```

## Constraints

- Features must not import other features — shared code goes in `core/`
- No auth route guards unless requested
- Do not invent backend endpoints
- Surgical diffs only

## Anti-patterns

- `StatefulWidget` + manual plumbing instead of Riverpod
- Calling Dio directly from widgets
- Skipping GoRouter registration
- Parallel "also fix backend/web" edits from this skill

## Done checklist

- [ ] Feature slice created
- [ ] Route registered
- [ ] API matches backend contract
- [ ] `flutter analyze` clean
- [ ] Learning entry recorded
