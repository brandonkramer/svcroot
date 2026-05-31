package svcroot

import "path/filepath"

//
// ────────────────────────────────────────
// layout paths.
//

// Layout names runtime files under a service home directory.
type Layout struct {
	SessionsDir       string
	SocketName        string
	ObserveSocketName string
	LockName          string
	PipePrefix        string
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

// Sessions returns home/sessions (or the configured sessions dir).
func Sessions(home string, layout *Layout) string {
	l := withCoreDefaults(layout)
	return filepath.Join(home, l.SessionsDir)
}

// Socket returns the Unix domain socket path for home.
func Socket(home string, layout *Layout) string {
	l := withCoreDefaults(layout)
	return filepath.Join(home, l.SessionsDir, l.SocketName)
}

// Lock returns the lock file path for home.
func Lock(home string, layout *Layout) string {
	l := withCoreDefaults(layout)
	return filepath.Join(home, l.SessionsDir, l.LockName)
}

// ObserveSocket returns the read-only observe HTTP socket path for home.
func ObserveSocket(home string, layout *Layout) string {
	l := withDefaults(layout)
	return filepath.Join(home, l.SessionsDir, l.ObserveSocketName)
}

// WithDefaults returns layout with empty core fields filled from DefaultLayout.
func (l *Layout) WithDefaults() Layout {
	if l == nil {
		return withDefaults(nil)
	}
	return withDefaults(l)
}

// RuntimeDir returns home joined with the runtime directory name.
func (l *Layout) RuntimeDir(home string) string {
	return Sessions(home, l)
}

// At joins path elements under home. When underRuntime is true, the runtime directory is inserted first.
func (l *Layout) At(home string, underRuntime bool, elem ...string) string {
	base := home
	if underRuntime {
		base = Sessions(home, l)
	}
	if len(elem) == 0 {
		return base
	}
	return filepath.Join(append([]string{base}, elem...)...)
}

// ObserveSocketPath returns the observe socket path when ObserveSocketName is set.
func (l *Layout) ObserveSocketPath(home string) (string, bool) {
	if l == nil || l.ObserveSocketName == "" {
		return "", false
	}
	core := withCoreDefaults(l)
	return filepath.Join(home, core.SessionsDir, l.ObserveSocketName), true
}

// PipePrefixName returns the configured Windows pipe prefix when set.
func (l *Layout) PipePrefixName() (string, bool) {
	if l == nil || l.PipePrefix == "" {
		return "", false
	}
	return l.PipePrefix, true
}
