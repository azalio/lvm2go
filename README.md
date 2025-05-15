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

In some advanced scenarios, you might need to execute LVM commands directly within the container's namespace, even if `lvm2go` detects a containerized environment. For this purpose, the `WithForceNoNsenter` context option is provided:

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
```

See the example at [`examples/force_no_nsenter/main.go`](examples/force_no_nsenter/main.go) for more details.

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
