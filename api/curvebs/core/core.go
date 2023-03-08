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
	USER_LOGIN           = "user.login"
	USER_LOGOUT          = "user.logout"
	USER_CREATE          = "user.create"
	USER_DELETE          = "user.delete"
	USER_UPDATE_PASSWORD = "user.update.password"
	USER_RESET_PASSWORD  = "user.reset.password"
	USER_UPDATE_INFO     = "user.update.info"
	USER_LIST            = "user.list"

	// manager
	STATUS_ETCD                = "status.etcd"
	STATUS_MDS                 = "status.mds"
	STATUS_SNAPSHOTCLONESERVER = "status.snapshotcloneserver"
	STATUS_CHUNKSERVER         = "status.chunkserver"
	STATUS_CLUSTER             = "status.cluster"
	SPACE_CLUSTER              = "space.cluster"
	PERFORMANCE_CLUSTER        = "performance.cluster"
	TOPO_LIST                  = "topo.list"
	TOPO_POOL_LIST             = "topo.pool.list"
	TOPO_POOL_GET              = "topo.pool.get"
	VOLUME_LIST                = "volume.list"
	VOLUME_GET                 = "volume.get"
	SNAPSHOT_LIST              = "snapshot.list"
	HOST_LIST                  = "host.list"
	HOST_GET                   = "host.get"
	DISK_LIST                  = "disk.list"
)

func Exit(r *pigeon.Request, code errno.Errno) bool {
	r.SendJSON(pigeon.JSON{
		"errorCode": strconv.Itoa(code.Code()),
		"errorMsg":  code.Description(),
	})
	if code != errno.OK {
		r.HeadersOut[comm.HEADER_ERROR_CODE] = strconv.Itoa(code.Code())
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
	return r.Exit(200)
}
