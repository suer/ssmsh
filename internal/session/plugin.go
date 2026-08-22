package session

import (
	"fmt"
	"os/exec"
)

const pluginBinary = "session-manager-plugin"

var lookPath = exec.LookPath

func CheckPluginInstalled() error {
	if _, err := lookPath(pluginBinary); err != nil {
		return fmt.Errorf(
			"%s not found in PATH, install it: https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html",
			pluginBinary,
		)
	}
	return nil
}
