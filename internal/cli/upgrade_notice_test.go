package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestLatestAvailableVersionPrefersLocalSource(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("ARLOX_HOME", "")

	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, "cmd", "arlox"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "go.mod"), []byte("module github.com/sureshmopidevi/arlox\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "cmd", "arlox", "main.go"), []byte("package main\nfunc main(){}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(repo, "internal", "version"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "internal", "version", "VERSION"), []byte("9.9.9\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writeArloxSourceConfig(repo); err != nil {
		t.Fatal(err)
	}

	got, ok := latestAvailableVersion()
	if !ok {
		t.Fatal("expected source version")
	}
	if got != "9.9.9" {
		t.Fatalf("got %q, want 9.9.9", got)
	}
}

func TestUpgradeNoticeShownTracksLatest(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if upgradeNoticeShown("0.9.2") {
		t.Fatal("expected unseen notice")
	}
	if err := markUpgradeNoticeShown("0.9.2"); err != nil {
		t.Fatal(err)
	}
	if !upgradeNoticeShown("0.9.2") {
		t.Fatal("expected notice to be marked shown")
	}
	if upgradeNoticeShown("0.9.3") {
		t.Fatal("expected different version to show again")
	}
}

func TestSkipUpgradeNotice(t *testing.T) {
	cmd := &cobra.Command{Use: "upgrade"}
	if !skipUpgradeNotice(cmd) {
		t.Fatal("expected upgrade to skip notice")
	}
	cmd.Use = "create"
	if skipUpgradeNotice(cmd) {
		t.Fatal("expected create to show notice")
	}
}
