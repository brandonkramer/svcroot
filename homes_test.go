package svcroot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/brandonkramer/svcroot"
)

func TestKnownHomesRegistry(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "missing file",
			run: func(t *testing.T) {
				t.Parallel()
				got := svcroot.LoadKnownHomes(filepath.Join(t.TempDir(), "missing.json"))
				if len(got.Homes) != 0 {
					t.Fatalf("got=%+v", got)
				}
			},
		},
		{
			name: "register and list",
			run: func(t *testing.T) {
				t.Parallel()
				dir := t.TempDir()
				reg := filepath.Join(dir, "known.json")
				home := filepath.Join(dir, "home")
				if err := os.MkdirAll(home, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := svcroot.RegisterKnownHome(reg, home, "t1"); err != nil {
					t.Fatal(err)
				}
				if err := svcroot.RegisterKnownHome(reg, home, "t2"); err != nil {
					t.Fatal(err)
				}
				file := svcroot.LoadKnownHomes(reg)
				if len(file.Homes) != 1 || file.Homes[0].LastSeen != "t2" {
					t.Fatalf("file=%+v", file)
				}
			},
		},
		{
			name: "candidate homes",
			run: func(t *testing.T) {
				t.Parallel()
				dir := t.TempDir()
				reg := filepath.Join(dir, "known.json")
				home := filepath.Join(dir, "service-home")
				if err := os.MkdirAll(home, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := svcroot.RegisterKnownHome(reg, home, "2026-01-01T00:00:00Z"); err != nil {
					t.Fatal(err)
				}
				got, err := svcroot.CandidateHomes(reg, "SVC_HOME", "", filepath.Join(dir, "default"))
				if err != nil || len(got) < 2 {
					t.Fatalf("got=%v err=%v", got, err)
				}
				got2, err := svcroot.CandidateHomes(reg, "SVC_HOME", dir+"/explicit", dir+"/default")
				if err != nil || len(got2) == 0 {
					t.Fatalf("got=%v err=%v", got2, err)
				}
				got3, err := svcroot.CandidateHomes("", "SVC_HOME", "/tmp/explicit-home", "/tmp/default-home")
				if err != nil || len(got3) == 0 {
					t.Fatalf("got=%v err=%v", got3, err)
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, tc.run)
	}
}

func TestRegisterKnownHomeErrors(t *testing.T) {
	t.Parallel()

	if err := svcroot.RegisterKnownHome("", "", "now"); err == nil {
		t.Fatal("expected error")
	}
}

func TestFindProjectRootMissing(t *testing.T) {
	t.Parallel()

	if _, ok := svcroot.FindProjectRoot(t.TempDir(), ".svc"); ok {
		t.Fatal("expected false")
	}
}

func TestDefaultHomeDir(t *testing.T) {
	t.Parallel()

	if _, err := svcroot.DefaultHomeDir(".svc"); err != nil {
		t.Fatal(err)
	}
}

func TestSessionsPartialLayout(t *testing.T) {
	t.Parallel()

	layout := svcroot.Layout{SessionsDir: "s"}
	got := svcroot.Sessions("/tmp/x", &layout)
	if got != "/tmp/x/s" {
		t.Fatalf("got=%q", got)
	}
}
