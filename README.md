# lvm2go (Alpha)

[![Go Reference](https://pkg.go.dev/badge/github.com/jakobmoellerdev/lvm2go.svg)](https://pkg.go.dev/github.com/jakobmoellerdev/lvm2go)
[![Test](https://github.com/azalio/lvm2go/actions/workflows/test.yaml/badge.svg)](https://github.com/azalio/lvm2go/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jakobmoellerdev/lvm2go)](https://goreportcard.com/report/github.com/jakobmoellerdev/lvm2go)
[![License](https://img.shields.io/github/license/jakobmoellerdev/lvm2go)](https://github.com/azalio/lvm2go)

Package lvm2go implements a Go API for the lvm2 command line tools.

_This project is in Alpha stage and should not be used in production installations. Not all commands have been properly implemented and tested._

The API is designed to be simple and easy to use, while still providing
access to the full functionality of the LVM2 command line tools.

Compared to a simple command line wrapper, lvm2go provides a more structured
way to interact with lvm2, and allows for more complex interactions while safeguarding typing
and allowing for fine-grained control over the input of various usually problematic parameters,
such as sizes (and their conversion), validation of input parameters, and caching of data.

A simple usage example is shown below:

```go
package main

import (
 "context"
 "errors"
 "log/slog"
 "os"

 . "github.com/jakobmoellerdev/lvm2go"
)

func main() {
 if os.Geteuid() != 0 {
  panic("panicking because lvm2 requires root privileges for most operations.")
 }
 if err := run(); err != nil {
  slog.Error(err.Error())
  os.Exit(1)
 }
}

func run() (err error) {
 ctx := context.Background()
 lvm := NewClient()
 vgName := VolumeGroupName("test")
 lvName := LogicalVolumeName("test")
 deviceSize := MustParseSize("1G")
 lvSize := MustParseSize("100M")

 var losetup LoopbackDevice
 if losetup, err = NewLoopbackDevice(deviceSize); err != nil {
  return
 }
 defer func() {
  err = errors.Join(err, losetup.Close())
 }()

 if err = lvm.VGCreate(ctx, vgName, PhysicalVolumesFrom(losetup.Device())); err != nil {
  return
 }
 defer func() {
  err = errors.Join(err, lvm.VGRemove(ctx, vgName))
 }()

 if err = lvm.LVCreate(ctx, vgName, lvName, lvSize); err != nil {
  return
 }
 defer func() {
  err = errors.Join(err, lvm.LVRemove(ctx, vgName, lvName))
 }()

 return
}
```

## Containerization Support

When running inside a container (e.g., Docker, Kubernetes), `lvm2go` automatically uses `nsenter` to execute LVM commands in the host's mount namespace. This ensures that LVM operations correctly target the host system.

In some advanced scenarios, you might need to execute LVM commands directly within the container's namespace, even if `lvm2go` detects a containerized environment. There are two ways to achieve this:

### 1. Using the WithForceNoNsenter Context Option

The `WithForceNoNsenter` context option allows you to bypass nsenter for individual commands:

```go
import (
 "context"
 "github.com/azalio/lvm2go"
)

// ...

// Standard behavior: uses nsenter if containerized
cmdNormal := lvm2go.CommandContext(ctx, "lvs")

// Force execution without nsenter, even if containerized
ctxNoNsenter := lvm2go.WithForceNoNsenter(ctx, true)
cmdNoNsenter := lvm2go.CommandContext(ctxNoNsenter, "lvs")

// For debugging purposes, you can check if nsenter will be used
willUseNsenter := lvm2go.WillUseNsenter(ctx)
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

See the example at [`examples/no_nsenter_client/main.go`](examples/no_nsenter_client/main.go) for more details on using the client wrapper.

### Use Cases for Bypassing nsenter

There are several scenarios where bypassing the automatic nsenter behavior might be useful:

1. **Container-specific LVM operations**: When you need to perform LVM operations on volumes that are specific to the container's namespace, not the host.

2. **Testing and development**: When developing or testing LVM functionality within a containerized environment without affecting the host system.

3. **Nested containerization**: In complex setups with nested containers where the namespace hierarchy requires direct command execution.

4. **Custom namespace setups**: When your container has a custom namespace configuration that doesn't require nsenter for proper LVM operation.

5. **Debugging**: When troubleshooting issues related to namespace differences between the container and host.

See the example at [`examples/force_no_nsenter/main.go`](examples/force_no_nsenter/main.go) for more details on how to use this feature.

## Implemented commands by tested feature set

This set of commands is implemented and tested to some extent. The tested feature set is described in the table below.

| Command    | State | E2E Testing | Special Use Cases |
|------------|-------|-------------|-------------------|
| lvcreate   | Alpha | Basic       | Thin              |
| lvremove   | Alpha | Basic       | Thin              |
| lvextend   | Alpha | Basic       | Extents & Sizes   |
| lvchange   | Alpha | Basic       | (De-)Activation   |
| lvrename   | Alpha | Basic       |                   |
| lvs        | Alpha | Basic       |                   |
| vgcreate   | Alpha | Basic       |                   |
| vgremove   | Alpha | Basic       |                   |
| vgextend   | Alpha | Basic       |                   |
| vgreduce   | Alpha | Basic       |                   |
| vgchange   | Alpha | Basic       |                   |
| vgrename   | Alpha | Basic       |                   |
| vgs        | Alpha | Basic       |                   |
| pvs        | Alpha | Basic       |                   |
| pvcreate   | Alpha | Basic       |                   |
| pvchange   | Alpha | Basic       |                   |
| pvremove   | Alpha | Basic       |                   |
| pvmove     | Alpha | Basic       |                   |
| lvmdevices | Alpha | Basic       |                   |
| version    | Alpha | Basic       |                   |
