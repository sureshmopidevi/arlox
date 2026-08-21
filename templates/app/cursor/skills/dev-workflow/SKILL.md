# Skill: dev-workflow

Standard development commands for this project.

## Trigger

Invoke when asked to run the app, install packages, analyze, test, or clean the build.

---

## Commands

### Initial setup

```bash
flutter pub get
dart run build_runner build --delete-conflicting-outputs
```

> Run `build_runner` any time you add or modify a `@riverpod`-annotated provider.

### Run

```bash
# Development (uses default API_URL from AppConstants)
flutter run

# With a custom API URL
flutter run --dart-define=API_URL=https://api.example.com
```

### Analyze

```bash
flutter analyze
dart run custom_lint
```

Zero warnings required before committing or declaring a task complete.

### Test

```bash
flutter test
```

### Clean build

```bash
flutter clean \
  && flutter pub get \
  && dart run build_runner build --delete-conflicting-outputs
```

Use this when you see stale generated files or unexplained build errors.

---

## Notes

- **`build_runner` is required** after any change to a file with `@riverpod`, `@JsonSerializable`, or other code-gen annotations.
- Pass `API_URL` via `--dart-define` to override `AppConstants.baseUrl` without modifying source.
- The `custom_lint` step catches Riverpod-specific lint warnings that `flutter analyze` misses.
