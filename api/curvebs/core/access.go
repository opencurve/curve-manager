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

	READ_PERM    = 4
	WRITE_PERM   = 2
	MANAGER_PERM = 1
)

var (
	// a switch whether to enable interface verification, default false
	enableCheck bool
	// api expire time, default 60s
	apiExpireSeconds int
	// login expire time, default 1800s
	loginExpireSeconds int

	// whether allow multiple user with write permission logined at same time
	enableMultipleWriteUserLogin bool

	// method to permission
	method2permission = map[string]int{
		USER_CREATE:                 READ_PERM + MANAGER_PERM,
		USER_DELETE:                 READ_PERM + MANAGER_PERM,
		USER_LIST:                   READ_PERM + MANAGER_PERM,
		USER_UPDATE_PERMISSION:      READ_PERM + MANAGER_PERM,
		USER_LOGIN:                  READ_PERM,
		USER_LOGOUT:                 READ_PERM,
		USER_GET:                    READ_PERM,
		USER_UPDATE_PASSWORD:        READ_PERM,
		USER_RESET_PASSWORD:         READ_PERM,
		USER_UPDATE_EMAIL:           READ_PERM,
		STATUS_ETCD:                 READ_PERM,
		STATUS_MDS:                  READ_PERM,
		STATUS_SNAPSHOTCLONESERVER:  READ_PERM,
		STATUS_CHUNKSERVER:          READ_PERM,
		STATUS_CLUSTER:              READ_PERM,
		SPACE_CLUSTER:               READ_PERM,
		SPACE_TREND_CLUSTER:         READ_PERM,
		PERFORMANCE_CLUSTER:         READ_PERM,
		TOPO_LIST:                   READ_PERM,
		TOPO_POOL_LIST:              READ_PERM,
		TOPO_POOL_GET:               READ_PERM,
		VOLUME_LIST:                 READ_PERM,
		VOLUME_GET:                  READ_PERM,
		SNAPSHOT_LIST:               READ_PERM,
		HOST_LIST:                   READ_PERM,
		HOST_GET:                    READ_PERM,
		DISK_LIST:                   READ_PERM,
		CLEAN_RECYCLEBIN:            READ_PERM + WRITE_PERM,
		CREATE_NAMESPACE:            READ_PERM + WRITE_PERM,
		CREATE_VOLUME:               READ_PERM + WRITE_PERM,
		EXTEND_VOLUME:               READ_PERM + WRITE_PERM,
		VOLUME_THROTTLE:             READ_PERM + WRITE_PERM,
		DELETE_VOLUME:               READ_PERM + WRITE_PERM,
		RECOVER_VOLUME:              READ_PERM + WRITE_PERM,
		CLONE_VOLUME:                READ_PERM + WRITE_PERM,
		CREATE_SNAPSHOT:             READ_PERM + WRITE_PERM,
		CANCEL_SNAPSHOT:             READ_PERM + WRITE_PERM,
		FLATTEN:                     READ_PERM + WRITE_PERM,
		DELETE_SNAPSHOT:             READ_PERM + WRITE_PERM,
		GET_SYSTEM_LOG:              READ_PERM,
		GET_SYSTEM_ALERT:            READ_PERM,
		GET_UNREAD_SYSTEM_ALERT_NUM: READ_PERM,
		UPDATE_READ_SYSTEM_ALERT_ID: READ_PERM,
		GET_ALERT_CONF:              READ_PERM,
		UPDATE_ALERT_CONF:           READ_PERM + MANAGER_PERM,
		GET_ALERT_CANDIDATE:         READ_PERM + MANAGER_PERM,
		UPDATE_ALERT_USER:           READ_PERM + MANAGER_PERM,
		DEPLOY_HOST_LIST:            READ_PERM + MANAGER_PERM,
		DEPLOY_HOST_COMMIT:          READ_PERM + MANAGER_PERM,
		DEPLOY_DISK_LIST:            READ_PERM + MANAGER_PERM,
		DEPLOY_DISK_COMMIT:          READ_PERM + MANAGER_PERM,
		DEPLOY_DISK_FORMAT_STATUS:   READ_PERM + MANAGER_PERM,
		DEPLOY_DISK_FORMAT:          READ_PERM + MANAGER_PERM,
		DEPLOY_CONFIG_SHOW:          READ_PERM + MANAGER_PERM,
		DEPLOY_CONFIG_COMMIT:        READ_PERM + MANAGER_PERM,
		DEPLOY_CLUSTER_ADD:          READ_PERM + MANAGER_PERM,
		DEPLOY_CLUSTER_DEPLOY:       READ_PERM + MANAGER_PERM,
	}
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

func IsLoginRequest(r *pigeon.Request) bool {
	return r.Args[METHOD] == USER_LOGIN
}

func IsResetPasswordRequest(r *pigeon.Request) bool {
	return r.Args[METHOD] == USER_RESET_PASSWORD
}

func NeedRecordLog(r *pigeon.Request) bool {
	method := r.Args[METHOD]
	return method == USER_LOGIN || method == USER_LOGOUT || method == USER_RESET_PASSWORD ||
		(method2permission[method]&(WRITE_PERM+MANAGER_PERM) > 0 && method != GET_SYSTEM_LOG)
}

func checkPermission(r *pigeon.Request, perm int) bool {
	method := r.Args[METHOD]
	if p, ok := method2permission[method]; ok {
		return p&perm == p
	}
	return false
}

func checkTimeOut(r *pigeon.Request) bool {
	argTime := r.HeadersIn[comm.HEADER_AUTH_TIMESTAMP]
	inTime, err := strconv.ParseInt(argTime, 10, 64)
	if err != nil {
		r.Logger().Error("checkTimeOut failed, invalid time argument",
			pigeon.Field("time", argTime),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return false
	}
	nowSec := time.Now().Unix()
	if inTime+int64(apiExpireSeconds) < nowSec {
		r.Logger().Error("checkTimeOut failed, time expired",
			pigeon.Field("inTime", inTime),
			pigeon.Field("ttl", apiExpireSeconds),
			pigeon.Field("now", nowSec),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return false
	}
	return true
}

/*
* algorithm：
* 1. String-Items: HTTP-Method; URI; Args; Timestamp; Token
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
	// v := reflect.ValueOf(data)
	// for i := 0; i < v.Elem().NumField(); i++ {
	// 	value := fmt.Sprintf("%+v", v.Elem().Field(i))
	// 	if value != "" {
	// 		stringItems = append(stringItems, value)
	// 	}
	// }
	sort.Strings(stringItems)
	signStr := strings.Join(stringItems, ":")
	sign := common.GetMd5Sum32Little(signStr)
	if inSign != sign {
		r.Logger().Error("checkSignature failed",
			pigeon.Field("signStr", signStr),
			pigeon.Field("sign", sign),
			pigeon.Field("inSign", inSign),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return false
	}
	return true
}

func checkToken(r *pigeon.Request) (bool, int) {
	token := r.HeadersIn[comm.HEADER_AUTH_TOKEN]
	return storage.CheckSession(token, loginExpireSeconds)
}

func AccessAllowed(r *pigeon.Request, data interface{}) errno.Errno {
	if !IsLoginRequest(r) && !IsResetPasswordRequest(r) {
		// check user token valied
		ok, perm := checkToken(r)
		if !ok {
			r.Logger().Error("checkToken failed",
				pigeon.Field("token", r.HeadersIn[comm.HEADER_AUTH_TOKEN]),
				pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
			return errno.USER_IS_UNAUTHORIZED
		}
		// check request method is allowed to this user
		if !checkPermission(r, perm) {
			r.Logger().Error("checkToken checkPermission failed",
				pigeon.Field("method", r.Args[METHOD]),
				pigeon.Field("user perm", perm),
				pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
			return errno.OPERATION_IS_NOT_PERMIT
		}
		if enableCheck {
			// check request ttl and sigenatrue
			if !checkTimeOut(r) || !checkSignature(r, data) {
				return errno.REQUEST_IS_DENIED_FOR_SIGNATURE
			}
		}
	}
	return errno.OK
}
