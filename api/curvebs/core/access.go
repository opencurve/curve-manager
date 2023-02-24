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

package core

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/curve-manager/internal/storage"

	"github.com/opencurve/pigeon"
)

const (
	ACCESS_API_ENABLE_CHECK     = "access.api.enable_check"
	ACCESS_API_EXPIRE_SECONDS   = "access.api.expire_seconds"
	ACCESS_LOGIN_EXPIRE_SECONDS = "access.login.expire_seconds"

	// method
	METHOD             = "method"
	METHOD_USER_LOGIN  = "user.login"
	METHOD_USER_CREATE = "user.create"
	METHOD_USER_DELETE = "user.delete"
	METHOD_USER_LIST   = "user.list"
)

var (
	// a switch whether to enable interface verification, default false
	enableCheck bool
	// api expire time, default 60s
	apiExpireSeconds int
	// login expire time, default 1800s
	loginExpireSeconds int
)

func InitAccess(cfg *pigeon.Configure) {
	enableCheck = cfg.GetConfig().GetBool(ACCESS_API_ENABLE_CHECK)
	apiExpireSeconds = cfg.GetConfig().GetInt(ACCESS_API_EXPIRE_SECONDS)
	if apiExpireSeconds <= 0 {
		apiExpireSeconds = 60
	}
	loginExpireSeconds = cfg.GetConfig().GetInt(ACCESS_LOGIN_EXPIRE_SECONDS)
	if loginExpireSeconds <= 0 {
		loginExpireSeconds = 1800
	}
}

func isLoginRequest(r *pigeon.Request) bool {
	return r.Args[METHOD] == METHOD_USER_LOGIN
}

func isAdminRequest(r *pigeon.Request) bool {
	method := r.Args[METHOD]
	return method == METHOD_USER_CREATE || method == METHOD_USER_DELETE || method == METHOD_USER_LIST
}

func checkTimeOut(r *pigeon.Request) bool {
	argTime := r.HeadersIn[comm.HEADER_AUTH_TIMESTAMP]
	inTime, err := strconv.ParseInt(argTime, 10, 64)
	if err != nil {
		r.Logger().Error("checkTimeOut failed, invalid time argument",
			pigeon.Field("time", argTime))
		return false
	}
	nowSec := time.Now().Unix()
	if inTime+int64(apiExpireSeconds) < nowSec {
		r.Logger().Error("checkTimeOut failed, time expired",
			pigeon.Field("inTime", inTime),
			pigeon.Field("ttl", apiExpireSeconds),
			pigeon.Field("now", nowSec))
		return false
	}
	return true
}

/*
* algorithm：
* 1. String-Items: HTTP-Method; URI; Args; QueryValue1; QueryValue2; ... QueryValuen; Timestamp; Token
* 2. Sorted-Items: sort String-Items based alphabetically
* 3. Sign-String: join Sorted-Items with ":"
* 4. Sign: MD532Little(Sign-String)
 */
func checkSignature(r *pigeon.Request, data interface{}) bool {
	token := r.HeadersIn[comm.HEADER_AUTH_TOKEN]
	inSign := r.HeadersIn[comm.HEADER_AUTH_SIGN]
	timeStamp := r.HeadersIn[comm.HEADER_AUTH_TIMESTAMP]
	stringItems := []string{r.Method, r.Uri, timeStamp, token}
	for _, v := range r.Args {
		stringItems = append(stringItems, v)
	}
	v := reflect.ValueOf(data)
	for i := 0; i < v.Elem().NumField(); i++ {
		stringItems = append(stringItems, fmt.Sprintf("%+v", v.Elem().Field(i)))
	}
	sort.Strings(stringItems)
	signStr := strings.Join(stringItems, ":")
	sign := common.GetMd5Sum32Little(signStr)
	r.Logger().Error("checkSignature",
		pigeon.Field("signStr", signStr),
		pigeon.Field("sign", sign))
	return inSign == sign
}

func checkToken(r *pigeon.Request) bool {
	token := r.HeadersIn[comm.HEADER_AUTH_TOKEN]
	if !storage.CheckSession(token, apiExpireSeconds) {
		return false
	}
	if isAdminRequest(r) {
		return storage.IsAdmin(token)
	}
	return true
}

func AccessAllowed(r *pigeon.Request, data interface{}) errno.Errno {
	if !isLoginRequest(r) {
		if !checkToken(r) {
			return errno.USER_IS_UNAUTHORIZED
		}
		if !checkTimeOut(r) || !checkSignature(r, data) {
			return errno.REQUEST_IS_DENIED_FOR_SIGNATURE
		}
	}
	return errno.OK
}
