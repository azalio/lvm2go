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

import (
	"context"
	"io"
)

// noNsenterClient is a client wrapper that applies the NoNsenter context to all operations.
// noNsenterClient is created using the WithNoNsenter function in client.go
type noNsenterClient struct {
	client Client
}

// applyNoNsenter applies the NoNsenter context to the given context.
func (c *noNsenterClient) applyNoNsenter(ctx context.Context) context.Context {
	return WithForceNoNsenter(ctx, true)
}

// Ensure noNsenterClient implements Client
var _ Client = (*noNsenterClient)(nil)

// Version implements MetaClient.
func (c *noNsenterClient) Version(ctx context.Context, opts ...VersionOption) (Version, error) {
	return c.client.Version(c.applyNoNsenter(ctx), opts...)
}

// RawConfig implements MetaClient.
func (c *noNsenterClient) RawConfig(ctx context.Context, opts ...ConfigOption) (RawConfig, error) {
	return c.client.RawConfig(c.applyNoNsenter(ctx), opts...)
}

// ReadAndDecodeConfig implements MetaClient.
func (c *noNsenterClient) ReadAndDecodeConfig(ctx context.Context, v any, opts ...ConfigOption) error {
	return c.client.ReadAndDecodeConfig(c.applyNoNsenter(ctx), v, opts...)
}

// WriteAndEncodeConfig implements MetaClient.
func (c *noNsenterClient) WriteAndEncodeConfig(ctx context.Context, v any, writer io.Writer) error {
	return c.client.WriteAndEncodeConfig(c.applyNoNsenter(ctx), v, writer)
}

// UpdateGlobalConfig implements MetaClient.
func (c *noNsenterClient) UpdateGlobalConfig(ctx context.Context, v any) error {
	return c.client.UpdateGlobalConfig(c.applyNoNsenter(ctx), v)
}

// UpdateLocalConfig implements MetaClient.
func (c *noNsenterClient) UpdateLocalConfig(ctx context.Context, v any) error {
	return c.client.UpdateLocalConfig(c.applyNoNsenter(ctx), v)
}

// UpdateProfileConfig implements MetaClient.
func (c *noNsenterClient) UpdateProfileConfig(ctx context.Context, v any, profile Profile) error {
	return c.client.UpdateProfileConfig(c.applyNoNsenter(ctx), v, profile)
}

// CreateProfile implements MetaClient.
func (c *noNsenterClient) CreateProfile(ctx context.Context, v any, profile Profile) (string, error) {
	return c.client.CreateProfile(c.applyNoNsenter(ctx), v, profile)
}

// RemoveProfile implements MetaClient.
func (c *noNsenterClient) RemoveProfile(ctx context.Context, profile Profile) error {
	return c.client.RemoveProfile(c.applyNoNsenter(ctx), profile)
}

// GetProfilePath implements MetaClient.
func (c *noNsenterClient) GetProfilePath(ctx context.Context, profile Profile) (string, error) {
	return c.client.GetProfilePath(c.applyNoNsenter(ctx), profile)
}

// GetProfileDirectory implements MetaClient.
func (c *noNsenterClient) GetProfileDirectory(ctx context.Context) (string, error) {
	return c.client.GetProfileDirectory(c.applyNoNsenter(ctx))
}

// VG implements VolumeGroupClient.
func (c *noNsenterClient) VG(ctx context.Context, opts ...VGsOption) (*VolumeGroup, error) {
	return c.client.VG(c.applyNoNsenter(ctx), opts...)
}

// VGs implements VolumeGroupClient.
func (c *noNsenterClient) VGs(ctx context.Context, opts ...VGsOption) ([]*VolumeGroup, error) {
	return c.client.VGs(c.applyNoNsenter(ctx), opts...)
}

// VGCreate implements VolumeGroupClient.
func (c *noNsenterClient) VGCreate(ctx context.Context, opts ...VGCreateOption) error {
	return c.client.VGCreate(c.applyNoNsenter(ctx), opts...)
}

// VGRemove implements VolumeGroupClient.
func (c *noNsenterClient) VGRemove(ctx context.Context, opts ...VGRemoveOption) error {
	return c.client.VGRemove(c.applyNoNsenter(ctx), opts...)
}

// VGExtend implements VolumeGroupClient.
func (c *noNsenterClient) VGExtend(ctx context.Context, opts ...VGExtendOption) error {
	return c.client.VGExtend(c.applyNoNsenter(ctx), opts...)
}

// VGReduce implements VolumeGroupClient.
func (c *noNsenterClient) VGReduce(ctx context.Context, opts ...VGReduceOption) error {
	return c.client.VGReduce(c.applyNoNsenter(ctx), opts...)
}

// VGRename implements VolumeGroupClient.
func (c *noNsenterClient) VGRename(ctx context.Context, opts ...VGRenameOption) error {
	return c.client.VGRename(c.applyNoNsenter(ctx), opts...)
}

