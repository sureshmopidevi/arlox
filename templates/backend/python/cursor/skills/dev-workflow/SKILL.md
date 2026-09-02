---
name: dev-workflow
description: Run and verify the FastAPI backend.
---

# Development Workflow

1. Install dependencies with `uv sync --extra dev`.
2. Start the API with `uv run uvicorn app.main:app --reload`.
3. Run `uv run pytest` after code changes.
4. For model changes, generate an Alembic revision, inspect it, and run
   `uv run alembic upgrade head`.
