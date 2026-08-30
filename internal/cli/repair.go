package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sureshmopidevi/arlox/internal/generate"
	"github.com/sureshmopidevi/arlox/internal/ui"
	"github.com/sureshmopidevi/arlox/internal/workspace"
)

func repairCmd() *cobra.Command {
	var (
		force bool
		deps  bool
	)

	cmd := &cobra.Command{
		Use:   "repair",
		Short: "Inspect and repair workspace files, rules, skills, and configs",
		Long: `Repair validates your workspace and restores missing files:
  • Re-syncs or recreates missing .code-workspace files
  • Restores missing root Makefile
  • Restores missing .cursor rules and skills
  • Restores missing .env and app.env configs from templates
  • Re-initializes missing git repositories in stacks
  • Optionally installs missing stack dependencies (go.sum, node_modules, pub get)`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			return runRepair(cwd, force, deps)
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "overwrite locally modified skills/rules with defaults")
	cmd.Flags().BoolVar(&deps, "deps", true, "install missing stack dependencies (go mod tidy, npm install, flutter pub get)")
	return cmd
}

type repairReport struct {
	repaired []string
	healthy  []string
	warnings []string
}

func (r *repairReport) fix(item string) {
	r.repaired = append(r.repaired, item)
}

func (r *repairReport) ok(item string) {
	r.healthy = append(r.healthy, item)
}

func (r *repairReport) warn(item string) {
	r.warnings = append(r.warnings, item)
}

