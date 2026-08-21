# Skill: Add Feature

Use this skill whenever adding a new section / page to the admin console.

## Steps

### 1. Create the feature folder

```
src/features/<name>/
  api/<name>Api.ts       — API call functions (typed, use apiClient)
  components/<Name>Home.tsx  — main component for this feature
  hooks/use<Name>Queries.ts  — TanStack Query hooks
  types.ts               — TypeScript interfaces
```

### 2. Create the page

```
src/pages/<Name>Page.tsx
```

Thin wrapper only:
```tsx
import { <Name>Home } from '@/features/<name>/components/<Name>Home'
export function <Name>Page() { return <<Name>Home /> }
```

### 3. Add the route

In `src/app/router.tsx`, inside the `<ProtectedRoute>` / `<AdminLayout>` block:
```tsx
<Route path="/<name>" element={<<Name>Page />} />
```
Import `<Name>Page` at the top.

### 4. Add the nav link

In `src/layouts/AdminSidebar.tsx`, append to `navItems`:
```tsx
{ label: '<Display Name>', to: '/<name>' }
```

### 5. Verify

```bash
npm run build
```

Must exit 0 before finishing.

### 6. Record learning

Append an entry to `.cursor/skills/add-feature/learned/README.md` using the format
in `.cursor/rules/learning.mdc`.
