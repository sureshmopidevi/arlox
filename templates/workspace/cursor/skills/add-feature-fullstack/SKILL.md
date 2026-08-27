---
name: add-feature-fullstack
description: >-
  Implements a feature across backend, web, and Flutter app in strict order
  (backend → web → app) with API contracts and verification between phases.
  Use when the user asks to add a full-stack feature, cross-stack feature,
  feature for backend and web, or anything spanning multiple vibeit stacks.
---

# Add Feature Across Stacks

Run **one stack at a time**. Never parallelize backend/web/app for the same feature.

## When to use

- User wants a feature in more than one of: API, admin web, Flutter app
- Phrases like: "full stack", "backend and web", "across stacks", "API + UI"

**Do not use** for single-stack work — use `add-feature-backend`, `add-feature-web`, or `add-feature-mobile` instead.

## Success criteria (all that apply)

| Stack | Must prove |
|-------|------------|
| Backend | `make test` green + `learned/<feature>.md` with API contract |
| Web | Types + API client + UI + route/nav + `npm run build` green |
| App | Feature slice + route + `flutter analyze` green |

If `web/` exists, **Phase 2 (web) is mandatory** — do not stop after backend.

## 0. Detect stacks

```text
backend/  → Phase 1
web/      → Phase 2   ← required when this folder exists
app/      → Phase 3
```

State the plan to the user before coding:

```text
1. backend → verify: make test + contract doc
2. web     → verify: npm run build
3. app     → verify: flutter analyze
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

**Verify:** `cd backend && make test`  
Do **not** start web until this passes and the contract file exists.

## 2. Phase — Web (if `web/` exists)

Work **only** under `web/`. Follow `web/.cursor/skills/add-feature-web/SKILL.md`.

**Must consume the backend contract** from  
`backend/.cursor/skills/add-feature-backend/learned/<feature>.md` (or ask if missing).

**Required deliverables:**

1. `src/features/<name>/types.ts` — match API shapes
2. `src/features/<name>/api/<name>Api.ts` — calls via `apiClient`
3. Query/mutation hooks
4. Feature UI component + thin page
5. Route in `src/app/router.tsx`
6. Nav item in `src/layouts/AdminSidebar.tsx` (unless user said otherwise)

**Verify:** `cd web && npm run build`  
Do **not** start app until this passes.

## 3. Phase — App (if `app/` exists)

Work **only** under `app/`. Follow `app/.cursor/skills/add-feature-mobile/SKILL.md`.

Consume the same backend contract. Do not invent endpoints.

**Verify:** `cd app && flutter analyze`

## Subagent pattern (preferred)

For each phase, use a Task/subagent with:

```text
Scope: only <stack>/ 
Read: backend → add-feature-backend | web → add-feature-web | app → add-feature-mobile
Contract: backend/.cursor/skills/add-feature-backend/learned/<feature>.md  (web/mobile)
Verify: <stack verify command>
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
