package svcroot_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/brandonkramer/svcroot"
)

func TestService(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "resolve and open",
			run: func(t *testing.T) {
				svc := &svcroot.Service{
					EnvVar:     "SVC_HOME",
					DefaultDir: ".svc",
					Layout: svcroot.Layout{
						SessionsDir: "run",
						SocketName:  "api.sock",
						LockName:    "api.lock",
					},
				}
				t.Setenv("SVC_HOME", "")
				home, err := svc.Resolve()
				if err != nil {
					t.Fatal(err)
				}
				h := svc.Open(home)
				if got := h.Root(); got != home {
					t.Fatalf("root=%q", got)
				}
				if got := h.Socket(); got != filepath.Join(home, "run", "api.sock") {
					t.Fatalf("socket=%q", got)
				}
				if got := h.Lock(); got != filepath.Join(home, "run", "api.lock") {
					t.Fatalf("lock=%q", got)
				}
				if _, ok := h.ObserveSocket(); ok {
					t.Fatal("expected observe socket to be disabled")
				}
				if _, ok := h.PipePrefix(); ok {
					t.Fatal("expected pipe prefix to be disabled")
				}
			},
		},
		{
			name: "registry",
			run: func(t *testing.T) {
				dir := t.TempDir()
				t.Setenv("HOME", dir)
				svc := &svcroot.Service{
					DefaultDir:   ".app",
					RegistryFile: "run/known.json",
					Layout:       svcroot.DefaultLayout,
				}
				reg, err := svc.RegistryPath()
				if err != nil {
					t.Fatal(err)
				}
				want := filepath.Join(dir, ".app", "run", "known.json")
				if reg != want {
					t.Fatalf("registry=%q want=%q", reg, want)
				}
				home := filepath.Join(dir, "daemon")
				if err := os.MkdirAll(home, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := svc.Register(home, "now"); err != nil {
					t.Fatal(err)
				}
				got, err := svc.Candidates("")
				if err != nil || len(got) == 0 {
					t.Fatalf("candidates=%v err=%v", got, err)
				}
			},
		},
		{
			name: "join paths",
			run: func(t *testing.T) {
				t.Parallel()
				h := (&svcroot.Service{Layout: svcroot.DefaultLayout}).Open("/data/home")
				if got := h.Join("tasks", "a.json"); got != "/data/home/tasks/a.json" {
					t.Fatalf("join=%q", got)
				}
				if got := h.RuntimeJoin("known.json"); got != "/data/home/sessions/known.json" {
					t.Fatalf("runtime join=%q", got)
				}
			},
		},
		{
			name: "registry not configured",
			run: func(t *testing.T) {
				t.Parallel()
				err := (&svcroot.Service{}).Register("/tmp/home", "now")
				if !errors.Is(err, svcroot.ErrRegistryNotConfigured) {
					t.Fatalf("err=%v", err)
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, tc.run)
	}
}
