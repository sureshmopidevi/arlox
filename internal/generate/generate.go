package generate

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sureshmopidevi/arlox/internal/designsystems"
	arloxexec "github.com/sureshmopidevi/arlox/internal/exec"
	tmplfs "github.com/sureshmopidevi/arlox/templates"
	"github.com/sureshmopidevi/arlox/internal/ui"
	"github.com/sureshmopidevi/arlox/internal/version"
	"github.com/sureshmopidevi/arlox/internal/workspace"
)

// Data holds the template variables for all stacks.
type Data struct {
	Name                  string
	DisplayName           string
	Module                string
	Package               string
	Org                   string
	APIURL                string
	ArloxVersion          string
	PostgresPort          int
	WebDesignSystem       string
	WebDesignSystemLabel  string
}

// HasTool reports whether name is available on PATH.
func HasTool(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// Stack generates the given stack under root using embedded templates.
// On fatal error during template scaffolding, the directory is cleaned up.
// Non-fatal post-scaffold issues (e.g. offline package installs) return a warning string.
func Stack(root string, stack workspace.Stack, data Data) (string, error) {
	switch stack {
	case workspace.Backend:
		if !HasTool("go") {
			return "", fmt.Errorf("go not found in PATH")
		}
	case workspace.Web:
		if !HasTool("node") || !HasTool("npm") {
			return "", fmt.Errorf("node and npm are required for the web stack")
		}
	case workspace.App:
		if !HasTool("flutter") {
			return "", fmt.Errorf("flutter not found in PATH")
		}
	}

	if !HasTool("git") {
		return "", fmt.Errorf("git not found in PATH")
	}

	stackName := string(stack)
	stackDir := filepath.Join(root, stackName)

	created := false
	if _, err := os.Stat(stackDir); os.IsNotExist(err) {
		if err := os.MkdirAll(stackDir, 0o755); err != nil {
			return "", err
		}
		created = true
	}

	cleanup := func(err error) error {
		if created {
			os.RemoveAll(stackDir)
		}
		return err
	}

	// For Flutter: run flutter create first, then overlay our templates.
	if stack == workspace.App {
		ui.Dim("app: running flutter create…")
		if err := arloxexec.RunInDir(stackDir, "flutter", "create",
			"--org", data.Org,
			"--project-name", data.Package,
			"."); err != nil {
			return "", cleanup(fmt.Errorf("flutter create: %w", err))
		}
	}

	if err := renderTemplates(tmplfs.FS, stackName, stackDir, data); err != nil {
		return "", cleanup(err)
	}

	if stack == workspace.Web {
		ds := data.WebDesignSystem
		if ds == "" {
			ds = designsystems.DefaultWebID
		}
		overlay := designsystems.WebOverlayPath(ds)
		if err := renderTemplates(tmplfs.FS, overlay, stackDir, data); err != nil {
			return "", cleanup(fmt.Errorf("design system overlay %s: %w", ds, err))
		}
		if err := writeDesignSystemManifest(stackDir, ds); err != nil {
			return "", cleanup(err)
		}
	}

	warn, err := finalizeStack(stack, stackDir, data)
	if err != nil {
		return "", cleanup(err)
	}

	if err := workspace.InitGit(stackDir); err != nil {
		return "", cleanup(fmt.Errorf("git init: %w", err))
	}

	if err := writeManifest(stackDir, stackName); err != nil {
		return "", cleanup(err)
	}

	return warn, nil
}

// WorkspaceRoot renders workspace-level templates (README, .cursor, etc.) into root.
func WorkspaceRoot(root string, data Data) error {
	return renderTemplates(tmplfs.FS, "workspace", root, data)
}

// UpdateSkills refreshes .cursor skills/rules from templates for the workspace
// root and each present stack. Files under learned/ and locally modified files
// are skipped unless force is true.
func UpdateSkills(root string, force bool) error {
	if err := updateCursorFromTemplate(root, "workspace", "", force); err != nil {
		return err
	}
	for _, stack := range []workspace.Stack{workspace.Backend, workspace.Web, workspace.App} {
		if !workspace.StackExists(root, stack) {
			continue
		}
		stackDir := filepath.Join(root, string(stack))
		if err := updateCursorFromTemplate(stackDir, string(stack), "", force); err != nil {
			return err
		}
	}
	return nil
}

func updateCursorFromTemplate(destRoot, srcDir, _ string, force bool) error {
	return fs.WalkDir(tmplfs.FS, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.Contains(path, "cursor") {
			return nil
		}

		rel, _ := filepath.Rel(srcDir, path)
		destPath := strings.TrimSuffix(filepath.Join(destRoot, mapTemplatePath(rel)), ".tmpl")

		if strings.Contains(destPath, filepath.Join("learned")) {
			return nil
		}

		if !force && locallyModified(destPath, path) {
			return nil
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
			return err
		}
		content, err := tmplfs.FS.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(destPath, content, 0o644)
	})
}

// renderTemplates walks the embedded FS under srcDir and writes rendered files to destDir.
// Files with a .tmpl suffix are parsed through Go's text/template engine;
// all other files are copied verbatim.
func renderTemplates(fsys fs.FS, srcDir, destDir string, data Data) error {
	return fs.WalkDir(fsys, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		rel, _ := filepath.Rel(srcDir, path)
		destRel := mapTemplatePath(rel)
		dest := strings.TrimSuffix(filepath.Join(destDir, destRel), ".tmpl")

		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}

		raw, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}

		// Only .tmpl files are rendered through the template engine.
		if !strings.HasSuffix(path, ".tmpl") {
			return os.WriteFile(dest, raw, 0o644)
		}

		tmpl, err := template.New(path).Delims("[[", "]]").Parse(string(raw))
		if err != nil {
			return err
		}

		f, err := os.Create(dest)
		if err != nil {
			return err
		}
		if err := tmpl.Execute(f, data); err != nil {
			f.Close()
			return err
		}
		return f.Close()
	})
}

// mapTemplatePath remaps embed-friendly names to dotfiles/dirs on disk.
func mapTemplatePath(rel string) string {
	parts := strings.Split(rel, string(filepath.Separator))
	for i, p := range parts {
		switch {
		case p == "cursor":
			parts[i] = ".cursor"
		case p == "gitignore":
			parts[i] = ".gitignore"
		case strings.HasPrefix(p, "env.example"):
			parts[i] = "." + p
		}
	}
	return filepath.Join(parts...)
}

// locallyModified reports whether the local file differs from the embedded template.
func locallyModified(localPath, templatePath string) bool {
	local, err := os.ReadFile(localPath)
	if err != nil {
		return false
	}
	tmpl, err := tmplfs.FS.ReadFile(templatePath)
	if err != nil {
		return false
	}
	return md5.Sum(local) != md5.Sum(tmpl)
}

type originManifest struct {
	Stack        string            `json:"stack"`
	ArloxVersion string            `json:"arloxVersion"`
	Hashes       map[string]string `json:"hashes"`
}

func writeManifest(dir, stack string) error {
	m := originManifest{
		Stack:        stack,
		ArloxVersion: version.Version,
		Hashes:       make(map[string]string),
	}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		if info.IsDir() || !info.Mode().IsRegular() {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		if shouldSkipManifestPath(rel) {
			return nil
		}
		sum, err := fileMD5(path)
		if err != nil {
			return err
		}
		m.Hashes[rel] = sum
		return nil
	})
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, ".origin-manifest.json"), data, 0o644)
}
