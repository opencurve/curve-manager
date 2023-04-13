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
	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/curve-manager/internal/snapshotclone"
	"github.com/opencurve/pigeon"
)

type Snapshot struct {
	UUID       string `json:"uuid" binding:"required"`
	User       string `json:"user" binding:"required"`
	VolumeName string `json:"volumeName" binding:"required"`
}

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

func CreateSnapshot(r *pigeon.Request, volumeName, user, snapshotName string) errno.Errno {
	err := snapshotclone.CreateSnapshot(volumeName, user, snapshotName)
	if err != nil {
		r.Logger().Error("CreateSnapshot failed",
			pigeon.Field("volumeName", volumeName),
			pigeon.Field("user", user),
			pigeon.Field("snapshotName", snapshotName),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return errno.CREATE_SNAPSHOT_FAILED
	}
	return errno.OK
}

func CancelSnapshot(r *pigeon.Request, snapshots []Snapshot) errno.Errno {
	size := len(snapshots)
	if size == 0 {
		return errno.OK
	}
	ret := make(chan common.QueryResult, size)
	for _, snapshot := range snapshots {
		go func(uuid, user, volumeName string) {
			err := snapshotclone.CancelSnapshot(uuid, user, volumeName)
			ret <- common.QueryResult{
				Key:    uuid,
				Err:    err,
				Result: nil,
			}
		}(snapshot.UUID, snapshot.User, snapshot.VolumeName)
	}

	count := 0
	success := true
	for res := range ret {
		if res.Err != nil {
			r.Logger().Error("CancelSnapshot failed",
				pigeon.Field("uuid", res.Key.(string)),
				pigeon.Field("error", res.Err),
				pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
			success = false
		}
		count += 1
		if count >= size {
			break
		}
	}
	if !success {
		return errno.CANCEL_SNAPSHOT_FAILED
	}
	return errno.OK
}

func DeleteSnapshot(r *pigeon.Request, fileName, user string, uuids []string, failed bool) errno.Errno {
	toDelete := []snapshotclone.SnapshotInfo{}
	if fileName == "" && user == "" && len(uuids) == 0 && !failed {
		return errno.OK
	}

	status := ""
	if failed {
		status = snapshotclone.STATUS_ERROR
	}
	if len(uuids) != 0 {
		for _, uuid := range uuids {
			info, err := snapshotclone.GetSnapshot(1, 1, uuid, user, fileName, status)
			if err != nil {
				r.Logger().Error("DeleteSnapshot GetSnapshot failed",
					pigeon.Field("fileName", fileName),
					pigeon.Field("user", user),
					pigeon.Field("uuids", uuids),
					pigeon.Field("status", status),
					pigeon.Field("uuid", uuid),
					pigeon.Field("error", err),
					pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
				return errno.GET_SNAPSHOT_FAILED
			}
			if info.Total > 0 {
				toDelete = append(toDelete, info.Info...)
			} else {
				r.Logger().Error("GetSnapshot by uuid failed",
					pigeon.Field("uuid", uuid),
					pigeon.Field("user", user),
					pigeon.Field("fileName", fileName),
					pigeon.Field("status", status),
					pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
				return errno.GET_SNAPSHOT_FAILED
			}
		}
	} else {
		snapshots, err := snapshotclone.GetAllSnapshot(user, fileName, status)
		if err != nil {
			r.Logger().Error("DeleteSnapshot GetAllSnapshot failed",
				pigeon.Field("fileName", fileName),
				pigeon.Field("user", user),
				pigeon.Field("status", status),
				pigeon.Field("error", err),
				pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
			return errno.GET_SNAPSHOT_FAILED
		}
		toDelete = append(toDelete, snapshots...)
	}
	for _, s := range toDelete {
		err := snapshotclone.DeleteSnapshot(s.UUID, s.File, s.User)
		if err != nil {
			r.Logger().Error("DeleteSnapshot failed",
				pigeon.Field("fileName", s.File),
				pigeon.Field("user", s.User),
				pigeon.Field("uuid", s.UUID),
				pigeon.Field("error", err),
				pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
			return errno.DELETE_SNAPSHOT_FAILED
		}
	}
	return errno.OK
}
