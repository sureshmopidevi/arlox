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
func finalizeStack(stack workspace.Stack, stackDir string, data Data) error {
	switch stack {
	case workspace.Backend:
		if err := runInDir(stackDir, "go", "mod", "tidy"); err != nil {
			return fmt.Errorf("go mod tidy: %w", err)
		}
		if err := copyIfMissing(
			filepath.Join(stackDir, "configs/local/app.env.example"),
			filepath.Join(stackDir, "configs/local/app.env"),
		); err != nil {
			return err
		}
		tryPostgresSetup(data.Name)
	case workspace.Web:
		if err := runInDir(stackDir, "npm", "install", "--no-audit", "--no-fund"); err != nil {
			return fmt.Errorf("npm install: %w", err)
		}
		if err := copyIfMissing(
			filepath.Join(stackDir, ".env.example"),
			filepath.Join(stackDir, ".env"),
		); err != nil {
			return err
		}
	case workspace.App:
		if err := runInDir(stackDir, "flutter", "pub", "get"); err != nil {
			return fmt.Errorf("flutter pub get: %w", err)
		}
	}
	return nil
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
