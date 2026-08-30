package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var noColor bool

func init() {
	_, noColorSet := os.LookupEnv("NO_COLOR")
	noColor = noColorSet || !term.IsTerminal(int(os.Stdout.Fd()))
}

func render(s lipgloss.Style, text string) string {
	if noColor {
		return text
	}
	return s.Render(text)
}

var (
	styleWorkspace = lipgloss.NewStyle().Bold(true)
	styleBackend   = lipgloss.NewStyle().Foreground(lipgloss.Color("#4B9CD3")).Bold(true)
	styleWeb       = lipgloss.NewStyle().Foreground(lipgloss.Color("#5BC8E5")).Bold(true)
	styleApp       = lipgloss.NewStyle().Foreground(lipgloss.Color("#C678DD")).Bold(true)
	styleSkip      = lipgloss.NewStyle().Foreground(lipgloss.Color("#E5C07B")).Faint(true)
	styleError     = lipgloss.NewStyle().Foreground(lipgloss.Color("#E06C75")).Bold(true)
	styleSuccess   = lipgloss.NewStyle().Foreground(lipgloss.Color("#98C379")).Bold(true)
	styleMuted     = lipgloss.NewStyle().Faint(true)
)

func stackStyle(stack string) lipgloss.Style {
	switch stack {
	case "backend":
		return styleBackend
	case "web":
		return styleWeb
	case "app":
		return styleApp
	case "workspace":
		return styleWorkspace
	default:
		return lipgloss.NewStyle()
	}
}

// Header prints a bold workspace action header.
func Header(action, name string) {
	fmt.Printf("\n  %s  %s\n\n", render(styleMuted, "arlox"), render(styleWorkspace, action+"  "+name))
}

// WorkspaceInfo prints where files will be written (always before generate).
func WorkspaceInfo(absPath, wsFile string) {
	fmt.Printf("  %s\n", render(styleMuted, "Workspace"))
	fmt.Printf("    %s\n", render(styleWorkspace, absPath))
	if wsFile != "" {
		fmt.Printf("    %s\n", render(styleMuted, filepath.Base(wsFile)))
	}
	fmt.Println()
}

// Progress prints a stack generation status line.
func Progress(stack, status string) {
	label := render(stackStyle(stack), fmt.Sprintf("%-10s", stack))
	var statusStr string
	switch status {
	case "generating":
		statusStr = render(styleMuted, "generating…")
	case "done":
		statusStr = render(styleSuccess, "done")
	case "skipped":
		statusStr = render(styleSkip, "skipped (already exists)")
	case "failed":
		statusStr = render(styleError, "failed")
	default:
		statusStr = status
	}
	fmt.Printf("  %s %s\n", label, statusStr)
}

// SummaryOpts configures Summary output.
type SummaryOpts struct {
	Action     string // created | updated
	Name       string
	Root       string
	Stacks     []string
	Skipped    []string
	Failed     []string
	OpenCmd    string
	RelativeCD bool
}

// Summary prints color-coded next-step commands.
func Summary(opts SummaryOpts) {
	action := opts.Action
	if action == "" {
		action = "created"
	}
	fmt.Println()
	fmt.Printf("  %s  %s\n\n", render(styleSuccess, action), render(styleWorkspace, opts.Name))

	wsFile := opts.Name + ".code-workspace"
	wsOpen := wsFile
	if !opts.RelativeCD {
		wsOpen = filepath.Join(opts.Name, wsFile)
	}

	if opts.OpenCmd != "" {
		fmt.Printf("  %s   %s\n", render(styleMuted, fmt.Sprintf("%-12s", "opened")), render(styleWorkspace, opts.OpenCmd))
	}
	// --classic opens Cursor's IDE editor; bare `cursor` launches the Agents window.
	fmt.Printf("  %s   %s\n", render(styleMuted, fmt.Sprintf("%-12s", "cursor")), render(styleWorkspace, "cursor --classic "+wsOpen))
	fmt.Printf("  %s   %s\n", render(styleMuted, fmt.Sprintf("%-12s", "vscode")), render(styleWorkspace, "code "+wsOpen))
	fmt.Printf("  %s   %s\n", render(styleMuted, fmt.Sprintf("%-12s", "antigravity")), render(styleWorkspace, "agy-ide "+wsOpen))

	prefix := opts.Name
	if opts.RelativeCD {
		prefix = "."
	}

	makeCmd := "make dev"
	if prefix != "." {
		makeCmd = "cd " + prefix + " && " + makeCmd
	}
	fmt.Printf("  %s   %s\n", render(styleMuted, fmt.Sprintf("%-12s", "make")), render(styleWorkspace, makeCmd))

	for _, s := range opts.Stacks {
		cmd := nextCommand(prefix, s)
		fmt.Printf("  %s   %s\n", render(styleMuted, fmt.Sprintf("%-12s", s)), render(stackStyle(s), cmd))
	}
	for _, s := range opts.Skipped {
		fmt.Printf("  %s   %s\n", render(styleMuted, fmt.Sprintf("%-12s", s)), render(styleSkip, "skipped (already exists)"))
	}
	for _, s := range opts.Failed {
		fmt.Printf("  %s   %s\n", render(styleMuted, fmt.Sprintf("%-12s", s)), render(styleError, "failed"))
	}

	fmt.Printf("\n  %s\n\n", render(styleMuted, "Open the .code-workspace file in Cursor, VS Code, or Antigravity IDE, then run make dev."))
}

func nextCommand(prefix, stack string) string {
	var cmd string
	switch stack {
	case "backend":
		cmd = "make backend.run"
	case "web":
		cmd = "make web.dev"
	case "app":
		cmd = "make app.run"
	default:
		cmd = "make help"
	}
	if prefix != "." {
		return "cd " + prefix + " && " + cmd
	}
	return cmd
}

func Error(msg string) {
	fmt.Fprintf(os.Stderr, "%s %s\n", render(styleError, "error"), msg)
}

func Success(msg string) {
	fmt.Printf("%s %s\n", render(styleSuccess, "ok"), msg)
}

func Dim(msg string) {
	fmt.Println(render(styleMuted, msg))
}

func StackLabel(stack string) string {
	return render(stackStyle(stack), stack)
}

func CommandBackend() string { return render(styleBackend, "make run") }
func CommandWeb() string     { return render(styleWeb, "npm run dev") }
func CommandApp() string     { return render(styleApp, "flutter run") }
func CommandOpen(wsFile string) string {
	return render(styleWorkspace, "cursor --classic "+wsFile)
}
