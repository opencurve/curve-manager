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
	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/api/curvebs/agent"
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/pigeon"
)

func GetClusterSpace(r *pigeon.Request, ctx *Context) bool {
	space, err := agent.GetClusterSpace(r.Logger(), r.HeadersIn[comm.HEADER_REQUEST_ID])
	if err != errno.OK {
		return core.Exit(r, err)
	}
	return core.ExitSuccessWithData(r, space)
}

func GetClusterSpaceTrend(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*GetClusterSpaceTrendRequest)
	space, err := agent.GetClusterSpaceTrend(r, data.Start, data.End, data.Interval)
	if err != errno.OK {
		return core.Exit(r, err)
	}
	return core.ExitSuccessWithData(r, space)
}

func GetClusterPerformance(r *pigeon.Request, ctx *Context) bool {
	performance, err := agent.GetClusterPerformance(r)
	if err != errno.OK {
		return core.Exit(r, err)
	}
	return core.ExitSuccessWithData(r, performance)
}

func ListTopology(r *pigeon.Request, ctx *Context) bool {
	topo, err := agent.ListTopology(r)
	if err != errno.OK {
		return core.Exit(r, err)
	}
	return core.ExitSuccessWithData(r, topo)
}

func ListLogicalPool(r *pigeon.Request, ctx *Context) bool {
	pools, err := agent.ListLogicalPool(r)
	if err != errno.OK {
		return core.Exit(r, err)
	}
	return core.ExitSuccessWithData(r, pools)
}

func GetLogicalPool(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*GetLogicalPoolRequest)
	pools, err := agent.GetLogicalPool(r, data.Id)
	if err != errno.OK {
		return core.Exit(r, err)
	}
	return core.ExitSuccessWithData(r, pools)
}
