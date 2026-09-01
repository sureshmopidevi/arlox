# Web design systems

When you include the **web** stack, arlox scaffolds a React admin console with a UI library of your choice. The selection is made at create/add time and persisted so AI agents build new features with the same components.

---

## Choosing a design system

### Interactive create

```bash
arlox create myapp
```

1. Pick stacks (include **web**).
2. Answer **Web design system** — six options, same style as the stack picker.

### Non-interactive (CI / scripts)

Pass `--web-ui` with `--web`. If omitted in a non-TTY environment, arlox defaults to **`tailwind`**.

```bash
arlox create myapp --web --web-ui shadcn
arlox add --web --web-ui antd
```

---

## Supported systems

| ID | Label | Best for |
|----|-------|----------|
| `tailwind` | Tailwind only | Minimal scaffold, full control, no component library |
| `shadcn` | shadcn/ui | Tailwind + accessible Radix primitives, copy-paste components |
| `antd` | Ant Design | Enterprise admin UIs, dense data tables, forms |
| `mui` | Material UI | Material Design, large ecosystem |
| `chakra` | Chakra UI | Accessible components, composable layout props |
| `mantine` | Mantine | Batteries-included forms, modals, hooks |

All options ship the same **feature skeleton**: login form, admin layout, sidebar, home page, auth hooks, and Vitest smoke test — implemented with that library’s components.

---

## What gets generated

### Two-layer templates

1. **Web core** (`templates/web/`) — router, API client, auth store, architecture rules, feature folder layout.
2. **Design system overlay** (`templates/design-systems/web/<id>/`) — `package.json`, styles, UI shell, `design-system.mdc`, `AGENTS.md`.

Overlays win on path conflicts, so each project gets one coherent UI stack.

### Manifest file

After scaffold, the choice is written to:

```text
web/.arlox/design-system.json
```

Example:

```json
{
  "id": "shadcn",
  "label": "shadcn/ui",
  "arloxVersion": "0.15.0"
}
```

Agents and skills read this file — do not guess from `package.json` alone.

### Cursor rules

Each overlay installs:

```text
web/.cursor/rules/design-system.mdc
```

Rules describe which imports to use, how to add components, and what **not** to mix (e.g. no antd + shadcn in one project).

---

## Adding features with the same system

Use the **`add-feature-web`** skill (or workspace **`add-feature-fullstack`**).

**Phase 0** in the skill:

1. Read `web/.arlox/design-system.json`.
2. Read `.cursor/rules/design-system.mdc`.
3. Build UI with **that library only**.

| `id` | Where to get components |
|------|-------------------------|
| `tailwind` | Native HTML + Tailwind utility classes |
| `shadcn` | `src/components/ui/*` — add more via `npx shadcn@latest add <name>` |
| `antd` | `import { Button, Table, Form } from 'antd'` |
| `mui` | `import { Button, TextField } from '@mui/material'` |
| `chakra` | `import { Button, Input } from '@chakra-ui/react'` |
| `mantine` | `import { Button, TextInput } from '@mantine/core'` |

---

## Examples

### shadcn admin with backend

```bash
arlox create myapp --backend --web --web-ui shadcn
cd myapp
make db-up && make backend.setup
make dev
```

Generated web includes `components.json` and pre-shipped `Button`, `Input`, `Label`, `Card`.

Add a table later:

```bash
cd web
npx shadcn@latest add table
```

### Ant Design web-only prototype

```bash
arlox create prototype --web --web-ui antd
cd prototype/web && npm run dev
```

### Add web to existing workspace with Mantine

```bash
cd existing-project
arlox add --web --web-ui mantine
```

---

## Architecture (for contributors)

| Path | Role |
|------|------|
| `internal/designsystems/catalog.go` | Valid IDs, labels, overlay paths |
| `internal/cli/web_design_system.go` | `--web-ui` flag + interactive prompt |
| `internal/generate/generate.go` | Two-pass render for `workspace.Web` |
| `internal/generate/design_system.go` | Writes `.arlox/design-system.json` |
| `templates/design-systems/web/catalog.json` | Human-readable catalog mirror |
| `templates/design-systems/web/<id>/` | Per-system overlay files |

Adding a new design system:

1. Add entry to `WebCatalog` in `catalog.go` and `catalog.json`.
2. Create overlay directory with the same shell files as existing systems.
3. Extend `verify.sh` or `design_system_test.go` to render/build the new overlay.
4. Bump `internal/version/VERSION` (MINOR).

---

## FAQ

**Can I change design system after create?**  
Not via CLI today. You would migrate dependencies and components manually, or regenerate the web stack in a fresh folder.

**Does mobile/app use the same design system?**  
No. Flutter app uses Material; web design systems apply only to `web/`.

**Why default to Tailwind in CI?**  
Non-interactive runs have no prompt; `tailwind` is the lightest overlay and matches the original scaffold.

**Do repair/skills update change my design system?**  
`arlox skills update` refreshes shared skills from `templates/web/` but not the overlay-specific `design-system.mdc` unless you regenerate web. The manifest in `.arlox/` is set at create time.
