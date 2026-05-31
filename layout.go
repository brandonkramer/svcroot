package svcroot

import "path/filepath"

//
// ────────────────────────────────────────
// layout paths.
//

// Layout names runtime files under a service root directory.
type Layout struct {
	// SessionsDir is the runtime state directory name under the service root.
	SessionsDir string
	// SocketName is the primary RPC socket file name inside SessionsDir.
	SocketName string
	// ObserveSocketName is the observe socket file name; leave empty to disable.
	ObserveSocketName string
	// LockName is the service lock file name inside SessionsDir.
	LockName string
	// PipePrefix is the Windows named-pipe prefix; leave empty to disable.
	PipePrefix string
}

const (
	defaultSessionsDir       = "sessions"
	defaultSocketName        = "rpc.sock"
	defaultObserveSocketName = "observe.sock"
	defaultLockName          = "service.lock"
	defaultPipePrefix        = "svc"
)

// DefaultLayout fills empty Layout fields with conventional names.
var DefaultLayout = Layout{
	SessionsDir:       defaultSessionsDir,
	SocketName:        defaultSocketName,
	ObserveSocketName: defaultObserveSocketName,
	LockName:          defaultLockName,
	PipePrefix:        defaultPipePrefix,
}

func layoutValue(layout *Layout) Layout {
	if layout == nil {
		return DefaultLayout
	}
	return *layout
}

func withCoreDefaults(l *Layout) Layout {
	v := layoutValue(l)
	if v.SessionsDir == "" {
		v.SessionsDir = defaultSessionsDir
	}
	if v.SocketName == "" {
		v.SocketName = defaultSocketName
	}
	if v.LockName == "" {
		v.LockName = defaultLockName
	}
	return v
}

func withDefaults(l *Layout) Layout {
	v := withCoreDefaults(l)
	if v.ObserveSocketName == "" {
		v.ObserveSocketName = defaultObserveSocketName
	}
	if v.PipePrefix == "" {
		v.PipePrefix = defaultPipePrefix
	}
	return v
}

// Sessions returns root joined with the configured runtime directory.
func Sessions(root string, layout *Layout) string {
	l := withCoreDefaults(layout)
	return filepath.Join(root, l.SessionsDir)
}

// Socket returns the Unix domain socket path under root.
func Socket(root string, layout *Layout) string {
	l := withCoreDefaults(layout)
	return filepath.Join(root, l.SessionsDir, l.SocketName)
}

// Lock returns the lock file path under root.
func Lock(root string, layout *Layout) string {
	l := withCoreDefaults(layout)
	return filepath.Join(root, l.SessionsDir, l.LockName)
}

// ObserveSocket returns the read-only observe HTTP socket path under root.
func ObserveSocket(root string, layout *Layout) string {
	l := withDefaults(layout)
	return filepath.Join(root, l.SessionsDir, l.ObserveSocketName)
}

// WithDefaults returns layout with empty core fields filled from DefaultLayout.
func (l *Layout) WithDefaults() Layout {
	if l == nil {
		return withDefaults(nil)
	}
	return withDefaults(l)
}

// RuntimeDir returns root joined with the runtime directory name.
func (l *Layout) RuntimeDir(root string) string {
	return Sessions(root, l)
}

// At joins path elements under root. When underRuntime is true, the runtime directory is inserted first.
func (l *Layout) At(root string, underRuntime bool, elem ...string) string {
	base := root
	if underRuntime {
		base = Sessions(root, l)
	}
	if len(elem) == 0 {
		return base
	}
	return filepath.Join(append([]string{base}, elem...)...)
}

// ObserveSocketPath returns the observe socket path when ObserveSocketName is set.
func (l *Layout) ObserveSocketPath(root string) (string, bool) {
	if l == nil || l.ObserveSocketName == "" {
		return "", false
	}
	core := withCoreDefaults(l)
	return filepath.Join(root, core.SessionsDir, l.ObserveSocketName), true
}

// PipePrefixName returns the configured Windows pipe prefix when set.
func (l *Layout) PipePrefixName() (string, bool) {
	if l == nil || l.PipePrefix == "" {
		return "", false
	}
	return l.PipePrefix, true
}