func runRepair(cwd string, force, installDeps bool) error {
	root, ok := workspace.DetectWorkspace(cwd)
	if !ok {
		ui.Header("repair", "system check")
		ui.Dim("no arlox workspace detected in current directory — running toolchain check")
		return runDoctor()
	}

	wsFile := workspace.FindWorkspaceFile(root)
	name := ""
	if wsFile != "" {
		name = strings.TrimSuffix(filepath.Base(wsFile), ".code-workspace")
	} else {
		name = filepath.Base(root)
	}

	ui.Header("repair", name)
	ui.Dim(fmt.Sprintf("workspace  %s", root))
	fmt.Println()

	var report repairReport
	data := buildData(name, stackFlags{})

	// 1. Workspace file (.code-workspace)
	if wsFile == "" {
		wsFile = filepath.Join(root, name+".code-workspace")
		if err := workspace.EnsureWorkspace(root, name); err != nil {
			report.warn("failed to create .code-workspace: " + err.Error())
		} else {
			report.fix("created " + filepath.Base(wsFile))
		}
	} else {
		report.ok(filepath.Base(wsFile) + " exists")
	}

	if err := workspace.SyncWorkspaceFolders(wsFile, root); err != nil {
		report.warn("failed to sync folders in .code-workspace: " + err.Error())
	} else {
		report.fix("synced folder list in " + filepath.Base(wsFile))
	}

	// 2. Root Makefile
	rootMakefile := filepath.Join(root, "Makefile")
	if _, err := os.Stat(rootMakefile); os.IsNotExist(err) {
		if err := generate.WorkspaceRoot(root, data); err != nil {
			report.warn("failed to restore root Makefile: " + err.Error())
		} else {
			report.fix("restored root Makefile")
		}
	} else {
		report.ok("root Makefile exists")
	}

	// 3. Cursor skills & rules
	if err := generate.UpdateSkills(root, force); err != nil {
		report.warn("skills update error: " + err.Error())
	} else {
		report.fix("updated .cursor rules and skills")
	}

	// 4. Per-stack repairs (configs, git, deps)
	presentStacks := workspace.ListPresentStacks(root)
	for _, stackName := range presentStacks {
		stack := workspace.Stack(stackName)
		stackDir := filepath.Join(root, stackName)

		// Git init if missing
		gitDir := filepath.Join(stackDir, ".git")
		if _, err := os.Stat(gitDir); os.IsNotExist(err) {
			if err := workspace.InitGit(stackDir); err != nil {
				report.warn(fmt.Sprintf("%s: git init failed: %v", stackName, err))
			} else {
				report.fix(fmt.Sprintf("%s: initialized git repo", stackName))
			}
		} else {
			report.ok(fmt.Sprintf("%s: git repo initialized", stackName))
		}

		switch stack {
		case workspace.Backend:
			// Config check
			envExample := filepath.Join(stackDir, "configs/local/app.env.example")
			envFile := filepath.Join(stackDir, "configs/local/app.env")
			if _, err := os.Stat(envFile); os.IsNotExist(err) {
				if _, err := os.Stat(envExample); err == nil {
					data, readErr := os.ReadFile(envExample)
					if readErr == nil && os.WriteFile(envFile, data, 0o644) == nil {
						report.fix("backend: seeded configs/local/app.env from example")
					}
				}
			} else {
				report.ok("backend: app.env exists")
			}

			// Dependencies
			if installDeps && generate.HasTool("go") {
				goSum := filepath.Join(stackDir, "go.sum")
				if _, err := os.Stat(goSum); os.IsNotExist(err) {
					ui.Dim("backend: running go mod tidy...")
					if err := runInDir(stackDir, "go", "mod", "tidy"); err != nil {
						report.warn("backend: go mod tidy had issues: " + err.Error())
					} else {
						report.fix("backend: resolved go.sum dependencies")
					}
				} else {
					report.ok("backend: go.sum dependencies present")
				}
			}

		case workspace.Web:
			// Config check
			envExample := filepath.Join(stackDir, ".env.example")
			envFile := filepath.Join(stackDir, ".env")
			if _, err := os.Stat(envFile); os.IsNotExist(err) {
				if _, err := os.Stat(envExample); err == nil {
					data, readErr := os.ReadFile(envExample)
					if readErr == nil && os.WriteFile(envFile, data, 0o644) == nil {
						report.fix("web: seeded .env from example")
					}
				}
			} else {
				report.ok("web: .env exists")
			}

			// Dependencies
			if installDeps && generate.HasTool("npm") {
				nodeModules := filepath.Join(stackDir, "node_modules")
				if _, err := os.Stat(nodeModules); os.IsNotExist(err) {
					ui.Dim("web: running npm install...")
					if err := runInDir(stackDir, "npm", "install", "--no-audit", "--no-fund"); err != nil {
						report.warn("web: npm install had issues: " + err.Error())
					} else {
						report.fix("web: installed node_modules")
					}
				} else {
					report.ok("web: node_modules present")
				}
			}

		case workspace.App:
			// Dependencies
			if installDeps && generate.HasTool("flutter") {
				dartTool := filepath.Join(stackDir, ".dart_tool")
				if _, err := os.Stat(dartTool); os.IsNotExist(err) {
					ui.Dim("app: running flutter pub get...")
					if err := runInDir(stackDir, "flutter", "pub", "get"); err != nil {
						report.warn("app: flutter pub get had issues: " + err.Error())
					} else {
						report.fix("app: fetched flutter packages")
					}
				} else {
					report.ok("app: flutter packages present")
				}
			}
		}
	}

	fmt.Println()
	if len(report.repaired) > 0 {
		fmt.Println("  Repairs Applied:")
		for _, f := range report.repaired {
			ui.Success(f)
		}
		fmt.Println()
	}

	if len(report.healthy) > 0 {
		fmt.Println("  Healthy Checks:")
		for _, h := range report.healthy {
			ui.Dim("  ✓ " + h)
		}
		fmt.Println()
	}

	if len(report.warnings) > 0 {
		fmt.Println("  Warnings / Action Needed:")
		for _, w := range report.warnings {
			ui.Warn(w)
		}
		fmt.Println()
	}

	ui.Success(fmt.Sprintf("workspace repair completed for %s", name))
	return nil
}
