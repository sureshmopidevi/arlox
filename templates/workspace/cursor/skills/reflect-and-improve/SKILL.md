---
name: reflect-and-improve
description: >-
  Review learned/README.md entries across workspace and stacks, run apply-pending
  where patterns should be promoted, and surface unapplied lessons. Use when
  maintaining agent guardrails, after several features, or when arlox skills status
  shows pending learnings.
---

# Reflect and Improve

## User intent

When done, proven patterns from `learned/README.md` are promoted into rules or skills (or marked reviewed), and the team knows what learning debt remains.

## Read first

1. **`.arlox/project.json`**
2. Run **`arlox skills status`** from workspace root (lists unapplied entries per stack).
3. Each present stack's **`learned/README.md`** and **`apply-pending/SKILL.md`**.

## Success criteria

| Check | Must be true |
|-------|----------------|
| Scan | Every present stack's `learned/README.md` reviewed |
| Promote | Reusable patterns applied to `.cursor/rules/*.mdc` or skill steps |
| Mark | Applied entries get `**Applied:** YYYY-MM-DD — <what changed>` |
| Status | `arlox skills status` shows fewer (or zero) unapplied entries |

## Steps

1. **Status** — `arlox skills status` from workspace root.
2. **Per stack** (backend, web, app as present):
   - Read `learned/README.md`.
   - Run that stack's **`apply-pending/SKILL.md`** for entries without `**Applied:**`.
3. **Cross-stack** — if the same pitfall appears in multiple stacks, consider a workspace-level note in `contracts/README.md` or root `.cursor/rules/`.
4. **Re-check** — `arlox skills status`.

## Verify

```bash
arlox skills status
```

Target: zero unapplied entries, or each remaining entry explicitly marked reviewed with reason.

## Learn

This skill *is* the maintenance loop. Do not append meta-entries unless you discovered a new process gap worth documenting.

## Anti-patterns

- Auto-editing canonical `SKILL.md` files with one-off feature details
- Deleting `learned/` entries (`arlox skills update` preserves them)
- Promoting unproven or single-use patterns to global rules
- Skipping stacks that have pending entries
