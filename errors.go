package svcroot

import "errors"

// ErrEmptyHome is returned when a home path resolves to an empty string.
var ErrEmptyHome = errors.New("svcroot: empty home path")

// ErrRegistryNotConfigured is returned when registry operations run without RegistryFile.
var ErrRegistryNotConfigured = errors.New("svcroot: registry file not configured")
