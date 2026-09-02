package cli

import (
	"path/filepath"
	"testing"
)

func TestRunUpgradeFallsBackToRemoteWithoutLocalSource(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("ARLOX_HOME", "")
	t.Chdir(t.TempDir())

	if _, err := findArloxSource(""); err == nil {
		t.Fatal("expected no local source in isolated HOME")
	}
}

func TestFindArloxSourceIgnoresInvalidSavedConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("ARLOX_HOME", "")

	missing := filepath.Join(home, "missing-arlox")
	if err := writeArloxSourceConfig(missing); err != nil {
		t.Fatal(err)
	}

	_, err := findArloxSource("")
	if err == nil {
		t.Fatal("expected error when saved source path is not a repo")
	}
}
