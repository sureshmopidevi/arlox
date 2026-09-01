package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunInitBrownfieldBackendOnly(t *testing.T) {
	root := filepath.Join(t.TempDir(), "brownfield")
	backendDir := filepath.Join(root, "backend")
	if err := os.MkdirAll(backendDir, 0o755); err != nil {
		t.Fatal(err)
	}

	originalGoMod := "module github.com/acme/legacy-api\n\ngo 1.22\n"
	if err := os.WriteFile(filepath.Join(backendDir, "go.mod"), []byte(originalGoMod), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := runInit(root, ".", "brownfield", false, false, stackFlags{}); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	wsFile := filepath.Join(root, "brownfield.code-workspace")
	if _, err := os.Stat(wsFile); os.IsNotExist(err) {
		t.Fatalf("expected %s", wsFile)
	}
	if _, err := os.Stat(filepath.Join(root, "Makefile")); os.IsNotExist(err) {
		t.Fatal("expected root Makefile")
	}
	if _, err := os.Stat(filepath.Join(root, "contracts/auth.md")); os.IsNotExist(err) {
		t.Fatal("expected contracts/auth.md")
	}

	got, err := os.ReadFile(filepath.Join(backendDir, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != originalGoMod {
		t.Fatalf("backend go.mod changed:\n%s", got)
	}
}

func TestRunInitRequiresStack(t *testing.T) {
	root := t.TempDir()
	if err := runInit(root, ".", "empty", false, false, stackFlags{}); err == nil {
		t.Fatal("expected error when no stacks present")
	}
}

func TestInferInitNameFromDirectory(t *testing.T) {
	root := filepath.Join(t.TempDir(), "myapp")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	name, err := inferInitName(root, "")
	if err != nil {
		t.Fatal(err)
	}
	if name != "myapp" {
		t.Fatalf("want myapp, got %q", name)
	}
}

func TestInferInitNameRequiresFlagForInvalidDir(t *testing.T) {
	root := filepath.Join(t.TempDir(), "MyApp")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := inferInitName(root, ""); err == nil {
		t.Fatal("expected error for uppercase directory name")
	}
}
