---
name: dev-workflow
description: Run and verify the Spring Boot backend.
---

# Development Workflow

1. Confirm Java 21 and Maven are available.
2. Start development with `mvn spring-boot:run`.
3. Run `mvn test` after code changes.
4. Run `mvn package` before release.
5. Configure overrides through environment variables referenced by
   `application.properties`.
