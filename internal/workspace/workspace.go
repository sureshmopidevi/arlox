package workspace

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Stack identifies a project stack.
type Stack string

const (
	Backend Stack = "backend"
	Web     Stack = "web"
	App     Stack = "app"
)

var allStacks = []Stack{Backend, Web, App}

var stackMarkers = map[Stack]string{
	Backend: "go.mod",
	Web:     "package.json",
	App:     "pubspec.yaml",
}

// workspaceFile mirrors the .code-workspace JSON format.
type workspaceFile struct {
	Folders    []workspaceFolder `json:"folders"`
	Settings   map[string]any    `json:"settings,omitempty"`
	Extensions *workspaceExts    `json:"extensions,omitempty"`
}

type workspaceFolder struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type workspaceExts struct {
	Recommendations []string `json:"recommendations"`
}

var nameRe = regexp.MustCompile(`^[a-z0-9_-]+$`)

// ValidateName returns an error if name is not lowercase [a-z0-9_-]+.
func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if !nameRe.MatchString(name) {
		return fmt.Errorf("name must contain only lowercase letters, numbers, hyphens, and underscores")
	}
	return nil
}

// DetectWorkspace returns the workspace root if a *.code-workspace file exists
// in cwd or any parent directory.
func DetectWorkspace(cwd string) (string, bool) {
	dir := cwd
	for {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return "", false
		}
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".code-workspace") {
				return dir, true
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", false
}

// FindWorkspaceFile returns the path to the *.code-workspace file in root, or "".
// Prefers {basename(root)}.code-workspace, otherwise the first match alphabetically.
func FindWorkspaceFile(root string) string {
	entries, err := os.ReadDir(root)
	if err != nil {
		return ""
	}
	preferred := filepath.Join(root, filepath.Base(root)+".code-workspace")
	if _, err := os.Stat(preferred); err == nil {
		return preferred
	}
	var found string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".code-workspace") {
			continue
		}
		candidate := filepath.Join(root, e.Name())
		if found == "" || e.Name() < filepath.Base(found) {
			found = candidate
		}
	}
	return found
}

// StackExists reports whether the stack directory contains its marker file.
func StackExists(root string, stack Stack) bool {
	marker := stackMarkers[stack]
	if marker == "" {
		return false
	}
	_, err := os.Stat(filepath.Join(root, string(stack), marker))
	return err == nil
}

// StackPresentButEmpty reports whether the stack dir exists but has no marker file.
func StackPresentButEmpty(root string, stack Stack) bool {
	dir := filepath.Join(root, string(stack))
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return false
	}
	return !StackExists(root, stack)
}

// EnsureWorkspace creates the workspace dir and writes an initial .code-workspace if missing.
func EnsureWorkspace(root, name string) error {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return err
	}
	wsFile := filepath.Join(root, name+".code-workspace")
	if _, err := os.Stat(wsFile); os.IsNotExist(err) {
		return WriteWorkspaceFile(wsFile, name, []string{"."})
	}
	return nil
}

// WriteWorkspaceFile writes a fresh .code-workspace JSON file.
// The project name is used as the display name for path "." (workspace root).
// Stack folders are labeled "{name}-backend", "{name}-web", "{name}-app".
func WriteWorkspaceFile(path, name string, folders []string) error {
	wf := workspaceFile{
		Folders: []workspaceFolder{},
		Settings: map[string]any{
			"files.exclude": map[string]any{"**/.dart_tool": true},
		},
		Extensions: &workspaceExts{
			Recommendations: []string{
				"golang.go",
				"dart-code.flutter",
				"dbaeumer.vscode-eslint",
			},
		},
	}
	for _, f := range folders {
		wf.Folders = append(wf.Folders, workspaceFolder{
			Name: FolderDisplayName(name, f),
			Path: f,
		})
	}
	data, err := json.MarshalIndent(wf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// FolderDisplayName returns the IDE sidebar label for a workspace folder path.
func FolderDisplayName(projectName, path string) string {
	if path == "." || path == "" {
		return projectName
	}
	switch path {
	case "backend", "web", "app":
		return projectName + "-" + path
	default:
		return path
	}
}

// ReadWorkspaceFolders returns the folder paths listed in the workspace file.
func ReadWorkspaceFolders(wsPath string) []string {
	data, err := os.ReadFile(wsPath)
	if err != nil {
		return nil
	}
	var wf workspaceFile
	if err := json.Unmarshal(data, &wf); err != nil {
		return nil
	}
	folders := make([]string, 0, len(wf.Folders))
	for _, f := range wf.Folders {
		folders = append(folders, f.Path)
	}
	return folders
}

// ListPresentStacks returns stack folder names that exist under root.
func ListPresentStacks(root string) []string {
	var present []string
	for _, s := range allStacks {
		if StackExists(root, s) {
			present = append(present, string(s))
		}
	}
	return present
}

// SyncWorkspaceFolders ensures the root "." folder and every present stack are listed.
func SyncWorkspaceFolders(wsPath, root string) error {
	folders := append([]string{"."}, ListPresentStacks(root)...)
	return MergeWorkspaceFolders(wsPath, folders)
}

// MergeWorkspaceFolders adds any new folders into the workspace file (union)
// and refreshes display names so stacks are labeled "{project}-backend", etc.
func MergeWorkspaceFolders(wsPath string, newFolders []string) error {
	data, err := os.ReadFile(wsPath)
	if err != nil {
		return err
	}
	var wf workspaceFile
	if err := json.Unmarshal(data, &wf); err != nil {
		return err
	}
	if wf.Folders == nil {
		wf.Folders = []workspaceFolder{}
	}

	rootName := strings.TrimSuffix(filepath.Base(wsPath), ".code-workspace")

	seen := make(map[string]bool)
	for i, f := range wf.Folders {
		seen[f.Path] = true
		wf.Folders[i].Name = FolderDisplayName(rootName, f.Path)
	}
	for _, f := range newFolders {
		if !seen[f] {
			wf.Folders = append(wf.Folders, workspaceFolder{
				Name: FolderDisplayName(rootName, f),
				Path: f,
			})
			seen[f] = true
		}
	}

	out, err := json.MarshalIndent(wf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(wsPath, out, 0o644)
}

// IsUnrelatedDir reports whether path is a non-empty directory without a .code-workspace.
func IsUnrelatedDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return false
	}
	entries, err := os.ReadDir(path)
	if err != nil || len(entries) == 0 {
		return false
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".code-workspace") {
			return false
		}
	}
	return true
}

// InitGit runs git init in dir if no .git directory is present.
func InitGit(dir string) error {
	if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
		return nil
	}
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	return cmd.Run()
}

// ListMissingStacks returns stacks that don't yet have their marker file.
func ListMissingStacks(root string) []Stack {
	var missing []Stack
	for _, s := range allStacks {
		if !StackExists(root, s) {
			missing = append(missing, s)
		}
	}
	return missing
}
