package generate

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sureshmopidevi/arlox/internal/workspace"
)

func TestDBNameFromProject(t *testing.T) {
	tests := map[string]string{
		"my-app":        "my_app",
		"issue_tracker": "issue_tracker",
		"a_b-c":         "a_b_c",
	}
	for in, want := range tests {
		if got := dbNameFromProject(in); got != want {
			t.Errorf("dbNameFromProject(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestCopyIfMissing(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "source.env")
	dst := filepath.Join(tmp, "sub", "target.env")

	if err := os.WriteFile(src, []byte("KEY=VALUE\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := copyIfMissing(src, dst); err != nil {
		t.Fatalf("copyIfMissing failed: %v", err)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "KEY=VALUE\n" {
		t.Fatalf("content mismatch: got %q", string(data))
	}

	// Idempotency check: should not overwrite existing dst
	if err := os.WriteFile(src, []byte("MODIFIED\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := copyIfMissing(src, dst); err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(dst)
	if string(data) != "KEY=VALUE\n" {
		t.Fatalf("expected existing file to be preserved, got %q", string(data))
	}
}

func TestFinalizeStackSeedsConfigsEvenIfToolsWarning(t *testing.T) {
	tmp := t.TempDir()
	backendDir := filepath.Join(tmp, "backend")
	cfgDir := filepath.Join(backendDir, "configs", "local")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	example := filepath.Join(cfgDir, "app.env.example")
	if err := os.WriteFile(example, []byte("PORT=8080\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// finalizeStack will attempt `go mod tidy` in an empty dir, which will produce a warning
	// but the config file MUST be created and no fatal error returned.
	warn, err := finalizeStack(workspace.GoGin, backendDir, Data{Name: "test"})
	if err != nil {
		t.Fatalf("expected non-fatal return, got err: %v", err)
	}
	_ = warn // warning is expected

	targetEnv := filepath.Join(cfgDir, "app.env")
	if _, err := os.Stat(targetEnv); os.IsNotExist(err) {
		t.Fatalf("expected %s to be created by finalizeStack", targetEnv)
	}
}

func TestMapTemplatePath(t *testing.T) {
	tests := map[string]string{
		"cursor/rules/karpathy.mdc":     ".cursor/rules/karpathy.mdc",
		"gitignore":                     ".gitignore",
		"env.example":                   ".env.example",
		"env.example.tmpl":              ".env.example.tmpl",
		"configs/local/app.env.example": "configs/local/app.env.example",
		"lib/main.dart.tmpl":            "lib/main.dart.tmpl",
	}
	for in, want := range tests {
		if got := mapTemplatePath(in); got != want {
			t.Errorf("mapTemplatePath(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestTemplateSourcePath(t *testing.T) {
	tests := map[workspace.Stack][]workspace.Variant{
		workspace.Backend: workspace.BackendVariants,
		workspace.Web:     workspace.WebVariants,
		workspace.App:     workspace.AppVariants,
	}
	for stack, variants := range tests {
		for _, variant := range variants {
			want := filepath.Join(string(stack), string(variant))
			if got := templateSourcePath(stack, variant); got != want {
				t.Errorf("templateSourcePath(%q, %q) = %q, want %q", stack, variant, got, want)
			}
		}
	}
}

func TestManifestRoundtrip(t *testing.T) {
	dir := t.TempDir()
	if err := writeManifest(dir, workspace.Web, workspace.ReactVite); err != nil {
		t.Fatal(err)
	}
	if got := manifestVariant(dir); got != workspace.ReactVite {
		t.Fatalf("manifestVariant() = %q, want %q", got, workspace.ReactVite)
	}
	data, err := os.ReadFile(filepath.Join(dir, ".origin-manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"variant": "react-vite"`) {
		t.Fatalf("manifest does not contain variant: %s", data)
	}
}

func TestRequiredToolsByVariant(t *testing.T) {
	tests := map[workspace.Variant][]string{
		workspace.GoGin:         {"go"},
		workspace.PyFastAPI:     {"uv"},
		workspace.NodeExpress:   {"node", "npm"},
		workspace.JavaSpring:    {"java", "mvn"},
		workspace.ReactVite:     {"node", "npm"},
		workspace.Flutter:       {"flutter"},
		workspace.ReactPWA:      {"node", "npm"},
		workspace.NativeIOS:     {"swift", "xcodebuild"},
		workspace.NativeAndroid: {"gradle"},
	}
	for variant, want := range tests {
		got := requiredTools(variant)
		if strings.Join(got, ",") != strings.Join(want, ",") {
			t.Errorf("requiredTools(%q) = %v, want %v", variant, got, want)
		}
	}
}

func TestFinalizeStackVariantRouting(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test uses Unix executable shims")
	}
	tests := map[workspace.Variant]string{
		workspace.GoGin:         "go",
		workspace.PyFastAPI:     "uv",
		workspace.NodeExpress:   "npm",
		workspace.NodeFastify:   "npm",
		workspace.JavaSpring:    "mvn",
		workspace.ReactVite:     "npm",
		workspace.NextJS:        "npm",
		workspace.VueVite:       "npm",
		workspace.SvelteVite:    "npm",
		workspace.Angular:       "npm",
		workspace.Nuxt:          "npm",
		workspace.Flutter:       "flutter",
		workspace.ReactPWA:      "npm",
		workspace.NativeIOS:     "swift",
		workspace.NativeAndroid: "gradle",
	}
	for variant, command := range tests {
		t.Run(string(variant), func(t *testing.T) {
			root := t.TempDir()
			bin := filepath.Join(root, "bin")
			if err := os.MkdirAll(bin, 0o755); err != nil {
				t.Fatal(err)
			}
			shim := filepath.Join(bin, command)
			if err := os.WriteFile(shim, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
				t.Fatal(err)
			}
			t.Setenv("PATH", bin)

			if variant == workspace.GoGin {
				cfg := filepath.Join(root, "configs", "local")
				if err := os.MkdirAll(cfg, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(cfg, "app.env.example"), nil, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			if variant == workspace.ReactVite {
				if err := os.WriteFile(filepath.Join(root, ".env.example"), nil, 0o644); err != nil {
					t.Fatal(err)
				}
			}

			if warning, err := finalizeStack(variant, root, Data{Name: "demo"}); err != nil || warning != "" {
				t.Fatalf("finalizeStack(%q) = warning %q, error %v", variant, warning, err)
			}
		})
	}
}
