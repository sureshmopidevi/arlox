package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sureshmopidevi/arlox/internal/workspace"
)

func TestRunRepairRecreatesMissingWorkspaceAndConfigs(t *testing.T) {
	name := "repairdemo"
	root := filepath.Join(t.TempDir(), name)
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create workspace and a backend folder with app.env.example
	wsFile := filepath.Join(root, name+".code-workspace")
	if err := workspace.WriteWorkspaceFile(wsFile, name, []string{"."}); err != nil {
		t.Fatal(err)
	}

	backendDir := filepath.Join(root, "backend")
	cfgDir := filepath.Join(backendDir, "configs", "local")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(backendDir, "go.mod"), []byte("module demo-backend\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cfgDir, "app.env.example"), []byte("PORT=8080\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Remove .code-workspace to simulate broken/missing file
	if err := os.Remove(wsFile); err != nil {
		t.Fatal(err)
	}

	// Run repair (disable deps install for fast isolated unit test)
	if err := runRepair(root, false, false); err != nil {
		t.Fatalf("runRepair failed: %v", err)
	}

	// 1. Verify .code-workspace is restored
	if _, err := os.Stat(wsFile); os.IsNotExist(err) {
		t.Fatalf("expected %s to be restored", wsFile)
	}

	// 2. Verify root Makefile is restored
	if _, err := os.Stat(filepath.Join(root, "Makefile")); os.IsNotExist(err) {
		t.Fatalf("expected root Makefile to be restored")
	}

	// 3. Verify backend app.env is seeded
	if _, err := os.Stat(filepath.Join(cfgDir, "app.env")); os.IsNotExist(err) {
		t.Fatalf("expected backend configs/local/app.env to be restored")
	}

	// 4. Verify backend git repo is initialized
	if _, err := os.Stat(filepath.Join(backendDir, ".git")); os.IsNotExist(err) {
		t.Fatalf("expected backend .git to be initialized")
	}
}
