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

type CreateUserRequest struct {
	UserName   string `json:"userName" binding:"required"`
	PassWord   string `json:"passWord" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Permission int    `json:"permission" default:"1"`
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

type UpdateUserInfoRequest struct {
	UserName   string `json:"userName" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Permission int    `json:"permission" binding:"required"`
}

type ListUserRequest struct{}

var requests = []Request{
	{
		"POST",
		"user.login",
		LoginRequest{},
		Login,
	},
	{
		"POST",
		"user.create",
		CreateUserRequest{},
		CreateUser,
	},
	{
		"POST",
		"user.delete",
		DeleteUserRequest{},
		DeleteUser,
	},
	{
		"POST",
		"user.update.password",
		ChangePassWordRequest{},
		ChangePassWord,
	},
	{
		"POST",
		"user.reset.password",
		ResetPassWordRequest{},
		ResetPassWord,
	},
	{
		"POST",
		"user.update.info",
		UpdateUserInfoRequest{},
		UpdateUserInfo,
	},
	{
		"GET",
		"user.list",
		ListUserRequest{},
		ListUser,
	},
}
