---
name: dev-workflow
description: Build, run, and test this SwiftUI package.
---

# Development Workflow

```bash
swift package resolve
swift test
swift run
```

Use `swift run` for the package's macOS SwiftUI host. Open `Package.swift` in
Xcode for iOS simulator work. Signing, capabilities, asset catalogs, previews,
and archives require an Xcode iOS app target and cannot be validated by
`swift test`.
