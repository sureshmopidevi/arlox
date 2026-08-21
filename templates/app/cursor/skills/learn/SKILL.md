# Skill: learn

Persist a decision or lesson to the appropriate `learned/README.md`.

## Trigger

Invoke when asked to **"learn \<something\>"** or after completing a significant task.

---

## Steps

1. Identify the target `learned/README.md` based on scope:

   | Work type | File |
   |-----------|------|
   | Feature work | `.cursor/skills/add-feature/learned/README.md` |
   | Core change | `.cursor/skills/add-feature/learned/README.md` (under `### Core`) |
   | Deferred/risky | Same file, prefix title with `[PENDING]` |

2. Append an entry at the top of the entries section (newest first):

   ```markdown
   ## YYYY-MM-DD — <short title>
   **What**: one sentence describing what was done.
   **Why**: the reason this approach was chosen.
   **Gotchas**: anything non-obvious a future agent should know.
   ```

3. For deferred items, prefix the title:

   ```markdown
   ## YYYY-MM-DD — [PENDING] <short title>
   ```

4. Confirm the entry was written by reading back the file.

---

## Notes

- Entries should be concise — 3–5 lines maximum per entry.
- Do not edit or remove existing entries; only append.
- The `apply-pending` skill processes `[PENDING]` entries.
