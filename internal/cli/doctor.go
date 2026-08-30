package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sureshmopidevi/arlox/internal/ui"
	"github.com/sureshmopidevi/arlox/internal/version"
	"github.com/sureshmopidevi/arlox/internal/workspace"
)

func doctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check system toolchains, PATH configuration, and workspace status",
		Long: `Doctor verifies that your development environment is properly configured for arlox.
Checks toolchains (Go, Node, npm, Flutter, Git) and current workspace status.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDoctor()
		},
	}
}

func runDoctor() error {
	ui.Header("doctor", "system & toolchain check")

	// 1. Check GOPATH bin on PATH
	gopathBin := ""
	if out, err := exec.Command("go", "env", "GOPATH").Output(); err == nil {
		gp := strings.TrimSpace(string(out))
		if gp != "" {
			gopathBin = filepath.Join(gp, "bin")
		}
	}
	if gopathBin == "" {
		if home, err := os.UserHomeDir(); err == nil {
			gopathBin = filepath.Join(home, "go", "bin")
		}
	}

	pathEnv := os.Getenv("PATH")
	pathHasGOPATH := strings.Contains(":"+pathEnv+":", ":"+gopathBin+":")

	fmt.Println("  Environment:")
	fmt.Printf("    • arlox version : %s\n", version.Version)
	if gopathBin != "" {
		if pathHasGOPATH {
			fmt.Printf("    • GOPATH bin    : %s (on PATH)\n", gopathBin)
		} else {
			fmt.Printf("    • GOPATH bin    : %s (NOT on PATH — add to ~/.zshrc or ~/.bashrc)\n", gopathBin)
		}
	}
	fmt.Println()

	// 2. Check toolchains
	type toolInfo struct {
		name    string
		version string
		ok      bool
	}

	tools := []string{"go", "node", "npm", "flutter", "git", "arlox"}
	fmt.Println("  Toolchains:")
	for _, t := range tools {
		path, err := exec.LookPath(t)
		if err != nil {
			fmt.Printf("    %s %-10s not found\n", "❌", t)
			continue
		}

		ver := ""
		switch t {
		case "go":
			if out, err := exec.Command(path, "version").Output(); err == nil {
				ver = strings.TrimSpace(string(out))
			}
		case "node", "npm":
			if out, err := exec.Command(path, "-v").Output(); err == nil {
				ver = strings.TrimSpace(string(out))
			}
		case "flutter":
			if out, err := exec.Command(path, "--version").Output(); err == nil {
				lines := strings.Split(strings.TrimSpace(string(out)), "\n")
				if len(lines) > 0 {
					ver = lines[0]
				}
			}
		case "git":
			if out, err := exec.Command(path, "--version").Output(); err == nil {
				ver = strings.TrimSpace(string(out))
			}
		case "arlox":
			ver = version.Version
		}

		if ver != "" {
			fmt.Printf("    %s %-10s %s (%s)\n", "✅", t, ver, path)
		} else {
			fmt.Printf("    %s %-10s %s\n", "✅", t, path)
		}
	}
	fmt.Println()

	// 3. Workspace check (if in workspace)
	if cwd, err := os.Getwd(); err == nil {
		if root, ok := workspace.DetectWorkspace(cwd); ok {
			wsFile := workspace.FindWorkspaceFile(root)
			name := strings.TrimSuffix(filepath.Base(wsFile), ".code-workspace")
			fmt.Println("  Current Workspace:")
			fmt.Printf("    • Project name : %s\n", name)
			fmt.Printf("    • Root         : %s\n", root)
			stacks := workspace.ListPresentStacks(root)
			fmt.Printf("    • Stacks       : %s\n", strings.Join(stacks, ", "))
			fmt.Println()
		}
	}

	return nil
}
