# Skill: Apply Pending

Review and promote staged patches from `.cursor/skills/.pending/` into live skill/rule files.

## Steps

1. List all files in `.cursor/skills/.pending/`.
2. For each patch file:
   a. Read the patch.
   b. Read the target file it references.
   c. Decide: **apply**, **reject**, or **defer**.
      - Apply: the change is correct, clear, and improves the target.
      - Reject: the change is wrong, out of scope, or already covered.
      - Defer: needs more context; leave in pending.
3. For each **apply** decision:
   - Edit the target skill/rule file with the change.
   - Delete the patch file from `.pending/`.
4. For each **reject** decision:
   - Delete the patch file.
   - Note the rejection reason in your response.
5. Report: how many applied, rejected, deferred.

## Rules
- Never apply a patch that breaks existing content without resolving the conflict.
- Keep skill files under ~60 lines unless the content genuinely requires more.
- After applying, verify the target file reads coherently end-to-end.
