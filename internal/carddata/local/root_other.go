//go:build !linux && !openbsd && !freebsd && !netbsd && !darwin && !windows

package localdata

import (
	"os"
	"path/filepath"
)

// SharedRoot returns a best-effort shared data directory on other platforms.
func SharedRoot() (string, error) {
	configDirectory, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDirectory, "fyne", SharedApplicationID), nil
}
