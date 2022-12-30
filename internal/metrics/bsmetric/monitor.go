package bsmetric

import (
	"fmt"

	"github.com/opencurve/curve-manager/internal/metrics/core"
)

const (
	ETCD_CLUSTER_VERSION_NAME  = "etcd_cluster_version"
	ETCD_SERVER_IS_LEADER_NAME = "etcd_server_is_leader"
	ETCD_CLUSTER_VERSION       = "cluster_version"

	LOGICAL_POOL_METRIC_PREFIX    = "topology_metric_logical_pool_"
	LOGICAL_POOL_LOGICAL_CAPACITY = "_logical_capacity"
	LOGICAL_POOL_LOGICAL_ALLOC    = "_logical_alloc"
)

type PoolSpace struct {
	Total string
	Used string
}

func parseVectorMetric(info *QueryResponseOfVector, key string, isValue bool) *map[string]string {
	ret := make(map[string]string)
	if info == nil {
		return &ret
	}
	for _, metric := range info.Data.Result {
		if !isValue {
			ret[metric.Metric["instance"]] = metric.Metric[key]
		} else {
			ret[metric.Metric["instance"]] = metric.Value[1].(string)
		}
	}
	return &ret
}

func GetEtcdStatus() (*[]EtcdStatus, error) {
	rquestSize := 2
	results := make(chan metricResult, rquestSize)
	// get version
	go func() {
		var ret QueryResponseOfVector
		err := core.GMetricClient.GetMetricFromPrometheus(
			VECTOR_METRIC_PATH, VECTOR_METRIC_QUERY_KEY, ETCD_CLUSTER_VERSION_NAME, &ret)
		results <- metricResult{
			key:    ETCD_CLUSTER_VERSION_NAME,
			err:    err,
			result: &ret,
		}

	}()

	// get leader
	go func() {
		var ret QueryResponseOfVector
		err := core.GMetricClient.GetMetricFromPrometheus(
			VECTOR_METRIC_PATH, VECTOR_METRIC_QUERY_KEY, ETCD_SERVER_IS_LEADER_NAME, &ret)
		results <- metricResult{
			key:    ETCD_SERVER_IS_LEADER_NAME,
			err:    err,
			result: &ret,
		}
	}()

	var ret []EtcdStatus
	retMap := make(map[string]*EtcdStatus)
	for _, addr := range core.GMetricClient.EtcdAddr {
		retMap[addr] = &EtcdStatus{
			Address: addr,
			Version: "",
			Leader:  false,
			Online:  false,
		}
	}

	count := 0
	for res := range results {
		if res.err != nil {
			return &ret, res.err
		}
		if res.key.(string) == ETCD_CLUSTER_VERSION_NAME {
			versions := parseVectorMetric(res.result.(*QueryResponseOfVector), ETCD_CLUSTER_VERSION, false)
			for k, v := range *versions {
				(*retMap[k]).Version = v
			}
		} else {
			leaders := parseVectorMetric(res.result.(*QueryResponseOfVector), "", true)
			for k, v := range *leaders {
				if v == "1" {
					(*retMap[k]).Leader = true
				} else {
					(*retMap[k]).Leader = false
				}
				(*retMap[k]).Online = true
			}
		}
		count = count + 1
		if count >= rquestSize {
			break
		}
	}
	for _, v := range retMap {
		ret = append(ret, *v)
	}
	return &ret, nil
}

func GetPoolSpace(name string) (*PoolSpace, error) {
	var space PoolSpace
	spaceName := formatToMetricName(name)

	rquestSize := 2
	results := make(chan metricResult, rquestSize)
	// get total space
	totalName := fmt.Sprintf("%s%s%s", LOGICAL_POOL_METRIC_PREFIX, spaceName, LOGICAL_POOL_LOGICAL_CAPACITY)
	go func() {
		var ret QueryResponseOfVector
		err := core.GMetricClient.GetMetricFromPrometheus(
			VECTOR_METRIC_PATH, VECTOR_METRIC_QUERY_KEY, totalName, &ret)
		results <- metricResult{
			key:    totalName,
			err:    err,
			result: &ret,
		}
	}()

	// get used space
	usedName := fmt.Sprintf("%s%s%s", LOGICAL_POOL_METRIC_PREFIX, spaceName, LOGICAL_POOL_LOGICAL_ALLOC)
	go func() {
		var ret QueryResponseOfVector
		err := core.GMetricClient.GetMetricFromPrometheus(
			VECTOR_METRIC_PATH, VECTOR_METRIC_QUERY_KEY, usedName, &ret)
		results <- metricResult{
			key:    usedName,
			err:    err,
			result: &ret,
		}
	}()

	count := 0
	for res := range results {
		if res.err != nil {
			return &space, res.err
		}
		if res.key.(string) == totalName {
			ret := parseVectorMetric(res.result.(*QueryResponseOfVector), "", true)
			for _, v := range *ret {
				space.Total = v
			}
		} else {
			ret := parseVectorMetric(res.result.(*QueryResponseOfVector), "", true)
			for _, v := range *ret {
				space.Used = v
			}
		}
		count = count + 1
		if count >= rquestSize {
			break
		}
	}
	return &space, nil
}
