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

package bsmetric

import (
	"fmt"
	"strconv"

	"github.com/opencurve/curve-manager/internal/common"
	metricomm "github.com/opencurve/curve-manager/internal/metrics/common"
	"github.com/opencurve/curve-manager/internal/metrics/core"
)

const (
	CLUSTER_METRIC_PREFIX    = "topology_metric_cluster_"
	CLUSTER_LOGICAL_CAPACITY = "topology_metric_cluster_logical_capacity"
	CLUSTER_LOGICAL_ALLOC    = "topology_metric_cluster_logical_alloc"

	LOGICAL_POOL_METRIC_PREFIX    = "topology_metric_logical_pool_"
	LOGICAL_POOL_LOGICAL_CAPACITY = "_logical_capacity"
	LOGICAL_POOL_LOGICAL_ALLOC    = "_logical_alloc"
	LOGICAL_POOL_SERVER_NUM       = "_server_num"
	LOGICAL_POOL_CHUNKSERVER_NUM  = "_chunkserver_num"
	LOGICAL_POOL_COPYSET_NUM      = "_copyset_num"

	FILE_PREFIX = "curve_client_"
)

type PoolItemNum struct {
	ServerNum      uint32
	ChunkServerNum uint32
	CopysetNum     uint32
}

func GetEtcdStatus() ([]ServiceStatus, string) {
	// init value
	ret := []ServiceStatus{}
	retMap := make(map[string]*ServiceStatus)
	for _, addr := range core.GMetricClient.EtcdAddr {
		retMap[addr] = &ServiceStatus{
			Address: addr,
			Version: "",
			Leader:  false,
			Online:  false,
		}
	}

	// version, leader
	requestSize := 2
	results := make(chan metricomm.MetricResult, requestSize)
	go metricomm.QueryInstantMetric(metricomm.ETCD_CLUSTER_VERSION_NAME, &results)
	go metricomm.QueryInstantMetric(metricomm.ETCD_SERVER_IS_LEADER_NAME, &results)

	count := 0
	for res := range results {
		if res.Err != nil {
			return ret, res.Err.Error()
		}
		if res.Key.(string) == metricomm.ETCD_CLUSTER_VERSION_NAME {
			versions := metricomm.ParseVectorMetric(res.Result.(*metricomm.QueryResponseOfVector), false)
			for k, v := range versions {
				if _, ok := retMap[k]; ok {
					(*retMap[k]).Version = v[metricomm.ETCD_CLUSTER_VERSION]
				}
			}
		} else {
			leaders := metricomm.ParseVectorMetric(res.Result.(*metricomm.QueryResponseOfVector), true)
			for k, v := range leaders {
				if _, ok := retMap[k]; ok {
					if v["value"] == "1" {
						(*retMap[k]).Leader = true
					} else {
						(*retMap[k]).Leader = false
					}
					(*retMap[k]).Online = true
				}
			}
		}
		count += 1
		if count >= requestSize {
			break
		}
	}
	for _, v := range retMap {
		ret = append(ret, *v)
	}
	return ret, ""
}

func GetPoolSpace(name string) (*metricomm.Space, error) {
	space := metricomm.Space{}
	poolName := metricomm.FormatToMetricName(name)

	// total, alloc
	requestSize := 2
	results := make(chan metricomm.MetricResult, requestSize)
	totalName := fmt.Sprintf("%s%s%s", LOGICAL_POOL_METRIC_PREFIX, poolName, LOGICAL_POOL_LOGICAL_CAPACITY)
	usedName := fmt.Sprintf("%s%s%s", LOGICAL_POOL_METRIC_PREFIX, poolName, LOGICAL_POOL_LOGICAL_ALLOC)

	go metricomm.QueryInstantMetric(totalName, &results)
	go metricomm.QueryInstantMetric(usedName, &results)

	count := 0
	for res := range results {
		if res.Err != nil {
			return &space, res.Err
		}
		ret := metricomm.ParseVectorMetric(res.Result.(*metricomm.QueryResponseOfVector), true)
		if res.Key.(string) == totalName {
			for _, v := range ret {
				total, ok := strconv.ParseUint(v["value"], 10, 64)
				if ok != nil {
					return nil, ok
				}
				space.Total = total / common.GiB
				break
			}
		} else {
			for _, v := range ret {
				used, ok := strconv.ParseUint(v["value"], 10, 64)
				if ok != nil {
					return nil, ok
				}
				space.Used = used / common.GiB
				break
			}
		}
		count += 1
		if count >= requestSize {
			break
		}
	}
	return &space, nil
}

