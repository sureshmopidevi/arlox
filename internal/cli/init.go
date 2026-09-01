package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sureshmopidevi/arlox/internal/exec"
	"github.com/sureshmopidevi/arlox/internal/generate"
	"github.com/sureshmopidevi/arlox/internal/ui"
	"github.com/sureshmopidevi/arlox/internal/workspace"
)

func initCmd() *cobra.Command {
	var (
		nameFlag string
		force    bool
		deps     bool
		f        stackFlags
	)

	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Add arlox workspace orchestration to an existing project",
		Long: `Initialize arlox workspace files in an existing directory with one or more stacks.

Detects backend/, web/, and app/ via marker files (go.mod, package.json, pubspec.yaml).
Only workspace-level files are written — stack source trees are never overwritten.

Examples:
  arlox init              # current directory
  arlox init . --name myapp
  arlox init ~/Projects/legacy-api`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			pathArg := "."
			if len(args) > 0 {
				pathArg = args[0]
			}
			return runInit(cwd, pathArg, nameFlag, force, deps, f)
		},
	}

	cmd.Flags().StringVar(&nameFlag, "name", "", "workspace project name (default: infer from .code-workspace or directory)")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite locally modified root skills/rules")
	cmd.Flags().BoolVar(&deps, "deps", false, "install missing stack dependencies (go mod tidy, npm install, flutter pub get)")
	bindAddFlags(cmd, &f)
	return cmd
}

func runInit(cwd, pathArg, nameFlag string, force, installDeps bool, f stackFlags) error {
	root, err := resolveInitPath(cwd, pathArg)
	if err != nil {
		return err
	}

	stacks := workspace.ListPresentStacks(root)
	if len(stacks) == 0 {
		ui.Error("no stacks detected — expected backend/go.mod, web/package.json, or app/pubspec.yaml")
		return fmt.Errorf("no stacks found in %s", root)
	}

	name, err := inferInitName(root, nameFlag)
	if err != nil {
		return err
	}

	needsSkills := cursorTreeMissing(root)
	for _, s := range stacks {
		if cursorTreeMissing(filepath.Join(root, s)) {
			needsSkills = true
			break
		}
	}

	ui.Header("init", name)
	ui.WorkspaceInfo(root, filepath.Join(root, name+".code-workspace"))
	ui.Dim(fmt.Sprintf("detected stacks: %s", strings.Join(stacks, ", ")))
	fmt.Println()

	data := buildData(name, f)

	if err := workspace.EnsureWorkspace(root, name); err != nil {
		return err
	}
	if err := generate.WorkspaceRoot(root, data); err != nil {
		return err
	}

	wsFile := workspace.FindWorkspaceFile(root)
	if wsFile == "" {
		wsFile = filepath.Join(root, name+".code-workspace")
	}
	if err := workspace.SyncWorkspaceFolders(wsFile, root); err != nil {
		return err
	}

	if needsSkills {
		if err := generate.UpdateSkills(root, force); err != nil {
			return err
		}
		ui.Success("installed .cursor rules and skills")
	}

	if installDeps {
		installInitDeps(root, stacks)
	}

	ui.Summary(ui.SummaryOpts{
		Action:     "initialized",
		Name:       name,
		Root:       root,
		Stacks:     stacks,
		RelativeCD: root == cwd,
	})
	return nil
}

func resolveInitPath(cwd, pathArg string) (string, error) {
	target := cwd
	if pathArg != "" && pathArg != "." {
		if filepath.IsAbs(pathArg) {
			target = pathArg
		} else {
			target = filepath.Join(cwd, pathArg)
		}
	}
	abs, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", abs)
	}
	return abs, nil
}

func inferInitName(root, nameFlag string) (string, error) {
	if nameFlag != "" {
		if err := workspace.ValidateName(nameFlag); err != nil {
			return "", err
		}
		return nameFlag, nil
	}
	if ws := workspace.FindWorkspaceFile(root); ws != "" {
		name := strings.TrimSuffix(filepath.Base(ws), ".code-workspace")
		if err := workspace.ValidateName(name); err == nil {
			return name, nil
		}
	}
	name := filepath.Base(root)
	if err := workspace.ValidateName(name); err != nil {
		return "", fmt.Errorf("could not infer project name from %q — use --name", filepath.Base(root))
	}
	return name, nil
}

func cursorTreeMissing(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".cursor"))
	return os.IsNotExist(err)
}

func installInitDeps(root string, stacks []string) {
	for _, stackName := range stacks {
		stackDir := filepath.Join(root, stackName)
		switch workspace.Stack(stackName) {
		case workspace.Backend:
			if generate.HasTool("go") {
				ui.Dim("backend: running go mod tidy…")
				if err := exec.RunInDir(stackDir, "go", "mod", "tidy"); err != nil {
					ui.Warn("backend: go mod tidy: " + err.Error())
				}
			}
		case workspace.Web:
			if generate.HasTool("npm") {
				ui.Dim("web: running npm install…")
				if err := exec.RunInDir(stackDir, "npm", "install", "--no-audit", "--no-fund"); err != nil {
					ui.Warn("web: npm install: " + err.Error())
				}
			}
		case workspace.App:
			if generate.HasTool("flutter") {
				ui.Dim("app: running flutter pub get…")
				if err := exec.RunInDir(stackDir, "flutter", "pub", "get"); err != nil {
					ui.Warn("app: flutter pub get: " + err.Error())
				}
			}
		}
	}
}
