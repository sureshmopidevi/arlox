package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/sureshmopidevi/arlox/internal/ui"
)

func uninstallCmd() *cobra.Command {
	var (
		yes bool
		all bool
	)

	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Remove arlox (and legacy binaries) from your system",
		Long: `Uninstall removes the arlox binary from $(go env GOPATH)/bin and other PATH locations.
Also cleans up legacy vibeit binaries when present.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUninstall(yes, all)
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation prompt")
	cmd.Flags().BoolVar(&all, "all", true, "also remove legacy vibeit binary if found")
	return cmd
}

func findInstalledBinaries(includeLegacy bool) []string {
	var targets []string
	names := []string{"arlox"}
	if includeLegacy {
		names = append(names, "vibeit")
	}

	seen := make(map[string]bool)

	// Check go env GOPATH / bin
	if out, err := exec.Command("go", "env", "GOPATH").Output(); err == nil {
		gopath := strings.TrimSpace(string(out))
		if gopath != "" {
			for _, name := range names {
				p := filepath.Join(gopath, "bin", name)
				if _, err := os.Stat(p); err == nil && !seen[p] {
					seen[p] = true
					targets = append(targets, p)
				}
			}
		}
	}

	// Check ~/go/bin
	if home, err := os.UserHomeDir(); err == nil {
		for _, name := range names {
			p := filepath.Join(home, "go", "bin", name)
			if _, err := os.Stat(p); err == nil && !seen[p] {
				seen[p] = true
				targets = append(targets, p)
			}
		}
	}

	// Check exec.LookPath
	for _, name := range names {
		if path, err := exec.LookPath(name); err == nil && !seen[path] {
			// Don't include binary in the source checkout ./bin/arlox
			if !strings.Contains(path, filepath.Join("arlox", "bin")) {
				seen[path] = true
				targets = append(targets, path)
			}
		}
	}

	return targets
}

func runUninstall(autoConfirm, includeLegacy bool) error {
	ui.Header("uninstall", "arlox")

	binaries := findInstalledBinaries(includeLegacy)
	if len(binaries) == 0 {
		ui.Dim("no installed arlox binaries found on PATH or in GOPATH/bin")
		return nil
	}

	fmt.Println("  Found installed binaries:")
	for _, b := range binaries {
		fmt.Printf("    • %s\n", b)
	}
	fmt.Println()

	if !autoConfirm {
		var confirm bool
		err := huh.NewForm(huh.NewGroup(
			huh.NewConfirm().
				Title("Remove these binaries?").
				Description("This deletes the compiled arlox executable(s).").
				Value(&confirm),
		)).Run()
		if err != nil {
			return err
		}
		if !confirm {
			ui.Dim("uninstall aborted")
			return nil
		}
	}

	removed := removeBinaries(binaries)

	fmt.Println()
	if removed > 0 {
		ui.Success(fmt.Sprintf("arlox uninstalled (%d binary removed)", removed))
		ui.Dim("To reinstall later: go install github.com/sureshmopidevi/arlox/cmd/arlox@latest")
	}
	return nil
}

func removeBinaries(binaries []string) int {
	removed := 0
	for _, b := range binaries {
		if err := os.Remove(b); err != nil {
			ui.Error(fmt.Sprintf("failed to remove %s: %v", b, err))
		} else {
			ui.Success(fmt.Sprintf("removed %s", b))
			removed++
		}
	}
	return removed
}
