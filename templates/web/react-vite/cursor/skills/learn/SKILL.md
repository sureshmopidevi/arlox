# Skill: Learn

Capture what was learned after completing a task.

## When to use

- After adding a feature (required by `add-feature-web` skill step 7).
- After debugging a non-obvious issue.
- After refactoring a core module.

## Steps

1. **Summarise** what was built or fixed in one sentence.
2. **List files** that were added or meaningfully changed.
3. **Identify decisions** — why did you choose this approach over alternatives?
4. **Note pitfalls** — what would trip up the next person?
5. **Append** to `.cursor/skills/add-feature-web/learned/README.md` using the format:

```markdown
## YYYY-MM-DD — <Feature or Fix Name>

**Pattern:** <one-line summary>
**Files touched:** <comma-separated list>
**Key decisions:**
- <decision 1>

**Pitfalls:**
- <pitfall 1>
```

## Scope

- Feature learnings → `learned/README.md` (always).
- Core module changes (api, stores, layouts) → also add a note in the relevant
  `.cursor/rules/*.mdc` if the rule needs updating.
