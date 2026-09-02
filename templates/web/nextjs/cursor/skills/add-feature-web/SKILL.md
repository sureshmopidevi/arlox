---
name: add-feature-web
description: Add a tested feature to this Next.js App Router application.
---

# Add Feature (Web)

1. Read the architecture, Tailwind, and learned guidance.
2. Define the route, data boundary, user states, and testable outcome.
3. Add the route in `app/`; keep it server-rendered unless interactivity requires a Client Component.
4. Put shared UI in `components/` and validate external data at its boundary.
5. Cover the primary behavior with Vitest and Testing Library.
6. Run `npm test` and `npm run build`.
7. Use the learn skill to record non-obvious decisions.

Avoid speculative state libraries, API wrappers, and component layers.
