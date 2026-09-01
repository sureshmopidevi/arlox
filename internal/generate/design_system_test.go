package generate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/sureshmopidevi/arlox/internal/designsystems"
	tmplfs "github.com/sureshmopidevi/arlox/templates"
)

func TestWebDesignSystemOverlaysRender(t *testing.T) {
	for _, sys := range designsystems.WebCatalog {
		t.Run(sys.ID, func(t *testing.T) {
			root := t.TempDir()
			data := Data{
				Name:                 "dsdemo",
				DisplayName:          "Dsdemo",
				APIURL:               "http://localhost:8080/api/v1",
				ArloxVersion:         "0.15.0",
				WebDesignSystem:      sys.ID,
				WebDesignSystemLabel: sys.Label,
			}
			webDir := filepath.Join(root, "web")
			if err := os.MkdirAll(webDir, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := renderTemplates(tmplfs.FS, "web", webDir, data); err != nil {
				t.Fatalf("web core: %v", err)
			}
			overlay := designsystems.WebOverlayPath(sys.ID)
			if err := renderTemplates(tmplfs.FS, overlay, webDir, data); err != nil {
				t.Fatalf("overlay: %v", err)
			}
			if err := writeDesignSystemManifest(webDir, sys.ID); err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(filepath.Join(webDir, "package.json")); err != nil {
				t.Fatal("package.json missing")
			}
			raw, err := os.ReadFile(filepath.Join(webDir, ".arlox", "design-system.json"))
			if err != nil {
				t.Fatal(err)
			}
			var manifest designSystemManifest
			if err := json.Unmarshal(raw, &manifest); err != nil {
				t.Fatal(err)
			}
			if manifest.ID != sys.ID {
				t.Fatalf("manifest id %q want %q", manifest.ID, sys.ID)
			}
		})
	}
}
