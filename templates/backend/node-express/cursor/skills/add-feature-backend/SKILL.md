---
name: add-feature-backend
description: Add a tested Express feature and record its API contract.
---

# Add Backend Feature

1. Read the architecture and karpathy rules plus `learned/`.
2. Confirm endpoints, validation, persistence, and success criteria.
3. Add routes, service logic, and Prisma persistence only as needed.
4. Add a Prisma migration when the schema changes.
5. Add focused Vitest tests; run `npm test` and `npm run build`.
6. Record endpoints, request/response shapes, and errors in
   `learned/<feature>.md`; update `learned/README.md`.

Do not refactor unrelated features or introduce another framework or ORM.
