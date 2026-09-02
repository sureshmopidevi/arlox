---
name: add-feature-web
description: Add a tested feature to this Angular application.
---

# Add Feature (Web)

1. Read the architecture, Angular, and learned guidance.
2. Define user states, route/data boundaries, and a testable outcome.
3. Add standalone feature components and typed services under `src/app/features/<name>/`.
4. Use signals for local state and keep templates declarative.
5. Cover primary behavior with Vitest and TestBed.
6. Run `npm test` and `npm run build`.
7. Use the learn skill for non-obvious decisions.

Avoid speculative shared modules, generic services, and state libraries.
