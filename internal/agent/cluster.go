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

	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/metrics/bsmetric"
	bsrpc "github.com/opencurve/curve-manager/internal/rpc/curvebs"
)

func GetClusterSpace() (interface{}, error) {
	var result Space
	// get logical pools form mds
	pools, err := bsrpc.GMdsClient.ListLogicalPool()
	if err != nil {
		return nil, fmt.Errorf("ListLogicalPool failed, %s", err)
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
		return nil, err
	}
	for _, info := range poolInfos {
		result.Total += info.Space.Total
		result.Alloc += info.Space.Alloc
		result.CanRecycled += info.Space.CanRecycled
	}
	return &result, nil
}

func GetClusterPerformance() (interface{}, error) {
	return bsmetric.GetClusterPerformance()
}

func GetClusterStatus() (interface{}, error) {
	retErr := fmt.Errorf("")
	clusterStatus := ClusterStatus{}
	// 1. get pool numbers in cluster
	pools, err := bsrpc.GMdsClient.ListLogicalPool()
	if err != nil {
		clusterStatus.Healthy = false
		clusterStatus.PoolNum = 0
		retErr = fmt.Errorf("ListLogicalPool failed: %s  |  ", err)
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
			retErr = fmt.Errorf("%sCheck %s service healthy failed: %s  |  ", retErr, res.Key.(string), res.Err)
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
		retErr = fmt.Errorf("%sCheck copysets in cluster failed: %s  |  ", retErr, err)
	}
	healthy = health && healthy
	clusterStatus.Healthy = healthy
	clusterStatus.CopysetNum.Total = cs.getCopysetTotalNum()
	clusterStatus.CopysetNum.Unhealthy = cs.getCopysetUnhealthyNum()
	return clusterStatus, retErr
}
