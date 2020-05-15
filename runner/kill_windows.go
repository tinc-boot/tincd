package runner

import (
	"os/exec"
)

func killProcess(cmd *exec.Cmd) {
	_ = cmd.Process.Kill()
}
