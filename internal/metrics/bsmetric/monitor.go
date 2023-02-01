package bsmetric

import (
	"fmt"
	"strconv"

	comm "github.com/opencurve/curve-manager/internal/metrics/common"
	"github.com/opencurve/curve-manager/internal/metrics/core"
)

const (
	CLUSTER_METRIC_PREFIX = "topology_metric_cluster_"

	LOGICAL_POOL_METRIC_PREFIX    = "topology_metric_logical_pool_"
	LOGICAL_POOL_LOGICAL_CAPACITY = "_logical_capacity"
	LOGICAL_POOL_LOGICAL_ALLOC    = "_logical_alloc"
	LOGICAL_POOL_SERVER_NUM       = "_server_num"
	LOGICAL_POOL_CHUNKSERVER_NUM  = "_chunkserver_num"
	LOGICAL_POOL_COPYSET_NUM      = "_copyset_num"
)

type Space struct {
	Total uint64
	Used  uint64
}

type PoolItemNum struct {
	ServerNum      uint32
	ChunkServerNum uint32
	CopysetNum     uint32
}

func GetEtcdStatus() ([]ServiceStatus, string) {
	// init value
	var ret []ServiceStatus
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
	rquestSize := 2
	results := make(chan comm.MetricResult, rquestSize)
	go comm.QueryInstantMetric(comm.ETCD_CLUSTER_VERSION_NAME, &results)
	go comm.QueryInstantMetric(comm.ETCD_SERVER_IS_LEADER_NAME, &results)

	count := 0
	for res := range results {
		if res.Err != nil {
			return ret, res.Err.Error()
		}
		if res.Key.(string) == comm.ETCD_CLUSTER_VERSION_NAME {
			versions := comm.ParseVectorMetric(res.Result.(*comm.QueryResponseOfVector), comm.ETCD_CLUSTER_VERSION, false)
			for k, v := range *versions {
				(*retMap[k]).Version = v
			}
		} else {
			leaders := comm.ParseVectorMetric(res.Result.(*comm.QueryResponseOfVector), "", true)
			for k, v := range *leaders {
				if v == "1" {
					(*retMap[k]).Leader = true
				} else {
					(*retMap[k]).Leader = false
				}
				(*retMap[k]).Online = true
			}
		}
		count += 1
		if count >= rquestSize {
			break
		}
	}
	for _, v := range retMap {
		ret = append(ret, *v)
	}
	return ret, ""
}

func GetPoolSpace(name string) (*Space, error) {
	var space Space
	poolName := comm.FormatToMetricName(name)

	// total, alloc
	rquestSize := 2
	results := make(chan comm.MetricResult, rquestSize)
	totalName := fmt.Sprintf("%s%s%s", LOGICAL_POOL_METRIC_PREFIX, poolName, LOGICAL_POOL_LOGICAL_CAPACITY)
	usedName := fmt.Sprintf("%s%s%s", LOGICAL_POOL_METRIC_PREFIX, poolName, LOGICAL_POOL_LOGICAL_ALLOC)

	go comm.QueryInstantMetric(totalName, &results)
	go comm.QueryInstantMetric(usedName, &results)

	count := 0
	for res := range results {
		if res.Err != nil {
			return &space, res.Err
		}
		ret := comm.ParseVectorMetric(res.Result.(*comm.QueryResponseOfVector), "", true)
		if res.Key.(string) == totalName {
			for _, v := range *ret {
				total, ok := strconv.ParseUint(v, 10, 64)
				if ok != nil {
					return nil, ok
				}
				space.Total = total / comm.GB
				break
			}
		} else {
			for _, v := range *ret {
				used, ok := strconv.ParseUint(v, 10, 64)
				if ok != nil {
					return nil, ok
				}
				space.Used = used / comm.GB
				break
			}
		}
		count += 1
		if count >= rquestSize {
			break
		}
	}
	return &space, nil
}

func GetPoolItemNum(name string) (*PoolItemNum, error) {
	var poolItemNum PoolItemNum
	poolName := comm.FormatToMetricName(name)

	// serverNUm, chunkserverNum, copysetNum
	rquestSize := 3
	results := make(chan comm.MetricResult, rquestSize)
	serverName := fmt.Sprintf("%s%s%s", LOGICAL_POOL_METRIC_PREFIX, poolName, LOGICAL_POOL_SERVER_NUM)
	chunkserverName := fmt.Sprintf("%s%s%s", LOGICAL_POOL_METRIC_PREFIX, poolName, LOGICAL_POOL_CHUNKSERVER_NUM)
	copysetName := fmt.Sprintf("%s%s%s", LOGICAL_POOL_METRIC_PREFIX, poolName, LOGICAL_POOL_COPYSET_NUM)

	go comm.QueryInstantMetric(serverName, &results)
	go comm.QueryInstantMetric(chunkserverName, &results)
	go comm.QueryInstantMetric(copysetName, &results)

	count := 0
	for res := range results {
		if res.Err != nil {
			return &poolItemNum, res.Err
		}
		ret := comm.ParseVectorMetric(res.Result.(*comm.QueryResponseOfVector), "", true)
		for _, v := range *ret {
			iv, ok := strconv.ParseUint(v, 10, 32)
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
		if count >= rquestSize {
			break
		}
	}
	return &poolItemNum, nil
}

func GetPoolPerformance(name string) ([]comm.Performance, error) {
	poolName := comm.FormatToMetricName(name)
	prefix := fmt.Sprintf("%s%s_", LOGICAL_POOL_METRIC_PREFIX, poolName)
	return comm.GetPerformance(prefix)
}

func GetClusterPerformance() ([]comm.Performance, error) {
	return comm.GetPerformance(CLUSTER_METRIC_PREFIX)
}
