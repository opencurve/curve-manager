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
	"github.com/opencurve/curve-manager/internal/metrics/bsmetric"
	bsrpc "github.com/opencurve/curve-manager/internal/rpc/curvebs"
	"github.com/opencurve/pigeon"
)

const (
	ETCD_SERVICE                  = "etcd"
	MDS_SERVICE                   = "mds"
	SNAPSHOT_CLONE_SERVER_SERVICE = "snapshotcloneserver"
)

type CopysetNum struct {
	Total     uint32 `json:"total" binding:"required"`
	Unhealthy uint32 `json:"unhealthy" binding:"required"`
}

type ClusterStatus struct {
	Healthy    bool       `json:"healthy" binding:"required"`
	PoolNum    uint32     `json:"poolNum" binding:"required"`
	CopysetNum CopysetNum `json:"copysetNum" binding:"required"`
}

func GetClusterSpace(r *pigeon.Request) (interface{}, errno.Errno) {
	var result Space
	// get logical pools form mds
	pools, err := bsrpc.GMdsClient.ListLogicalPool()
	if err != nil {
		r.Logger().Error("GetClusterSpace bsrpc.ListLogicalPool failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.LIST_POOL_FAILED
	}

	var poolInfos []PoolInfo
	for _, pool := range pools {
		var info PoolInfo
		info.Name = pool.Name
		info.Id = pool.Id
		poolInfos = append(poolInfos, info)
	}

	err = getPoolSpace(&poolInfos)
	if err != nil {
		r.Logger().Error("GetClusterSpace getPoolSpace failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_POOL_SPACE_FAILED
	}
	for _, info := range poolInfos {
		result.Total += info.Space.Total
		result.Alloc += info.Space.Alloc
		result.CanRecycled += info.Space.CanRecycled
	}
	return &result, errno.OK
}

func GetClusterPerformance(r *pigeon.Request) (interface{}, errno.Errno) {
	performance, err := bsmetric.GetClusterPerformance()
	if err != nil {
		r.Logger().Error("GetClusterPerformance failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_CLUSTER_PERFORMANCE_FAILED
	}
	return performance, errno.OK
}

func GetClusterStatus(r *pigeon.Request) interface{} {
	clusterStatus := ClusterStatus{}
	// 1. get pool numbers in cluster
	pools, err := bsrpc.GMdsClient.ListLogicalPool()
	if err != nil {
		clusterStatus.Healthy = false
		clusterStatus.PoolNum = 0
		r.Logger().Error("GetClusterStatus bsrpc.ListLogicalPool failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
	}
	clusterStatus.PoolNum = uint32(len(pools))

	healthy := true
	// 2. check service status
	// etcd, mds, snapshotcloneserver
	size := 3
	ret := make(chan common.QueryResult, size)

	go func() {
		ret <- checkServiceHealthy(ETCD_SERVICE)
	}()

	go func() {
		ret <- checkServiceHealthy(MDS_SERVICE)
	}()

	go func() {
		ret <- checkServiceHealthy(SNAPSHOT_CLONE_SERVER_SERVICE)
	}()
	count := 0
	for res := range ret {
		if res.Err != nil {
			r.Logger().Error("GetClusterStatus check service status failed",
				pigeon.Field("service", res.Key.(string)),
				pigeon.Field("error", res.Err),
				pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		}
		healthy = healthy && res.Result.(bool)
		count += 1
		if count >= size {
			break
		}
	}

	// 3. check copyset in cluster
	cs := NewCopyset()
	health, err := cs.checkCopysetsInCluster()
	if err != nil {
		r.Logger().Error("GetClusterStatus checkCopysetsInCluster failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
	}
	healthy = health && healthy
	clusterStatus.Healthy = healthy
	clusterStatus.CopysetNum.Total = cs.getCopysetTotalNum()
	clusterStatus.CopysetNum.Unhealthy = cs.getCopysetUnhealthyNum()
	return clusterStatus
}
