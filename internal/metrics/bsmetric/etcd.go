package bsmetric

import (
	"sync"

	"github.com/opencurve/curve-manager/internal/metrics/core"
)

func GetEtcdStatus() (*[]EtcdStatus, error) {
	var verResult, leaderResult QueryResponseOfVector
	var verErr, leaderErr error
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		verErr = core.GMetricClient.GetMetricFromPrometheus(
			VECTOR_METRIC_PATH, VECTOR_METRIC_QUERY_KEY, ETCD_CLUSTER_VERSION_NAME, &verResult)
	}()

	go func() {
		defer wg.Done()
		leaderErr = core.GMetricClient.GetMetricFromPrometheus(
			VECTOR_METRIC_PATH, VECTOR_METRIC_QUERY_KEY, ETCD_SERVER_IS_LEADER_NAME, &leaderResult)
	}()

	wg.Wait()
	var ret []EtcdStatus
	retMap := make(map[string]*EtcdStatus)
	for index, addr := range core.GMetricClient.EtcdAddr {
		ret = append(ret, EtcdStatus{
			Address: addr,
			Version: "",
			Leader:  false,
			Online:  false,
		})
		retMap[addr] = &ret[index]
	}

	if verErr != nil {
		return &ret, verErr
	}
	// get etcd version
	versions := parseVectorMetric(&verResult, ETCD_CLUSTER_VERSION, false)
	for k, v := range *versions {
		(*retMap[k]).Version = v
	}

	if leaderErr != nil {
		return &ret, leaderErr
	}
	// get etcd leader
	leaders := parseVectorMetric(&leaderResult, "", true)
	for k, v := range *leaders {
		if v == "1" {
			(*retMap[k]).Leader = true
		} else {
			(*retMap[k]).Leader = false
		}
		(*retMap[k]).Online = true
	}
	return &ret, nil
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
