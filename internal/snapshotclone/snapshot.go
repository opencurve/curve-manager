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

package snapshotclone

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/opencurve/curve-manager/internal/common"
)

const (
	ACTION_CREATE_SANPSHOT   = "CreateSnapshot"
	ACTION_CANCEL_SNAPSHOT   = "CancelSnapshot"
	ACTION_GET_SNAPSHOT_LIST = "GetFileSnapshotList"
	ACTION_DELETE_SNAPSHOT   = "DeleteSnapshot"

	// snapshot status
	STATUS_DONE           = "done"
	STATUS_IN_PROCESS     = "inProcess"
	STATUS_DELETING       = "deleting"
	STATUS_ERROR_DELETING = "errorDeleting"
	STATUS_CANCELING      = "canceling"
	STATUS_ERROR          = "error"
	STATUS_NOT_SUPPORT    = "notSupport"

	LIST_SNAPSHOT_LIMIT = 100
)

func getNumericStatus(status string) int {
	switch status {
	case STATUS_DONE:
		return 0
	case STATUS_IN_PROCESS:
		return 1
	case STATUS_DELETING:
		return 2
	case STATUS_ERROR_DELETING:
		return 3
	case STATUS_CANCELING:
		return 4
	case STATUS_ERROR:
		return 5
	default:
		return -1
	}
}

func getStrStatus(status int) string {
	switch status {
	case 0:
		return STATUS_DONE
	case 1:
		return STATUS_IN_PROCESS
	case 2:
		return STATUS_DELETING
	case 3:
		return STATUS_ERROR_DELETING
	case 4:
		return STATUS_CANCELING
	case 5:
		return STATUS_ERROR
	default:
		return STATUS_NOT_SUPPORT
	}
}

func transferSnapshotInfo(in *GetSnapShotResponse, out *ListSnapshotInfo) {
	out.Total = int(in.TotalCount)
	for _, info := range in.Snapshots {
		var item SnapshotInfo
		item.UUID = info.UUID
		item.Name = info.Name
		item.User = info.User
		item.File = info.File
		item.SeqNum = info.SeqNum
		item.Ctime = time.Unix(int64(info.Ctime/1000000), 0).Format(common.TIME_FORMAT)
		item.FileLength = info.FileLength / common.GiB
		item.Status = getStrStatus(info.Status)
		item.Progress = fmt.Sprintf("%d%%", info.Progress)
		out.Info = append(out.Info, item)
	}
}

func GetSnapshot(size, page uint32, uuid, user, fileName, status string) (ListSnapshotInfo, error) {
	var snapshotInfo GetSnapShotResponse
	listSnapshotInfo := ListSnapshotInfo{
		Info: []SnapshotInfo{},
	}
	params := fmt.Sprintf("Action=%s&Version=%s&Limit=%d&Offset=%d",
		ACTION_GET_SNAPSHOT_LIST, SNAPSHOT_CLONE_VERSION, size, (page-1)*size)
	if uuid != "" {
		params = fmt.Sprintf("%s&UUID=%s", params, uuid)
	}
	if user != "" {
		params = fmt.Sprintf("%s&User=%s", params, user)
	}
	if fileName != "" {
		params = fmt.Sprintf("%s&File=%s", params, fileName)
	}
	if status != "" {
		s := getNumericStatus(status)
		if s < 0 {
			return listSnapshotInfo, fmt.Errorf("status not support")
		}
		params = fmt.Sprintf("%s&Status=%d", params, s)
	}

	resp, err := GSnapshotCloneClient.sendHttp2SnapshotClone(params)
	if err != nil {
		return listSnapshotInfo, err
	}

	err = json.Unmarshal([]byte(resp), &snapshotInfo)
	if err != nil {
		return listSnapshotInfo, err
	}

	if snapshotInfo.Code != ERROR_CODE_SUCCESS {
		return listSnapshotInfo, fmt.Errorf(snapshotInfo.Message)
	}

	transferSnapshotInfo(&snapshotInfo, &listSnapshotInfo)
	return listSnapshotInfo, nil
}

func GetAllSnapshot(user, fileName, status string) ([]SnapshotInfo, error) {
	result := []SnapshotInfo{}
	page := 1
	limit := LIST_SNAPSHOT_LIMIT
	for true {
		infos, err := GetSnapshot(uint32(limit), uint32(page), "", user, fileName, status)
		if err != nil {
			return nil, err
		}
		result = append(result, infos.Info...)
		if infos.Total < limit {
			break
		}
		page++
	}
	return result, nil
}

func CreateSnapshot(volumeName, user, snapshotName string) error {
	params := fmt.Sprintf("Action=%s&Version=%s&User=%s&File=%s&Name=%s",
		ACTION_CREATE_SANPSHOT, SNAPSHOT_CLONE_VERSION, user, volumeName, snapshotName)
	resp, err := GSnapshotCloneClient.sendHttp2SnapshotClone(params)
	if err != nil {
		return fmt.Errorf("err:%v, params:%s", err, params)
	}
	var response CreateSnapshotCloneResponse
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return fmt.Errorf("err:%v, params:%s", err, params)
	}
	if response.Code != ERROR_CODE_SUCCESS {
		return fmt.Errorf("err:%s, params:%s", response.Message, params)
	}
	return nil
}

func CancelSnapshot(uuid, user, volumeName string) error {
	params := fmt.Sprintf("Action=%s&Version=%s&User=%s&File=%s&UUID=%s",
		ACTION_CANCEL_SNAPSHOT, SNAPSHOT_CLONE_VERSION, user, volumeName, uuid)
	resp, err := GSnapshotCloneClient.sendHttp2SnapshotClone(params)
	if err != nil {
		return fmt.Errorf("err:%v, params:%s", err, params)
	}
	var response SnapshotCloneResponse
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return fmt.Errorf("err:%v, params:%s", err, params)
	}
	if response.Code != ERROR_CODE_SUCCESS {
		return fmt.Errorf("err:%s, params:%s", response.Message, params)
	}
	return nil
}

func DeleteSnapshot(uuid, volumeName, user string) error {
	params := fmt.Sprintf("Action=%s&Version=%s&User=%s&File=%s&UUID=%s",
		ACTION_DELETE_SNAPSHOT, SNAPSHOT_CLONE_VERSION, user, volumeName, uuid)
	resp, err := GSnapshotCloneClient.sendHttp2SnapshotClone(params)
	if err != nil {
		return fmt.Errorf("err:%v, params:%s", err, params)
	}
	var response SnapshotCloneResponse
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return fmt.Errorf("err:%v, params:%s", err, params)
	}
	if response.Code != ERROR_CODE_SUCCESS {
		return fmt.Errorf("err:%s, params:%s", response.Message, params)
	}
	return nil
}
