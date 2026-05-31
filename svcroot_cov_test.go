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
	got := Sessions("/tmp/home", nil)
	if got == "" {
		t.Fatal("expected sessions path")
	}
}

func TestRegisterKnownHomeEmptyPath(t *testing.T) {
	stubPathAbs(t, func(string) (string, error) { return "", nil })
	if !errors.Is(RegisterKnownHome("", "", "now"), ErrEmptyHome) {
		t.Fatal("expected ErrEmptyHome")
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
	path, err := (&Service{}).RegistryPath()
	if err != nil || path != "" {
		t.Fatalf("path=%q err=%v", path, err)
	}
}

func TestHomeNilGuards(t *testing.T) {
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
	old := registryWithLock
	t.Cleanup(func() { registryWithLock = old })
	registryWithLock = func(path string, fn func() error) error {
		if cfg.mkdirErr != nil {
			return cfg.mkdirErr
		}
		if cfg.openErr != nil {
			return cfg.openErr
		}
		if cfg.flockErr != nil {
			return cfg.flockErr
		}
		return old(path, fn)
	}
}

func TestLayoutDefaultsAndNilReceivers(t *testing.T) {
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
	if _, ok := FindProjectRoot("\x00", ".svc"); ok {
		t.Fatal("expected false")
	}
}

func stubPathAbs(t *testing.T, fn func(string) (string, error)) {
	t.Helper()
	old := pathAbs
	t.Cleanup(func() { pathAbs = old })
	pathAbs = fn
}

func stubUserHomeDir(t *testing.T, fn func() (string, error)) {
	t.Helper()
	old := userHomeDir
	t.Cleanup(func() { userHomeDir = old })
	userHomeDir = fn
}

func stubRegistryDataMkdirAll(t *testing.T, fn func(string, os.FileMode) error) {
	t.Helper()
	old := registryDataMkdirAll
	t.Cleanup(func() { registryDataMkdirAll = old })
	registryDataMkdirAll = fn
}

func stubProjectPathAbs(t *testing.T, fn func(string) (string, error)) {
	t.Helper()
	old := projectPathAbs
	t.Cleanup(func() { projectPathAbs = old })
	projectPathAbs = fn
}

func TestRegisterKnownHomeAbsErrors(t *testing.T) {
	errAbs := errors.New("abs failed")
	stubPathAbs(t, func(string) (string, error) { return "", errAbs })
	if err := RegisterKnownHome(t.TempDir()+"/known.json", "home", "now"); !errors.Is(err, errAbs) {
		t.Fatalf("err=%v", err)
	}

	stubPathAbs(t, func(string) (string, error) { return "", nil })
	if !errors.Is(RegisterKnownHome("", "", "now"), ErrEmptyHome) {
		t.Fatal("expected ErrEmptyHome")
	}
}

func TestRegisterKnownHomeRegistryErrors(t *testing.T) {
	home := t.TempDir()
	if err := RegisterKnownHome("/dev/null/sub/known.json", home, "now"); err == nil {
		t.Fatal("expected mkdir error")
	}

	dir := t.TempDir()
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })
	if err := RegisterKnownHome(filepath.Join(dir, "known.json"), home, "now"); err == nil {
		t.Fatal("expected write error")
	}
}

func TestCandidateHomesBranches(t *testing.T) {
	dir := t.TempDir()
	reg := filepath.Join(dir, "known.json")
	home := filepath.Join(dir, "home")
	explicit := filepath.Join(dir, "explicit")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := RegisterKnownHome(reg, home, "now"); err != nil {
		t.Fatal(err)
	}

	errAbs := errors.New("abs failed")
	stubPathAbs(t, func(path string) (string, error) {
		if path == "bad" {
			return "", errAbs
		}
		return filepath.Abs(path)
	})
	got, err := CandidateHomes(reg, "SVC_HOME", "bad", explicit)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) < 2 {
		t.Fatalf("got=%v", got)
	}

	t.Setenv("SVC_HOME", explicit)
	got, err = CandidateHomes(reg, "SVC_HOME", "", home)
	if err != nil || len(got) < 2 {
		t.Fatalf("got=%v err=%v", got, err)
	}

	got, err = CandidateHomes(reg, "SVC_HOME", home, home)
	if err != nil || len(got) != 1 {
		t.Fatalf("got=%v err=%v", got, err)
	}
}

