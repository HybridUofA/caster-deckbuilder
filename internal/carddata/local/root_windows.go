//go:build windows

package localdata

import (
	"os"
	"path/filepath"
)

// SharedRoot returns Fyne's existing per-user data directory on Windows.
func SharedRoot() (string, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDirectory, "AppData", "Roaming", "fyne", SharedApplicationID), nil
}
