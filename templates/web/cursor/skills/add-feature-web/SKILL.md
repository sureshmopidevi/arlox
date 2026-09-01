---
name: add-feature-web
description: >-
  Ship an admin feature reachable at a new route with types, API client, TanStack
  Query hooks, and sidebar nav; npm run build green. Use when adding a web page,
  admin section, or frontend resource that calls the backend API.
---

# Add Feature (Web)

## User intent

When done, the user can open a new route in the admin console, see the feature UI built with the project's design system, and data loads from the backend per the shared contract.

## Read first

1. **`.arlox/project.json`** (workspace root) — `stacks.web.packageName`, naming conventions.
2. **`web/.arlox/design-system.json`** and **`.cursor/rules/design-system.mdc`** — use only that UI library.
3. **`contracts/<feature>.md`** (workspace root, preferred) or `../backend/.cursor/skills/add-feature-backend/learned/<feature>.md` as fallback.
4. **`.cursor/rules/architecture.mdc`** and **`.cursor/rules/karpathy.mdc`** (if present).
5. **`.cursor/skills/add-feature-web/learned/README.md`**.

**Do not use** alone for multi-stack work — prefer workspace `add-feature-fullstack`.

| `design-system id` | UI imports |
|--------------------|------------|
| `tailwind` | Native HTML + Tailwind utility classes |
| `shadcn` | `@/components/ui/*` — add via `npx shadcn@latest add <component>` |
| `antd` | `antd` |
| `mui` | `@mui/material` |
| `chakra` | `@chakra-ui/react` |
| `mantine` | `@mantine/core` |

## Success criteria

| Check | Must be true |
|-------|----------------|
| Contract | Types and API paths match `contracts/<feature>.md` |
| UI | Route + sidebar nav; feature reachable |
| Build | `npm run build` exits 0 |
| Learn | Entry appended via `learn/SKILL.md` |

If the contract is missing and backend exists in this workspace, stop and implement or document backend first (or ask the user).

## Steps

1. **Feature folder** — `src/features/<name>/` with `types.ts`, `api/<name>Api.ts`, hooks, components.
2. **Thin page** — `src/pages/<Name>Page.tsx` imports feature component only.
3. **Route** — register in `src/app/router.tsx` under `<AdminLayout>` (and `<ProtectedRoute>` if auth required).
4. **Sidebar** — append to `navItems` in `src/layouts/AdminSidebar.tsx`.
5. **Verify** — `npm run build`.
6. **Learn** — run **`.cursor/skills/learn/SKILL.md`** (mandatory).

## Verify

```bash
npm run build
```

## Learn

Run **`.cursor/skills/learn/SKILL.md`** after verify passes. Append to **`learned/README.md`**.

## Anti-patterns

- Business logic or API calls inside `pages/`
- Skipping route or sidebar (feature unreachable)
- Inventing endpoints the backend does not expose
- Mixing UI libraries or ignoring `design-system.json`
- Drive-by CSS or architecture refactors
