package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {

	cmd := exec.Command("/bin/sh", "./container/cgroups/memory/mount_and_exhaust_mem.sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println("cmd start error")
		os.Exit(1)
	}

	memoryCgroupSetup(cmd.Process.Pid)

	if err := cmd.Wait(); err != nil {
		fmt.Println("cmd wait error")
		os.Exit(1)
	}
}

func memoryCgroupSetup(pid int) {
	cpath := "/sys/fs/cgroup/memory/mycontainer/"
	if err := os.MkdirAll(cpath, 0644); err != nil {
		fmt.Println("failed to create memory cgroup")
	}
	addProcessToCgroup(cpath+"cgroup.procs", pid)
}

func addProcessToCgroup(fpath string, pid int) {
	file, err := os.OpenFile(fpath, os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("%d", pid)); err != nil {
		fmt.Println("failer to write pid")
		panic(err)
	}
}
