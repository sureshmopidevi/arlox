package generate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
		if err := runInDir(stackDir, "go", "mod", "tidy"); err != nil {
			warning = "go mod tidy skipped/failed (offline or network error) — run: make backend.tidy"
		}
		tryPostgresSetup(data.Name)
		return warning, nil

	case workspace.Web:
		if err := copyIfMissing(
			filepath.Join(stackDir, ".env.example"),
			filepath.Join(stackDir, ".env"),
		); err != nil {
			return "", err
		}
		var warning string
		if err := runInDir(stackDir, "npm", "install", "--no-audit", "--no-fund"); err != nil {
			warning = "npm install skipped/failed (offline or network error) — run: make web.install"
		}
		return warning, nil

	case workspace.App:
		var warning string
		if err := runInDir(stackDir, "flutter", "pub", "get"); err != nil {
			warning = "flutter pub get skipped/failed (offline or network error) — run: make app.get"
		}
		return warning, nil
	}
	return "", nil
}

func runInDir(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w\n%s", err, out)
	}
	return nil
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
func tryPostgresSetup(projectName string) {
	if _, err := exec.LookPath("createdb"); err != nil {
		return
	}
	_ = exec.Command("brew", "services", "start", "postgresql@16").Run()
	_ = exec.Command("brew", "services", "start", "postgresql").Run()
	dbName := dbNameFromProject(projectName)
	_ = exec.Command("createdb", dbName).Run()
}

func dbNameFromProject(name string) string {
	return strings.ReplaceAll(name, "-", "_")
}
