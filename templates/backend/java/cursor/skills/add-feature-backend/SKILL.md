---
name: add-feature-backend
description: Add a tested Spring Boot feature and record its API contract.
---

# Add Backend Feature

1. Read the architecture and karpathy rules plus `learned/`.
2. Confirm endpoints, validation, persistence, and success criteria.
3. Add a controller, service, entity, and repository only as needed.
4. Add focused Spring tests and run `mvn test`.
5. Record endpoints, request/response shapes, and errors in
   `learned/<feature>.md`; update `learned/README.md`.

Use constructor injection. Do not refactor unrelated features or introduce
another framework or ORM.