func GetPoolItemNum(name string) (*PoolItemNum, error) {
	poolItemNum := PoolItemNum{}
	poolName := metricomm.FormatToMetricName(name)

	// serverNUm, chunkserverNum, copysetNum
	requestSize := 3
	results := make(chan metricomm.MetricResult, requestSize)
	serverName := fmt.Sprintf("%s%s%s", LOGICAL_POOL_METRIC_PREFIX, poolName, LOGICAL_POOL_SERVER_NUM)
	chunkserverName := fmt.Sprintf("%s%s%s", LOGICAL_POOL_METRIC_PREFIX, poolName, LOGICAL_POOL_CHUNKSERVER_NUM)
	copysetName := fmt.Sprintf("%s%s%s", LOGICAL_POOL_METRIC_PREFIX, poolName, LOGICAL_POOL_COPYSET_NUM)

	go metricomm.QueryInstantMetric(serverName, &results)
	go metricomm.QueryInstantMetric(chunkserverName, &results)
	go metricomm.QueryInstantMetric(copysetName, &results)

	count := 0
	for res := range results {
		if res.Err != nil {
			return &poolItemNum, res.Err
		}
		ret := metricomm.ParseVectorMetric(res.Result.(*metricomm.QueryResponseOfVector), true)
		for _, v := range ret {
			iv, ok := strconv.ParseUint(v["value"], 10, 32)
			if ok != nil {
				return &poolItemNum, ok
			}
			if res.Key.(string) == serverName {
				poolItemNum.ServerNum = uint32(iv)
			} else if res.Key.(string) == chunkserverName {
				poolItemNum.ChunkServerNum = uint32(iv)
			} else {
				poolItemNum.CopysetNum = uint32(iv)
			}
			break
		}
		count += 1
		if count >= requestSize {
			break
		}
	}
	return &poolItemNum, nil
}

func GetPoolPerformance(name string, start, end, interval uint64) ([]metricomm.Performance, error) {
	poolName := metricomm.FormatToMetricName(name)
	prefix := fmt.Sprintf("%s%s_", LOGICAL_POOL_METRIC_PREFIX, poolName)
	return metricomm.GetPerformance(prefix, start, end, interval)
}

func GetClusterPerformance(start, end, interval uint64) ([]metricomm.Performance, error) {
	return metricomm.GetPerformance(CLUSTER_METRIC_PREFIX, start, end, interval)
}

func GetVolumePerformance(volumeName string, start, end, interval uint64) ([]metricomm.UserPerformance, error) {
	name := metricomm.FormatToMetricName(volumeName)
	prefix := fmt.Sprintf("%s%s_", FILE_PREFIX, name)
	return metricomm.GetUserPerformance(prefix, start, end, interval)
}

func GetClusterSpace(start, end, interval uint64) ([]metricomm.SpaceTrend, error) {
	spaces := []metricomm.SpaceTrend{}
	retMap := make(map[float64]*metricomm.SpaceTrend)

	// total, alloc
	requestSize := 2
	results := make(chan metricomm.MetricResult, requestSize)
	totalName := fmt.Sprintf("%s&start=%d&end=%d&step=%d", CLUSTER_LOGICAL_CAPACITY, start, end, interval)
	usedName := fmt.Sprintf("%s&start=%d&end=%d&step=%d", CLUSTER_LOGICAL_ALLOC, start, end, interval)
	go metricomm.QueryRangeMetric(totalName, &results)
	go metricomm.QueryRangeMetric(usedName, &results)

	count := 0
	for res := range results {
		if res.Err != nil {
			return nil, res.Err
		}
		ret := metricomm.ParseMatrixMetric(res.Result.(*metricomm.QueryResponseOfMatrix), metricomm.INSTANCE)
		if res.Key.(string) == totalName {
			for _, v := range ret {
				for _, item := range v {
					total, e := strconv.ParseUint(item.Value, 10, 64)
					if e != nil {
						return nil, e
					}
					if _, ok := retMap[item.Timestamp]; ok {
						retMap[item.Timestamp].Total = total / common.GiB
					} else {
						retMap[item.Timestamp] = &metricomm.SpaceTrend{
							Timestamp: item.Timestamp,
							Total:     total / common.GiB,
						}
					}
				}
			}
		} else if res.Key.(string) == usedName {
			for _, v := range ret {
				for _, item := range v {
					used, e := strconv.ParseUint(item.Value, 10, 64)
					if e != nil {
						return nil, e
					}
					if _, ok := retMap[item.Timestamp]; ok {
						retMap[item.Timestamp].Used = used / common.GiB
					} else {
						retMap[item.Timestamp] = &metricomm.SpaceTrend{
							Timestamp: item.Timestamp,
							Used:      used / common.GiB,
						}
					}
				}
			}
		}
		count += 1
		if count >= requestSize {
			break
		}
	}
	for _, v := range retMap {
		spaces = append(spaces, *v)
	}
	return spaces, nil
}
