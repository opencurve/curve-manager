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

package user

import (
	"github.com/mcuadros/go-defaults"
	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/curve-manager/internal/agent"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/pigeon"
)

func Login(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*LoginRequest)
	user, err := agent.Login(data.UserName, data.PassWord)
	if err != nil {
		r.Logger().Error("user login failed",
			pigeon.Field("userName", data.UserName),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.USER_LOGIN_FAILED)
	}
	return core.ExitSuccessWithData(r, user)
}

func CreateUser(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*CreateUserRequest)
	defaults.SetDefaults(data)
	err := agent.CreateUser(data.UserName, data.PassWord, data.Email, data.Permission)
	if err != nil {
		r.Logger().Error("create user failed",
			pigeon.Field("userName", data.UserName),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.CREATE_USER_FAILED)
	}
	return core.Exit(r, errno.OK)
}

func DeleteUser(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*DeleteUserRequest)
	err := agent.DeleteUser(data.UserName)
	if err != nil {
		r.Logger().Error("delete user failed",
			pigeon.Field("userName", data.UserName),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.DELETE_USER_FAILED)
	}
	return core.Exit(r, errno.OK)
}

func ChangePassWord(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*ChangePassWordRequest)
	err := agent.ChangePassWord(data.UserName, data.OldPassWord, data.NewPassWord)
	if err != nil {
		r.Logger().Error("change password failed",
			pigeon.Field("userName", data.UserName),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.CHANGE_PASSWORD_FAILED)
	}
	return core.Exit(r, errno.OK)
}

func ResetPassWord(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*ResetPassWordRequest)
	err := agent.ResetPassWord(data.UserName)
	if err != nil {
		r.Logger().Error("reset password failed",
			pigeon.Field("userName", data.UserName),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.RESET_PASSWORD_FAILED)
	}
	return core.Exit(r, errno.OK)
}

func UpdateUserInfo(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*UpdateUserInfoRequest)
	err := agent.UpdateUserInfo(data.UserName, data.Email, data.Permission)
	if err != nil {
		r.Logger().Error("update user info failed",
			pigeon.Field("userName", data.UserName),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.UPDATE_USER_INFO_FAILED)
	}
	return core.Exit(r, errno.OK)
}

func ListUser(r *pigeon.Request, ctx *Context) bool {
	users, err := agent.ListUser()
	if err != nil {
		r.Logger().Error("list user failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.LIST_USER_FAILED)
	}
	return core.ExitSuccessWithData(r, users)
}
