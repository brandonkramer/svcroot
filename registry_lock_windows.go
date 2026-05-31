//go:build windows

package svcroot

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

//
// ────────────────────────────────────────
// registry locking (windows).
//

var (
	registryMkdirAll = os.MkdirAll
	registryOpenFile = func(lockPath string) (*os.File, error) {
		return os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	}
)

func withRegistryLock(path string, fn func() error) error {
	lockPath := path + ".lock"
	if err := registryMkdirAll(filepath.Dir(lockPath), 0o755); err != nil {
		return fmt.Errorf("svcroot: create lock dir: %w", err)
	}
	f, err := registryOpenFile(lockPath)
	if err != nil {
		return fmt.Errorf("svcroot: open registry lock: %w", err)
	}
	defer f.Close()
	h := windows.Handle(f.Fd())
	var ol windows.Overlapped
	if err := windows.LockFileEx(h, windows.LOCKFILE_EXCLUSIVE_LOCK, 0, 1, 0, &ol); err != nil {
		return fmt.Errorf("svcroot: acquire registry lock: %w", err)
	}
	defer func() { _ = windows.UnlockFileEx(h, 0, 1, 0, &ol) }()
	return fn()
}
