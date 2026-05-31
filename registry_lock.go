package svcroot

import (
	"context"

	"github.com/brandonkramer/filelock"
)

//
// ────────────────────────────────────────
// known-homes registry locking.
//

var registryWithLock = func(path string, fn func() error) error {
	return filelock.WithSidecar(context.Background(), path, filelock.DefaultSidecar, fn)
}

func withRegistryLock(path string, fn func() error) error {
	return registryWithLock(path, fn)
}
