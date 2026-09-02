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
func finalizeStack(variant workspace.Variant, stackDir string, data Data) (string, error) {
	switch variant {
	case workspace.GoGin:
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

	case workspace.ReactVite:
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

	case workspace.NodeExpress, workspace.NodeFastify,
		workspace.NextJS, workspace.VueVite, workspace.SvelteVite,
		workspace.Angular, workspace.Nuxt, workspace.ReactPWA:
		if err := copyIfPresent(
			filepath.Join(stackDir, ".env.example"),
			filepath.Join(stackDir, ".env"),
		); err != nil {
			return "", err
		}
		if err := runInDir(stackDir, "npm", "install", "--no-audit", "--no-fund"); err != nil {
			return "npm install skipped/failed (offline or network error) — run: npm install", nil
		}
		return "", nil

	case workspace.PyFastAPI:
		if err := runInDir(stackDir, "uv", "sync", "--extra", "dev"); err != nil {
			return "Python dependencies skipped/failed — run: uv sync --extra dev", nil
		}
		return "", nil

	case workspace.JavaSpring:
		if err := runInDir(stackDir, "mvn", "dependency:go-offline"); err != nil {
			return "Maven dependency resolution skipped/failed — run: mvn dependency:go-offline", nil
		}
		return "", nil

	case workspace.Flutter:
		var warning string
		if err := runInDir(stackDir, "flutter", "pub", "get"); err != nil {
			warning = "flutter pub get skipped/failed (offline or network error) — run: make app.get"
		}
		return warning, nil

	case workspace.NativeIOS:
		if err := runInDir(stackDir, "swift", "package", "resolve"); err != nil {
			return "Swift package resolution skipped/failed — run: swift package resolve", nil
		}
		return "", nil

	case workspace.NativeAndroid:
		if err := runInDir(stackDir, "gradle", "dependencies"); err != nil {
			return "Gradle dependency resolution skipped/failed — run: gradle dependencies", nil
		}
		return "", nil
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

func copyIfPresent(src, dst string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil
	}
	return copyIfMissing(src, dst)
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
