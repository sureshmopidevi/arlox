# Skill: add-feature

Scaffold a new vertical feature slice following the project architecture.

## Trigger

Invoke when asked to **"add feature \<name\>"** or **"create feature \<name\>"**.

---

## Steps

### 0. Prepare

1. Read `.cursor/rules/architecture.mdc`.
2. Read `.cursor/skills/add-feature/learned/README.md` for prior decisions relevant to this feature.
3. Confirm the feature name with the user if ambiguous.

### 1. Scaffold files

Create under `lib/features/<name>/`:

```
data/
  datasources/<name>_remote_datasource.dart
  models/<name>_dto.dart
  repositories/<name>_repository_impl.dart
domain/
  entities/<name>.dart
  repositories/<name>_repository.dart
  usecases/get_<name>.dart          # or appropriate verb (create_, update_, delete_)
presentation/
  screens/<name>_screen.dart        # Scaffold + AppBar + body
  providers/<name>_providers.dart   # manual Riverpod Provider / NotifierProvider
```

### 2. Datasource

- Inject `DioClient` via constructor parameter (do not access `dioClientProvider` directly inside the datasource — let the provider wire it).
- Catch `DioException`, convert with `DioClient.mapError`, rethrow as `ApiException`.

### 3. Repository

- Domain interface in `domain/repositories/<name>_repository.dart` — return types use domain entities, not DTOs.
- Implementation in `data/repositories/<name>_repository_impl.dart`.

### 4. Use case

```dart
class Get<Name> {
  const Get<Name>(this._repository);
  final <Name>Repository _repository;
  Future<List<<Name>>> call() => _repository.get<Name>s();
}
```

### 5. Provider

Feature-scoped provider in `presentation/providers/<name>_provider.dart`:

```dart
@riverpod
Future<List<<Name>>> <name>(Ref ref) async {
  final client = ref.watch(dioClientProvider);
  final datasource = <Name>RemoteDatasource(client);
  final repo = <Name>RepositoryImpl(datasource);
  return Get<Name>(repo).call();
}
```

### 6. Screen

- Use `ConsumerWidget` or `ConsumerStatefulWidget`.
- Handle `AsyncValue` states: loading → `CircularProgressIndicator`, error → `EmptyStateView`, data → your list/content.
- Material widgets only. No custom tokens.

### 7. Register route

Add a `GoRoute` entry to `lib/core/router/app_router.dart`.

### 8. Analyze

```bash
flutter analyze
dart run custom_lint
```

Fix **all** warnings before finishing.

### 9. Learn

Append an entry to `.cursor/skills/add-feature/learned/README.md`:

```markdown
## YYYY-MM-DD — add <name> feature
**What**: scaffolded <name> feature with remote datasource and GoRouter route.
**Why**: <rationale for any non-obvious decisions>.
**Gotchas**: <anything a future agent should know>.
```

---

## Constraints

- Material widgets only — no custom theme files, no Verdant tokens.
- `ConsumerWidget` or `ConsumerStatefulWidget` for all screens — never `StatefulWidget` + manual state.
- No auth guard on routes unless explicitly requested.
- Keep use cases single-responsibility — one public `call` method per class.
- Features must not import each other. Use `core/` for anything shared.
