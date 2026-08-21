package workspace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

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
	if len(wf.Folders) != 3 {
		t.Fatalf("expected 3 folders, got %d", len(wf.Folders))
	}
}
