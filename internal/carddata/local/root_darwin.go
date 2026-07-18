//go:build darwin

package localdata

import (
	"os"
	"path/filepath"
)

// SharedRoot returns Fyne's existing per-user data directory on macOS.
func SharedRoot() (string, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDirectory, "Library", "Preferences", "fyne", SharedApplicationID), nil
}
