---
name: add-feature-android
description: Add a tested MVVM/repository feature ready for Compose and Hilt.
---

# Add Feature (Android)

1. Read the architecture rule and `learned/README.md`.
2. Define immutable UI state, actions, repository needs, and success criteria.
3. Add framework-free repository contracts and ViewModel logic under
   `app/src/main/java/`.
4. Render immutable state with Compose and bind implementations with Hilt
   constructor injection.
5. Do not invent API endpoints or leak Android types into domain contracts.
6. Add JUnit coverage for state transitions.
7. Run `gradle test lint assembleDebug`; list device checks still required.
8. Record the decision and tooling impact in `learned/README.md`.
