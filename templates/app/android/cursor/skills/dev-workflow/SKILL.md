---
name: dev-workflow
description: Build, install, lint, and test this Jetpack Compose app.
---

# Development Workflow

```bash
gradle test
gradle assembleDebug
gradle lint
gradle installDebug
```

Run `gradle wrapper` once from a trusted Gradle 9.3+ installation, commit the
generated wrapper, and then use `./gradlew` for every command. `ANDROID_HOME`
must point to an SDK containing API 37. Before declaring work complete, run
`./gradlew test lint assembleDebug`.
