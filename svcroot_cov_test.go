package svcroot

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultHomeDirError(t *testing.T) {
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")
	if _, err := DefaultHomeDir(".svc"); err == nil {
		t.Skip("user home unavailable on this system")
	}
}

func TestSessionsWithPartialLayout(t *testing.T) {
	t.Parallel()

	got := Sessions("/tmp/home", nil)
	if got == "" {
		t.Fatal("expected sessions path")
	}
}

func TestRegisterKnownHomeEmptyPath(t *testing.T) {
	t.Parallel()

	err := RegisterKnownHome("", "", "now")
	if !errors.Is(err, ErrEmptyHome) && err == nil {
		t.Fatalf("err=%v", err)
	}
}

func TestRegistryLockInjectedErrors(t *testing.T) {
	dir := t.TempDir()
	reg := filepath.Join(dir, "known.json")
	errMkdir := errors.New("mkdir failed")
	errOpen := errors.New("open failed")
	errFlock := errors.New("flock failed")

	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "mkdir",
			run: func(t *testing.T) {
				stubRegistryLock(t, registryLockStub{mkdirErr: errMkdir})
				if err := RegisterKnownHome(reg, dir, "now"); !errors.Is(err, errMkdir) {
					t.Fatalf("err=%v", err)
				}
			},
		},
		{
			name: "open",
			run: func(t *testing.T) {
				stubRegistryLock(t, registryLockStub{openErr: errOpen})
				if err := RegisterKnownHome(reg, dir, "now"); !errors.Is(err, errOpen) {
					t.Fatalf("err=%v", err)
				}
			},
		},
		{
			name: "flock",
			run: func(t *testing.T) {
				stubRegistryLock(t, registryLockStub{flockErr: errFlock})
				if err := RegisterKnownHome(reg, dir, "now"); !errors.Is(err, errFlock) {
					t.Fatalf("err=%v", err)
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, tc.run)
	}
}

func TestServiceRegistryPathEmpty(t *testing.T) {
	t.Parallel()

	path, err := (&Service{}).RegistryPath()
	if err != nil || path != "" {
		t.Fatalf("path=%q err=%v", path, err)
	}
}

func TestHomeNilGuards(t *testing.T) {
	t.Parallel()

	var h *Home
	if _, ok := h.ObserveSocket(); ok {
		t.Fatal("expected false")
	}
	if _, ok := h.PipePrefix(); ok {
		t.Fatal("expected false")
	}
}

type registryLockStub struct {
	mkdirErr error
	openErr  error
	flockErr error
}

func stubRegistryLock(t *testing.T, cfg registryLockStub) {
	t.Helper()
	oldMkdir := registryMkdirAll
	oldOpen := registryOpenFile
	oldFlock := registryFlock
	t.Cleanup(func() {
		registryMkdirAll = oldMkdir
		registryOpenFile = oldOpen
		registryFlock = oldFlock
	})
	if cfg.mkdirErr != nil {
		registryMkdirAll = func(string, os.FileMode) error { return cfg.mkdirErr }
	}
	if cfg.openErr != nil {
		registryOpenFile = func(string) (*os.File, error) { return nil, cfg.openErr }
	}
	if cfg.flockErr != nil {
		registryFlock = func(int, int) error { return cfg.flockErr }
	}
}

func TestLayoutDefaultsAndNilReceivers(t *testing.T) {
	t.Parallel()

	var nilLayout *Layout
	if got := nilLayout.WithDefaults().PipePrefix; got != "svc" {
		t.Fatalf("pipe=%q", got)
	}
	if got := nilLayout.RuntimeDir("/home"); got != "/home/sessions" {
		t.Fatalf("runtime=%q", got)
	}
	if got := nilLayout.At("/home", false, "x"); got != "/home/x" {
		t.Fatalf("at=%q", got)
	}
	partial := Layout{SessionsDir: "run", SocketName: "s.sock", LockName: "l.lock"}
	if got := ObserveSocket("/home", &partial); got != "/home/run/observe.sock" {
		t.Fatalf("observe=%q", got)
	}
}

func TestServiceDefaultHomeError(t *testing.T) {
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")
	_, err := (&Service{DefaultDir: ".app", RegistryFile: "known.json"}).RegistryPath()
	if err == nil {
		t.Skip("user home unavailable on this system")
	}
	_, err = (&Service{DefaultDir: ".app", RegistryFile: "known.json"}).Candidates("")
	if err == nil {
		t.Skip("user home unavailable on this system")
	}
}

func TestFindProjectRootInvalidFrom(t *testing.T) {
	t.Parallel()

	if _, ok := FindProjectRoot("\x00", ".svc"); ok {
		t.Fatal("expected false")
	}
}
