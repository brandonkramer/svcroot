package svcroot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//
// ────────────────────────────────────────
// root resolution.
//

// ResolveHome returns the envVar value when set, otherwise ~/.defaultDirName.
func ResolveHome(envVar, defaultDirName string) (string, error) {
	if h := strings.TrimSpace(os.Getenv(envVar)); h != "" {
		return h, nil
	}
	return DefaultHomeDir(defaultDirName)
}

// DefaultHomeDir returns ~/.defaultDirName under the user home directory.
func DefaultHomeDir(defaultDirName string) (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("svcroot: resolve user home: %w", err)
	}
	return filepath.Join(userHome, defaultDirName), nil
}
