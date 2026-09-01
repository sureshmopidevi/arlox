package cli

import (
	"bytes"
	"fmt"
	"os"
	osexec "os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	arloxexec "github.com/sureshmopidevi/arlox/internal/exec"
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
  2. $ARLOX_HOME
  3. current directory (if it is the arlox repo)
  4. ~/arlox

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

	if _, err := osexec.LookPath("go"); err != nil {
		ui.Error("go not found on PATH — install Go first: https://go.dev/dl/")
		return fmt.Errorf("go required")
	}

	gopathBin, err := goPathBin()
	if err != nil {
		ui.Error(err.Error())
		return err
	}
	binaryPath := filepath.Join(gopathBin, "arlox")

	root, err := findArloxSource(sourceFlag)
	if err != nil {
		ui.Error(err.Error())
		return err
	}
	ui.Dim("source  " + root)

	if pull {
		if isGitCheckout(root) {
			ui.Dim("git     pull --ff-only")
			if err := runGitPullFF(root); err != nil {
				ui.Error("git pull failed: " + err.Error())
				return err
			}
			ui.Success("git pull complete")
		} else {
			ui.Dim("git     skipped (not a git checkout)")
		}
	} else {
		ui.Dim("git     skipped (--no-pull)")
	}

	ui.Dim("build   go build → bin/arlox")
	if err := arloxexec.RunInDir(root, "go", "build", "-o", "bin/arlox", "./cmd/arlox"); err != nil {
		ui.Error("build failed: " + err.Error())
		return err
	}
	ui.Success("build complete")

	ui.Dim("install go install → " + gopathBin)
	if err := arloxexec.RunInDir(root, "go", "install", "./cmd/arlox"); err != nil {
		ui.Error("install failed: " + err.Error())
		return err
	}
	ui.Success("installed to " + binaryPath)
	ui.Dim("binary  " + binaryPath)

	after := binaryVersionAt(binaryPath)
	if after == "" {
		after = "(unknown — run: arlox version)"
	}

	if after == before {
		ui.Success(fmt.Sprintf("reinstalled  %s", before))
	} else if strings.HasPrefix(after, "(") {
		ui.Success(fmt.Sprintf("reinstalled  %s", before))
	} else {
		ui.Success(fmt.Sprintf("upgraded  %s → %s", before, after))
	}

	ui.Dim("hint    run: arlox version")
	ui.Dim("hint    in workspace: arlox skills update")
	fmt.Println()
	return nil
}

func goPathBin() (string, error) {
	out, err := osexec.Command("go", "env", "GOPATH").Output()
	if err != nil {
		return "", fmt.Errorf("go env GOPATH failed: %w", err)
	}
	gopath := strings.TrimSpace(string(out))
	if gopath == "" {
		return "", fmt.Errorf("GOPATH is empty")
	}
	return filepath.Join(gopath, "bin"), nil
}

func findArloxSource(explicit string) (string, error) {
	var candidates []string
	if explicit != "" {
		candidates = append(candidates, explicit)
	}
	if v := strings.TrimSpace(os.Getenv("ARLOX_HOME")); v != "" {
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

func isGitCheckout(root string) bool {
	_, err := os.Stat(filepath.Join(root, ".git"))
	return err == nil
}

func runGitPullFF(root string) error {
	cmd := osexec.Command("git", "pull", "--ff-only")
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func binaryVersionAt(path string) string {
	out, err := osexec.Command(path, "version").Output()
	if err != nil {
		return ""
	}
	fields := strings.Fields(strings.TrimSpace(string(out)))
	if len(fields) >= 2 {
		return fields[len(fields)-1]
	}
	return strings.TrimSpace(string(out))
}
