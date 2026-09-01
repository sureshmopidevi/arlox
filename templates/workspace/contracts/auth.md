# auth — API contract

Login and session endpoints used by web and Flutter app.

## Endpoints

| Method | Path | Auth | Request | Response |
|--------|------|------|---------|----------|
| POST | /api/v1/auth/login | no | `{ "email": string, "password": string }` | `{ data: { token: string, user: { id, email, name? } } }` |
| GET | /api/v1/auth/me | yes | — | `{ data: { id: number, email: string, name?: string } }` |

## Errors

| Status | When |
|--------|------|
| 400 | invalid request body |
| 401 | invalid credentials (login) or missing/invalid token (me) |
| 404 | user not found (me) |

## Notes for web/app

- Envelope: unwrap `data` from JSON responses
- Auth header: `Authorization: Bearer <token>`
- Web: `src/features/auth/` — Zustand + React Query
- App: `lib/features/auth/` — Riverpod + secure storage
