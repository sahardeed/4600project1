package builtins

import (
	"fmt"
	"os"
	"os/exec"
)

func handleTcsh(args ...string) error {
	if len(args) < 1 {
		fmt.Println("Usage: tcsh [command]")
		return nil
	}

	command := args[0]
	cmd := exec.Command("tcsh", "-c", command)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
