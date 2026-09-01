package exec

import (
	"os"
	osexec "os/exec"
)

// RunInDir runs a command in dir with stdout/stderr streamed to the terminal.
func RunInDir(dir, name string, args ...string) error {
	cmd := osexec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
