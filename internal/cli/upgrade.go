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

func upgradeCmd() *cobra.Command {
	var (
		noPull bool
		source string
	)
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Update arlox from the local source repo and reinstall",
		Long: `Rebuild and reinstall arlox from your local checkout.

Looks for the source repo in this order:
  1. --source <path>
  2. $ARLOX_HOME (or $VIBEIT_HOME)
  3. current directory (if it is the arlox repo)
  4. ~/arlox (or ~/vibeit)

By default runs git pull --ff-only before building. Use --no-pull to skip.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpgrade(source, !noPull)
		},
	}
	cmd.Flags().BoolVar(&noPull, "no-pull", false, "skip git pull before building")
	cmd.Flags().StringVar(&source, "source", "", "path to arlox source repo (default: $ARLOX_HOME or ~/arlox)")
	return cmd
}

func runUpgrade(sourceFlag string, pull bool) error {
	before := version.Version
	ui.Header("upgrade", "arlox "+before)

	if _, err := exec.LookPath("go"); err != nil {
		ui.Error("go not found on PATH — install Go first: https://go.dev/dl/")
		return fmt.Errorf("go required")
	}

	root, err := findArloxSource(sourceFlag)
	if err != nil {
		ui.Error(err.Error())
		return err
	}
	ui.Dim("source  " + root)

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
	if v := strings.TrimSpace(os.Getenv("VIBEIT_HOME")); v != "" {
		candidates = append(candidates, v)
	}
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates, cwd)
	}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, "arlox"), filepath.Join(home, "vibeit"))
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

	return "", fmt.Errorf("arlox source repo not found — set ARLOX_HOME or pass --source /path/to/arlox")
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
