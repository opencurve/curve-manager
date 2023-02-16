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

func GetEtcdStatus(r *pigeon.Request, ctx *Context) bool {
	status, err := agent.GetEtcdStatus(r)
	if err != errno.OK {
		return core.Exit(r, err)
	}
	return core.ExitSuccessWithData(r, status)
}

func GetMdsStatus(r *pigeon.Request, ctx *Context) bool {
	status, err := agent.GetMdsStatus(r)
	if err != errno.OK {
		return core.Exit(r, err)
	}
	return core.ExitSuccessWithData(r, status)
}

func GetSnapShotCloneServerStatus(r *pigeon.Request, ctx *Context) bool {
	status, err := agent.GetSnapShotCloneServerStatus(r)
	if err != errno.OK {
		return core.Exit(r, err)
	}
	return core.ExitSuccessWithData(r, status)
}

func GetChunkServerStatus(r *pigeon.Request, ctx *Context) bool {
	status, err := agent.GetChunkServerStatus(r)
	if err != errno.OK {
		return core.Exit(r, err)
	}
	return core.ExitSuccessWithData(r, status)
}

func GetClusterStatus(r *pigeon.Request, ctx *Context) bool {
	status := agent.GetClusterStatus(r)
	return core.ExitSuccessWithData(r, status)
}
