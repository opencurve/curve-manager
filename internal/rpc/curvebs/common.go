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

package curvebs

type PhysicalPool struct {
	Id   uint32 `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc"`
}

type LogicalPool struct {
	Id             uint32 `json:"id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	PhysicalPoolId uint32 `json:"physicalPoolId" binding:"required"`
	Type           string `json:"type" binding:"required"`
	CreateTime     string `json:"createTime" binding:"required"`
	AllocateStatus string `json:"allocateStatus" binding:"required"`
	ScanEnable     bool   `json:"scanEnable"`
}

type Zone struct {
	Id               uint32 `json:"id" binding:"required"`
	Name             string `json:"name" binding:"required"`
	PhysicalPoolId   uint32 `json:"physicalPoolId" binding:"required"`
	PhysicalPoolName string `json:"physicalName" binding:"required"`
	Desc             string `json:"desc"`
}

type Server struct {
	Id               uint32 `json:"id" binding:"required"`
	HostName         string `json:"hostName" binding:"required"`
	InternalIp       string `json:"internalIp" binding:"required"`
	InternalPort     uint32 `json:"internalPort" binding:"required"`
	ExternalIp       string `json:"externalIp" binding:"required"`
	ExternalPort     uint32 `json:"externalPort" binding:"required"`
	ZoneId           uint32 `json:"zoneId" binding:"required"`
	ZoneName         string `json:"zoneName" binding:"required"`
	PhysicalPoolId   uint32 `json:"physicalPoolId" binding:"required"`
	PhysicalPoolName string `json:"physicalName" binding:"required"`
	Desc             string `json:"desc"`
}

type ChunkServer struct {
	Id           uint32 `json:"id" binding:"required"`
	DiskType     string `json:"diskType" binding:"required"`
	HostIp       string `json:"hostIp" binding:"required"`
	Port         uint32 `json:"port" binding:"required"`
	Status       string `json:"status" binding:"required"`
	DiskStatus   string `json:"diskStatus" binding:"required"`
	OnlineStatus string `json:"onlineStatus" binding:"required"`
	MountPoint   string `json:"mountPoint" binding:"required"`
	DiskCapacity string `json:"diskCapacity" binding:"required"`
	DiskUsed     string `json:"diskUsed" binding:"required"`
	ExternalIp   string `json:"externalIp"`
}

type ThrottleParams struct {
	Type        string `json:"type"`
	Limit       uint64 `json:"limit"`
	Burst       uint64 `json:"burst"`
	BurstLength uint64 `json:"burstLength"`
}

type FileInfo struct {
	Id                   uint64           `json:"id"`
	FileName             string           `json:"fileName"`
	ParentId             uint64           `json:"parentId"`
	FileType             string           `json:"fileType"`
	Owner                string           `json:"owner"`
	ChunkSize            uint32           `json:"chunkSize"`
	SegmentSize          uint32           `json:"segmentSize"`
	Length               uint64           `json:"length"`
	AllocateSize         uint64           `json:"alloc"`
	Ctime                string           `json:"ctime"`
	SeqNum               uint64           `json:"seqNum"`
	FileStatus           string           `json:"fileStatus"`
	OriginalFullPathName string           `json:"originalFullPathName"`
	CloneSource          string           `json:"cloneSource"`
	CloneLength          uint64           `json:"cloneLength"`
	StripeUnit           uint64           `json:"stripeUnit"`
	StripeCount          uint64           `json:"stripeCount"`
	ThrottleParams       []ThrottleParams `json:"throttleParams"`
	Epoch                uint64           `json:"epoch"`
}

type CopySetInfo struct {
	LogicalPoolId      uint32 `json:"logicalPoolId" binding:"required"`
	CopysetId          uint32 `json:"copysetId" binding:"required"`
	Scanning           bool   `json:"scanning"`
	LastScanSec        uint64 `json:"lastScanSec"`
	LastScanConsistent bool   `json:"lastScanConsistent"`
}

type ChunkServerLocation struct {
	ChunkServerId uint32 `json:"chunkServerId" binding:"required"`
	HostIp        string `json:"hostIp" binding:"required"`
	Port          uint32 `josn:"port" binding:"required"`
	ExternalIp    string `json:"externalIp"`
}
type CopySetServerInfo struct {
	CopysetId uint32                `json:"copysetId" binding:"required"`
	CsLocs    []ChunkServerLocation `json:"csLocs" binding:"required"`
}
