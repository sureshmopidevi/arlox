package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindInstalledBinaries(t *testing.T) {
	targets := findInstalledBinaries(true)
	_ = targets
}

func TestRemoveBinaries(t *testing.T) {
	tmp := t.TempDir()
	fakeBin := filepath.Join(tmp, "arlox")
	if err := os.WriteFile(fakeBin, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	removed := removeBinaries([]string{fakeBin})
	if removed != 1 {
		t.Fatalf("expected 1 binary removed, got %d", removed)
	}

	if _, err := os.Stat(fakeBin); !os.IsNotExist(err) {
		t.Fatalf("expected fake binary to be deleted")
	}
}
