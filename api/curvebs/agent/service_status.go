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
	"fmt"

	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/curve-manager/internal/metrics/bsmetric"
	bsrpc "github.com/opencurve/curve-manager/internal/rpc/curvebs"
	"github.com/opencurve/pigeon"
)

type VersionNum struct {
	Version string `json:"version"`
	Number  int    `json:"number"`
}

type ChunkServerStatus struct {
	TotalNum  int          `json:"totalNum"`
	OnlineNum int          `json:"onlineNum"`
	Versions  []VersionNum `json:"versions"`
}

func GetEtcdStatus(r *pigeon.Request) (interface{}, errno.Errno) {
	status, err := bsmetric.GetEtcdStatus()
	if err != "" {
		r.Logger().Error("GetEtcdStatus failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_ETCD_STATUS_FAILED
	}
	return status, errno.OK
}

func GetMdsStatus(r *pigeon.Request) (interface{}, errno.Errno) {
	status, err := bsmetric.GetMdsStatus()
	if err != "" {
		r.Logger().Error("GetMdsStatus failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_MDS_STATUS_FAILED
	}
	return status, errno.OK
}

func GetSnapShotCloneServerStatus(r *pigeon.Request) (interface{}, errno.Errno) {
	status, err := bsmetric.GetSnapShotCloneServerStatus()
	if err != "" {
		r.Logger().Error("GetSnapShotCloneServerStatus failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_SNAPSHOT_CLONE_STATUS_FAILED
	}
	return status, errno.OK
}

func GetChunkServerStatus(r *pigeon.Request) (interface{}, errno.Errno) {
	var result ChunkServerStatus
	// get chunkserver form mds
	chunkservers, err := bsrpc.GMdsClient.GetChunkServerInCluster()
	if err != nil {
		r.Logger().Error("GetChunkServerStatus bsrpc.GetChunkServerInCluster failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_CHUNKSERVER_IN_CLUSTER_FAILED
	}

	online := 0
	var endponits []string
	for _, cs := range chunkservers {
		if cs.OnlineStatus == bsrpc.ONLINE_STATUS {
			online += 1
		}
		endponits = append(endponits, fmt.Sprintf("%s:%d", cs.HostIp, cs.Port))
	}
	result.TotalNum = len(chunkservers)
	result.OnlineNum = online

	// get version form metric
	versions, err := bsmetric.GetChunkServerVersion(&endponits)
	if err != nil {
		r.Logger().Error("GetChunkServerStatus bsmetric.GetChunkServerVersion failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_CHUNKSERVER_VERSION_FAILED
	}
	for k, v := range versions {
		result.Versions = append(result.Versions, VersionNum{
			Version: k,
			Number:  v,
		})
	}
	return &result, errno.OK
}

func checkServiceHealthy(name string) common.QueryResult {
	var ret common.QueryResult
	var status []bsmetric.ServiceStatus
	var err string
	ret.Key = name
	switch name {
	case ETCD_SERVICE:
		status, err = bsmetric.GetEtcdStatus()
	case MDS_SERVICE:
		status, err = bsmetric.GetMdsStatus()
	case SNAPSHOT_CLONE_SERVER_SERVICE:
		status, err = bsmetric.GetSnapShotCloneServerStatus()
	default:
		ret.Result = false
		ret.Err = fmt.Errorf("Invalid service name")
		return ret
	}

	if err != "" {
		ret.Err = fmt.Errorf(err)
		ret.Result = false
		return ret
	}

	leaderNum := 0
	hasOffline := false
	for _, s := range status {
		if s.Leader {
			leaderNum += 1
		}
		if !s.Online {
			hasOffline = true
		}
	}

	ret.Result = leaderNum == 1 && !hasOffline
	return ret
}
