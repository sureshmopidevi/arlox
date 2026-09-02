---
name: add-feature-backend
description: Add a tested FastAPI feature and record its API contract.
---

# Add Backend Feature

1. Read the architecture and karpathy rules plus `learned/`.
2. Confirm endpoints, validation, persistence, and success criteria.
3. Add a thin router, service logic, and SQLAlchemy persistence only as needed.
4. Add an Alembic migration when the model changes.
5. Add focused tests and run `uv run pytest`.
6. Record endpoints, request/response shapes, and errors in
   `learned/<feature>.md`; update `learned/README.md`.

Do not refactor unrelated features or introduce another framework or ORM.
