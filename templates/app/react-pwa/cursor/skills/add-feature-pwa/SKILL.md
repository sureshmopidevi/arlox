---
name: add-feature-pwa
description: Add a tested, offline-aware vertical feature to this React PWA.
---

# Add Feature (PWA)

1. Read `.cursor/rules/architecture.mdc` and `learned/README.md`.
2. Define the user-visible behavior and whether it must work offline.
3. Add the minimum slice under `src/features/<name>/`.
4. Keep API functions separate from components and do not invent endpoints.
5. Add a React Testing Library test for the primary behavior.
6. If caching changes, document freshness, invalidation, and privacy effects.
7. Run `npm test` and `npm run build`.
8. Append the decision and offline impact to `learned/README.md`.
