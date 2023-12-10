package builtins

import (
	"fmt"
	"os"
	"os/exec"
)

func handleZsh(args ...string) error {
	if len(args) < 1 {
		fmt.Println("Usage: zsh [command]")
		return nil
	}

	command := args[0]
	cmd := exec.Command("zsh", "-c", command)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
