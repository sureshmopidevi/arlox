package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/sureshmopidevi/arlox/internal/generate"
	"github.com/sureshmopidevi/arlox/internal/ui"
	"github.com/sureshmopidevi/arlox/internal/version"
	"github.com/sureshmopidevi/arlox/internal/workspace"
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
		Short: "Create a new arlox workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			args, err = resolveVariantArgs(args, &f, true)
			if err != nil {
				return err
			}
			return runCreate(args, f)
		},
	}
	bindStackFlags(create, &f)

	root := &cobra.Command{
		Use:     "arlox",
		Short:   "Scaffold and manage multi-stack workspaces",
		Long:    "arlox create workspaces with backend, web, and/or Flutter stacks plus Cursor rules and skills.",
		Version: version.Version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			maybeShowUpgradeNotice(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	root.SetVersionTemplate("arlox {{.Version}}\n")
	root.AddCommand(create)
	root.AddCommand(addCmd())
	root.AddCommand(skillsCmd())
	root.AddCommand(repairCmd())
	root.AddCommand(doctorCmd())
	root.AddCommand(upgradeCmd())
	root.AddCommand(uninstallCmd())
	root.AddCommand(versionCmd())
	return root
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print arlox version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("arlox " + version.Version)
		},
	}
}

type stackFlags struct {
	backend workspace.Variant
	web     workspace.Variant
	app     workspace.Variant
	open    bool
	module  string
	org     string
	out     string
}

type variantFlag struct {
	value *workspace.Variant
}

func (f variantFlag) String() string {
	if f.value == nil {
		return ""
	}
	return string(*f.value)
}

func (f variantFlag) Set(value string) error {
	*f.value = workspace.Variant(value)
	return nil
}

func (variantFlag) Type() string { return "variant" }

func bindStackFlags(cmd *cobra.Command, f *stackFlags) {
	bindAddFlags(cmd, f)
	cmd.Flags().StringVar(&f.module, "module", "", "Go module path (default: github.com/example/<name>-backend)")
	cmd.Flags().StringVar(&f.org, "org", "com.example", "Flutter org identifier")
	cmd.Flags().StringVar(&f.out, "out", "", "parent directory for workspace (default: current directory)")
}

// bindAddFlags binds only the flags relevant to the add subcommand.
func bindAddFlags(cmd *cobra.Command, f *stackFlags) {
	cmd.Flags().Var(variantFlag{value: &f.backend}, "backend", "include backend stack (go, python, node-express, node-fastify, java)")
	cmd.Flags().Var(variantFlag{value: &f.web}, "web", "include web stack (react-vite, nextjs, vue, svelte, angular, nuxt)")
	cmd.Flags().Var(variantFlag{value: &f.app}, "app", "include app stack (flutter, react-pwa, ios, android)")
	cmd.Flags().Lookup("backend").NoOptDefVal = string(workspace.GoGin)
	cmd.Flags().Lookup("web").NoOptDefVal = string(workspace.ReactVite)
	cmd.Flags().Lookup("app").NoOptDefVal = string(workspace.Flutter)
	cmd.Flags().BoolVar(&f.open, "open", false, "open workspace in Cursor, VS Code, or Antigravity IDE after creation")
}

type stackSelection struct {
	stack   workspace.Stack
	variant workspace.Variant
}

func selectedVariantsFromFlags(f stackFlags) []stackSelection {
	var selections []stackSelection
	if f.backend != "" {
		selections = append(selections, stackSelection{workspace.Backend, f.backend})
	}
	if f.web != "" {
		selections = append(selections, stackSelection{workspace.Web, f.web})
	}
	if f.app != "" {
		selections = append(selections, stackSelection{workspace.App, f.app})
	}
	return selections
}

func validateSelections(selections []stackSelection) error {
	for _, selection := range selections {
		if !workspace.ValidVariant(selection.stack, selection.variant) {
			return fmt.Errorf("invalid --%s variant %q", selection.stack, selection.variant)
		}
	}
	return nil
}

// resolveVariantArgs handles Cobra's optional-value flag behavior. With
// NoOptDefVal set for legacy bare flags, pflag leaves a space-separated value
// in the positional arguments; move recognized values back to their flag.
func resolveVariantArgs(args []string, f *stackFlags, hasName bool) ([]string, error) {
	keep := 0
	if hasName && len(args) > 0 {
		keep = 1
	}
	positional := append([]string(nil), args[:keep]...)
	for _, arg := range args[keep:] {
		var target *workspace.Variant
		for _, candidate := range []struct {
			stack workspace.Stack
			value *workspace.Variant
		}{
			{workspace.Backend, &f.backend},
			{workspace.Web, &f.web},
			{workspace.App, &f.app},
		} {
			if *candidate.value == workspace.DefaultVariant(candidate.stack) &&
				workspace.ValidVariant(candidate.stack, workspace.Variant(arg)) {
				target = candidate.value
				break
			}
		}
		if target == nil {
			return nil, fmt.Errorf("unexpected argument %q (variant values must follow a matching stack flag)", arg)
		}
		*target = workspace.Variant(arg)
	}
	if len(positional) > 1 {
		return nil, fmt.Errorf("accepts at most 1 arg, received %d", len(positional))
	}
	return positional, nil
}

