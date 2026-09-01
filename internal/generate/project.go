package generate

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sureshmopidevi/arlox/internal/naming"
	"github.com/sureshmopidevi/arlox/internal/version"
	"github.com/sureshmopidevi/arlox/internal/workspace"
)

type projectManifest struct {
	Name        string                 `json:"name"`
	Kebab       string                 `json:"kebab"`
	Snake       string                 `json:"snake"`
	DisplayName string                 `json:"displayName"`
	ArloxVersion string                `json:"arloxVersion"`
	Stacks      projectManifestStacks  `json:"stacks"`
}

type projectManifestStacks struct {
	Backend *projectBackendStack `json:"backend,omitempty"`
	Web     *projectWebStack     `json:"web,omitempty"`
	App     *projectAppStack     `json:"app,omitempty"`
}

type projectBackendStack struct {
	Module string `json:"module"`
	DBName string `json:"dbName"`
}

type projectWebStack struct {
	PackageName  string `json:"packageName"`
	DesignSystem string `json:"designSystem,omitempty"`
}

type projectAppStack struct {
	Package string `json:"package"`
	Org     string `json:"org"`
}

// WriteProjectManifest writes or updates .arlox/project.json at the workspace root.
func WriteProjectManifest(root string, data Data, present []workspace.Stack) error {
	n := naming.FromSlug(data.Name, data.Module, data.Org)
	// Prefer explicit Data fields when populated (buildData sets them).
	if data.KebabName != "" {
		n.Kebab = data.KebabName
		n.WebPackage = data.WebPackageName
	}
	if data.SnakeName != "" {
		n.Snake = data.SnakeName
		n.FlutterPkg = data.SnakeName
		n.DBName = data.DBName
	}
	if data.Module != "" {
		n.Module = data.Module
	}

	m := projectManifest{
		Name:         n.Name,
		Kebab:        n.Kebab,
		Snake:        n.Snake,
		DisplayName:  data.DisplayName,
		ArloxVersion: version.Version,
	}

	for _, s := range present {
		switch s {
		case workspace.Backend:
			m.Stacks.Backend = &projectBackendStack{
				Module: n.Module,
				DBName: n.DBName,
			}
		case workspace.Web:
			m.Stacks.Web = &projectWebStack{
				PackageName:  n.WebPackage,
				DesignSystem: data.WebDesignSystem,
			}
		case workspace.App:
			m.Stacks.App = &projectAppStack{
				Package: n.FlutterPkg,
				Org:     data.Org,
			}
		}
	}

	raw, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Join(root, ".arlox")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "project.json"), append(raw, '\n'), 0o644)
}
