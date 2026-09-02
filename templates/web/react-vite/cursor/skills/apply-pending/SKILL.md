# Skill: Apply Pending

Review the learned index and apply unapplied lessons to the codebase.

## Steps

1. Read `.cursor/skills/add-feature-web/learned/README.md`.
2. Identify entries **not yet marked as applied**.
3. For each pending entry:
   a. Determine if the pattern should be generalised into core code or rules.
   b. If yes — propose the change to the user before applying.
   c. If no — mark as reviewed and note why it stays feature-local.
4. After applying, add `**Applied:** YYYY-MM-DD — <what changed>` to the entry.

## What counts as "applied"

- A learned pattern extracted into a shared hook/util.
- A rule updated in `.cursor/rules/*.mdc`.
- A skill step updated in `.cursor/skills/add-feature-web/SKILL.md`.

## What stays feature-local

- Patterns that are specific to one domain and unlikely to recur.
- Experiments that didn't pan out (still worth documenting, not worth spreading).
