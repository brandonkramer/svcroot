# svcroot

Service instance roots, runtime file layout, project-root discovery, and known-instance registries for long-running local services.

## Install

```bash
go get github.com/brandonkramer/svcroot
```

## Configure once with Service

```go
svc := &svcroot.Service{
    EnvVar:       "MYAPP_HOME",
    DefaultDir:   ".myapp",
    RegistryFile: "sessions/known-instances.json", // optional
    Layout: svcroot.Layout{
        SessionsDir: "sessions",
        SocketName:  "rpc.sock",
        LockName:    "service.lock",
    },
}

root, err := svc.Resolve()
h := svc.Open(root)

socket := h.Socket()
lock := h.Lock()
tasksDir := h.Join("tasks")
registry := h.RuntimeJoin("state.json")
```

## Legacy helpers

```go
layout := &svcroot.DefaultLayout
socket := svcroot.Socket(root, layout)
lock := svcroot.Lock(root, layout)
```

## Development

```bash
make check          # race tests + linux compile + golangci-lint
make install-hooks  # lefthook pre-commit/pre-push
```
