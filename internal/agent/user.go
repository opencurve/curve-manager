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
	"fmt"

	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/email"
	"github.com/opencurve/curve-manager/internal/storage"
)

func Login(name, passwd string) (interface{}, error) {
	return storage.Login(name, passwd)
}

func CreateUser(name, passwd, email string, permission int) error {
	return storage.CreateUser(name, passwd, email, permission)
}

func DeleteUser(name string) error {
	return storage.DeleteUser(name)
}

func ChangePassWord(name, oldPassword, newPassword string) error {
	passwd, err := storage.GetUserPassword(name)
	if err != nil {
		return err
	}
	if passwd != oldPassword {
		return fmt.Errorf("old passwd not matched, storedPasswd=%s, inPasswd=%s", passwd, oldPassword)
	}
	return storage.ChangePassWord(name, newPassword)
}

func ResetPassWord(name string) error {
	emailAddr, err := storage.GetUserEmail(name)
	if err != nil {
		return err
	}

	if emailAddr == "" {
		return fmt.Errorf("no email address")
	}
	passwd := storage.GetNewPassWord()
	err = storage.ChangePassWord(name, common.GetMd5Sum32Little(passwd))
	if err != nil {
		return err
	}

	err = email.SendNewPassWord(name, emailAddr, passwd)
	return err
}

func UpdateUserInfo(name, email string, permission int) error {
	return storage.UpdateUserInfo(name, email, permission)
}

func ListUser() (interface{}, error) {
	return storage.ListUser()
}
