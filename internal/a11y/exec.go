package a11y

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func RunAppleScript(ctx context.Context, script string) (string, error) {
	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("osascript error: %s\n%s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func IsSupported() bool {
	return runtime.GOOS == "darwin"
}
