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
	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/curve-manager/internal/snapshotclone"
	"github.com/opencurve/pigeon"
)

func GetSnapshot(r *pigeon.Request, size, page uint32, uuid, user, fileName, status string) (interface{}, errno.Errno) {
	snapshots, err := snapshotclone.GetSnapshot(size, page, uuid, user, fileName, status)
	if err != nil {
		r.Logger().Error("GetSnapshot failed",
			pigeon.Field("size", size),
			pigeon.Field("page", page),
			pigeon.Field("uuid", uuid),
			pigeon.Field("user", user),
			pigeon.Field("fileName", fileName),
			pigeon.Field("status", status),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.LIST_SNAPSHOT_FAILED
	}
	return snapshots, errno.OK
}
