package generate

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sureshmopidevi/arlox/internal/designsystems"
	"github.com/sureshmopidevi/arlox/internal/version"
)

type designSystemManifest struct {
	ID           string `json:"id"`
	Label        string `json:"label"`
	ArloxVersion string `json:"arloxVersion"`
}

func writeDesignSystemManifest(webDir, id string) error {
	manifest := designSystemManifest{
		ID:           id,
		Label:        designsystems.WebLabel(id),
		ArloxVersion: version.Version,
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Join(webDir, ".arlox")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "design-system.json"), append(data, '\n'), 0o644)
}