func buildData(name string, f stackFlags) generate.Data {
	module := f.module
	if module == "" {
		module = "github.com/example/" + name + "-backend"
	}
	// Flutter package names must be snake_case (no hyphens); match project name like web's package.json.
	pkg := strings.ReplaceAll(name, "-", "_")
	data := generate.Data{
		Name:         name,
		DisplayName:  toDisplayName(name),
		Module:       module,
		Package:      pkg,
		Org:          f.org,
		APIURL:       "http://localhost:8080/api/v1",
		ArloxVersion: version.Version,
	}
	if selected := selectedVariantsFromFlags(f); len(selected) > 0 {
		data.Variant = selected[0].variant
	}
	return data
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
			Description("lowercase letters, numbers, hyphens, and underscores").
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
		workspace.Backend: "backend",
		workspace.Web:     "web",
		workspace.App:     "app",
	}

	opts := make([]huh.Option[string], 0, len(available))
	for _, s := range available {
		opts = append(opts, huh.NewOption(labels[s], string(s)))
	}

	var selected []string
	err := huh.NewForm(huh.NewGroup(
		huh.NewMultiSelect[string]().
			Title("What should arlox generate?").
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

func promptVariant(stack workspace.Stack) (workspace.Variant, error) {
	if stack == workspace.App {
		var target string
		if err := huh.NewForm(huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which app target?").
				Options(
					huh.NewOption("Flutter", string(workspace.Flutter)),
					huh.NewOption("React PWA", string(workspace.ReactPWA)),
					huh.NewOption("Native", "native"),
				).
				Value(&target),
		)).Run(); err != nil {
			return "", err
		}
		if target != "native" {
			return workspace.Variant(target), nil
		}

		var native workspace.Variant
		err := huh.NewForm(huh.NewGroup(
			huh.NewSelect[workspace.Variant]().
				Title("Which native platform?").
				Options(
					huh.NewOption("iOS (SwiftUI)", workspace.NativeIOS),
					huh.NewOption("Android (Kotlin + Compose)", workspace.NativeAndroid),
				).
				Value(&native),
		)).Run()
		return native, err
	}

	labels := map[workspace.Variant]string{
		workspace.GoGin:       "Go + Gin",
		workspace.PyFastAPI:   "Python + FastAPI",
		workspace.NodeExpress: "Node.js + Express",
		workspace.NodeFastify: "Node.js + Fastify",
		workspace.JavaSpring:  "Java + Spring Boot",
		workspace.ReactVite:   "React + Vite",
		workspace.NextJS:      "Next.js",
		workspace.VueVite:     "Vue + Vite",
		workspace.SvelteVite:  "Svelte + Vite",
		workspace.Angular:     "Angular",
		workspace.Nuxt:        "Nuxt",
	}
	options := make([]huh.Option[workspace.Variant], 0, len(workspace.VariantsForStack(stack)))
	for _, variant := range workspace.VariantsForStack(stack) {
		options = append(options, huh.NewOption(labels[variant], variant))
	}
	var selected workspace.Variant
	err := huh.NewForm(huh.NewGroup(
		huh.NewSelect[workspace.Variant]().
			Title(fmt.Sprintf("Which %s framework?", stack)).
			Options(options...).
			Value(&selected),
	)).Run()
	return selected, err
}

type stackResult struct {
	created  []string
	warnings map[string]string
	skipped  []string
	failed   []string
}

func executeStacks(root string, selections []stackSelection, data generate.Data) stackResult {
	r := stackResult{
		warnings: make(map[string]string),
	}
	for _, selection := range selections {
		s := selection.stack
		if workspace.StackExists(root, s) {
			ui.Progress(string(s), "skipped")
			r.skipped = append(r.skipped, string(s))
			continue
		}
		if workspace.StackPresentButEmpty(root, s) {
			ui.Progress(string(s), "failed")
			ui.Error(fmt.Sprintf("%s/ exists but is incomplete — remove it or run arlox add --%s again after fixing", s, s))
			r.failed = append(r.failed, string(s))
			continue
		}
		ui.Progress(string(s), "generating")
		warn, err := generate.Stack(root, s, selection.variant, data)
		if err != nil {
			ui.Progress(string(s), "failed")
			ui.Error(err.Error())
			r.failed = append(r.failed, string(s))
			continue
		}
		if !workspace.StackExists(root, s) {
			ui.Progress(string(s), "failed")
			ui.Error(fmt.Sprintf("%s finished but stack files are missing — try again or run arlox add --%s", s, s))
			r.failed = append(r.failed, string(s))
			continue
		}
		if warn != "" {
			ui.Progress(string(s), "warning")
			r.warnings[string(s)] = warn
		} else {
			ui.Progress(string(s), "done")
		}
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

// generateOpts captures the options for the shared generate workflow.
type generateOpts struct {
	action     string // "created" or "updated"
	relativeCD bool   // true for add (already inside workspace)
}

// runGenerate is the shared orchestration for create and add: select stacks,
// generate them, sync the workspace file, and print the summary.
func runGenerate(root, name string, f stackFlags, opts generateOpts) error {
	if stop, err := offerUpgradeBeforeGenerate(); err != nil {
		return err
	} else if stop {
		return nil
	}

	data := buildData(name, f)
	missing := workspace.ListMissingStacks(root)

	selections := selectedVariantsFromFlags(f)
	if err := validateSelections(selections); err != nil {
		return err
	}
	if len(selections) == 0 {
		var err error
		stacks, err := promptStacks(missing)
		if err != nil {
			return err
		}
		if len(stacks) == 0 {
			wsFile := workspace.FindWorkspaceFile(root)
			if wsFile != "" {
				if err := workspace.SyncWorkspaceFolders(wsFile, root); err != nil {
					ui.Error("failed to update workspace file: " + err.Error())
					return err
				}
			}
			ui.Dim("nothing selected — workspace folder is ready for arlox add")
			return nil
		}
		for _, stack := range stacks {
			variant, err := promptVariant(stack)
			if err != nil {
				return err
			}
			selections = append(selections, stackSelection{stack: stack, variant: variant})
		}
	}

	result := executeStacks(root, selections, data)

	wsFile := workspace.FindWorkspaceFile(root)
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
		Action:     opts.action,
		Name:       name,
		Root:       root,
		Stacks:     result.created,
		Warnings:   result.warnings,
		Skipped:    result.skipped,
		Failed:     result.failed,
		OpenCmd:    openCmd,
		RelativeCD: opts.relativeCD,
	})
	if len(result.failed) > 0 {
		return fmt.Errorf("%d stack(s) failed", len(result.failed))
	}
	return nil
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
		ui.Error(fmt.Sprintf("%s exists and is not an arlox workspace", target))
		return fmt.Errorf("aborting: unrelated directory")
	}

	ui.Header("create", name)
	ui.WorkspaceInfo(target, filepath.Join(target, name+".code-workspace"))

	wsFile := filepath.Join(target, name+".code-workspace")
	_, wsErr := os.Stat(wsFile)
	isNewWorkspace := os.IsNotExist(wsErr)

	if err := workspace.EnsureWorkspace(target, name); err != nil {
		return err
	}
	if isNewWorkspace {
		if err := generate.WorkspaceRoot(target, buildData(name, f)); err != nil {
			return err
		}
	}

	return runGenerate(target, name, f, generateOpts{action: "created"})
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
			remaining, err := resolveVariantArgs(args, &f, false)
			if err != nil {
				return err
			}
			if len(remaining) > 0 {
				return fmt.Errorf("add accepts no arguments")
			}
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			root, ok := workspace.DetectWorkspace(cwd)
			if !ok {
				ui.Error("no arlox workspace found — run from a workspace root or stack folder")
				return fmt.Errorf("run arlox create first")
			}

			wsFile := workspace.FindWorkspaceFile(root)
			if wsFile == "" {
				ui.Error("workspace directory found but .code-workspace file is missing")
				return fmt.Errorf("missing .code-workspace file in %s", root)
			}
			name := strings.TrimSuffix(filepath.Base(wsFile), ".code-workspace")

			ui.Header("add", name)

			if err := workspace.SyncWorkspaceFolders(wsFile, root); err != nil {
				ui.Error("failed to update workspace file: " + err.Error())
				return err
			}

			if len(workspace.ListMissingStacks(root)) == 0 {
				ui.Dim("all stacks already exist")
				return nil
			}

			return runGenerate(root, name, f, generateOpts{action: "updated", relativeCD: true})
		},
	}

	bindAddFlags(cmd, &f)
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
				ui.Error("no arlox workspace found — run from a workspace root or stack folder")
				return fmt.Errorf("run arlox create first")
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
