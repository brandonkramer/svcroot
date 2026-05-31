package svcroot

import (
	"fmt"
	"path/filepath"
)

//
// ────────────────────────────────────────
// service configuration.
//

// Service describes how one long-running service resolves its root directory,
// lays out runtime files, and optionally tracks known instances on disk.
type Service struct {
	// EnvVar overrides DefaultDir when set in the process environment.
	EnvVar string
	// DefaultDir is appended to the user home when EnvVar is unset (e.g. ".myapp").
	DefaultDir string
	// Layout names runtime files under each service root.
	Layout Layout
	// RegistryFile is joined with DefaultHome() for known-instance registry paths.
	// Leave empty when the service does not persist known instances.
	RegistryFile string
}

// Resolve returns the configured service root directory.
func (s *Service) Resolve() (string, error) {
	return ResolveHome(s.EnvVar, s.DefaultDir)
}

// DefaultHome returns ~/.DefaultDir regardless of EnvVar.
func (s *Service) DefaultHome() (string, error) {
	return DefaultHomeDir(s.DefaultDir)
}

// Open returns path helpers for an explicit service root.
func (s *Service) Open(root string) *Home {
	layout := s.Layout
	return &Home{root: root, layout: &layout}
}

// RegistryPath returns the absolute known-instance registry path, or ("", nil) when unset.
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

// Register records root in this service's known-instance registry.
func (s *Service) Register(root, seenAt string) error {
	path, err := s.RegistryPath()
	if err != nil {
		return err
	}
	if path == "" {
		return fmt.Errorf("svcroot: register home: %w", ErrRegistryNotConfigured)
	}
	return RegisterKnownHome(path, root, seenAt)
}

// Candidates returns deduplicated service roots from explicitRoot, EnvVar, DefaultHome, and the registry.
func (s *Service) Candidates(explicitRoot string) ([]string, error) {
	path, err := s.RegistryPath()
	if err != nil {
		return nil, err
	}
	defaultHome, err := s.DefaultHome()
	if err != nil {
		return nil, err
	}
	return CandidateHomes(path, s.EnvVar, explicitRoot, defaultHome)
}
