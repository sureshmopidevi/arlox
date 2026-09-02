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

	"github.com/sureshmopidevi/arlox/internal/workspace"
	tmplfs "github.com/sureshmopidevi/arlox/templates"
)

// Data holds the template variables for all stacks.
type Data struct {
	Name         string
	DisplayName  string
	Module       string
	Package      string
	Org          string
	APIURL       string
	ArloxVersion string
	Variant      workspace.Variant
}

// HasTool reports whether name is available on PATH.
func HasTool(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// Stack generates the given stack under root using embedded templates.
// On fatal error during template scaffolding, the directory is cleaned up.
// Non-fatal post-scaffold issues (e.g. offline package installs) return a warning string.
func Stack(root string, stack workspace.Stack, variant workspace.Variant, data Data) (string, error) {
	if !workspace.ValidVariant(stack, variant) {
		return "", fmt.Errorf("variant %q is not valid for the %s stack", variant, stack)
	}
	for _, tool := range requiredTools(variant) {
		if !HasTool(tool) {
			return "", fmt.Errorf("%s not found in PATH (required for %s)", tool, variant)
		}
	}

	if !HasTool("git") {
		return "", fmt.Errorf("git not found in PATH")
	}

	data.Variant = variant
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
	if variant == workspace.Flutter {
		cmd := exec.Command("flutter", "create",
			"--org", data.Org,
			"--project-name", data.Package,
			".")
		cmd.Dir = stackDir
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", cleanup(fmt.Errorf("flutter create: %w\n%s", err, out))
		}
	}

	source := templateSourcePath(stack, variant)
	if err := renderTemplates(tmplfs.FS, source, stackDir, data); err != nil {
		return "", cleanup(err)
	}

	warn, err := finalizeStack(variant, stackDir, data)
	if err != nil {
		return "", cleanup(err)
	}

	if err := workspace.InitGit(stackDir); err != nil {
		return "", cleanup(fmt.Errorf("git init: %w", err))
	}

	if err := writeManifest(stackDir, stack, variant); err != nil {
		return "", cleanup(err)
	}

	return warn, nil
}

func templateSourcePath(stack workspace.Stack, variant workspace.Variant) string {
	return filepath.Join(string(stack), string(variant))
}

func requiredTools(variant workspace.Variant) []string {
	switch variant {
	case workspace.GoGin:
		return []string{"go"}
	case workspace.PyFastAPI:
		return []string{"uv"}
	case workspace.NodeExpress, workspace.NodeFastify,
		workspace.ReactVite, workspace.NextJS, workspace.VueVite,
		workspace.SvelteVite, workspace.Angular, workspace.Nuxt,
		workspace.ReactPWA:
		return []string{"node", "npm"}
	case workspace.JavaSpring:
		return []string{"java", "mvn"}
	case workspace.Flutter:
		return []string{"flutter"}
	case workspace.NativeIOS:
		return []string{"swift", "xcodebuild"}
	case workspace.NativeAndroid:
		return []string{"gradle"}
	default:
		return nil
	}
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
		variant, ok := workspace.DetectVariant(root, stack)
		if !ok {
			variant = workspace.DefaultVariant(stack)
		}
		if err := updateCursorFromTemplate(stackDir, templateSourcePath(stack, variant), "", force); err != nil {
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
	Stack   workspace.Stack   `json:"stack"`
	Variant workspace.Variant `json:"variant"`
	Hashes  map[string]string `json:"hashes"`
}

func writeManifest(dir string, stack workspace.Stack, variant workspace.Variant) error {
	source := templateSourcePath(stack, variant)
	m := originManifest{
		Stack:   stack,
		Variant: variant,
		Hashes:  make(map[string]string),
	}
	_ = fs.WalkDir(tmplfs.FS, source, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		raw, readErr := tmplfs.FS.ReadFile(path)
		if readErr != nil {
			return nil
		}
		sum := md5.Sum(raw)
		m.Hashes[path] = fmt.Sprintf("%x", sum)
		return nil
	})
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, ".origin-manifest.json"), data, 0o644)
}

func manifestVariant(dir string) workspace.Variant {
	data, err := os.ReadFile(filepath.Join(dir, ".origin-manifest.json"))
	if err != nil {
		return ""
	}
	var manifest originManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return ""
	}
	return manifest.Variant
}
