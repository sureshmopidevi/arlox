package workspace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteWorkspaceFileNeverNullFolders(t *testing.T) {
	root := t.TempDir()
	wsPath := filepath.Join(root, "demo.code-workspace")

	if err := WriteWorkspaceFile(wsPath, "demo", []string{"."}); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(wsPath)
	if err != nil {
		t.Fatal(err)
	}
	var top map[string]any
	if err := json.Unmarshal(raw, &top); err != nil {
		t.Fatal(err)
	}
	folders, ok := top["folders"].([]any)
	if !ok {
		t.Fatalf("folders must be a JSON array, got %T (%s)", top["folders"], raw)
	}
	if len(folders) != 1 {
		t.Fatalf("expected 1 folder, got %d", len(folders))
	}

	var wf workspaceFile
	if err := json.Unmarshal(raw, &wf); err != nil {
		t.Fatal(err)
	}
	if wf.Folders[0].Path != "." || wf.Folders[0].Name != "demo" {
		t.Fatalf("expected root folder name=demo path=., got %+v", wf.Folders[0])
	}
}

func TestSyncWorkspaceFolders(t *testing.T) {
	root := t.TempDir()
	wsPath := filepath.Join(root, "demo.code-workspace")

	if err := WriteWorkspaceFile(wsPath, "demo", []string{"backend"}); err != nil {
		t.Fatal(err)
	}
	for _, stack := range []string{"backend", "web", "app"} {
		dir := filepath.Join(root, stack)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		marker := stackMarkers[Stack(stack)]
		if err := os.WriteFile(filepath.Join(dir, marker), []byte("ok"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if err := SyncWorkspaceFolders(wsPath, root); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(wsPath)
	if err != nil {
		t.Fatal(err)
	}
	var wf workspaceFile
	if err := json.Unmarshal(data, &wf); err != nil {
		t.Fatal(err)
	}
	// root "." + backend + web + app
	if len(wf.Folders) != 4 {
		t.Fatalf("expected 4 folders, got %d: %+v", len(wf.Folders), wf.Folders)
	}

	paths := make(map[string]string)
	for _, f := range wf.Folders {
		paths[f.Path] = f.Name
	}
	if paths["."] != "demo" {
		t.Fatalf("expected root folder named demo, got %q", paths["."])
	}
	want := map[string]string{
		"backend": "demo-backend",
		"web":     "demo-web",
		"app":     "demo-app",
	}
	for path, name := range want {
		if paths[path] != name {
			t.Fatalf("folder %s: want name %q, got %q", path, name, paths[path])
		}
	}
}
