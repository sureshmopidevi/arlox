package cli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sureshmopidevi/arlox/internal/ui"
	"github.com/sureshmopidevi/arlox/internal/version"
)

const arloxModule = "github.com/sureshmopidevi/arlox"
const arloxInstallModule = arloxModule + "/cmd/arlox@latest"

func upgradeCmd() *cobra.Command {
	var (
		noPull bool
		source string
	)
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Update arlox to the latest version",
		Long: `Update arlox to the latest version.

If a local source checkout is found, rebuilds from that repo (and runs
git pull --ff-only unless --no-pull is set).

Otherwise upgrades via:

  go install github.com/sureshmopidevi/arlox/cmd/arlox@latest

This is the same path used by the curl installer, so arlox upgrade works
after a remote install with no local clone.

Local source lookup order:
  1. --source <path>
  2. $ARLOX_HOME
  3. ~/.config/arlox/source (saved by install.sh or a prior upgrade)
  4. current directory (if it is the arlox repo)
  5. ~/arlox`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpgrade(source, !noPull)
		},
	}
	cmd.Flags().BoolVar(&noPull, "no-pull", false, "skip git pull before building from a local source repo")
	cmd.Flags().StringVar(&source, "source", "", "path to local arlox source repo")
	return cmd
}

func runUpgrade(sourceFlag string, pull bool) error {
	before := version.Version
	ui.Header("upgrade", "arlox "+before)

	if _, err := exec.LookPath("go"); err != nil {
		ui.Error("go not found on PATH — install Go first: https://go.dev/dl/")
		return fmt.Errorf("go required")
	}

	if sourceFlag != "" {
		root, err := findArloxSource(sourceFlag)
		if err != nil {
			ui.Error(err.Error())
			return err
		}
		return runLocalUpgrade(root, pull, before)
	}

	root, err := findArloxSource("")
	if err != nil {
		return runRemoteUpgrade(before)
	}
	return runLocalUpgrade(root, pull, before)
}

func runLocalUpgrade(root string, pull bool, before string) error {
	ui.Dim("source  " + root)

	if err := writeArloxSourceConfig(root); err != nil {
		ui.Dim("note    could not save source path: " + err.Error())
	}

	if pull {
		ui.Dim("git     pull --ff-only")
		if err := gitPullFF(root); err != nil {
			ui.Error("git pull failed: " + err.Error())
			return err
		}
	} else {
		ui.Dim("git     skipped (--no-pull)")
	}

	ui.Dim("build   go build → bin/arlox")
	if err := runInDir(root, "go", "build", "-o", "bin/arlox", "./cmd/arlox"); err != nil {
		ui.Error("build failed: " + err.Error())
		return err
	}

	ui.Dim("install go install → $(go env GOPATH)/bin")
	if err := runInDir(root, "go", "install", "./cmd/arlox"); err != nil {
		ui.Error("install failed: " + err.Error())
		return err
	}

	return finishUpgrade(before)
}

func runRemoteUpgrade(before string) error {
	ui.Dim("install go install " + arloxInstallModule)
	if err := runInDir("", "go", "install", arloxInstallModule); err != nil {
		ui.Error("install failed: " + err.Error())
		return err
	}
	return finishUpgrade(before)
}

func finishUpgrade(before string) error {
	after := installedVersion()
	if after == "" {
		after = "(unknown — run: arlox version)"
	}
	ui.Success(fmt.Sprintf("upgraded  %s → %s", before, after))
	fmt.Println()
	return nil
}

func findArloxSource(explicit string) (string, error) {
	var candidates []string
	if explicit != "" {
		candidates = append(candidates, explicit)
	}
	if v := strings.TrimSpace(os.Getenv("ARLOX_HOME")); v != "" {
		candidates = append(candidates, v)
	}
	if v := readArloxSourceConfig(); v != "" {
		candidates = append(candidates, v)
	}
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates, cwd)
	}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, "arlox"))
	}

	seen := map[string]bool{}
	for _, c := range candidates {
		abs, err := filepath.Abs(c)
		if err != nil {
			continue
		}
		if seen[abs] {
			continue
		}
		seen[abs] = true
		if isArloxRepo(abs) {
			return abs, nil
		}
	}

	return "", fmt.Errorf("arlox source repo not found — run from the repo, set ARLOX_HOME, pass --source /path/to/arlox, or reinstall with ./install.sh")
}

func arloxSourceConfigFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "arlox", "source"), nil
}

func readArloxSourceConfig() string {
	path, err := arloxSourceConfigFile()
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func writeArloxSourceConfig(root string) error {
	path, err := arloxSourceConfigFile()
	if err != nil {
		return err
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(abs+"\n"), 0o644)
}

func isArloxRepo(dir string) bool {
	modPath := filepath.Join(dir, "go.mod")
	data, err := os.ReadFile(modPath)
	if err != nil {
		return false
	}
	if !bytes.Contains(data, []byte("module "+arloxModule)) {
		return false
	}
	if _, err := os.Stat(filepath.Join(dir, "cmd", "arlox")); err != nil {
		return false
	}
	return true
}

func gitPullFF(root string) error {
	if _, err := os.Stat(filepath.Join(root, ".git")); err != nil {
		return nil // not a git checkout; skip quietly
	}
	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runInDir(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installedVersion() string {
	path, err := exec.LookPath("arlox")
	if err != nil {
		return ""
	}
	out, err := exec.Command(path, "version").Output()
	if err != nil {
		return ""
	}
	// "arlox 0.1.0" → "0.1.0"
	fields := strings.Fields(strings.TrimSpace(string(out)))
	if len(fields) >= 2 {
		return fields[len(fields)-1]
	}
	return strings.TrimSpace(string(out))
}
