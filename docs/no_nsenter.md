# NoNsenter Feature in lvm2go

This document explains the NoNsenter feature in lvm2go, which allows users to bypass the automatic nsenter behavior when running in a containerized environment.

## Background

When running inside a container (e.g., Docker, Kubernetes), `lvm2go` automatically uses `nsenter` to execute LVM commands in the host's mount namespace. This ensures that LVM operations correctly target the host system.

However, in some scenarios, you might need to execute LVM commands directly within the container's namespace, even if `lvm2go` detects a containerized environment.

## Use Cases for Bypassing nsenter

There are several scenarios where bypassing the automatic nsenter behavior might be useful:

1. **Container-specific LVM operations**: When you need to perform LVM operations on volumes that are specific to the container's namespace, not the host.

2. **Testing and development**: When developing or testing LVM functionality within a containerized environment without affecting the host system.

3. **Nested containerization**: In complex setups with nested containers where the namespace hierarchy requires direct command execution.

4. **Custom namespace setups**: When your container has a custom namespace configuration that doesn't require nsenter for proper LVM operation.

5. **Debugging**: When troubleshooting issues related to namespace differences between the container and host.

## Implementation Options

lvm2go provides two ways to bypass the automatic nsenter behavior:

### 1. Using the WithForceNoNsenter Context Option

The `WithForceNoNsenter` context option allows you to bypass nsenter for individual commands:

```go
import (
    "context"
    "github.com/azalio/lvm2go"
)

// Standard behavior: uses nsenter if containerized
ctx := context.Background()
cmdNormal := lvm2go.CommandContext(ctx, "lvs")

// Force execution without nsenter, even if containerized
ctxNoNsenter := lvm2go.WithForceNoNsenter(ctx, true)
cmdNoNsenter := lvm2go.CommandContext(ctxNoNsenter, "lvs")
```

This approach requires modifying the context for each command, which can be cumbersome when using the client interface:

```go
// Need to create a new context for each operation
ctxNoNsenter := lvm2go.WithForceNoNsenter(ctx, true)
if err := client.VGRemove(ctxNoNsenter, name, lvm2go.Force(force)); err != nil {
    // Handle error
}
```

### 2. Using the WithNoNsenter Client Wrapper

For a more convenient approach, especially when using the client in a struct like a Kubernetes controller, you can use the `WithNoNsenter` client wrapper:

```go
import (
    "context"
    "github.com/azalio/lvm2go"
)

// Create a standard client
standardClient := lvm2go.NewClient()

// Create a client that will never use nsenter
noNsenterClient := lvm2go.WithNoNsenter(standardClient)

// All operations with this client will bypass nsenter
vgs, err := noNsenterClient.VGs(ctx)
```

This is particularly useful in structures like Kubernetes controllers:

```go
type VolumeGroupReconciler struct {
    client.Client
    LVM          lvm2go.Client
    NodeName     string
}

func NewReconciler(mgr manager.Manager) *VolumeGroupReconciler {
    return &VolumeGroupReconciler{
        Client: mgr.GetClient(),
        // Create an LVM client that never uses nsenter
        LVM: lvm2go.WithNoNsenter(lvm2go.NewClient()),
    }
}

func (r *VolumeGroupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // All LVM operations will bypass nsenter
    if err := r.LVM.VGRemove(ctx, name, lvm2go.Force(force)); err != nil {
        return ctrl.Result{}, err
    }
    return ctrl.Result{}, nil
}
```

## How It Works

The `WithNoNsenter` function returns a new client that wraps the original client and automatically applies the NoNsenter context to all operations. This is implemented as a client wrapper that intercepts all method calls, applies the NoNsenter context, and then delegates to the original client.

Internally, it uses the `WithForceNoNsenter` context option to set the NoNsenter flag on the context before passing it to the original client.

## Examples

See the following examples for more details:

- [`examples/force_no_nsenter/main.go`](../examples/force_no_nsenter/main.go) - Using the WithForceNoNsenter context option
- [`examples/no_nsenter_client/main.go`](../examples/no_nsenter_client/main.go) - Using the WithNoNsenter client wrapper