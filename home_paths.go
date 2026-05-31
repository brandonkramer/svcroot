package svcroot

import "path/filepath"

//
// ────────────────────────────────────────
// home path helpers.
//

// Home resolves paths under one service home root.
type Home struct {
	root   string
	layout *Layout
}

// Root returns the service home directory.
func (h *Home) Root() string { return h.root }

// Runtime returns the runtime state directory under home (typically home/sessions).
func (h *Home) Runtime() string { return Sessions(h.root, h.layout) }

// Join joins path elements directly under the home root.
func (h *Home) Join(elem ...string) string {
	return filepath.Join(append([]string{h.root}, elem...)...)
}

// RuntimeJoin joins path elements under the runtime directory.
func (h *Home) RuntimeJoin(elem ...string) string {
	return filepath.Join(append([]string{h.Runtime()}, elem...)...)
}

// Socket returns the primary RPC socket path when configured.
func (h *Home) Socket() string { return Socket(h.root, h.layout) }

// Lock returns the service lock file path when configured.
func (h *Home) Lock() string { return Lock(h.root, h.layout) }

// ObserveSocket returns the observe socket path and whether this layout defines one.
func (h *Home) ObserveSocket() (string, bool) {
	if h == nil || h.layout == nil {
		return "", false
	}
	return h.layout.ObserveSocketPath(h.root)
}

// PipePrefix returns the Windows named-pipe prefix and whether this layout defines one.
func (h *Home) PipePrefix() (string, bool) {
	if h == nil || h.layout == nil {
		return "", false
	}
	return h.layout.PipePrefixName()
}
