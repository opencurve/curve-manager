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
	"time"

	"github.com/opencurve/curve-manager/internal/common"
)

type sessionItem struct {
	userName   string
	permission int
	timestamp  int64
}

func AddSession(userInfo *UserInfo) {
	now := time.Now().Unix()
	tokenStr := fmt.Sprintf("username=%s&password=%s&timestamp=%d", userInfo.UserName, userInfo.PassWord, now)
	token := common.GetMd5Sum32Little(tokenStr)
	gStorage.mutex.Lock()
	defer gStorage.mutex.Unlock()
	// check if have logined, if so will disconnect old one
	if oldToken, ok := gStorage.loginOnce[userInfo.UserName]; ok {
		delete(gStorage.session, oldToken)
	}
	gStorage.session[token] = sessionItem{
		userName:   userInfo.UserName,
		permission: userInfo.Permission,
		timestamp:  now,
	}
	gStorage.loginOnce[userInfo.UserName] = token
	if userInfo.Permission&WRITE_PERM == WRITE_PERM {
		gStorage.loginedWriteUser = userInfo.UserName
	}
	userInfo.Token = token
}

func CheckSession(s string, expireSec int) (bool, int) {
	now := time.Now().Unix()
	gStorage.mutex.Lock()
	defer gStorage.mutex.Unlock()
	if item, ok := gStorage.session[s]; ok {
		if item.timestamp+int64(expireSec) < now {
			delete(gStorage.session, s)
			delete(gStorage.loginOnce, item.userName)
			if item.permission&WRITE_PERM == WRITE_PERM {
				gStorage.loginedWriteUser = ""
			}
			return false, 0
		}
		item.timestamp = now
		gStorage.session[s] = item
		return true, item.permission
	}
	return false, 0
}

func GetLoginWriteUser() string {
	gStorage.mutex.Lock()
	defer gStorage.mutex.Unlock()
	return gStorage.loginedWriteUser
}

func Logout(name string) {
	gStorage.mutex.Lock()
	defer gStorage.mutex.Unlock()
	token, _ := gStorage.loginOnce[name]
	delete(gStorage.session, token)
	delete(gStorage.loginOnce, name)
	if name == gStorage.loginedWriteUser {
		gStorage.loginedWriteUser = ""
	}
}

func GetLoginUserByToken(token string) string {
	gStorage.mutex.Lock()
	defer gStorage.mutex.Unlock()
	if item, ok := gStorage.session[token]; ok {
		return item.userName
	}
	return ""
}
