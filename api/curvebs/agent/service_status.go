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

	"github.com/SeanHai/curve-go-rpc/rpc/curvebs"
	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/curve-manager/internal/metrics/bsmetric"
	bsrpc "github.com/opencurve/curve-manager/internal/rpc/curvebs"
	"github.com/opencurve/pigeon"
)

const (
	SERVICE_ETCD                  = "etcd"
	SERVICE_MDS                   = "mds"
	SERVICE_CHUNKSERVER           = "chunkserver"
	SERVICE_SNAPSHOT_CLONE_SERVER = "snapshotcloneserver"
)

type VersionNum struct {
	Version string `json:"version"`
	Number  int    `json:"number"`
}

type ChunkServerStatus struct {
	TotalNum   int          `json:"totalNum"`
	OnlineNum  int          `json:"onlineNum"`
	Versions   []VersionNum `json:"versions"`
	NotOnlines []string     `json:"-"`
}

type serviceStatus struct {
	healthy bool
	detail  string
}

func GetEtcdStatus(r *pigeon.Request) (interface{}, errno.Errno) {
	status, err := bsmetric.GetEtcdStatus()
	if err != "" {
		r.Logger().Warn("GetEtcdStatus failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
	}
	return status, errno.OK
}

func GetMdsStatus(r *pigeon.Request) (interface{}, errno.Errno) {
	status, err := bsmetric.GetMdsStatus()
	if err != "" {
		r.Logger().Warn("GetMdsStatus failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
	}
	return status, errno.OK
}

func GetSnapShotCloneServerStatus(r *pigeon.Request) (interface{}, errno.Errno) {
	status, err := bsmetric.GetSnapShotCloneServerStatus()
	if err != "" {
		r.Logger().Warn("GetSnapShotCloneServerStatus failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
	}
	return status, errno.OK
}

func GetChunkServerStatus(l *pigeon.Logger, rId string) (interface{}, errno.Errno) {
	var result ChunkServerStatus
	// get chunkserver form mds
	chunkservers, err := bsrpc.GMdsClient.GetChunkServerInCluster()
	if err != nil {
		l.Error("GetChunkServerStatus bsrpc.GetChunkServerInCluster failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", rId))
		return nil, errno.GET_CHUNKSERVER_IN_CLUSTER_FAILED
	}

	online := 0
	var endponits []string
	for _, cs := range chunkservers {
		endpoint := fmt.Sprintf("%s:%d", cs.HostIp, cs.Port)
		if cs.OnlineStatus == curvebs.ONLINE_STATUS {
			online += 1
			endponits = append(endponits, endpoint)
		} else {
			result.NotOnlines = append(result.NotOnlines, endpoint)
		}
	}
	result.TotalNum = len(chunkservers)
	result.OnlineNum = online

	// get version form metric
	versions, err := bsmetric.GetChunkServerVersion(&endponits)
	if err != nil {
		l.Warn("GetChunkServerStatus bsmetric.GetChunkServerVersion failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", rId))
	} else {
		for k, v := range versions {
			result.Versions = append(result.Versions, VersionNum{
				Version: k,
				Number:  v,
			})
		}
	}
	return &result, errno.OK
}

func checkServiceHealthy(name string) common.QueryResult {
	ret := common.QueryResult{}
	var status []bsmetric.ServiceStatus
	var err string
	ret.Key = name
	switch name {
	case SERVICE_ETCD:
		status, err = bsmetric.GetEtcdStatus()
	case SERVICE_MDS:
		status, err = bsmetric.GetMdsStatus()
	case SERVICE_SNAPSHOT_CLONE_SERVER:
		status, err = bsmetric.GetSnapShotCloneServerStatus()
	default:
		ret.Result = serviceStatus{
			healthy: false,
		}
		ret.Err = fmt.Errorf("Invalid service name")
		return ret
	}

	if err != "" {
		ret.Err = fmt.Errorf(err)
	}

	healthy := true
	detail := ""
	leaderNum := 0
	leaderVec := []string{}
	offline := 0
	offlineVec := []string{}
	for _, s := range status {
		if s.Leader {
			leaderNum++
			leaderVec = append(leaderVec, s.Address)
		}
		if !s.Online {
			offline++
			offlineVec = append(offlineVec, s.Address)
		}
	}

	if leaderNum != 1 {
		healthy = false
		detail = fmt.Sprintf("%s leader number = %d, leaders: %v", name, leaderNum, leaderVec)
	}
	if offline > 0 {
		healthy = false
		detail = fmt.Sprintf("%s %s offline number = %d, offlines: %v", detail, name, offline, offlineVec)
	}
	ret.Result = serviceStatus{
		healthy: healthy,
		detail:  detail,
	}
	return ret
}
