// +build !windows,!darwin

package internal

import (
	"context"
	"os/exec"
)

func DetectTincBinary() (string, error) {
	return exec.LookPath("tincd")
}

func Preload(ctx context.Context) error {
	return nil
}
