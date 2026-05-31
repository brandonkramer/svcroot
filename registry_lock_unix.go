//go:build unix

package svcroot

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

//
// ────────────────────────────────────────
// registry locking (unix).
//

var (
	registryMkdirAll = os.MkdirAll
	registryOpenFile = func(lockPath string) (*os.File, error) {
		return os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	}
	registryFlock = syscall.Flock
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
	if err := registryFlock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		return fmt.Errorf("svcroot: acquire registry lock: %w", err)
	}
	defer func() { _ = registryFlock(int(f.Fd()), syscall.LOCK_UN) }()
	return fn()
}
