package generate

import (
	"fmt"
	"os"
	osexec "os/exec"
	"path/filepath"
	"github.com/sureshmopidevi/arlox/internal/naming"

	arloxexec "github.com/sureshmopidevi/arlox/internal/exec"
	"github.com/sureshmopidevi/arlox/internal/ui"
	"github.com/sureshmopidevi/arlox/internal/workspace"
)

// finalizeStack installs dependencies and runs one-time setup after templates render.
// Non-fatal issues (like offline package install failure) return a warning string without failing.
func finalizeStack(stack workspace.Stack, stackDir string, data Data) (string, error) {
	switch stack {
	case workspace.Backend:
		if err := copyIfMissing(
			filepath.Join(stackDir, "configs/local/app.env.example"),
			filepath.Join(stackDir, "configs/local/app.env"),
		); err != nil {
			return "", err
		}
		var warning string
		ui.Dim("backend: running go mod tidy…")
		if err := arloxexec.RunInDir(stackDir, "go", "mod", "tidy"); err != nil {
			warning = "go mod tidy skipped/failed (offline or network error) — run: make backend.tidy"
		}
		tryPostgresSetup(stackDir, data.Name)
		return warning, nil

	case workspace.Web:
		if err := copyIfMissing(
			filepath.Join(stackDir, ".env.example"),
			filepath.Join(stackDir, ".env"),
		); err != nil {
			return "", err
		}
		var warning string
		ui.Dim("web: running npm install…")
		if err := arloxexec.RunInDir(stackDir, "npm", "install", "--no-audit", "--no-fund"); err != nil {
			warning = "npm install skipped/failed (offline or network error) — run: make web.install"
		}
		return warning, nil

	case workspace.App:
		var warning string
		ui.Dim("app: running flutter pub get…")
		if err := arloxexec.RunInDir(stackDir, "flutter", "pub", "get"); err != nil {
			warning = "flutter pub get skipped/failed (offline or network error) — run: make app.get"
		}
		return warning, nil
	}
	return "", nil
}

func copyIfMissing(src, dst string) error {
	if _, err := os.Stat(dst); err == nil {
		return nil
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("copy %s: %w", filepath.Base(src), err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}

// tryPostgresSetup starts local Postgres (best effort) and creates the project DB.
func tryPostgresSetup(stackDir, projectName string) {
	workspaceRoot := filepath.Dir(stackDir)
	composeFile := filepath.Join(workspaceRoot, "docker-compose.yml")
	if _, err := os.Stat(composeFile); err == nil {
		if _, err := osexec.LookPath("docker"); err == nil {
			ui.Dim("backend: starting postgres via docker compose…")
			_ = arloxexec.RunInDir(workspaceRoot, "docker", "compose", "up", "-d", "postgres")
			return
		}
	}

	if _, err := osexec.LookPath("createdb"); err != nil {
		return
	}
	_ = osexec.Command("brew", "services", "start", "postgresql@16").Run()
	_ = osexec.Command("brew", "services", "start", "postgresql").Run()
	dbName := naming.Snake(projectName)
	_ = osexec.Command("createdb", dbName).Run()
}

func dbNameFromProject(name string) string {
	return naming.Snake(name)
}
