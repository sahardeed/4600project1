package builtins

import (
	"fmt"
	"os"
	"os/exec"
)

func handleBash(args ...string) error {
	if len(args) < 1 {
		fmt.Println("Usage: bash [command]")
		return nil
	}

	command := args[0]
	cmd := exec.Command("bash", "-c", command)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
