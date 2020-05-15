//+build linux darwin

package runner

import (
	"os/exec"
	"syscall"
)

func killProcess(cmd *exec.Cmd) {
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
		_ = cmd.Process.Kill()
	}
}
