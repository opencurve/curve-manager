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

package agent

import (
	"sort"

	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/email"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/curve-manager/internal/storage"
	"github.com/opencurve/pigeon"
)

type ListUserInfo struct {
	UserName   string `json:"userName" binding:"required"`
	Email      string `json:"email"`
	Permission int    `json:"permission" binding:"required"`
}

func Login(r *pigeon.Request, name, passwd string) (interface{}, errno.Errno) {
	userInfo, err := storage.GetUser(name)
	if err != nil {
		r.Logger().Error("Login failed",
			pigeon.Field("userName", name),
			pigeon.Field("passWord", passwd),
			pigeon.Field("error", err))
		return nil, errno.GET_USER_FAILED
	}
	if passwd != userInfo.PassWord {
		r.Logger().Error("Login failed",
			pigeon.Field("userName", name),
			pigeon.Field("inPassWord", passwd),
			pigeon.Field("storedPassword", userInfo.PassWord),
			pigeon.Field("error", err))
		return nil, errno.USER_PASSWORD_NOT_MATCH
	}
	storage.AddSession(&userInfo)
	return userInfo, errno.OK
}

func CreateUser(r *pigeon.Request, name, passwd, email string, permission int) errno.Errno {
	err := storage.SetUser(name, passwd, email, permission)
	if err != nil {
		r.Logger().Error("Create user failed",
			pigeon.Field("userName", name),
			pigeon.Field("passWord", passwd),
			pigeon.Field("email", email),
			pigeon.Field("permission", permission))
		return errno.CREATE_USER_FAILED
	}
	return errno.OK
}

func DeleteUser(r *pigeon.Request, name string) errno.Errno {
	err := storage.DeleteUser(name)
	if err != nil {
		r.Logger().Error("Delete user failed",
			pigeon.Field("userName", name),
			pigeon.Field("err", err))
		return errno.DELETE_USER_FAILED
	}
	return errno.OK
}

func ChangePassWord(r *pigeon.Request, name, oldPassword, newPassword string) errno.Errno {
	passwd, err := storage.GetUserPassword(name)
	if err != nil {
		r.Logger().Error("GetUserPassword failed",
			pigeon.Field("userName", name),
			pigeon.Field("error", err))
		return errno.GET_USER_PASSWORD_FAILED
	}
	if passwd != oldPassword {
		r.Logger().Error("ChangePassWord failed, old password not match",
			pigeon.Field("userName", name),
			pigeon.Field("inPassword", oldPassword),
			pigeon.Field("stored password", passwd))
		return errno.USER_PASSWORD_NOT_MATCH
	}
	err = storage.UpdateUserPassWord(name, newPassword)
	if err != nil {
		r.Logger().Error("UpdateUserPassWord failed",
			pigeon.Field("userName", name),
			pigeon.Field("error", err))
		return errno.UPDATE_USER_PASSWORD_FAILED
	}
	return errno.OK
}

func ResetPassWord(r *pigeon.Request, name string) errno.Errno {
	emailAddr, err := storage.GetUserEmail(name)
	if err != nil {
		r.Logger().Error("GetUserEmail failed",
			pigeon.Field("userName", name),
			pigeon.Field("error", err))
		return errno.GET_USER_EMAIL_FAILED
	}

	if emailAddr == "" {
		r.Logger().Error("ResetPassWord failed, email is empty",
			pigeon.Field("userName", name))
		return errno.USER_EMAIL_EMPTY
	}
	passwd := storage.GetNewPassWord()
	err = storage.UpdateUserPassWord(name, common.GetMd5Sum32Little(passwd))
	if err != nil {
		r.Logger().Error("UpdateUserPassWord failed",
			pigeon.Field("userName", name),
			pigeon.Field("passWord", common.GetMd5Sum32Little(passwd)),
			pigeon.Field("error", err))
		return errno.UPDATE_USER_PASSWORD_FAILED
	}

	err = email.SendNewPassWord(name, emailAddr, passwd)
	if err != nil {
		r.Logger().Error("Email sendNewPassWord failed",
			pigeon.Field("userName", name),
			pigeon.Field("emailAddr", emailAddr),
			pigeon.Field("newPassword", passwd),
			pigeon.Field("error", err))
		return errno.SEND_USER_PASSWORD_FAILED
	}
	return errno.OK
}

func UpdateUserInfo(r *pigeon.Request, name, email string, permission int) errno.Errno {
	err := storage.UpdateUserInfo(name, email, permission)
	if err != nil {
		r.Logger().Error("UpdateUserInfo failed",
			pigeon.Field("userName", name),
			pigeon.Field("error", err))
		return errno.UPDATE_USER_INFO_FAILED
	}
	return errno.OK
}

func sortUser(users []storage.UserInfo) {
	sort.Slice(users, func(i, j int) bool {
		return users[i].UserName < users[j].UserName
	})
}

func ListUser(r *pigeon.Request, size, page uint32, userName string) (interface{}, errno.Errno) {
	users, err := storage.ListUser(userName)
	if err != nil {
		r.Logger().Error("ListUser failed",
			pigeon.Field("error", err))
		return nil, errno.LIST_USER_FAILED
	}
	sortUser(*users)
	length := uint32(len(*users))
	start := (page - 1) * size
	end := common.MinUint32(page*size, length)

	infos := []ListUserInfo{}
	for _, user := range (*users)[start:end] {
		var info ListUserInfo
		info.UserName = user.UserName
		info.Email = user.Email
		info.Permission = user.Permission
		infos = append(infos, info)
	}
	return infos, errno.OK
}
