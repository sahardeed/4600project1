package builtins

import (
	"fmt"
	"os"
	"os/exec"
)

func handleCsh(args ...string) error {
	if len(args) < 1 {
		fmt.Println("Usage: csh [command]")
		return nil
	}

	command := args[0]
	cmd := exec.Command("csh", "-c", command)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
