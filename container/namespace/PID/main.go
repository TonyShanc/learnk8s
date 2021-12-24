package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("/bin/sh")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		// UTS namespace
		Cloneflags: syscall.CLONE_NEWPID,
	}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running the exec.Command - %s\n", err)
		os.Exit(1)
	}
}