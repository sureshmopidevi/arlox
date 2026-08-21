package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/lt-sureshmopidevi/vibeit/internal/generate"
	"github.com/lt-sureshmopidevi/vibeit/internal/ui"
	"github.com/lt-sureshmopidevi/vibeit/internal/version"
	"github.com/lt-sureshmopidevi/vibeit/internal/workspace"
)

// Execute builds and runs the root cobra command.
func Execute() {
	if err := buildRoot().Execute(); err != nil {
		os.Exit(1)
	}
}

func buildRoot() *cobra.Command {
	var f stackFlags
	create := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new vibeit workspace",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(args, f)
		},
	}
	bindStackFlags(create, &f)

	root := &cobra.Command{
		Use:     "vibeit",
		Short:   "Scaffold and manage multi-stack workspaces",
		Long:    "vibeit create workspaces with backend, web, and/or Flutter stacks plus Cursor rules and skills.",
		Version: version.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	root.SetVersionTemplate("vibeit {{.Version}}\n")
	root.AddCommand(create)
	root.AddCommand(addCmd())
	root.AddCommand(skillsCmd())
	root.AddCommand(versionCmd())
	root.AddCommand(upgradeCmd())
	return root
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print vibeit version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("vibeit " + version.Version)
		},
	}
}

type stackFlags struct {
	backend bool
	web     bool
	app     bool
	open    bool
	module  string
	org     string
	out     string
}

func bindStackFlags(cmd *cobra.Command, f *stackFlags) {
	cmd.Flags().BoolVar(&f.backend, "backend", false, "include backend (Go) stack")
	cmd.Flags().BoolVar(&f.web, "web", false, "include web (React/TS) stack")
	cmd.Flags().BoolVar(&f.app, "app", false, "include app (Flutter) stack")
	cmd.Flags().BoolVar(&f.open, "open", false, "open workspace in Cursor, VS Code, or Antigravity IDE after creation")
	cmd.Flags().StringVar(&f.module, "module", "", "Go module path (default: github.com/example/<name>-backend)")
	cmd.Flags().StringVar(&f.org, "org", "com.example", "Flutter org identifier")
	cmd.Flags().StringVar(&f.out, "out", "", "parent directory for workspace (default: current directory)")
}

func selectedStacks(f stackFlags) []workspace.Stack {
	var stacks []workspace.Stack
	if f.backend {
		stacks = append(stacks, workspace.Backend)
	}
	if f.web {
		stacks = append(stacks, workspace.Web)
	}
	if f.app {
		stacks = append(stacks, workspace.App)
	}
	return stacks
}

func buildData(name string, f stackFlags) generate.Data {
	module := f.module
	if module == "" {
		module = "github.com/example/" + name + "-backend"
	}
	// Flutter package names must be snake_case (no hyphens); match project name like web's package.json.
	pkg := strings.ReplaceAll(name, "-", "_")
	return generate.Data{
		Name:        name,
		DisplayName: toDisplayName(name),
		Module:      module,
		Package:     pkg,
		Org:         f.org,
		APIURL:      "http://localhost:8080/api/v1",
	}
}

func toDisplayName(name string) string {
	parts := strings.Split(name, "-")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}

func promptName() (string, error) {
	var name string
	err := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Project name").
			Description("lowercase letters, numbers, and hyphens").
			Validate(workspace.ValidateName).
			Value(&name),
	)).Run()
	return name, err
}

func promptStacks(available []workspace.Stack) ([]workspace.Stack, error) {
	if len(available) == 0 {
		return nil, nil
	}

	labels := map[workspace.Stack]string{
		workspace.Backend: "backend   Go + Gin API",
		workspace.Web:     "web       Vite + React + Tailwind",
		workspace.App:     "app       Flutter (Material home)",
	}

	opts := make([]huh.Option[string], 0, len(available))
	for _, s := range available {
		opts = append(opts, huh.NewOption(labels[s], string(s)))
	}

	var selected []string
	err := huh.NewForm(huh.NewGroup(
		huh.NewMultiSelect[string]().
			Title("What should vibeit generate?").
			Description("space toggle · enter confirm").
			Options(opts...).
			Value(&selected),
	)).Run()
	if err != nil {
		return nil, err
	}

	stacks := make([]workspace.Stack, 0, len(selected))
	for _, s := range selected {
		stacks = append(stacks, workspace.Stack(s))
	}
	return stacks, nil
}

type stackResult struct {
	created []string
	skipped []string
	failed  []string
}

func runStacks(root string, stacks []workspace.Stack, data generate.Data) stackResult {
	var r stackResult
	for _, s := range stacks {
		if workspace.StackExists(root, s) {
			ui.Progress(string(s), "skipped")
			r.skipped = append(r.skipped, string(s))
			continue
		}
		ui.Progress(string(s), "generating")
		if err := generate.Stack(root, s, data); err != nil {
			ui.Progress(string(s), "failed")
			ui.Error(err.Error())
			r.failed = append(r.failed, string(s))
			continue
		}
		if !workspace.StackExists(root, s) {
			ui.Progress(string(s), "failed")
			ui.Error(fmt.Sprintf("%s finished but stack files are missing — try again or run vibeit add --%s", s, s))
			r.failed = append(r.failed, string(s))
			continue
		}
		ui.Progress(string(s), "done")
		r.created = append(r.created, string(s))
	}
	return r
}

// editorLaunch describes how to open a .code-workspace in a given CLI.
type editorLaunch struct {
	bin  string
	args []string // extra args before the workspace path (e.g. --classic for Cursor IDE)
}

