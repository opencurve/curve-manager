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
	"strconv"

	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/api/curvebs/agent"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/pigeon"
)

const (
	// http method
	HTTP_POST = "POST"
	HTTP_GET  = "GET"

	// method
	METHOD = "method"

	// user
	USER_LOGIN             = "user.login"
	USER_LOGOUT            = "user.logout"
	USER_CREATE            = "user.create"
	USER_DELETE            = "user.delete"
	USER_UPDATE_PASSWORD   = "user.update.password"
	USER_RESET_PASSWORD    = "user.reset.password"
	USER_UPDATE_EMAIL      = "user.update.email"
	USER_UPDATE_PERMISSION = "user.update.permission"
	USER_LIST              = "user.list"
	USER_GET               = "user.get"

	// manager
	STATUS_ETCD                 = "status.etcd"
	STATUS_MDS                  = "status.mds"
	STATUS_SNAPSHOTCLONESERVER  = "status.snapshotcloneserver"
	STATUS_CHUNKSERVER          = "status.chunkserver"
	STATUS_CLUSTER              = "status.cluster"
	SPACE_CLUSTER               = "space.cluster"
	SPACE_TREND_CLUSTER         = "space.trend.cluster"
	PERFORMANCE_CLUSTER         = "performance.cluster"
	TOPO_LIST                   = "topo.list"
	TOPO_POOL_LIST              = "topo.pool.list"
	TOPO_POOL_GET               = "topo.pool.get"
	VOLUME_LIST                 = "volume.list"
	VOLUME_GET                  = "volume.get"
	SNAPSHOT_LIST               = "snapshot.list"
	HOST_LIST                   = "host.list"
	HOST_GET                    = "host.get"
	DISK_LIST                   = "disk.list"
	CLEAN_RECYCLEBIN            = "recyclebin.clean"
	CREATE_NAMESPACE            = "namespace.create"
	CREATE_VOLUME               = "volume.create"
	EXTEND_VOLUME               = "volume.extend"
	VOLUME_THROTTLE             = "volume.throttle"
	DELETE_VOLUME               = "volume.delete"
	RECOVER_VOLUME              = "volume.recover"
	CLONE_VOLUME                = "volume.clone"
	CREATE_SNAPSHOT             = "snapshot.create"
	CANCEL_SNAPSHOT             = "snapshot.cancel"
	DELETE_SNAPSHOT             = "snapshot.delete"
	FLATTEN                     = "volume.flatten"
	GET_SYSTEM_LOG              = "syslog.get"
	GET_SYSTEM_ALERT            = "alert.get"
	GET_UNREAD_SYSTEM_ALERT_NUM = "alert.unread.num.get"
	UPDATE_READ_SYSTEM_ALERT_ID = "alert.read.id.update"
	GET_ALERT_CONF              = "alert.conf.get"
	UPDATE_ALERT_CONF           = "alert.conf.update"
	GET_ALERT_CANDIDATE         = "alert.candidate.get"
	UPDATE_ALERT_USER           = "alert.user.update"
)

func Exit(r *pigeon.Request, code errno.Errno) bool {
	r.SendJSON(pigeon.JSON{
		"errorCode": strconv.Itoa(code.Code()),
		"errorMsg":  code.Description(),
	})
	if code != errno.OK {
		r.HeadersOut[comm.HEADER_ERROR_CODE] = strconv.Itoa(code.Code())
	}
	if NeedRecordLog(r) {
		agent.WriteSystemLog(r.Context.RemoteIP(), r.HeadersIn[comm.HEADER_LOG_USER], BELONG[r.Args[METHOD]],
			r.Args[METHOD], code.Description(), r.HeadersIn[comm.HEADER_LOG_CONTENT], code.Code())
	}
	return r.Exit(code.HTTPCode())
}

func Default(r *pigeon.Request) bool {
	r.Logger().Warn("unupport request uri", pigeon.Field("uri", r.Uri))
	return Exit(r, errno.UNSUPPORT_REQUEST_URI)
}

func ExitSuccessWithData(r *pigeon.Request, data interface{}) bool {
	r.SendJSON(pigeon.JSON{
		"data":      data,
		"errorCode": "0",
		"errorMsg":  "success",
	})
	if NeedRecordLog(r) {
		agent.WriteSystemLog(r.Context.RemoteIP(), r.HeadersIn[comm.HEADER_LOG_USER], BELONG[r.Args[METHOD]],
			r.Args[METHOD], "success", r.HeadersIn[comm.HEADER_LOG_CONTENT], 0)
	}
	return r.Exit(200)
}
