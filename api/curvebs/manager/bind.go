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
}