// knownEditors are VS Code–compatible CLIs that can open .code-workspace files.
// Order is preference for --open auto-launch.
// Cursor needs --classic so it opens the IDE editor instead of the Agents window.
var knownEditors = []editorLaunch{
	{bin: "cursor", args: []string{"--classic"}},
	{bin: "code"},
	{bin: "agy-ide"},
	{bin: "antigravity-ide"},
}

func tryOpen(wsFile string) string {
	for _, ed := range knownEditors {
		path, err := exec.LookPath(ed.bin)
		if err != nil {
			continue
		}
		args := append(append([]string{}, ed.args...), wsFile)
		cmd := exec.Command(path, args...)
		_ = cmd.Start()
		display := ed.bin
		if len(ed.args) > 0 {
			display = ed.bin + " " + strings.Join(ed.args, " ")
		}
		return display + " " + filepath.Base(wsFile)
	}
	return ""
}

func runCreate(args []string, f stackFlags) error {
	name := ""
	if len(args) > 0 {
		name = args[0]
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if name == "" {
		name, err = promptName()
		if err != nil {
			return err
		}
	} else if err := workspace.ValidateName(name); err != nil {
		return err
	}

	target := filepath.Join(resolveOutDir(cwd, f.out), name)

	if workspace.IsUnrelatedDir(target) {
		ui.Error(fmt.Sprintf("%s exists and is not a vibeit workspace", target))
		return fmt.Errorf("aborting: unrelated directory")
	}

	ui.Header("create", name)
	ui.WorkspaceInfo(target, filepath.Join(target, name+".code-workspace"))

	data := buildData(name, f)

	if err := workspace.EnsureWorkspace(target, name); err != nil {
		return err
	}
	if err := generate.WorkspaceRoot(target, data); err != nil {
		return err
	}

	stacks := selectedStacks(f)
	if len(stacks) == 0 {
		stacks, err = promptStacks(workspace.ListMissingStacks(target))
		if err != nil {
			return err
		}
		if len(stacks) == 0 {
			wsFile := workspace.FindWorkspaceFile(target)
			if wsFile != "" {
				if err := workspace.SyncWorkspaceFolders(wsFile, target); err != nil {
					ui.Error("failed to update workspace file: " + err.Error())
					return err
				}
			}
			ui.Dim("nothing selected — workspace folder is ready for vibeit add")
			return nil
		}
	}

	result := runStacks(target, stacks, data)

	wsFile := workspace.FindWorkspaceFile(target)
	if wsFile != "" {
		if err := workspace.SyncWorkspaceFolders(wsFile, target); err != nil {
			ui.Error("failed to update workspace file: " + err.Error())
			return err
		}
	}

	openCmd := ""
	if f.open && wsFile != "" {
		openCmd = tryOpen(wsFile)
	}

	ui.Summary(ui.SummaryOpts{
		Name:    name,
		Root:    target,
		Stacks:  result.created,
		Skipped: result.skipped,
		Failed:  result.failed,
		OpenCmd: openCmd,
	})
	return nil
}

func resolveOutDir(cwd, out string) string {
	if out == "" {
		return cwd
	}
	if filepath.IsAbs(out) {
		return out
	}
	abs, err := filepath.Abs(filepath.Join(cwd, out))
	if err != nil {
		return filepath.Join(cwd, out)
	}
	return abs
}

func addCmd() *cobra.Command {
	var f stackFlags

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add stacks to an existing workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			root, ok := workspace.DetectWorkspace(cwd)
			if !ok {
				ui.Error("no vibeit workspace found in current directory")
				return fmt.Errorf("run vibeit create first")
			}

			wsFile := workspace.FindWorkspaceFile(root)
			name := strings.TrimSuffix(filepath.Base(wsFile), ".code-workspace")

			ui.Header("add", name)

			data := buildData(name, f)
			missing := workspace.ListMissingStacks(root)
			if len(missing) == 0 {
				ui.Dim("all stacks already exist")
				return nil
			}

			stacks := selectedStacks(f)
			if len(stacks) == 0 {
				stacks, err = promptStacks(missing)
				if err != nil {
					return err
				}
				if len(stacks) == 0 {
					ui.Dim("nothing selected")
					return nil
				}
			}

			result := runStacks(root, stacks, data)

			if wsFile != "" {
				if err := workspace.SyncWorkspaceFolders(wsFile, root); err != nil {
					ui.Error("failed to update workspace file: " + err.Error())
					return err
				}
			}

			openCmd := ""
			if f.open && wsFile != "" {
				openCmd = tryOpen(wsFile)
			}

			ui.Summary(ui.SummaryOpts{
				Action:     "updated",
				Name:       name,
				Root:       root,
				Stacks:     result.created,
				Skipped:    result.skipped,
				Failed:     result.failed,
				OpenCmd:    openCmd,
				RelativeCD: true,
			})
			return nil
		},
	}

	bindStackFlags(cmd, &f)
	return cmd
}

func skillsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skills",
		Short: "Manage workspace Cursor skills",
	}

	var force bool
	update := &cobra.Command{
		Use:   "update",
		Short: "Update .cursor skills from templates (skips learned/ and locally edited)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			root, ok := workspace.DetectWorkspace(cwd)
			if !ok {
				root = cwd
			}

			ui.Header("skills update", root)
			if err := generate.UpdateSkills(root, force); err != nil {
				ui.Error(err.Error())
				return err
			}
			ui.Success("skills updated")
			return nil
		},
	}
	update.Flags().BoolVar(&force, "force", false, "overwrite locally modified skills")

	cmd.AddCommand(update)
	return cmd
}
