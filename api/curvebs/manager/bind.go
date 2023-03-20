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

package manager

import (
	"github.com/opencurve/curve-manager/api/curvebs/agent"
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/pigeon"
)

const (
	MODULE = "manager"
)

var METHOD_REQUEST map[string]Request

type (
	HandlerFunc func(r *pigeon.Request, ctx *Context) bool

	Context struct {
		Data interface{}
	}

	Request struct {
		httpMethod string
		method     string
		vType      interface{}
		handler    HandlerFunc
	}
)

func init() {
	METHOD_REQUEST = map[string]Request{}
	for _, request := range requests {
		METHOD_REQUEST[request.method] = request
		core.BELONG[request.method] = MODULE
	}
}

type GetEtcdStatusRequest struct{}

type GetMdsStatusRequest struct{}

type GetSnapShotCloneServerStatusRequest struct{}

type ListTopologyRequest struct{}

type ListLogicalPoolRequest struct{}

type GetLogicalPoolRequest struct {
	Id uint32 `json:"id" binding:"required"`
}

type GetChunkServerStatusRequest struct{}

type GetClusterStatusRequest struct{}

type GetClusterSpaceRequest struct{}

type GetClusterSpaceTrendRequest struct {
	Start    uint64 `json:"start" binding:"required"`
	End      uint64 `json:"end" binding:"required"`
	Interval uint64 `json:"interval" binding:"required"`
}

type GetClusterPerformanceRequest struct{}

type ListHostRequest struct {
	Size uint32 `json:"size" binding:"required"`
	Page uint32 `json:"page" binding:"required"`
}

type ListVolumeRequest struct {
	Size          uint32 `json:"size" binding:"required"`
	Page          uint32 `json:"page" binding:"required"`
	Path          string `json:"path" default:"/"`
	SortKey       string `json:"sortKey" default:"id"`
	SortDirection int    `json:"sortDirection" default:"1"`
}

type GetVolumeRequest struct {
	VolumeName string `json:"volumeName" binding:"required"`
}

type ListSnapshotRequest struct {
	Size     uint32 `json:"size" binding:"required"`
	Page     uint32 `json:"page" binding:"required"`
	UUID     string `json:"uuid"`
	User     string `json:"user"`
	FileName string `json:"fileName"`
	Status   string `json:"status"`
}

type GetHostRequest struct {
	HostName string `json:"hostName" binding:"required"`
}

type ListDiskRequest struct {
	Size     uint32 `json:"size" binding:"required"`
	Page     uint32 `json:"page" binding:"required"`
	HostName string `json:"hostName"`
}

type CleanRecycleBinRequest struct {
	Expiration uint64 `json:"expiration" default:"0"`
}

type CreateNameSpaceRequest struct {
	Name     string `json:"name" binding:"required"`
	User     string `json:"user" binding:"required"`
	PassWord string `json:"password"`
}

type CreateVolumeRequest struct {
	VolumeName string `json:"volumeName" binding:"required"`
	User       string `json:"user" binding:"required"`
	Length     uint64 `json:"length" binding:"required"`
	PassWord   string `json:"password"`
	StripUnit  uint64 `json:"stripUnit"`
	StripCount uint64 `json:"stripCount"`
}

type ExtendVolumeRequest struct {
	VolumeName string `json:"volumeName" binding:"required"`
	Length     uint64 `json:"length" binding:"required"`
}

type VolumeThrottleRequest struct {
	VolumeName   string `json:"volumeName" binding:"required"`
	ThrottleType string `json:"throttleType" binding:"required"`
	Limit        uint64 `json:"limit" binding:"required"`
	Burst        uint64 `json:"burst"`
	BurstLength  uint64 `json:"burstLength"`
}

type DeleteVolumeRequest struct {
	VolumeNames map[string]string `json:"volumeNames" binding:"required"`
}

type RecoverVolumeRequest struct {
	VolumeIds map[string]uint64 `json:"volumeIds" binding:"required"`
}

type CloneVolumeRequest struct {
	Src  string `json:"src" binding:"required"`
	Dest string `json:"dest" binding:"required"`
	User string `json:"user" binding:"required"`
	Lazy *bool  `json:"lazy" binding:"required"`
}

type CreateSnapshotRequest struct {
	VolumeName   string `json:"volumeName" binding:"required"`
	User         string `json:"user" binding:"required"`
	SnapshotName string `json:"snapshotName" binding:"required"`
}

type CancelSnapshotRequest struct {
	Snapshots []agent.Snapshot `json:"snapshots" binding:"required"`
}

type DeleteSnapshotRequest struct {
	FileName string   `json:"fileName"`
	User     string   `json:"user"`
	UUIDs    []string `json:"uuids"`
	Failed   bool     `json:"failed" default:"false"`
}

type FlattenRequest struct {
	VolumeName string `json:"volumeName" binding:"required"`
	User       string `json:"user" binding:"required"`
}