// VGChange implements VolumeGroupClient.
func (c *noNsenterClient) VGChange(ctx context.Context, opts ...VGChangeOption) error {
	return c.client.VGChange(c.applyNoNsenter(ctx), opts...)
}

// LV implements LogicalVolumeClient.
func (c *noNsenterClient) LV(ctx context.Context, opts ...LVsOption) (*LogicalVolume, error) {
	return c.client.LV(c.applyNoNsenter(ctx), opts...)
}

// LVs implements LogicalVolumeClient.
func (c *noNsenterClient) LVs(ctx context.Context, opts ...LVsOption) ([]*LogicalVolume, error) {
	return c.client.LVs(c.applyNoNsenter(ctx), opts...)
}

// LVCreate implements LogicalVolumeClient.
func (c *noNsenterClient) LVCreate(ctx context.Context, opts ...LVCreateOption) error {
	return c.client.LVCreate(c.applyNoNsenter(ctx), opts...)
}

// LVRemove implements LogicalVolumeClient.
func (c *noNsenterClient) LVRemove(ctx context.Context, opts ...LVRemoveOption) error {
	return c.client.LVRemove(c.applyNoNsenter(ctx), opts...)
}

// LVResize implements LogicalVolumeClient.
func (c *noNsenterClient) LVResize(ctx context.Context, opts ...LVResizeOption) error {
	return c.client.LVResize(c.applyNoNsenter(ctx), opts...)
}

// LVExtend implements LogicalVolumeClient.
func (c *noNsenterClient) LVExtend(ctx context.Context, opts ...LVExtendOption) error {
	return c.client.LVExtend(c.applyNoNsenter(ctx), opts...)
}

// LVReduce implements LogicalVolumeClient.
func (c *noNsenterClient) LVReduce(ctx context.Context, opts ...LVReduceOption) error {
	return c.client.LVReduce(c.applyNoNsenter(ctx), opts...)
}

// LVRename implements LogicalVolumeClient.
func (c *noNsenterClient) LVRename(ctx context.Context, opts ...LVRenameOption) error {
	return c.client.LVRename(c.applyNoNsenter(ctx), opts...)
}

// LVChange implements LogicalVolumeClient.
func (c *noNsenterClient) LVChange(ctx context.Context, opts ...LVChangeOption) error {
	return c.client.LVChange(c.applyNoNsenter(ctx), opts...)
}

// PVs implements PhysicalVolumeClient.
func (c *noNsenterClient) PVs(ctx context.Context, opts ...PVsOption) ([]*PhysicalVolume, error) {
	return c.client.PVs(c.applyNoNsenter(ctx), opts...)
}

// PVCreate implements PhysicalVolumeClient.
func (c *noNsenterClient) PVCreate(ctx context.Context, opts ...PVCreateOption) error {
	return c.client.PVCreate(c.applyNoNsenter(ctx), opts...)
}

// PVRemove implements PhysicalVolumeClient.
func (c *noNsenterClient) PVRemove(ctx context.Context, opts ...PVRemoveOption) error {
	return c.client.PVRemove(c.applyNoNsenter(ctx), opts...)
}

// PVResize implements PhysicalVolumeClient.
func (c *noNsenterClient) PVResize(ctx context.Context, opts ...PVResizeOption) error {
	return c.client.PVResize(c.applyNoNsenter(ctx), opts...)
}

// PVChange implements PhysicalVolumeClient.
func (c *noNsenterClient) PVChange(ctx context.Context, opts ...PVChangeOption) error {
	return c.client.PVChange(c.applyNoNsenter(ctx), opts...)
}

// PVMove implements PhysicalVolumeClient.
func (c *noNsenterClient) PVMove(ctx context.Context, opts ...PVMoveOption) error {
	return c.client.PVMove(c.applyNoNsenter(ctx), opts...)
}

// DevList implements DevicesClient.
func (c *noNsenterClient) DevList(ctx context.Context, opts ...DevListOption) ([]DeviceListEntry, error) {
	return c.client.DevList(c.applyNoNsenter(ctx), opts...)
}

// DevCheck implements DevicesClient.
func (c *noNsenterClient) DevCheck(ctx context.Context, opts ...DevCheckOption) error {
	return c.client.DevCheck(c.applyNoNsenter(ctx), opts...)
}

// DevUpdate implements DevicesClient.
func (c *noNsenterClient) DevUpdate(ctx context.Context, opts ...DevUpdateOption) error {
	return c.client.DevUpdate(c.applyNoNsenter(ctx), opts...)
}

// DevModify implements DevicesClient.
func (c *noNsenterClient) DevModify(ctx context.Context, opts ...DevModifyOption) error {
	return c.client.DevModify(c.applyNoNsenter(ctx), opts...)
}