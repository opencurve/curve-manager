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

	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/email"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/curve-manager/internal/storage"
	"github.com/opencurve/pigeon"
)

var (
	READ_PERM    = 4
	WRITE_PERM   = 2
	MANAGER_PERM = 1
)

type UserInfo struct {
	UserName   string `json:"userName" binding:"required"`
	Email      string `json:"email"`
	Permission int    `json:"permission" binding:"required"`
}

type ListUserInfo struct {
	Total int        `json:"total" binding:"required"`
	Info  []UserInfo `json:"info" binding:"required"`
}

const (
	ENABLE_MULTIPLE_WRITER_USER_LOGIN = "enable.multiple.write.user.login"
)

func Login(r *pigeon.Request, name, passwd string) (interface{}, errno.Errno) {
	userInfo, err := storage.GetUser(name)
	if err != nil {
		r.Logger().Error("Login failed",
			pigeon.Field("userName", name),
			pigeon.Field("passWord", passwd),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_USER_FAILED
	}
	if passwd != userInfo.PassWord {
		r.Logger().Error("Login failed",
			pigeon.Field("userName", name),
			pigeon.Field("inPassWord", passwd),
			pigeon.Field("storedPassword", userInfo.PassWord),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.USER_PASSWORD_NOT_MATCH
	}
	// check multiple write user login
	if userInfo.Permission&WRITE_PERM == WRITE_PERM &&
		!r.GetConfig().GetBool(ENABLE_MULTIPLE_WRITER_USER_LOGIN) {
		user := storage.GetLoginWriteUser()
		if user != "" && userInfo.UserName != user {
			r.Logger().Error("Login failed, there is already a write user logined",
				pigeon.Field("userName", name),
				pigeon.Field("permission", userInfo.Permission),
				pigeon.Field("login write user", user))
			return nil, errno.WRITE_USER_LOGIN_FAILED
		}
	}
	storage.AddSession(&userInfo)
	return userInfo, errno.OK
}

func Logout(r *pigeon.Request) errno.Errno {
	token := r.HeadersIn[comm.HEADER_AUTH_TOKEN]
	user := storage.GetLoginUserByToken(token)
	storage.Logout(user)
	return errno.OK
}

func CreateUser(r *pigeon.Request, name, passwd, email string, permission int) errno.Errno {
	err := storage.CreateUser(name, passwd, email, permission)
	if err != nil {
		r.Logger().Error("Create user failed",
			pigeon.Field("userName", name),
			pigeon.Field("passWord", passwd),
			pigeon.Field("email", email),
			pigeon.Field("permission", permission),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return errno.CREATE_USER_FAILED
	}
	return errno.OK
}

func DeleteUser(r *pigeon.Request, name string) errno.Errno {
	err := storage.DeleteUser(name)
	if err != nil {
		r.Logger().Error("Delete user failed",
			pigeon.Field("userName", name),
			pigeon.Field("err", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return errno.DELETE_USER_FAILED
	}
	return errno.OK
}

func ChangePassWord(r *pigeon.Request, name, oldPassword, newPassword string) errno.Errno {
	passwd, err := storage.GetUserPassword(name)
	if err != nil {
		r.Logger().Error("GetUserPassword failed",
			pigeon.Field("userName", name),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return errno.GET_USER_PASSWORD_FAILED
	}
	if passwd != oldPassword {
		r.Logger().Error("ChangePassWord failed, old password not match",
			pigeon.Field("userName", name),
			pigeon.Field("inPassword", oldPassword),
			pigeon.Field("stored password", passwd),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return errno.USER_PASSWORD_NOT_MATCH
	}
	err = storage.UpdateUserPassWord(name, newPassword)
	if err != nil {
		r.Logger().Error("UpdateUserPassWord failed",
			pigeon.Field("userName", name),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return errno.UPDATE_USER_PASSWORD_FAILED
	}
	return errno.OK
}

func ResetPassWord(r *pigeon.Request, name string) errno.Errno {
	emailAddr, err := storage.GetUserEmail(name)
	if err != nil {
		r.Logger().Error("GetUserEmail failed",
			pigeon.Field("userName", name),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return errno.GET_USER_EMAIL_FAILED
	}

	if emailAddr == "" {
		r.Logger().Error("ResetPassWord failed, email is empty",
			pigeon.Field("userName", name),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return errno.USER_EMAIL_EMPTY
	}
	passwd := storage.GetNewPassWord()
	err = storage.UpdateUserPassWord(name, common.GetMd5Sum32Little(passwd))
	if err != nil {
		r.Logger().Error("UpdateUserPassWord failed",
			pigeon.Field("userName", name),
			pigeon.Field("passWord", common.GetMd5Sum32Little(passwd)),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return errno.UPDATE_USER_PASSWORD_FAILED
	}

	err = email.SendNewPassWord(name, emailAddr, passwd)
	if err != nil {
		r.Logger().Error("Email sendNewPassWord failed",
			pigeon.Field("userName", name),
			pigeon.Field("emailAddr", emailAddr),
			pigeon.Field("newPassword", passwd),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return errno.SEND_USER_PASSWORD_FAILED
	}
	return errno.OK
}

func UpdateUserEmail(r *pigeon.Request, name, email string) errno.Errno {
	err := storage.UpdateUserEmail(name, email)
	if err != nil {
		r.Logger().Error("UpdateUserEmail failed",
			pigeon.Field("userName", name),
			pigeon.Field("email", email),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return errno.UPDATE_USER_EMAIL_FAILED
	}
	return errno.OK
}

func UpdateUserPermission(r *pigeon.Request, name string, permission int) errno.Errno {
	err := storage.UpdateUserPermission(name, permission)
	if err != nil {
		r.Logger().Error("UpdateUserPermission failed",
			pigeon.Field("userName", name),
			pigeon.Field("permission", permission),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return errno.UPDATE_USER_PERMISSION_FAILED
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
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.LIST_USER_FAILED
	}
	sortUser(*users)
	length := uint32(len(*users))
	start := (page - 1) * size
	end := common.MinUint32(page*size, length)

	info := ListUserInfo{
		Info: []UserInfo{},
	}
	info.Total = len(*users)
	if start >= length {
		return info, errno.OK
	}
	for _, user := range (*users)[start:end] {
		var item UserInfo
		item.UserName = user.UserName
		item.Email = user.Email
		item.Permission = user.Permission
		info.Info = append(info.Info, item)
	}
	return info, errno.OK
}

func GetUser(r *pigeon.Request) (interface{}, errno.Errno) {
	token := r.HeadersIn[comm.HEADER_AUTH_TOKEN]
	userName := storage.GetLoginUserByToken(token)
	info, err := storage.GetUser(userName)
	if err != nil {
		r.Logger().Error("GetUser failed",
			pigeon.Field("token", token),
			pigeon.Field("userName", userName),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_USER_FAILED
	}
	return info, errno.OK
}