func TestWithCoreDefaultsFields(t *testing.T) {
	cases := []struct {
		layout Layout
		check  func(Layout) bool
	}{
		{Layout{SocketName: "s", LockName: "l"}, func(v Layout) bool { return v.SessionsDir == defaultSessionsDir }},
		{Layout{SessionsDir: "s", LockName: "l"}, func(v Layout) bool { return v.SocketName == defaultSocketName }},
		{Layout{SessionsDir: "s", SocketName: "s"}, func(v Layout) bool { return v.LockName == defaultLockName }},
	}
	for _, tc := range cases {
		got := withCoreDefaults(&tc.layout)
		if !tc.check(got) {
			t.Fatalf("layout=%+v got=%+v", tc.layout, got)
		}
	}
}

func TestFindProjectRootBranches(t *testing.T) {
	errAbs := errors.New("abs failed")
	stubProjectPathAbs(t, func(string) (string, error) { return "", errAbs })
	if _, ok := FindProjectRoot("x", ".svc"); ok {
		t.Fatal("expected false")
	}

	base := t.TempDir()
	markerFile := filepath.Join(base, ".svc")
	if err := os.WriteFile(markerFile, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, ok := FindProjectRoot(base, ".svc"); ok {
		t.Fatal("expected false for file marker")
	}
}

func TestServiceRegisterSuccess(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	svc := &Service{
		DefaultDir:   ".app",
		RegistryFile: "known.json",
	}
	root := filepath.Join(dir, "daemon")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := svc.Register(root, "now"); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterKnownHomeWriteError(t *testing.T) {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatal(err)
	}
	reg := filepath.Join(dir, "known.json")
	if err := os.Mkdir(reg, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := RegisterKnownHome(reg, home, "now"); err == nil {
		t.Fatal("expected write error")
	}
}

func TestRegisterKnownHomeSortsMultipleEntries(t *testing.T) {
	dir := t.TempDir()
	reg := filepath.Join(dir, "known.json")
	homeB := filepath.Join(dir, "b-home")
	homeA := filepath.Join(dir, "a-home")
	for _, home := range []string{homeB, homeA} {
		if err := os.MkdirAll(home, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := RegisterKnownHome(reg, homeB, "t1"); err != nil {
		t.Fatal(err)
	}
	if err := RegisterKnownHome(reg, homeA, "t2"); err != nil {
		t.Fatal(err)
	}
	file := LoadKnownHomes(reg)
	if len(file.Homes) != 2 || file.Homes[0].Path >= file.Homes[1].Path {
		t.Fatalf("file=%+v", file)
	}
}

func TestRegisterKnownHomeDataMkdirError(t *testing.T) {
	errMkdir := errors.New("data mkdir failed")
	stubRegistryDataMkdirAll(t, func(string, os.FileMode) error { return errMkdir })
	if err := RegisterKnownHome(filepath.Join(t.TempDir(), "known.json"), t.TempDir(), "now"); !errors.Is(err, errMkdir) {
		t.Fatalf("err=%v", err)
	}
}

func TestServiceCandidatesDefaultHomeError(t *testing.T) {
	calls := 0
	stubUserHomeDir(t, func() (string, error) {
		calls++
		if calls == 1 {
			return t.TempDir(), nil
		}
		return "", errors.New("default home failed")
	})
	_, err := (&Service{DefaultDir: ".app", RegistryFile: "known.json"}).Candidates("")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestServiceCandidatesRegistryPathError(t *testing.T) {
	stubUserHomeDir(t, func() (string, error) { return "", errors.New("home failed") })
	_, err := (&Service{DefaultDir: ".app", RegistryFile: "known.json"}).Candidates("")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestServiceRegisterRegistryPathError(t *testing.T) {
	stubUserHomeDir(t, func() (string, error) { return "", errors.New("home failed") })
	err := (&Service{DefaultDir: ".app", RegistryFile: "known.json"}).Register(t.TempDir(), "now")
	if err == nil {
		t.Fatal("expected error")
	}
}
