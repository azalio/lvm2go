/*
 Copyright 2024 The lvm2go Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package lvm2go_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	. "github.com/azalio/lvm2go"
)

func TestLVMDevices(t *testing.T) {
	t.Parallel()
	SkipOrFailTestIfNotRoot(t)

	_, err := exec.LookPath("lvmdevices")
	if err != nil {
		t.Skip("Skipping test because lvmdevices command is not found")
	}

	clnt := NewClient()
	ctx := context.Background()

	losetup := MakeTestLoopbackDevice(t, MustParseSize("1M"))

	devFile := DevicesFile(strings.ToLower(t.Name()))
	t.Cleanup(func() {
		if err := os.Remove(fmt.Sprintf("/etc/lvm/devices/%s", devFile)); err != nil && !os.IsNotExist(err) {
			t.Fatalf("Failed to remove devices file: %s", err)
		}
	})

	if err := clnt.DevModify(ctx, AddDevice(losetup.Device()), devFile); err != nil {
		t.Fatalf("Failed to add device to devices file: %s", err)
	}

	devs, err := clnt.DevList(ctx, devFile)
	if err != nil {
		t.Fatalf("Failed to list devices: %s", err)
	} else if len(devs) != 1 {
		t.Fatalf("Expected 1 device, got %d", len(devs))
	}

	verifyDevListEntryForLoopBackDevice(t, devs[0], losetup)

	if err := clnt.DevModify(ctx, DelDevice(losetup.Device()), devFile); err != nil {
		t.Fatalf("Failed to remove device from devices file: %s", err)
	}

	devs, err = clnt.DevList(ctx, devFile)
	if err != nil {
		t.Fatalf("Failed to list devices: %s", err)
	} else if len(devs) != 0 {
		t.Fatalf("Expected 0 devices, got %d", len(devs))
	}
}

func verifyDevListEntryForLoopBackDevice(
	t *testing.T,
	dev DeviceListEntry,
	losetup LoopbackDevice,
) {
	if dev.IDType != DeviceIDTypeLoopFile {
		t.Fatalf("Expected ID type %q, got %q", DeviceIDTypeLoopFile, dev.IDType)
	}
	if dev.IDName != losetup.File() {
		t.Fatalf("Expected ID name %q, got %q", losetup.File(), dev.IDName)
	}
	if dev.DevName != losetup.Device() {
		t.Fatalf("Expected dev name %q, got %q", losetup.Device(), dev.DevName)
	}
	if ePVID := "none"; dev.PVID != ePVID {
		t.Fatalf("Expected PVID %q, got %q", ePVID, dev.PVID)
	}
}
