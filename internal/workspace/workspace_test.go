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
	exclude, ok := wf.Folders[0].Settings["files.exclude"].(map[string]any)
	if !ok || exclude["backend"] != true || exclude["web"] != true || exclude["app"] != true {
		t.Fatalf("expected root folder to hide stack dirs, got settings=%+v", wf.Folders[0].Settings)
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
		marker := stackMarkerCandidates[Stack(stack)][0]
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
	for _, f := range wf.Folders {
		if f.Path != "." {
			continue
		}
		exclude, ok := f.Settings["files.exclude"].(map[string]any)
		if !ok || exclude["backend"] != true {
			t.Fatalf("expected synced root folder to hide stacks, got settings=%+v", f.Settings)
		}
	}
}

func TestDetectWorkspaceWalksUp(t *testing.T) {
	root := t.TempDir()
	wsPath := filepath.Join(root, "demo.code-workspace")
	if err := WriteWorkspaceFile(wsPath, "demo", []string{"."}); err != nil {
		t.Fatal(err)
	}
	stackDir := filepath.Join(root, "backend")
	if err := os.MkdirAll(stackDir, 0o755); err != nil {
		t.Fatal(err)
	}

	found, ok := DetectWorkspace(stackDir)
	if !ok || found != root {
		t.Fatalf("DetectWorkspace from stack dir: got (%q, %v), want (%q, true)", found, ok, root)
	}

	_, ok = DetectWorkspace(t.TempDir())
	if ok {
		t.Fatal("expected no workspace in empty temp dir")
	}
}

func TestValidateName(t *testing.T) {
	valid := []string{"myapp", "my-app", "my_app", "app2", "a_b-c"}
	for _, name := range valid {
		if err := ValidateName(name); err != nil {
			t.Errorf("ValidateName(%q): unexpected error: %v", name, err)
		}
	}
	invalid := []string{"", "MyApp", "my app", "my.app", "my@app"}
	for _, name := range invalid {
		if err := ValidateName(name); err == nil {
			t.Errorf("ValidateName(%q): expected error", name)
		}
	}
}

func TestFindWorkspaceFilePrefersMatchingName(t *testing.T) {
	root := t.TempDir()
	other := filepath.Join(root, "other.code-workspace")
	preferred := filepath.Join(root, "demo.code-workspace")
	for _, p := range []string{other, preferred} {
		if err := os.WriteFile(p, []byte("{}"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	got := FindWorkspaceFile(root)
	if got != preferred {
		t.Fatalf("FindWorkspaceFile: got %q, want %q", got, preferred)
	}
}

func TestVariantConstants(t *testing.T) {
	all := append(append(append([]Variant{}, BackendVariants...), WebVariants...), AppVariants...)
	seen := make(map[Variant]bool, len(all))
	for _, variant := range all {
		if variant == "" {
			t.Fatal("variant constants must not be empty")
		}
		if seen[variant] {
			t.Fatalf("variant %q appears more than once", variant)
		}
		seen[variant] = true
	}
	if len(all) != 15 {
		t.Fatalf("expected 15 variants, got %d", len(all))
	}
}

func TestVariantsByStack(t *testing.T) {
	tests := map[Stack][]Variant{
		Backend: {GoGin, PyFastAPI, NodeExpress, NodeFastify, JavaSpring},
		Web:     {ReactVite, NextJS, VueVite, SvelteVite, Angular, Nuxt},
		App:     {Flutter, ReactPWA, NativeIOS, NativeAndroid},
	}
	for stack, want := range tests {
		got := VariantsForStack(stack)
		if len(got) != len(want) {
			t.Fatalf("%s: got %d variants, want %d", stack, len(got), len(want))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("%s variant %d = %q, want %q", stack, i, got[i], want[i])
			}
		}
	}
}

func TestStackExistsMultiMarker(t *testing.T) {
	for stack, markers := range stackMarkerCandidates {
		for _, marker := range markers {
			t.Run(string(stack)+"/"+marker, func(t *testing.T) {
				root := t.TempDir()
				dir := filepath.Join(root, string(stack))
				if err := os.MkdirAll(dir, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(dir, marker), nil, 0o644); err != nil {
					t.Fatal(err)
				}
				if !StackExists(root, stack) {
					t.Fatalf("StackExists(%s) = false with marker %s", stack, marker)
				}
			})
		}
	}
}

func TestDetectVariant(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, string(Backend))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := []byte(`{"stack":"backend","variant":"python"}`)
	if err := os.WriteFile(filepath.Join(dir, ".origin-manifest.json"), manifest, 0o644); err != nil {
		t.Fatal(err)
	}
	got, ok := DetectVariant(root, Backend)
	if !ok || got != PyFastAPI {
		t.Fatalf("DetectVariant() = (%q, %v), want (%q, true)", got, ok, PyFastAPI)
	}
	if _, ok := DetectVariant(root, Web); ok {
		t.Fatal("expected no web variant")
	}
}
