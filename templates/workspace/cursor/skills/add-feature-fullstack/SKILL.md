---
name: add-feature-fullstack
description: >-
  Implements a feature across backend, web, and app variants in strict order
  (backend → web → app) with API contracts and verification between phases.
  Use when the user asks to add a full-stack feature, cross-stack feature,
  feature for backend and web, or anything spanning multiple arlox stacks.
---

# Add Feature Across Stacks

Run **one stack at a time**. Never parallelize backend/web/app for the same feature.

## When to use

- User wants a feature in more than one of: API, admin web, Flutter app
- Phrases like: "full stack", "backend and web", "across stacks", "API + UI"

**Do not use** for single-stack work — use the add-feature skill inside that stack instead.

## Success criteria (all that apply)

| Stack | Must prove |
|-------|------------|
| Backend | `make backend.test` green + `learned/<feature>.md` with API contract |
| Web | Variant-appropriate feature + `make web.build` green |
| App | Variant-appropriate feature + `make app.verify` green |

If `web/` exists, **Phase 2 (web) is mandatory** — do not stop after backend.

## 0. Detect stacks

```text
backend/  → Phase 1
web/      → Phase 2   ← required when this folder exists
app/      → Phase 3
```

State the plan to the user before coding:

```text
1. backend → verify: make backend.test + contract doc
2. web     → verify: make web.build
3. app     → verify: make app.verify
```

Skip only missing folders. Never skip an existing stack.

## 1. Phase — Backend (if `backend/` exists)

Work **only** under `backend/`. Follow `backend/.cursor/skills/add-feature-backend/SKILL.md`.

**Required output before leaving this phase:**

Create `backend/.cursor/skills/add-feature-backend/learned/<feature>.md` with this contract:

```markdown
# <feature> — API contract

## Endpoints
| Method | Path | Auth | Request | Response |
|--------|------|------|---------|----------|
| POST | /api/v1/<resource> | yes | `{ ... }` | `{ data: { ... } }` |

## Errors
| Status | When |
|--------|------|
| 400 | validation |
| 401 | missing/invalid token |
| 404 | not found |

## Notes for web/app
- Envelope: `{ data: T }` unwrapped by clients
- Auth header: `Authorization: Bearer <token>`
```

Read `backend/.cursor/skills/dev-workflow/SKILL.md` for variant-specific commands.

**Verify:** `make backend.test` from the workspace root.
Do **not** start web until this passes and the contract file exists.

## 2. Phase — Web (if `web/` exists)

Work **only** under `web/`. Follow `web/.cursor/skills/add-feature-web/SKILL.md`.

**Must consume the backend contract** from  
`backend/.cursor/skills/add-feature-backend/learned/<feature>.md` (or ask if missing).

Follow the installed web variant's architecture and add-feature skill. Required
file layout and routing differ between React, Next.js, Vue, Svelte, Angular,
and Nuxt; do not impose one framework's conventions on another.

**Verify:** `make web.build` from the workspace root.
Do **not** start app until this passes.

## 3. Phase — App (if `app/` exists)

Work **only** under `app/`. Follow the add-feature skill present under
`app/.cursor/skills/`.

Consume the same backend contract. Do not invent endpoints.

Read `app/.cursor/skills/dev-workflow/SKILL.md` for variant-specific commands.

**Verify:** `make app.verify` from the workspace root.

## Subagent pattern (preferred)

For each phase, use a Task/subagent with:

```text
Scope: only <stack>/ 
Read: the stack's architecture rule, dev-workflow, and add-feature skill
Contract: backend/.cursor/skills/add-feature-backend/learned/<feature>.md  (web/mobile)
Verify: make backend.test | make web.build | make app.verify
Return: files changed + verify output summary
```

Wait for completion + green verify before launching the next phase.

## Anti-patterns

- Building backend only and calling the feature "done" when `web/` exists
- Starting web/app before backend contract is written
- Parallel subagents on the same feature
- Editing multiple stacks in one agent turn
- Drive-by refactors outside the feature
- Hardcoding API URLs instead of using stack clients/config

## Example

User: "Add issues list with create form (full stack)"

```text
1. Backend: modules/issues + routes + tests → learned/issues.md
2. Web: features/issues + IssuesPage + /issues nav → npm run build
3. App: lib/features/issues + GoRoute → flutter analyze
```

## Done checklist

- [ ] Every existing stack was implemented (not skipped)
- [ ] Backend contract doc exists (if backend ran)
- [ ] Web UI is reachable via route + nav (if web ran)
- [ ] Each phase verify command passed
- [ ] No cross-stack drive-by edits
