package svcroot

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/brandonkramer/jsonfile"
)

//
// ────────────────────────────────────────
// known-instance registry.
//

var (
	pathAbs              = filepath.Abs
	registryDataMkdirAll = os.MkdirAll
)

// HomeEntry records one known service root.
type HomeEntry struct {
	// Path is the absolute service root directory.
	Path string `json:"path"`
	// LastSeen is an optional RFC3339 timestamp from the last registration.
	LastSeen string `json:"last_seen,omitempty"`
}

// KnownHomes is the on-disk known-instance registry.
type KnownHomes struct {
	// Homes lists previously registered service roots.
	Homes []HomeEntry `json:"homes"`
}

// LoadKnownHomes reads the registry at path, returning an empty file when missing.
func LoadKnownHomes(path string) KnownHomes {
	file, err := jsonfile.Read[KnownHomes](path)
	if err != nil {
		return KnownHomes{}
	}
	return file
}

// RegisterKnownHome records root in the registry at registryPath.
func RegisterKnownHome(registryPath, home, now string) error {
	abs, err := pathAbs(home)
	if err != nil {
		return fmt.Errorf("svcroot: resolve home path: %w", err)
	}
	if abs == "" {
		return fmt.Errorf("svcroot: resolve home path: %w", ErrEmptyHome)
	}
	err = withRegistryLock(registryPath, func() error {
		if err := registryDataMkdirAll(filepath.Dir(registryPath), 0o755); err != nil {
			return fmt.Errorf("svcroot: create registry dir: %w", err)
		}
		file := LoadKnownHomes(registryPath)
		found := false
		for i := range file.Homes {
			if file.Homes[i].Path == abs {
				file.Homes[i].LastSeen = now
				found = true
				break
			}
		}
		if !found {
			file.Homes = append(file.Homes, HomeEntry{Path: abs, LastSeen: now})
		}
		sort.Slice(file.Homes, func(i, j int) bool { return file.Homes[i].Path < file.Homes[j].Path })
		if err := jsonfile.Write(registryPath, file); err != nil {
			return fmt.Errorf("svcroot: write registry %s: %w", registryPath, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("svcroot: register known home %s: %w", registryPath, err)
	}
	return nil
}

// CandidateHomes returns deduplicated service roots from envVar, envHome, defaultHome, and the registry file.
func CandidateHomes(registryPath, envVar, envHome, defaultHome string) ([]string, error) {
	seen := map[string]struct{}{}
	var out []string
	add := func(home string) {
		home = strings.TrimSpace(home)
		if home == "" {
			return
		}
		abs, err := pathAbs(home)
		if err != nil {
			return
		}
		if _, ok := seen[abs]; ok {
			return
		}
		seen[abs] = struct{}{}
		out = append(out, abs)
	}
	if envHome != "" {
		add(envHome)
	} else if envVar != "" {
		add(os.Getenv(envVar))
	}
	add(defaultHome)
	for _, entry := range LoadKnownHomes(registryPath).Homes {
		add(entry.Path)
	}
	sort.Strings(out)
	return out, nil
}
