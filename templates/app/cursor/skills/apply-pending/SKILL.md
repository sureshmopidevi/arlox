# Skill: apply-pending

Apply deferred learnings from `learned/README.md` to the codebase.

## Trigger

Invoke when asked to **"apply pending learnings"** or **"apply pending"**.

---

## Steps

1. Read `.cursor/skills/add-feature-mobile/learned/README.md`.

2. Identify all entries whose title starts with `[PENDING]`.

3. For each pending entry (oldest first — apply in chronological order):

   a. Understand the change described.  
   b. Check whether it still applies (the codebase may have changed).  
   c. If it still applies: implement the change in the relevant files.  
   d. Run `flutter analyze` — fix any warnings introduced by the change.  
   e. Remove the `[PENDING]` prefix from the entry title in `learned/README.md`.  
   f. If it no longer applies: remove `[PENDING]` and append a note:  
      `**Resolution**: superseded by <reason> — no action taken.`

4. Report a summary of what was applied and what was skipped.

---

## Notes

- Do not apply entries that have no `[PENDING]` prefix — they are already resolved.
- If an entry is ambiguous or risky, ask for clarification before applying.
- After applying, always run `flutter analyze` to confirm zero warnings.
