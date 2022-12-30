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
	"encoding/json"
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

type SnapShotCloneServerStatus struct {
	Address string `json:"address"`
	Version string `json:"version"`
	Leader  bool   `json:"leader"`
	Online  bool   `json:"online"`
}

type metricResult struct {
	key    interface{}
	result interface{}
	err    error
}

type bvarConfMetric struct {
	ConfName  string `json:"conf_name"`
	ConfValue string `json:"conf_value"`
}

func parseBvarMetric(value string) (*map[string]string, error) {
	ret := make(map[string]string)
	lines := strings.Split(value, "\n")
	for _, line := range lines {
		items := strings.Split(line, " : ")
		if len(items) != 2 {
			return nil, fmt.Errorf("parseBvarMetric failed, line: %s", line)
		}
		ret[strings.TrimSpace(items[0])] = strings.Trim(strings.TrimSpace(items[1]), "\"")
	}
	return &ret, nil
}

func getBvarConfMetricValue(metric string) string {
	var conf bvarConfMetric
	err := json.Unmarshal([]byte(metric), &conf)
	if err != nil {
		return ""
	}
	return conf.ConfValue
}

func getBvarMetric(addrs []string, name string, results *chan metricResult) {
	for _, host := range addrs {
		go func(addr string) {
			resp, err := core.GMetricClient.GetMetricFromBvar(addr,
				fmt.Sprintf("%s/%s", BVAR_METRIC_PATH, name))
			*results <- metricResult{
				key:    addr,
				err:    err,
				result: resp,
			}
		}(host)
	}
}

// eg: LogicalPool1 -> logical_pool1
func formatToMetricName(name string) string {
	var target []string
	pos := 0
	for index, ch := range name {
		if ch >= 65 && ch <= 90 && index != 0 {
			target = append(target, strings.ToLower(name[pos:index]))
			pos = index
		}
	}
	target = append(target, strings.ToLower(name[pos:len(name)]))
	return strings.Join(target, "_")
}
