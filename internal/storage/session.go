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
	userName  string
	timestamp int64
}

func AddSession(userInfo *UserInfo) {
	now := time.Now().Unix()
	sigStr := fmt.Sprintf("username=%s&password=%s&timestamp=%d", userInfo.UserName, userInfo.PassWord, now)
	sig := common.GetMd5Sum32Little(sigStr)
	gStorage.sessionMutex.Lock()
	defer gStorage.sessionMutex.Unlock()
	gStorage.session[sig] = sessionItem{
		userName:  userInfo.UserName,
		timestamp: now,
	}
	userInfo.Token = sig
}

func CheckSession(s string, expireSec int) bool {
	now := time.Now().Unix()
	gStorage.sessionMutex.Lock()
	defer gStorage.sessionMutex.Unlock()
	if item, ok := gStorage.session[s]; ok {
		if item.timestamp+int64(expireSec) < now {
			delete(gStorage.session, s)
			return false
		}
		item.timestamp = now
		gStorage.session[s] = item
		return true
	}
	return false
}

func IsAdmin(s string) bool {
	gStorage.sessionMutex.Lock()
	defer gStorage.sessionMutex.Unlock()
	if item, ok := gStorage.session[s]; ok {
		if item.userName == USER_ADMIN_NAME {
			return true
		}
		return false
	}
	return false
}
