package svcroot

import (
	"os"
	"path/filepath"
)

//
// ────────────────────────────────────────
// project root discovery.
//

// FindProjectRoot walks upward from from to find markerDir as a subdirectory.
func FindProjectRoot(from, markerDir string) (string, bool) {
	dir, err := filepath.Abs(from)
	if err != nil {
		return "", false
	}
	for {
		candidate := filepath.Join(dir, markerDir)
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}
