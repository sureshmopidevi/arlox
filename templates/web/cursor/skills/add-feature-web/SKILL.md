---
name: add-feature-web
description: >-
  Adds an admin-console feature: types, API client, TanStack Query hooks,
  page/components, router entry, and sidebar nav. Verifies with npm run build.
  Use when adding a web page, admin section, UI feature, or frontend resource
  that talks to the backend API.
---

# Add Feature (Web)

Add a vertical slice to the React admin console. Keep pages thin; put logic in `features/`.

## When to use

- New admin section / page
- New UI that calls backend APIs
- "Add X to the web app / dashboard / console"

**Do not use** alone for multi-stack work — prefer workspace `add-feature-fullstack` so backend contracts exist first.

## Before coding

0. **Design system** — read `.arlox/design-system.json` and `.cursor/rules/design-system.mdc`.
   Use **only** that library for UI (do not mix antd + shadcn, etc.).

| `id` | UI imports |
|------|------------|
| `tailwind` | Native HTML + Tailwind utility classes |
| `shadcn` | `@/components/ui/*` — add via `npx shadcn@latest add <component>` |
| `antd` | `antd` — Form, Table, Button, Input, … |
| `mui` | `@mui/material` |
| `chakra` | `@chakra-ui/react` |
| `mantine` | `@mantine/core` |

1. Read `.cursor/rules/architecture.mdc` and `.cursor/rules/karpathy.mdc` (if present).
2. Read `.cursor/skills/add-feature-web/learned/README.md`.
3. If backend exists, **read the API contract**:
   - `../../contracts/<feature>.md` (preferred)
   - Or `../backend/.cursor/skills/add-feature-backend/learned/<feature>.md`
4. Success criteria: route works, types match API, `npm run build` exits 0.

If the contract is missing and backend is in this workspace, stop and implement/document backend first (or ask the user).

## Steps

### 1. Feature folder

```text
src/features/<name>/
  api/<name>Api.ts
  components/<Name>Home.tsx
  hooks/use<Name>Queries.ts
  types.ts
```

- `types.ts` — mirror backend JSON (after envelope unwrap)
- `api/<name>Api.ts` — only via `@/api` client / shared axios instance; typed returns
- hooks — TanStack Query `useQuery` / `useMutation`; handle errors with existing helpers
- components — UI only; use the project's design system (see Phase 0); no raw axios

### 2. Thin page

```text
src/pages/<Name>Page.tsx
```

```tsx
import { <Name>Home } from '@/features/<name>/components/<Name>Home'

export function <Name>Page() {
  return <<Name>Home />
}
```

### 3. Route

In `src/app/router.tsx`, under `<AdminLayout>` (and `<ProtectedRoute>` if auth required):

```tsx
<Route path="/<name>" element={<<Name>Page />} />
```

Import the page at the top.

### 4. Sidebar nav

In `src/layouts/AdminSidebar.tsx`, append to `navItems`:

```tsx
{ label: '<Display Name>', to: '/<name>' }
```

Use `envConfig.appName` for branding — do not hardcode product names.

### 5. Wire to API contract

- Paths must match backend (`/api/v1/...` as configured in env)
- Request/response fields must match `learned/<feature>.md`
- Reuse auth from `authStore` / existing interceptors — do not invent token storage

### 6. Verify

```bash
npm run build
```

Must exit 0. Fix type errors before finishing.

### 7. Record learning

Append to `.cursor/skills/add-feature-web/learned/README.md`:

```markdown
## YYYY-MM-DD — <name>
**What**: added <name> feature (api, hooks, page, route, nav).
**API**: consumed <endpoints>.
**Gotchas**: <non-obvious decisions>.
```

## Anti-patterns

- Business logic or API calls inside `pages/`
- Skipping route or sidebar (feature unreachable)
- Inventing endpoints that backend does not expose
- New state libraries when Zustand/React Query already cover the case
- Drive-by CSS/architecture refactors

## Done checklist

- [ ] Feature folder + thin page
- [ ] Route registered
- [ ] Sidebar link added
- [ ] Types/API match backend contract
- [ ] `npm run build` passes
- [ ] Learning entry recorded
