package generate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sureshmopidevi/arlox/internal/workspace"
)

func TestDBNameFromProject(t *testing.T) {
	tests := map[string]string{
		"my-app":        "my_app",
		"issue_tracker": "issue_tracker",
		"a_b-c":         "a_b_c",
	}
	for in, want := range tests {
		if got := dbNameFromProject(in); got != want {
			t.Errorf("dbNameFromProject(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestCopyIfMissing(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "source.env")
	dst := filepath.Join(tmp, "sub", "target.env")

	if err := os.WriteFile(src, []byte("KEY=VALUE\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := copyIfMissing(src, dst); err != nil {
		t.Fatalf("copyIfMissing failed: %v", err)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "KEY=VALUE\n" {
		t.Fatalf("content mismatch: got %q", string(data))
	}

	// Idempotency check: should not overwrite existing dst
	if err := os.WriteFile(src, []byte("MODIFIED\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := copyIfMissing(src, dst); err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(dst)
	if string(data) != "KEY=VALUE\n" {
		t.Fatalf("expected existing file to be preserved, got %q", string(data))
	}
}

func TestFinalizeStackSeedsConfigsEvenIfToolsWarning(t *testing.T) {
	tmp := t.TempDir()
	backendDir := filepath.Join(tmp, "backend")
	cfgDir := filepath.Join(backendDir, "configs", "local")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	example := filepath.Join(cfgDir, "app.env.example")
	if err := os.WriteFile(example, []byte("PORT=8080\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// finalizeStack will attempt `go mod tidy` in an empty dir, which will produce a warning
	// but the config file MUST be created and no fatal error returned.
	warn, err := finalizeStack(workspace.Backend, backendDir, Data{Name: "test"})
	if err != nil {
		t.Fatalf("expected non-fatal return, got err: %v", err)
	}
	_ = warn // warning is expected

	targetEnv := filepath.Join(cfgDir, "app.env")
	if _, err := os.Stat(targetEnv); os.IsNotExist(err) {
		t.Fatalf("expected %s to be created by finalizeStack", targetEnv)
	}
}

func TestMapTemplatePath(t *testing.T) {
	tests := map[string]string{
		"cursor/rules/karpathy.mdc":      ".cursor/rules/karpathy.mdc",
		"gitignore":                      ".gitignore",
		"env.example":                    ".env.example",
		"env.example.tmpl":               ".env.example.tmpl",
		"configs/local/app.env.example":  "configs/local/app.env.example",
		"lib/main.dart.tmpl":             "lib/main.dart.tmpl",
	}
	for in, want := range tests {
		if got := mapTemplatePath(in); got != want {
			t.Errorf("mapTemplatePath(%q) = %q, want %q", in, got, want)
		}
	}
}
