---
name: dev-workflow
description: Run and verify the Express backend.
---

# Development Workflow

1. Install dependencies with `npm install`.
2. Generate Prisma client with `npm run db:generate`.
3. Start development with `npm run dev`.
4. Run `npm test` and `npm run build` after changes.
5. For schema changes, run `npm run db:migrate` and commit the migration.
