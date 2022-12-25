/*
prometheus http api resonse data struct

{
  "status": "success" | "error",
  "data": <data>,

  // Only set if status is "error". The data field may still hold additional data.
  "errorType": "<string>",
  "error": "<string>"
}
*/

package bsmetric

import (
	"fmt"
	"strings"

	"github.com/opencurve/curve-manager/internal/metrics/core"
)

const (
	VECTOR_METRIC_PATH      = "/api/v1/query"
	VECTOR_METRIC_QUERY_KEY = "query"

	BVAR_METRIC_PATH = "/vars"

	// service version
	CURVEBS_VERSION = "curve_version"

	// etcd
	ETCD_CLUSTER_VERSION_NAME  = "etcd_cluster_version"
	ETCD_SERVER_IS_LEADER_NAME = "etcd_server_is_leader"

	ETCD_CLUSTER_VERSION = "cluster_version"

	//mds
	MDS_STATUS   = "mds_status"
	MDS_LEADER   = "leader"
	MDS_FOLLOWER = "follower"
)

// data struct: https://prometheus.io/docs/prometheus/latest/querying/api/#expression-query-result-formats
type QueryResponseOfVector struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

type EtcdStatus struct {
	Address string `json:"address"`
	Version string `json:"version"`
	Leader  bool   `json:"leader"`
	Online  bool   `json:"online"`
}

type MdsStatus struct {
	Address string `json:"address"`
	Version string `json:"version"`
	Leader  bool   `json:"leader"`
	Online  bool   `json:"online"`
}

type metricResult struct {
	addr   string
	result interface{}
	err    error
}

func parseBvarMetric(value string) (*map[string]string, error) {
	ret := make(map[string]string)
	lines := strings.Split(value, "\n")
	for _, line := range lines {
		items := strings.Split(line, ":")
		if len(items) != 2 {
			return nil, fmt.Errorf("parseBvarMetric failed, line: %s", line)
		}
		ret[strings.TrimSpace(items[0])] = strings.Replace(strings.TrimSpace(items[1]), "\"", "", -1)
	}
	return &ret, nil
}

func getBvarMetric(addrs []string, name string, results *chan metricResult) {
	for _, host := range addrs {
		go func(addr string) {
			resp, err := core.GMetricClient.GetMetricFromBvar(addr,
				fmt.Sprintf("%s/%s", BVAR_METRIC_PATH, name))
			*results <- metricResult{
				addr:   addr,
				err:    err,
				result: resp,
			}
		}(host)
	}
}
