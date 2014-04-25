package main

import "os"
import "os/exec"

func ExecGit(args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)
	cmd.Dir = "c:\\dev\\Core"
	cmd.Stderr = os.Stderr
	return cmd
}
