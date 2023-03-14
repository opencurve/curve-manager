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

const (
	SERVICE_ROUTER         = "/SnapshotCloneService"
	SNAPSHOT_CLONE_VERSION = "0.0.6"

	ERROR_CODE_SUCCESS = "0"
)

// struct received from snapshot clone server
type GetSnapShotResponse struct {
	Code       string     `json:"Code" binding:"required"`
	Message    string     `json:"Message" binding:"required"`
	RequestId  string     `json:"RequestId" binding:"required"`
	Snapshots  []snapshot `json:"Snapshots"`
	TotalCount uint64     `json:"TotalCount"`
}

type snapshot struct {
	UUID       string `json:"UUID" binding:"required"`
	User       string `json:"User" binding:"required"`
	File       string `json:"File" binding:"required"`
	Name       string `json:"Name" binding:"required"`
	SeqNum     uint64 `json:"SeqNum" binding:"required"`
	Ctime      uint64 `json:"Time" binding:"required"`
	FileLength uint64 `json:"FileLength" binding:"required"`
	Status     int    `json:"Status" binding:"required"`
	Progress   uint32 `json:"Progress" binding:"required"`
}

// struct return to uplayer
type ListSnapshotInfo struct {
	Total int            `json:"total" binding:"required"`
	Info  []SnapshotInfo `json:"info" binding:"required"`
}

type SnapshotInfo struct {
	UUID       string `json:"uuid" binding:"required"`
	User       string `json:"user" binding:"required"`
	File       string `json:"file" binding:"required"`
	Name       string `json:"name" binding:"required"`
	SeqNum     uint64 `json:"seqNum" binding:"required"`
	Ctime      string `json:"time" binding:"required"`
	FileLength uint64 `json:"length" binding:"required"`
	Status     string `json:"status" binding:"required"`
	Progress   string `json:"Progress" binding:"required"`
}

type CreateSnapshotCloneResponse struct {
	Code      string `json:"Code" binding:"required"`
	Message   string `json:"Message" binding:"required"`
	RequestId string `json:"RequestId" binding:"required"`
	UUID      string `json:"UUID"`
}

type SnapshotCloneResponse struct {
	Code      string `json:"Code" binding:"required"`
	Message   string `json:"Message" binding:"required"`
	RequestId string `json:"RequestId" binding:"required"`
}

type TaskInfo struct {
	File       string `json:"File"`
	FileType   int    `json:"FileType"`
	IsLazy     bool   `json:"IsLazy"`
	NextStep   int    `json:"NextStep"`
	Progress   int    `json:"Progress"`
	Src        string `json:"Src"`
	TaskStatus int    `json:"TaskStatus"`
	TaskType   int    `json:"TaskType"`
	Time       uint64 `json:"Time"`
	UUID       string `json:"UUID"`
	User       string `json:"User"`
}

type GetCloneTasksResponse struct {
	Code       string     `json:"Code" binding:"required"`
	Message    string     `json:"Message" binding:"required"`
	RequestId  string     `json:"RequestId" binding:"required"`
	TaskInfos  []TaskInfo `json:"TaskInfos"`
	TotalCount uint64     `json:"TotalCount"`
}
