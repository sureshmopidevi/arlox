# Skill: Make Workflow

Reference for all `make` targets in this project.

| Target | What it does |
|--------|-------------|
| `make setup` | `db-up` + `db-create` + `go mod tidy` + copy env example |
| `make run` | `go run ./cmd/server` with env file |
| `make build` | compile to `bin/<app>` |
| `make test` | `go test ./... -count=1` |
| `make lint` | `golangci-lint run ./...` |
| `make tidy` | `go mod tidy` |
| `make doctor` | prints go, psql, lint versions + env file status |
| `make db-up` | `brew services start postgresql` |
| `make db-create` | `createdb <db_name>` |

## Rules
- Always use `make` — never raw `go run` or `go build`.
- After dependency changes: `make tidy`.
- After any code change: `make test`.
- `make doctor` is your first debugging step when something doesn't work.
