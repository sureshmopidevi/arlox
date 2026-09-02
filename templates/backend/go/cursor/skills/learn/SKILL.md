# Skill: Learn

Stage a patch to a skill or rule for later review.

## When to use
- You discover the canonical SKILL.md or a rule needs a correction or addition.
- An architecture decision was made that future agents should know.
- A repeated mistake should be prevented.

## Steps

1. Identify what needs to change (skill file path + the change).
2. Write the patch to `.cursor/skills/.pending/<timestamp>-<slug>.md`:

```markdown
# Pending patch: <slug>

**Target file:** `.cursor/skills/<path>/SKILL.md`  
**Type:** correction | addition | removal  

## Change
<describe the change concisely>

## Proposed diff / new content
<show old → new, or just the addition>

## Reason
<why this is needed>
```

3. Do NOT modify the live skill/rule file directly.
4. Tell the user: "Staged patch to `.cursor/skills/.pending/<file>`. Run apply-pending to review."

## Notes
- One patch file per change.
- Keep patches small and focused.
- The apply-pending skill handles promotion.
