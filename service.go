package svcroot

import (
	"fmt"
	"path/filepath"
)

//
// ────────────────────────────────────────
// service configuration.
//

// Service describes how one long-running service resolves its home directory,
// lays out runtime files, and optionally tracks known homes on disk.
type Service struct {
	// EnvVar overrides DefaultDir when set in the process environment.
	EnvVar string
	// DefaultDir is appended to the user home when EnvVar is unset (e.g. ".myapp").
	DefaultDir string
	Layout     Layout
	// RegistryFile is joined with DefaultHome() for known-home registry paths.
	// Leave empty when the service does not persist known homes.
	RegistryFile string
}

// Resolve returns the configured home directory for this service.
func (s *Service) Resolve() (string, error) {
	return ResolveHome(s.EnvVar, s.DefaultDir)
}

// DefaultHome returns ~/.DefaultDir regardless of EnvVar.
func (s *Service) DefaultHome() (string, error) {
	return DefaultHomeDir(s.DefaultDir)
}

// Open returns path helpers for an explicit home root.
func (s *Service) Open(home string) *Home {
	layout := s.Layout
	return &Home{root: home, layout: &layout}
}

// RegistryPath returns the absolute known-homes registry path, or ("", nil) when unset.
func (s *Service) RegistryPath() (string, error) {
	if s.RegistryFile == "" {
		return "", nil
	}
	home, err := s.DefaultHome()
	if err != nil {
		return "", fmt.Errorf("svcroot: resolve registry home: %w", err)
	}
	return filepath.Join(home, s.RegistryFile), nil
}

// Register records home in this service's known-homes registry.
func (s *Service) Register(home, seenAt string) error {
	path, err := s.RegistryPath()
	if err != nil {
		return err
	}
	if path == "" {
		return fmt.Errorf("svcroot: register home: %w", ErrRegistryNotConfigured)
	}
	return RegisterKnownHome(path, home, seenAt)
}

// Candidates returns deduplicated homes from explicitHome, EnvVar, DefaultHome, and the registry.
func (s *Service) Candidates(explicitHome string) ([]string, error) {
	path, err := s.RegistryPath()
	if err != nil {
		return nil, err
	}
	defaultHome, err := s.DefaultHome()
	if err != nil {
		return nil, err
	}
	return CandidateHomes(path, s.EnvVar, explicitHome, defaultHome)
}
