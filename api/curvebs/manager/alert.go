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
* Created Date: 2023-03-22
* Author: wanghai (SeanHai)
 */

package manager

import (
	"github.com/opencurve/curve-manager/api/curvebs/agent"
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/pigeon"
)

func GetSysAlert(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*GetSysAlertRequest)
	logs, err := agent.GetSysAlert(r, data.Start, data.End, data.Page, data.Size, data.Name, data.Level, data.Content)
	if err != errno.OK {
		return core.Exit(r, err)
	}
	return core.ExitSuccessWithData(r, logs)
}

func GetUnreadSysAlertNum(r *pigeon.Request, ctx *Context) bool {
	number, err := agent.GetUnreadSysAlertNum(r)
	if err != errno.OK {
		return core.Exit(r, err)
	}
	return core.ExitSuccessWithData(r, number)
}

func UpdateReadSysAlertId(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*UpdateReadSysAlertIdRequest)
	err := agent.UpdateReadSysAlertId(r, data.Id)
	return core.Exit(r, err)
}

func GetAlertConf(r *pigeon.Request, ctx *Context) bool {
	confs, err := agent.GetAlertConf(r)
	if err != errno.OK {
		return core.Exit(r, err)
	}
	return core.ExitSuccessWithData(r, confs)
}

func UpdateAlertConf(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*UpdateAlertConfRequest)
	err := agent.UpdateAlertConf(r, *data.Enable, data.Interval, data.Times,
		data.Rule, data.Name, data.AlertUsers)
	return core.Exit(r, err)
}

func GetAlertCandidate(r *pigeon.Request, ctx *Context) bool {
	users, err := agent.GetAlertCandidate(r)
	if err != errno.OK {
		return core.Exit(r, err)
	}
	return core.ExitSuccessWithData(r, users)
}

func UpdateAlertUser(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*UpdateAlertUserRequest)
	err := agent.UpdateAlertUser(r, data.Alert, data.User)
	return core.Exit(r, err)
}
