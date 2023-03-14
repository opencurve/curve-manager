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
* Created Date: 2023-03-10
* Author: wanghai (SeanHai)
 */

package snapshotclone

import (
	"encoding/json"
	"fmt"

	"github.com/opencurve/pigeon"
)

const (
	ACTION_GET_CLONE_TASKS = "GetCloneTasks"
	ACTION_CLONE           = "Clone"
	ACTION_FLATTEN         = "Flatten"

	CLONE_STATUS_DONE           = "done"
	CLONE_STATUS_CLONING        = "cloning"
	CLONE_STATUS_RECOVERING     = "recovering"
	CLONE_STATUS_CLEANING       = "cleaning"
	CLONE_STATUS_ERROR_CLEANING = "errorCleaning"
	CLONE_STATUS_ERROR          = "error"
	CLONE_STATUS_RETRYING       = "retrying"
	CLONE_STATUS_METAINSTALLED  = "metaInstalled"
)

var cloneStatus map[int]string = map[int]string{
	0: CLONE_STATUS_DONE,
	1: CLONE_STATUS_CLONING,
	2: CLONE_STATUS_RECOVERING,
	3: CLONE_STATUS_CLEANING,
	4: CLONE_STATUS_ERROR_CLEANING,
	5: CLONE_STATUS_ERROR,
	6: CLONE_STATUS_RETRYING,
	7: CLONE_STATUS_METAINSTALLED,
}

func CreateClone(src, dest, user string, lazy bool) error {
	params := fmt.Sprintf("Action=%s&Version=%s&User=%s&Source=%s&Destination=%s&Lazy=%t",
		ACTION_CLONE, SNAPSHOT_CLONE_VERSION, user, src, dest, lazy)
	resp, err := GSnapshotCloneClient.sendHttp2SnapshotClone(params)
	if err != nil {
		return err
	}
	var response CreateSnapshotCloneResponse
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return err
	}
	if response.Code != ERROR_CODE_SUCCESS {
		return fmt.Errorf(response.Message)
	}
	return nil
}

func GetCloneTaskNeedFlatten(r *pigeon.Request, volumeName, user string) ([]string, error) {
	uuids := []string{}
	params := fmt.Sprintf("Action=%s&Version=%s&User=%s&File=%s",
		ACTION_GET_CLONE_TASKS, SNAPSHOT_CLONE_VERSION, user, volumeName)
	resp, err := GSnapshotCloneClient.sendHttp2SnapshotClone(params)
	if err != nil {
		return uuids, err
	}
	var response GetCloneTasksResponse
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return uuids, err
	}
	if response.Code != ERROR_CODE_SUCCESS {
		return uuids, fmt.Errorf(response.Message)
	}
	for _, v := range response.TaskInfos {
		if cloneStatus[v.TaskStatus] == CLONE_STATUS_METAINSTALLED {
			uuids = append(uuids, v.UUID)
		}
	}
	return uuids, nil
}

func Flatten(user, uuid string) error {
	params := fmt.Sprintf("Action=%s&Version=%s&User=%s&UUID=%s",
		ACTION_FLATTEN, SNAPSHOT_CLONE_VERSION, user, uuid)
	resp, err := GSnapshotCloneClient.sendHttp2SnapshotClone(params)
	if err != nil {
		return err
	}
	var response SnapshotCloneResponse
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return err
	}
	if response.Code != ERROR_CODE_SUCCESS {
		return fmt.Errorf(response.Message)
	}
	return nil
}
