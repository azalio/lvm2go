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

func main() {
	ctx := context.Background()

	// Check if we're running in a containerized environment
	isContainerized := lvm2go.IsContainerized(ctx)
	fmt.Printf("Running in containerized environment: %v\n", isContainerized)
	
	// Check if nsenter will be used with the default context
	willUseNsenter := lvm2go.WillUseNsenter(ctx)
	fmt.Printf("Will use nsenter with default context: %v\n", willUseNsenter)

	// Example 1: Standard behavior (uses nsenter if containerized)
	fmt.Println("\nExample 1: Standard behavior")
	cmd1 := lvm2go.CommandContext(ctx, "ls", "-l", "/")
	fmt.Printf("Command path: %s\n", cmd1.Path)
	fmt.Printf("Command args: %v\n", cmd1.Args)

	// Run the command
	output1, err := cmd1.CombinedOutput()
	if err != nil {
		log.Fatalf("Command failed: %v", err)
	}
	fmt.Printf("Command output:\n%s\n", string(output1))

	// Example 2: Force no nsenter, even if containerized
	fmt.Println("\nExample 2: Force no nsenter")
	ctxNoNsenter := lvm2go.WithForceNoNsenter(ctx, true)
	
	// Check if nsenter will be used with the modified context
	willUseNsenterWithModifiedCtx := lvm2go.WillUseNsenter(ctxNoNsenter)
	fmt.Printf("Will use nsenter with modified context: %v\n", willUseNsenterWithModifiedCtx)
	
	cmd2 := lvm2go.CommandContext(ctxNoNsenter, "ls", "-l", "/")
	fmt.Printf("Command path: %s\n", cmd2.Path)
	fmt.Printf("Command args: %v\n", cmd2.Args)

	// Run the command
	output2, err := cmd2.CombinedOutput()
	if err != nil {
		log.Fatalf("Command failed: %v", err)
	}
	fmt.Printf("Command output:\n%s\n", string(output2))

	// If we're in a container, the outputs should be different
	if isContainerized {
		fmt.Println("\nNote: Since we're in a containerized environment, the outputs should be different:")
		fmt.Println("- Example 1 shows the host's root directory (via nsenter)")
		fmt.Println("- Example 2 shows the container's root directory (direct execution)")
	} else {
		fmt.Println("\nNote: Since we're not in a containerized environment, both commands executed directly")
		fmt.Println("and the outputs should be identical.")
	}

	// Exit with success
	os.Exit(0)
}
