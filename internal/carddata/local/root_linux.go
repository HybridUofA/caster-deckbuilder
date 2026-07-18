//go:build linux || openbsd || freebsd || netbsd

package localdata

import (
	"os"
	"path/filepath"
)

// SharedRoot returns Fyne's existing per-user data directory on XDG desktops.
func SharedRoot() (string, error) {
	configDirectory, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDirectory, "fyne", SharedApplicationID), nil
}
