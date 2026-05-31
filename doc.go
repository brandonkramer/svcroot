// Package svcroot resolves service home directories, runtime file layout,
// project-root discovery, and optional known-home registries.
//
// Configure a [Service] once, then call [Service.Resolve] and [Service.Open]
// for path helpers. Use [Layout] directly when you prefer functional helpers.
//
// # Examples
//
//	svc := svcroot.Service{
//	    EnvVar:     "MYAPP_HOME",
//	    DefaultDir: ".myapp",
//	    Layout:     svcroot.DefaultLayout,
//	}
//	home, err := svc.Resolve()
//	socket := svc.Open(home).Socket()
package svcroot
