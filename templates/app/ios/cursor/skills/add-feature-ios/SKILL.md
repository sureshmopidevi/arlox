---
name: add-feature-ios
description: Add a tested SwiftUI MVVM feature with repository boundaries.
---

# Add Feature (iOS)

1. Read the architecture rule and `learned/README.md`.
2. Define view states, user actions, repository needs, and success criteria.
3. Add the smallest SwiftUI view and `@MainActor` view model.
4. Add a repository protocol only at an actual I/O boundary.
5. Inject dependencies; do not use hidden global state.
6. Add XCTest coverage for view-model behavior.
7. Run `swift test`; identify any checks that still require Xcode.
8. Record the decision and Xcode impact in `learned/README.md`.
