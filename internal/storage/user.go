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
* Created Date: 2023-02-20
* Author: wanghai (SeanHai)
 */

package storage

import (
	"fmt"

	"github.com/opencurve/curve-manager/internal/common"
)

const (
	USER_ADMIN_NAME     = "admin"
	USER_ADMIN_PASSWORD = "curve"

	ADMIN_PERM = 7
	WRITE_PERM = 2

	NEW_PASSWORD_LENGTH = 8
)

type UserInfo struct {
	UserName   string `json:"userName" binding:"required"`
	PassWord   string `json:"-"`
	Email      string `json:"email"`
	Permission int    `json:"permission" binding:"required"`
	Token      string `json:"token,omitempty" binding:"required"`
}

func createAdminUser() error {
	passwd := common.GetMd5Sum32Little(USER_ADMIN_PASSWORD)
	return gStorage.execSQL(CREATE_ADMIN, USER_ADMIN_NAME, passwd, "", ADMIN_PERM)
}

func GetUser(name string) (UserInfo, error) {
	var user UserInfo
	rows, err := gStorage.querySQL(GET_USER, name)
	if err != nil {
		return user, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&user.UserName, &user.PassWord, &user.Email, &user.Permission)
		if err != nil {
			return user, err
		}
	} else {
		return user, fmt.Errorf("user not exist")
	}
	return user, nil
}

func CreateUser(name, passwd, email string, permission int) error {
	return gStorage.execSQL(CREATE_USER, name, passwd, email, permission)
}

func DeleteUser(name string) error {
	return gStorage.execSQL(DELETE_USER, name)
}

func UpdateUserPassWord(name, passwd string) error {
	return gStorage.execSQL(UPDATE_USER_PASSWORD, passwd, name)
}

func UpdateUserEmail(name, email string) error {
	return gStorage.execSQL(UPDATE_USER_EMAIL, email, name)
}

func UpdateUserPermission(name string, perm int) error {
	return gStorage.execSQL(UPDATE_USER_PERMISSION, perm, name)
}

func ListUser(userName string) (*[]UserInfo, error) {
	var users []UserInfo
	var sql string
	var params []interface{}
	if userName == "" {
		sql = LIST_USER
		params = append(params, USER_ADMIN_NAME)
	} else if userName == USER_ADMIN_NAME {
		return &users, nil
	} else {
		sql = GET_USER
		params = append(params, userName)
	}
	rows, err := gStorage.querySQL(sql, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user UserInfo
		err = rows.Scan(&user.UserName, &user.PassWord, &user.Email, &user.Permission)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return &users, nil
}

func ListUserWithEmail() ([]string, error) {
	emial := ""
	users := []string{}
	rows, err := gStorage.querySQL(LIST_USER_WITH_EMAIL, emial)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user string
		err = rows.Scan(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetUserEmail(name string) (string, error) {
	rows, err := gStorage.querySQL(GET_USER_EMAIL, name)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var email string
	if rows.Next() {
		err = rows.Scan(&email)
		if err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf("user not exist")
	}
	return email, nil
}

func GetUserPassword(name string) (string, error) {
	rows, err := gStorage.querySQL(GET_USER_PASSWORD, name)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var passwd string
	if rows.Next() {
		err = rows.Scan(&passwd)
		if err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf("user not exist")
	}
	return passwd, nil
}

func GetNewPassWord() string {
	return common.GetRandString(NEW_PASSWORD_LENGTH)
}
