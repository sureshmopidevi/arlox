package generate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/sureshmopidevi/arlox/internal/workspace"
)

func TestWriteProjectManifest(t *testing.T) {
	root := t.TempDir()
	data := Data{
		Name:            "my-cool_app",
		DisplayName:     "My Cool App",
		Module:          "github.com/example/my-cool-app-backend",
		Org:             "com.example",
		KebabName:       "my-cool-app",
		SnakeName:       "my_cool_app",
		WebPackageName:  "my-cool-app",
		DBName:          "my_cool_app",
		WebDesignSystem: "tailwind",
		ArloxVersion:    "0.16.0",
	}
	if err := WriteProjectManifest(root, data, []workspace.Stack{workspace.Backend, workspace.Web}); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(root, ".arlox", "project.json"))
	if err != nil {
		t.Fatal(err)
	}
	var m projectManifest
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	if m.Kebab != "my-cool-app" || m.Snake != "my_cool_app" {
		t.Fatalf("manifest names: kebab=%q snake=%q", m.Kebab, m.Snake)
	}
	if m.Stacks.Backend == nil || m.Stacks.Backend.DBName != "my_cool_app" {
		t.Fatal("backend stack missing or wrong dbName")
	}
	if m.Stacks.Web == nil || m.Stacks.Web.PackageName != "my-cool-app" {
		t.Fatal("web stack missing or wrong packageName")
	}
}
