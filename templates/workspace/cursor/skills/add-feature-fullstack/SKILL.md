# Add Feature Across Stacks (Full-Stack)

**name:** add-feature-fullstack

**description:** Add a feature across backend, web, and app sequentially using subagents. Use when implementing a feature that spans multiple stacks in a vibeit workspace.

## When to Use

Invoke this skill when you need to implement a feature across multiple stacks (backend, web, app) in a vibeit workspace. This skill ensures correct sequencing, API contracts flow downward, and each phase is verified before the next begins.

Do not use this skill for single-stack changes; use the stack-specific `add-feature` skill instead.

## How It Works

### 1. Detect Stacks

Check which stacks exist in the workspace:
- `backend/` → backend stack exists
- `web/` → web stack exists
- `app/` → app stack exists

### 2. Route by Stack Count

- **Single stack only**: Launch that stack's native `add-feature` skill and stop. No sequencing needed.
- **Multiple stacks**: Proceed to sequential phases below.

### 3. Phase 1: Backend (if exists)

Create a Task subagent scoped to `backend/`:

```
Full Repository Path: [backend/ absolute path]
Use: backend/.cursor/skills/add-feature/SKILL.md
Constraint: Work only in backend/ folder
Verify: Run `make test` before returning
Document: Ensure learned/<feature>.md documents API contracts
```

Wait for Phase 1 to complete and verify green. **Do not start Phase 2 until Phase 1 passes.**

### 4. Phase 2: Web (if exists)

After Phase 1 completes, create a Task subagent scoped to `web/`:

```
Full Repository Path: [web/ absolute path]
Use: web/.cursor/skills/add-feature/SKILL.md
API Contract: Consume backend API from backend/learned/<feature>.md or backend module definitions
Constraint: Work only in web/ folder
Verify: Run `npm run build` before returning
```

Wait for Phase 2 to complete and verify green. **Do not start Phase 3 until Phase 2 passes.**

### 5. Phase 3: App (if exists)

After Phase 2 completes, create a Task subagent scoped to `app/`:

```
Full Repository Path: [app/ absolute path]
Use: app/.cursor/skills/add-feature/SKILL.md
API Contract: Consume backend API from backend/learned/<feature>.md or backend module definitions
Constraint: Work only in app/ folder
Verify: Run `flutter analyze` before returning
```

Wait for Phase 3 to complete and verify green.

## Surgical Scope

Each subagent must:
- Touch only its assigned stack folder (`backend/`, `web/`, or `app/`)
- Not refactor or modify unrelated code in other stacks
- Not make "drive-by" improvements across folders
- Respect existing code style and patterns within its stack

## Verification Checkpoints

| Phase | Command | Success Criteria |
|-------|---------|------------------|
| Backend | `make test` | All tests pass, learned/ doc exists |
| Web | `npm run build` | Build succeeds, no errors |
| App | `flutter analyze` | No analysis errors |

If any phase fails verification, fix it before advancing. Do not skip or defer.

## Example Workflow

```
User: "Add export to CSV feature across the full stack"

1. Check: backend/, web/, app/ all exist → proceed to phases
2. Phase 1: Launch subagent to add CSV export to backend API
   - Backend returns API endpoint contract in learned/csv-export.md
3. Phase 2: Launch subagent to add CSV export UI to web
   - Web consumes backend API, builds UI, verifies npm run build
4. Phase 3: Launch subagent to add CSV export to Flutter app
   - App consumes backend API, implements feature, verifies flutter analyze
5. Done: Feature complete across all stacks in correct order
```

## Notes

- Always run phases in strict order: backend → web → app
- Never parallelize phases; wait for completion + verification before next
- Backend API changes must be documented so web/app can consume them
- If a stack doesn't exist, skip that phase and continue
- If only one stack exists, use its native add-feature skill instead
