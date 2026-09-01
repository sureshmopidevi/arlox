package generate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCheckDrift_modifiedFile(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	sum, err := fileMD5(filepath.Join(dir, "main.go"))
	if err != nil {
		t.Fatal(err)
	}

	manifest := originManifest{
		Stack: "backend",
		Hashes: map[string]string{
			"main.go": sum,
		},
	}
	raw, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".origin-manifest.json"), raw, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main // changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := CheckDrift(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Modified) != 1 || report.Modified[0] != "main.go" {
		t.Fatalf("expected main.go modified, got %+v", report)
	}
}

func TestCheckDrift_missingFile(t *testing.T) {
	dir := t.TempDir()

	manifest := originManifest{
		Stack: "web",
		Hashes: map[string]string{
			"package.json": "abc123",
		},
	}
	raw, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".origin-manifest.json"), raw, 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := CheckDrift(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Missing) != 1 || report.Missing[0] != "package.json" {
		t.Fatalf("expected package.json missing, got %+v", report)
	}
}
