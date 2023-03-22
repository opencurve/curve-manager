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
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/pigeon"
)

const (
	MODULE = "user"
)

var METHOD_REQUEST map[string]Request

type (
	HandlerFunc func(r *pigeon.Request, ctx *Context) bool

	Context struct {
		Data interface{}
	}

	Request struct {
		httpMethod string
		method     string
		vType      interface{}
		handler    HandlerFunc
	}
)

func init() {
	METHOD_REQUEST = map[string]Request{}
	for _, request := range requests {
		METHOD_REQUEST[request.method] = request
		core.BELONG[request.method] = MODULE
	}
}

type LoginRequest struct {
	UserName string `json:"userName" binding:"required"`
	PassWord string `json:"passWord" binding:"required"`
}

type LogoutRequest struct {
	UserName string `json:"userName" binding:"required"`
}

type CreateUserRequest struct {
	UserName   string `json:"userName" binding:"required"`
	PassWord   string `json:"passWord" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Permission int    `json:"permission" default:"4"`
}

type DeleteUserRequest struct {
	UserName string `json:"userName" binding:"required"`
}

type ChangePassWordRequest struct {
	UserName    string `json:"userName" binding:"required"`
	OldPassWord string `json:"oldPassword" binding:"required"`
	NewPassWord string `json:"newPassword" binding:"required"`
}

type ResetPassWordRequest struct {
	UserName string `json:"userName" binding:"required"`
}

type UpdateUserEmailRequest struct {
	UserName string `json:"userName" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UpdateUserPermissionRequest struct {
	UserName   string `json:"userName" binding:"required"`
	Permission int    `json:"permission" binding:"required"`
}

type ListUserRequest struct {
	Size     uint32 `json:"size" binding:"required"`
	Page     uint32 `json:"page" binding:"required"`
	UserName string `json:"userName"`
}

type GetUserRequest struct{}

var requests = []Request{
	{
		core.HTTP_POST,
		core.USER_LOGIN,
		LoginRequest{},
		Login,
	},
	{
		core.HTTP_POST,
		core.USER_LOGOUT,
		LogoutRequest{},
		Logout,
	},
	{
		core.HTTP_POST,
		core.USER_CREATE,
		CreateUserRequest{},
		CreateUser,
	},
	{
		core.HTTP_POST,
		core.USER_DELETE,
		DeleteUserRequest{},
		DeleteUser,
	},
	{
		core.HTTP_POST,
		core.USER_UPDATE_PASSWORD,
		ChangePassWordRequest{},
		ChangePassWord,
	},
	{
		core.HTTP_POST,
		core.USER_RESET_PASSWORD,
		ResetPassWordRequest{},
		ResetPassWord,
	},
	{
		core.HTTP_POST,
		core.USER_UPDATE_EMAIL,
		UpdateUserEmailRequest{},
		UpdateUserEmail,
	},
	{
		core.HTTP_POST,
		core.USER_UPDATE_PERMISSION,
		UpdateUserPermissionRequest{},
		UpdateUserPermission,
	},
	{
		core.HTTP_POST,
		core.USER_LIST,
		ListUserRequest{},
		ListUser,
	},
	{
		core.HTTP_GET,
		core.USER_GET,
		GetUserRequest{},
		GetUser,
	},
}
