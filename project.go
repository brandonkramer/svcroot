package svcroot

import (
	"os"
	"path/filepath"
)

//
// ────────────────────────────────────────
// project root discovery.
//

var projectPathAbs = filepath.Abs

// FindProjectRoot walks upward from from to find markerDir as a subdirectory.
func FindProjectRoot(from, markerDir string) (string, bool) {
	dir, err := projectPathAbs(from)
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
