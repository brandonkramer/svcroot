# svcroot

Service root directories, runtime file layout, project-root discovery, and known-instance registries

## Configure once with Service

```go
svc := &svcroot.Service{
    EnvVar:       "MYAPP_HOME",
    DefaultDir:   ".myapp",
    RegistryFile: "sessions/known-homes.json", // optional
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
