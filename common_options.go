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

package lvm2go

type CommonOptions struct {
	Devices
	DevicesFile
	Profile
	Verbose
	RequestConfirm
}

func (opts CommonOptions) ApplyToArgs(args Arguments) error {
	for _, arg := range []Argument{
		opts.Devices,
		opts.DevicesFile,
		opts.Verbose,
		opts.RequestConfirm,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}

type RequestConfirm bool

func (opt RequestConfirm) ApplyToArgs(args Arguments) error {
	if !opt {
		args.AddOrReplace("--yes")
	}
	return nil
}

type Verbose bool

func (opt Verbose) ApplyToArgs(args Arguments) error {
	if opt {
		args.AddOrReplace("--verbose")
	}
	return nil
}