type GetSysLogRequest struct {
	Start  int64  `json:"start" default:"0"`
	End    int64  `json:"end" default:"0"`
	Page   uint32 `json:"page" binding:"required"`
	Size   uint32 `json:"size" binding:"required"`
	Filter string `json:"filter"`
}

var requests = []Request{
	{
		core.HTTP_GET,
		core.STATUS_ETCD,
		GetEtcdStatusRequest{},
		GetEtcdStatus,
	},
	{
		core.HTTP_GET,
		core.STATUS_MDS,
		GetMdsStatusRequest{},
		GetMdsStatus,
	},
	{
		core.HTTP_GET,
		core.STATUS_SNAPSHOTCLONESERVER,
		GetSnapShotCloneServerStatusRequest{},
		GetSnapShotCloneServerStatus,
	},
	{
		core.HTTP_GET,
		core.STATUS_CHUNKSERVER,
		GetChunkServerStatusRequest{},
		GetChunkServerStatus,
	},
	{
		core.HTTP_GET,
		core.STATUS_CLUSTER,
		GetClusterStatusRequest{},
		GetClusterStatus,
	},
	{
		core.HTTP_GET,
		core.SPACE_CLUSTER,
		GetClusterSpaceRequest{},
		GetClusterSpace,
	},
	{
		core.HTTP_POST,
		core.SPACE_TREND_CLUSTER,
		GetClusterSpaceTrendRequest{},
		GetClusterSpaceTrend,
	},
	{
		core.HTTP_GET,
		core.PERFORMANCE_CLUSTER,
		GetClusterPerformanceRequest{},
		GetClusterPerformance,
	},
	{
		core.HTTP_GET,
		core.TOPO_LIST,
		ListTopologyRequest{},
		ListTopology,
	},
	{
		core.HTTP_GET,
		core.TOPO_POOL_LIST,
		ListLogicalPoolRequest{},
		ListLogicalPool,
	},
	{
		core.HTTP_POST,
		core.TOPO_POOL_GET,
		GetLogicalPoolRequest{},
		GetLogicalPool,
	},
	{
		core.HTTP_POST,
		core.VOLUME_LIST,
		ListVolumeRequest{},
		ListVolume,
	},
	{
		core.HTTP_POST,
		core.VOLUME_GET,
		GetVolumeRequest{},
		GetVolume,
	},
	{
		core.HTTP_POST,
		core.SNAPSHOT_LIST,
		ListSnapshotRequest{},
		ListSnapshot,
	},
	{
		core.HTTP_POST,
		core.HOST_LIST,
		ListHostRequest{},
		ListHost,
	},
	{
		core.HTTP_POST,
		core.HOST_GET,
		GetHostRequest{},
		GetHost,
	},
	{
		core.HTTP_POST,
		core.DISK_LIST,
		ListDiskRequest{},
		ListDisk,
	},
	{
		core.HTTP_POST,
		core.CLEAN_RECYCLEBIN,
		CleanRecycleBinRequest{},
		CleanRecycleBin,
	},
	{
		core.HTTP_POST,
		core.CREATE_NAMESPACE,
		CreateNameSpaceRequest{},
		CreateNameSpace,
	},
	{
		core.HTTP_POST,
		core.CREATE_VOLUME,
		CreateVolumeRequest{},
		CreateVolume,
	},
	{
		core.HTTP_POST,
		core.EXTEND_VOLUME,
		ExtendVolumeRequest{},
		ExtendVolume,
	},
	{
		core.HTTP_POST,
		core.VOLUME_THROTTLE,
		VolumeThrottleRequest{},
		VolumeThrottle,
	},
	{
		core.HTTP_POST,
		core.DELETE_VOLUME,
		DeleteVolumeRequest{},
		DeleteVolume,
	},
	{
		core.HTTP_POST,
		core.RECOVER_VOLUME,
		RecoverVolumeRequest{},
		RecoverVolume,
	},
	{
		core.HTTP_POST,
		core.CLONE_VOLUME,
		CloneVolumeRequest{},
		CloneVolume,
	},
	{
		core.HTTP_POST,
		core.CREATE_SNAPSHOT,
		CreateSnapshotRequest{},
		CreateSnapshot,
	},
	{
		core.HTTP_POST,
		core.CANCEL_SNAPSHOT,
		CancelSnapshotRequest{},
		CancelSnapshot,
	},
	{
		core.HTTP_POST,
		core.DELETE_SNAPSHOT,
		DeleteSnapshotRequest{},
		DeleteSnapshot,
	},
	{
		core.HTTP_POST,
		core.FLATTEN,
		FlattenRequest{},
		Flatten,
	},
	{
		core.HTTP_POST,
		core.GET_SYSTEM_LOG,
		GetSysLogRequest{},
		GetSysLog,
	},
}
