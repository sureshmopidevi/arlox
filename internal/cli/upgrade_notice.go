package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/sureshmopidevi/arlox/internal/ui"
	"github.com/sureshmopidevi/arlox/internal/version"
)

func maybeShowUpgradeNotice(cmd *cobra.Command) {
	if os.Getenv("ARLOX_SKIP_UPGRADE_NOTICE") != "" {
		return
	}
	if skipUpgradeNotice(cmd) {
		return
	}
	latest, ok := latestAvailableVersion()
	if !ok || !version.Greater(latest, version.Version) {
		return
	}
	if upgradeNoticeShown(latest) {
		return
	}
	ui.UpgradeAvailable(version.Version, latest)
	_ = markUpgradeNoticeShown(latest)
}

func offerUpgradeBeforeGenerate() (stop bool, err error) {
	if os.Getenv("ARLOX_SKIP_UPGRADE_NOTICE") != "" {
		return false, nil
	}
	latest, ok := latestAvailableVersion()
	if !ok || !version.Greater(latest, version.Version) {
		return false, nil
	}
	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		return false, nil
	}

	var upgrade bool
	if err := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title(fmt.Sprintf("arlox %s is available (you have %s)", latest, version.Version)).
			Description("Upgrade now?").
			Value(&upgrade),
	)).Run(); err != nil {
		return false, err
	}
	if !upgrade {
		return false, nil
	}

	if err := runUpgrade("", true); err != nil {
		return true, err
	}
	ui.Dim("Re-run the same command to continue with the new version.")
	return true, nil
}

func skipUpgradeNotice(cmd *cobra.Command) bool {
	switch cmd.Name() {
	case "upgrade", "version", "uninstall", "help":
		return true
	default:
		return false
	}
}

func latestAvailableVersion() (string, bool) {
	if v, ok := localSourceVersion(); ok {
		return v, true
	}
	v, err := version.RemoteLatest("")
	if err != nil {
		return "", false
	}
	return v, true
}

func localSourceVersion() (string, bool) {
	root, err := findArloxSource("")
	if err != nil {
		return "", false
	}
	path := filepath.Join(root, "internal", "version", "VERSION")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	v := strings.TrimSpace(string(data))
	if v == "" {
		return "", false
	}
	return v, true
}

func upgradeNoticeFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "arlox", "upgrade-notice"), nil
}

func upgradeNoticeShown(latest string) bool {
	path, err := upgradeNoticeFile()
	if err != nil {
		return false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(data)) == latest
}

func markUpgradeNoticeShown(latest string) error {
	path, err := upgradeNoticeFile()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(latest+"\n"), 0o644)
}
