package generate

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DriftReport lists files that differ from the scaffold origin manifest.
type DriftReport struct {
	Missing  []string
	Modified []string
}

// CheckDrift compares on-disk files against .origin-manifest.json in stackDir.
func CheckDrift(stackDir string) (DriftReport, error) {
	manifestPath := filepath.Join(stackDir, ".origin-manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return DriftReport{}, err
	}

	var m originManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return DriftReport{}, err
	}

	report := DriftReport{}
	for rel, want := range m.Hashes {
		path := filepath.Join(stackDir, rel)
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				report.Missing = append(report.Missing, rel)
				continue
			}
			return DriftReport{}, err
		}
		if info.IsDir() {
			continue
		}
		got, err := fileMD5(path)
		if err != nil {
			return DriftReport{}, err
		}
		if got != want {
			report.Modified = append(report.Modified, rel)
		}
	}
	return report, nil
}

func fileMD5(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := md5.Sum(raw)
	return fmt.Sprintf("%x", sum), nil
}

func shouldSkipManifestPath(rel string) bool {
	parts := strings.Split(rel, string(filepath.Separator))
	for _, p := range parts {
		switch p {
		case ".git", "node_modules", ".dart_tool", "bin", "build", ".idea", ".vscode", "ios", "android", "macos", "linux", "windows", "web", "ephemeral":
			return true
		}
	}
	base := filepath.Base(rel)
	if base == ".origin-manifest.json" || base == "go.sum" {
		return true
	}
	return false
}
