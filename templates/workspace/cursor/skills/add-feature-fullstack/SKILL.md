---
name: add-feature-fullstack
description: >-
  Ship a feature across backend, web, and Flutter in order (contract → backend
  → web → app) with green verify between phases. Use when the user asks for a
  full-stack, cross-stack, or multi-stack feature spanning API and UI.
---

# Add Feature Across Stacks

## User intent

When done, every **existing** stack implements the same feature: API documented in `contracts/<feature>.md`, backend tests pass, web route works, and Flutter screen is wired — with no parallel cross-stack edits.

## Read first

1. **`.arlox/project.json`** — naming (`kebab`, `snake`), stack presence.
2. **`contracts/README.md`** — contract format.
3. Per-phase stack skills: `add-feature-backend`, `add-feature-web`, `add-feature-mobile`.

Run **one stack at a time**. Never parallelize backend/web/app for the same feature.

## Success criteria

| Stack | Must prove |
|-------|------------|
| Contract | `contracts/<feature>.md` exists before backend coding |
| Backend | `make test` green |
| Web | Route + nav + `npm run build` green |
| App | Feature slice + `flutter analyze` green |
| Learn | Each stack that ran must append via its `learn/SKILL.md` |

If `web/` exists, **Phase 2 (web) is mandatory** — do not stop after backend.

## Steps

### 0. Detect stacks

```text
backend/  → Phase 1
web/      → Phase 2   ← required when this folder exists
app/      → Phase 3
```

State the plan before coding:

```text
1. contracts → write contracts/<feature>.md
2. backend    → verify: make test
3. web        → verify: npm run build
4. app        → verify: flutter analyze
5. learn      → each stack runs learn/SKILL.md
```

Skip only **missing** folders. Never skip an existing stack.

### Phase 0 — Contract (always first)

Create **`contracts/<feature>.md`** using [`contracts/README.md`](../../contracts/README.md). Do **not** start backend until it exists.

### Phase 1 — Backend (if `backend/` exists)

Work **only** under `backend/`. Follow `backend/.cursor/skills/add-feature-backend/SKILL.md`.  
**Verify:** `cd backend && make test`

### Phase 2 — Web (if `web/` exists)

Work **only** under `web/`. Follow `web/.cursor/skills/add-feature-web/SKILL.md`.  
Read `web/.arlox/design-system.json`.  
**Verify:** `cd web && npm run build`

### Phase 3 — App (if `app/` exists)

Work **only** under `app/`. Follow `app/.cursor/skills/add-feature-mobile/SKILL.md`.  
**Verify:** `cd app && flutter analyze`

### Phase 4 — Learn (mandatory)

Each stack that ran must run its **`learn/SKILL.md`** and append to `learned/README.md`.  
If a pattern should become a reusable rule or skill step, run that stack's **`apply-pending/SKILL.md`** or workspace **`reflect-and-improve/SKILL.md`**.

## Verify

Run each stack's verify command in order; do not start the next phase until the current one is green.

## Learn

Mandatory after each phase. See Phase 4 above.

## Anti-patterns

- Building backend only when `web/` exists
- Starting backend before `contracts/<feature>.md`
- Starting web/app before prior phase verify passes
- Parallel subagents on the same feature
- API contracts only in `learned/` instead of `contracts/`
- Hardcoding URLs instead of stack clients/config
