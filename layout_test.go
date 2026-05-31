package svcroot

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLayoutPaths(t *testing.T) {
	t.Parallel()

	home := filepath.Join(t.TempDir(), "home")
	cases := []struct {
		name string
		fn   func(string, *Layout) string
		want string
	}{
		{
			name: "socket default",
			fn:   Socket,
			want: filepath.Join(home, "sessions", "rpc.sock"),
		},
		{
			name: "lock default",
			fn:   Lock,
			want: filepath.Join(home, "sessions", "service.lock"),
		},
		{
			name: "observe default",
			fn:   ObserveSocket,
			want: filepath.Join(home, "sessions", "observe.sock"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.fn(home, nil); got != tc.want {
				t.Fatalf("got=%q want=%q", got, tc.want)
			}
		})
	}
}

func TestResolveHome(t *testing.T) {
	cases := []struct {
		name    string
		env     string
		envVal  string
		defDir  string
		want    string
		wantErr bool
	}{
		{
			name:   "default dir",
			env:    "SVC_HOME",
			defDir: ".svc",
			want:   mustDefaultHome(t, ".svc"),
		},
		{
			name:   "env override",
			env:    "SVC_HOME",
			envVal: "/tmp/custom-home",
			defDir: ".svc",
			want:   "/tmp/custom-home",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(tc.env, tc.envVal)
			got, err := ResolveHome(tc.env, tc.defDir)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("got=%q want=%q", got, tc.want)
			}
		})
	}
}

func TestFindProjectRoot(t *testing.T) {
	t.Parallel()

	base := t.TempDir()
	project := filepath.Join(base, "repo", ".svc")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(base, "repo", "pkg", "inner")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}

	got, ok := FindProjectRoot(sub, ".svc")
	if !ok || got != project {
		t.Fatalf("FindProjectRoot()=%q ok=%v want=%q", got, ok, project)
	}
}

func TestLayoutHelpers(t *testing.T) {
	t.Parallel()

	layout := Layout{SessionsDir: "state", ObserveSocketName: "observe.sock", PipePrefix: "svc"}
	if got := layout.RuntimeDir("/home"); got != "/home/state" {
		t.Fatalf("runtime dir=%q", got)
	}
	if got := layout.At("/home", true, "cache", "x"); got != "/home/state/cache/x" {
		t.Fatalf("at=%q", got)
	}
	if got := layout.At("/home", false); got != "/home" {
		t.Fatalf("at base=%q", got)
	}
	if got := layout.WithDefaults().SocketName; got != "rpc.sock" {
		t.Fatalf("with defaults socket=%q", got)
	}
	if path, ok := (&Layout{ObserveSocketName: "observe.sock", SessionsDir: "run"}).ObserveSocketPath("/home"); !ok || path != "/home/run/observe.sock" {
		t.Fatalf("observe path=%q ok=%v", path, ok)
	}
	if _, ok := (&Layout{}).ObserveSocketPath("/home"); ok {
		t.Fatal("expected observe disabled")
	}
	if prefix, ok := layout.PipePrefixName(); !ok || prefix != "svc" {
		t.Fatalf("prefix=%q ok=%v", prefix, ok)
	}
	if _, ok := (*Layout)(nil).PipePrefixName(); ok {
		t.Fatal("expected nil layout pipe disabled")
	}
}

func mustDefaultHome(t *testing.T, dir string) string {
	t.Helper()
	got, err := DefaultHomeDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	return got
}
