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
	sigStr := fmt.Sprintf("username=%s&password=%s&timestamp=%d", userInfo.UserName, userInfo.PassWord, now)
	sig := common.GetMd5Sum32Little(sigStr)
	gStorage.sessionMutex.Lock()
	defer gStorage.sessionMutex.Unlock()
	// check if have logined, if so will disconnect old one
	if oldSig, ok := gStorage.loginOnce[userInfo.UserName]; ok {
		delete(gStorage.session, oldSig)
	}
	gStorage.session[sig] = sessionItem{
		userName:   userInfo.UserName,
		permission: userInfo.Permission,
		timestamp:  now,
	}
	gStorage.loginOnce[userInfo.UserName] = sig
	if userInfo.Permission&WRITE_PERM == WRITE_PERM {
		gStorage.loginedWriteUser = userInfo.UserName
	}
	userInfo.Token = sig
}

func CheckSession(s string, expireSec int) (bool, int) {
	now := time.Now().Unix()
	gStorage.sessionMutex.Lock()
	defer gStorage.sessionMutex.Unlock()
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
	gStorage.sessionMutex.Lock()
	defer gStorage.sessionMutex.Unlock()
	return gStorage.loginedWriteUser
}

func Logout(name string) {
	gStorage.sessionMutex.Lock()
	defer gStorage.sessionMutex.Unlock()
	sig, _ := gStorage.loginOnce[name]
	delete(gStorage.session, sig)
	delete(gStorage.loginOnce, name)
	if name == gStorage.loginedWriteUser {
		gStorage.loginedWriteUser = ""
	}
}
