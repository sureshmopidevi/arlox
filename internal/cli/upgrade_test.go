package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindArloxSourceUsesSavedConfig(t *testing.T) {
	t.Chdir(t.TempDir())

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("ARLOX_HOME", "")

	repo := t.TempDir()
	if err := os.WriteFile(filepath.Join(repo, "go.mod"), []byte("module github.com/sureshmopidevi/arlox\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(repo, "cmd", "arlox"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "cmd", "arlox", "main.go"), []byte("package main\nfunc main(){}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := findArloxSource("")
	if err == nil || got != "" {
		t.Fatalf("expected missing source, got %q err=%v", got, err)
	}

	if err := writeArloxSourceConfig(repo); err != nil {
		t.Fatal(err)
	}

	got, err = findArloxSource("")
	if err != nil {
		t.Fatal(err)
	}
	if got != repo {
		t.Fatalf("got %q, want %q", got, repo)
	}
}

func TestFindArloxSourcePrefersExplicitFlag(t *testing.T) {
	t.Chdir(t.TempDir())

	home := t.TempDir()
	t.Setenv("HOME", home)

	repo := t.TempDir()
	other := t.TempDir()
	for _, dir := range []string{repo, other} {
		if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module github.com/sureshmopidevi/arlox\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(dir, "cmd", "arlox"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "cmd", "arlox", "main.go"), []byte("package main\nfunc main(){}\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if err := writeArloxSourceConfig(other); err != nil {
		t.Fatal(err)
	}

	got, err := findArloxSource(repo)
	if err != nil {
		t.Fatal(err)
	}
	if got != repo {
		t.Fatalf("got %q, want %q", got, repo)
	}
}
