package builtins

import (
	"fmt"
	"os"
	"os/exec"
)

func handleSh(args ...string) error {
	if len(args) < 1 {
		fmt.Println("Usage: sh [command]")
		return nil
	}

	command := args[0]
	cmd := exec.Command(command, args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
