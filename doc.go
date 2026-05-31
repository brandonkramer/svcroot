// Package svcroot resolves service instance roots, runtime file layout,
// project-root discovery, and optional known-instance registries.
//
// Configure a [Service] once, then call [Service.Resolve] and [Service.Open]
// for path helpers. Use [Layout] directly when you prefer functional helpers.
//
// # Examples
//
//	svc := &svcroot.Service{
//	    EnvVar:     "MYAPP_HOME",
//	    DefaultDir: ".myapp",
//	    Layout:     svcroot.DefaultLayout,
//	}
//	root, err := svc.Resolve()
//	socket := svc.Open(root).Socket()
package svcroot
