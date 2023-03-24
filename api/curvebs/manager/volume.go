/*
*  Copyright (c) 2023 NetEase Inc.
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
 */

/*
* Project: Curve-Manager
* Created Date: 2023-02-11
* Author: wanghai (SeanHai)
 */

package manager

import (
	"github.com/opencurve/curve-manager/api/curvebs/agent"
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/pigeon"
)

func ListVolume(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*ListVolumeRequest)
	volumes, err := agent.ListVolume(r, data.Size, data.Page, data.Path, data.SortKey, data.SortDirection)
	if err != errno.OK {
		return core.Exit(r, err)
	}
	return core.ExitSuccessWithData(r, volumes)
}

func GetVolume(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*GetVolumeRequest)
	volume, err := agent.GetVolume(r, data.VolumeName)
	if err != errno.OK {
		return core.Exit(r, err)
	}
	return core.ExitSuccessWithData(r, volume)
}

func CleanRecycleBin(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*CleanRecycleBinRequest)
	err := agent.CleanRecycleBin(r, data.Expiration)
	return core.Exit(r, err)
}

func CreateNameSpace(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*CreateNameSpaceRequest)
	err := agent.CreateNameSpace(r, data.Name, data.User, data.PassWord)
	return core.Exit(r, err)
}

func CreateVolume(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*CreateVolumeRequest)
	err := agent.CreateVolume(r, data.VolumeName, data.User, data.PassWord, data.Length, data.StripUnit, data.StripCount)
	return core.Exit(r, err)
}

func ExtendVolume(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*ExtendVolumeRequest)
	err := agent.ExtendVolume(r, data.VolumeName, data.Length)
	return core.Exit(r, err)
}

func VolumeThrottle(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*VolumeThrottleRequest)
	err := agent.VolumeThrottle(r, data.VolumeName, data.ThrottleType, data.Limit, data.Burst, data.BurstLength)
	return core.Exit(r, err)
}

func DeleteVolume(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*DeleteVolumeRequest)
	err := agent.DeleteVolume(r, data.VolumeNames)
	return core.Exit(r, err)
}

func RecoverVolume(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*RecoverVolumeRequest)
	err := agent.RecoverVolume(r, data.VolumeIds)
	return core.Exit(r, err)
}

func CloneVolume(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*CloneVolumeRequest)
	err := agent.CloneVolume(r, data.Src, data.Dest, data.User, *data.Lazy)
	return core.Exit(r, err)
}

func Flatten(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*FlattenRequest)
	err := agent.Flatten(r, data.VolumeName, data.User)
	return core.Exit(r, err)
}
