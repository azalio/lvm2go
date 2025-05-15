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

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/azalio/lvm2go"
)

// This example demonstrates how to use the WithNoNsenter client wrapper
// to create a client that will never use nsenter, even in a containerized environment.
func main() {
	ctx := context.Background()

	// Check if we're running in a containerized environment
	isContainerized := lvm2go.IsContainerized(ctx)
	fmt.Printf("Running in containerized environment: %v\n", isContainerized)

	// Create a standard client
	standardClient := lvm2go.NewClient()

	// Create a client that will never use nsenter
	noNsenterClient := lvm2go.WithNoNsenter(standardClient)

	// Example 1: Using the standard client
	fmt.Println("\nExample 1: Using standard client")
	// This will use nsenter if in a containerized environment
	vgs1, err := standardClient.VGs(ctx)
	if err != nil {
		log.Fatalf("Failed to list volume groups with standard client: %v", err)
	}
	fmt.Printf("Found %d volume groups with standard client\n", len(vgs1))

	// Example 2: Using the NoNsenter client
	fmt.Println("\nExample 2: Using NoNsenter client")
	// This will never use nsenter, even in a containerized environment
	vgs2, err := noNsenterClient.VGs(ctx)
	if err != nil {
		log.Fatalf("Failed to list volume groups with NoNsenter client: %v", err)
	}
	fmt.Printf("Found %d volume groups with NoNsenter client\n", len(vgs2))

	// Example 3: Using the NoNsenter client in a reconciler-like structure
	fmt.Println("\nExample 3: Using NoNsenter client in a reconciler-like structure")
	reconciler := &VolumeGroupReconciler{
		LVM: noNsenterClient,
	}
	if err := reconciler.reconcileVolumeGroup(ctx, "example-vg"); err != nil {
		fmt.Printf("Reconciliation error: %v\n", err)
	} else {
		fmt.Println("Reconciliation successful")
	}

	// If we're in a container, explain the difference
	if isContainerized {
		fmt.Println("\nNote: Since we're in a containerized environment:")
		fmt.Println("- Example 1 uses nsenter to access the host's LVM")
		fmt.Println("- Example 2 directly accesses the container's LVM (without nsenter)")
		fmt.Println("- Example 3 demonstrates how to use the NoNsenter client in a reconciler structure")
	} else {
		fmt.Println("\nNote: Since we're not in a containerized environment, both clients behave the same")
	}

	// Exit with success
	os.Exit(0)
}

// VolumeGroupReconciler is an example structure similar to what users might have
// in a Kubernetes operator or other reconciliation-based system.
type VolumeGroupReconciler struct {
	// In a real application, this might embed a controller-runtime Client
	// client.Client
	
	// LVM client for LVM operations
	LVM lvm2go.Client
	
	// Other fields a reconciler might have
	NodeName     string
	SyncInterval string
}

// reconcileVolumeGroup is an example method that demonstrates how the NoNsenter client
// can be used in a reconciler structure.
func (r *VolumeGroupReconciler) reconcileVolumeGroup(ctx context.Context, name string) error {
	// Check if the volume group exists
	vg, err := r.LVM.VG(ctx, lvm2go.VolumeGroupName(name))
	if err != nil {
		if err == lvm2go.ErrVolumeGroupNotFound {
			fmt.Printf("Volume group %s not found\n", name)
			return nil
		}
		return fmt.Errorf("failed to get volume group: %w", err)
	}

	fmt.Printf("Found volume group: %s\n", vg.Name)

	// Example of removing a volume group with force option
	// This is commented out to prevent accidental deletion
	/*
	force := true
	if err := r.LVM.VGRemove(ctx, lvm2go.VolumeGroupName(name), lvm2go.Force(force)); err != nil {
		return fmt.Errorf("failed to remove volume group: %w", err)
	}
	*/

	return nil
}