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

package common

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/metrics/core"
)

const (
	INSTANCE   = "instance"
	DEVICE     = "device"
	FSTYPE     = "fstype"
	MOUNTPOINT = "mountpoint"

	INSTANT_METRIC_PATH     = "/api/v1/query"
	VECTOR_METRIC_QUERY_KEY = "query"

	RANGE_METRIC_PATH = "/api/v1/query_range"
	RANGE_START       = "start"
	RANGE_END         = "end"
	RANGE_STEP        = "step"

	BVAR_METRIC_PATH = "/vars"
	RAFT_STATUS_PATH = "/raft_stat"

	CONF_VALUE = "conf_value"

	CURVEBS_VERSION = "curve_version"

	ETCD_CLUSTER_VERSION_NAME  = "etcd_cluster_version"
	ETCD_SERVER_IS_LEADER_NAME = "etcd_server_is_leader"
	ETCD_CLUSTER_VERSION       = "cluster_version"

	WRITE_IOPS = "write_iops"
	WRITE_RATE = "write_rate"
	READ_IOPS  = "read_iops"
	READ_REAT  = "read_rate"

	WRITE_QPS = "write_qps"
	WRITE_BPS = "write_bps"
	READ_QPS  = "read_qps"
	READ_BPS  = "read_bps"
)

type MetricResult common.QueryResult

type Space struct {
	Total uint64
	Used  uint64
}

type SpaceTrend struct {
	Timestamp float64 `json:"timestamp"`
	Total     uint64  `json:"total"`
	Used      uint64  `json:"alloc"`
}

type FileSystemInfo struct {
	Device     string
	FsType     string
	MountPoint string
	SpaceTotal uint64
	SpaceAvail uint64
}

/*
prometheus http api resonse data struct

{
  "status": "success" | "error",
  "data": <data>,

  // Only set if status is "error". The data field may still hold additional data.
  "errorType": "<string>",
  "error": "<string>"
}

data struct: https://prometheus.io/docs/prometheus/latest/querying/api/#expression-query-result-formats
*/
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

type QueryResponseOfMatrix struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Values [][]interface{}   `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

type ServiceStatus struct {
	Address string `json:"address"`
	Version string `json:"version"`
	Leader  bool   `json:"leader"`
	Online  bool   `json:"online"`
}

type bvarConfMetric struct {
	ConfName  string `json:"conf_name"`
	ConfValue string `json:"conf_value"`
}

type RangeMetricItem struct {
	Timestamp float64 `json:"timestamp"`
	Value     string  `json:"value"`
}

type Performance struct {
	Timestamp float64 `json:"timestamp" binding:"required"`
	WriteIOPS string  `json:"writeIOPS" binding:"required"`
	WriteBPS  string  `json:"writeBPS" binding:"required"`
	ReadIOPS  string  `json:"readIOPS" binding:"required"`
	ReadBPS   string  `json:"readBPS" binding:"required"`
}

type UserPerformance struct {
	Timestamp float64 `json:"timestamp" binding:"required"`
	WriteQPS  string  `json:"writeQPS" binding:"required"`
	WriteBPS  string  `json:"writeBPS" binding:"required"`
	ReadQPS   string  `json:"readQPS" binding:"required"`
	ReadBPS   string  `json:"readBPS" binding:"required"`
}

func GetNodeCPUUtilizationName(instance string, interval uint64) string {
	return fmt.Sprintf("(sum+by(instance)+(irate(node_cpu_seconds_total{instance=%q,mode!=%q}[%ds]))",
		instance, NODE_CPU_IDLE, interval) +
		fmt.Sprintf("/on(instance)+group_left+sum+by+(instance)((irate(node_cpu_seconds_total[%ds]))))*100",
			interval)
}

func GetNodeMemUtilizationName(instance string) string {
	return fmt.Sprintf("100-((node_memory_MemAvailable_bytes{instance=%q}*100)/node_memory_MemTotal_bytes)",
		instance)
}

func GetNodeDiskPerformanceName(typeName, instance string, interval uint64) string {
	return fmt.Sprintf("irate(%s{instance=%q}[%ds])", typeName, instance, interval)
}

func GetNodeNetWorkReveiveName(typeName, instance string, interval uint64) string {
	return fmt.Sprintf("irate(%s{instance=%q,device!~%q}[%ds])",
		typeName, instance, NODE_NETWORK_DEVICE_FILTER, interval)
}

func ParseBvarMetric(value string) (*map[string]string, error) {
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

func GetBvarConfMetricValue(metric string) string {
	var conf bvarConfMetric
	err := json.Unmarshal([]byte(metric), &conf)
	if err != nil {
		return ""
	}
	return conf.ConfValue
}

func GetBvarMetric(addrs []string, name string, results *chan MetricResult) {
	for _, host := range addrs {
		go func(addr string) {
			resp, err := core.GMetricClient.GetMetricFromService(addr,
				fmt.Sprintf("%s/%s", BVAR_METRIC_PATH, name))
			*results <- MetricResult{
				Key:    addr,
				Err:    err,
				Result: resp,
			}
		}(host)
	}
}

/*
[8589934645]
peer_id: 10.166.24.22:8200:0\r\n
state: LEADER\r\n
readonly: 0\r\n
term: 19\r\n
conf_index: 11244429\r\n
peers: 10.166.24.22:8200:0 10.166.24.27:8218:0 10.166.24.29:8206:0\r\n
changing_conf: NO    stage: STAGE_NONE\r\n
election_timer: timeout(1000ms) STOPPED\r\n
vote_timer: timeout(1000ms) STOPPED\r\n
stepdown_timer: timeout(1000ms) SCHEDULING(in 577ms)\r\n
snapshot_timer: timeout(1800000ms) SCHEDULING(in 277280ms)\r\n
storage: [11243647, 11245778]\n
disk_index: 11245778\n
known_applied_index: 11245778\n
last_log_id: (index=11245778,term=19)\n
state_machine: Idle\n
last_committed_index: 11245778r\n
last_snapshot_index: 11244429
last_snapshot_term: 19
snapshot_status: IDLE
replicator_82304458296217@10.166.24.27:8218:0: next_index=11245779  flying_append_entries_size=0 idle hc=17738777 ac=206514 ic=0
replicator_80702435493905@10.166.24.29:8206:0: next_index=11245779  flying_append_entries_size=0 idle hc=17738818 ac=206282 ic=0
\r\n\r\n
[8589934712]
peer_id: 10.166.24.22:8200:0
state: FOLLOWER
readonly: 0
term: 16
conf_index: 15368827
peers: 10.166.24.22:8200:0 10.166.24.29:8212:0 10.166.24.30:8219:0
leader: 10.166.24.29:8212:0
last_msg_to_now: 48
election_timer: timeout(1000ms) SCHEDULING(in 719ms)
vote_timer: timeout(1000ms) STOPPED
stepdown_timer: timeout(1000ms) STOPPED
snapshot_timer: timeout(1800000ms) SCHEDULING(in 422640ms)
storage: [15367732, 15370070]
disk_index: 15370070
known_applied_index: 15370070
last_log_id: (index=15370070,term=16)
state_machine: Idle
last_committed_index: 15370070
last_snapshot_index: 15368827
last_snapshot_term: 16
snapshot_status: IDLE
*/
func ParseRaftStatusMetric(addr string, value string) ([]map[string]string, error) {
	var ret []map[string]string
	items := strings.Split(value, "\r\n\r\n")
	for _, item := range items {
		tmap := make(map[string]string)
		lines := strings.Split(item, "\r\n")
		raplicatorIndex := 0
		for index, line := range lines {
			if index == 0 {
				start := strings.Index(line, "[") + 1
				end := strings.Index(line, "]")
				if start >= end {
					return nil, fmt.Errorf(fmt.Sprintf("format error1: %s, %s", line, addr))
				}
				tmap[common.RAFT_STATUS_KEY_GROUPID] = line[start:end]
				continue
			}
			hit := strings.Count(line, ": ")
			if hit == 1 {
				c := strings.Split(line, ": ")
				if len(c) != 2 {
					return nil, fmt.Errorf(fmt.Sprintf("format error2: %s, %s", line, addr))
				}
				if strings.Contains(c[0], common.RAFT_STATUS_KEY_REPLICATOR) {
					c[0] = fmt.Sprintf("%s%d", common.RAFT_STATUS_KEY_REPLICATOR, raplicatorIndex)
					raplicatorIndex += 1
				}
				tmap[c[0]] = c[1]
			} else if hit == 2 {
				// line: [changing_conf: NO    stage: STAGE_NONE]
				v := strings.Split(line, "    ")
				if len(v) != 2 {
					return nil, fmt.Errorf(fmt.Sprintf("format error3: %s, %s", line, addr))
				}
				for _, i := range v {
					j := strings.Split(i, ": ")
					if len(j) != 2 {
						return nil, fmt.Errorf(fmt.Sprintf("format error4: %s, %s", i, addr))
					}
					tmap[j[0]] = j[1]
				}
			} else if strings.Contains(line, common.RAFT_STATUS_KEY_STORAGE) {
				storageItems := strings.Split(line, "\n")
				for _, sitem := range storageItems {
					sitemArr := strings.Split(sitem, ": ")
					if len(sitemArr) != 2 {
						return nil, fmt.Errorf(fmt.Sprintf("format error5: %s, %s", sitem, addr))
					}
					tmap[sitemArr[0]] = sitemArr[1]
				}
			}
		}
		ret = append(ret, tmap)
	}
	return ret, nil
}

func GetRaftStatusMetric(addrs []string, results *chan MetricResult) {
	for _, host := range addrs {
		go func(addr string) {
			resp, err := core.GMetricClient.GetMetricFromService(addr, RAFT_STATUS_PATH)
			*results <- MetricResult{
				Key:    addr,
				Err:    err,
				Result: resp,
			}
		}(host)
	}
}

// bvar related: https://github.com/apache/brpc/blob/84ac073a6c2c6b15d6d754686255c26acc917bfd/src/bvar/variable.cpp#L946
func FormatToMetricName(name string) string {
	var target []string
	pos := 0
	for index, ch := range name {
		if ch >= 'A' && ch <= 'Z' {
			if name[pos:index] != "" {
				target = append(target, strings.ToLower(name[pos:index]))
			}
			pos = index
		} else if !(ch >= '0' && ch <= '9') && !(ch >= 'a' && ch <= 'z') {
			if name[pos:index] != "" {
				target = append(target, strings.ToLower(name[pos:index]))
			}
			pos = index + 1
		}
	}
	if name[pos:] != "" {
		target = append(target, strings.ToLower(name[pos:]))
	}
	return strings.Join(target, "_")
}

// @return map[string]map[string]string, key: instance, value.key: meticKey or "value"
func ParseVectorMetric(info *QueryResponseOfVector, isValue bool) map[string]map[string]string {
	ret := make(map[string]map[string]string)
	if info == nil {
		return ret
	}
	for _, item := range info.Data.Result {
		if !isValue {
			ret[item.Metric[INSTANCE]] = item.Metric
		} else {
			tmap := make(map[string]string)
			tmap["value"] = item.Value[1].(string)
			ret[item.Metric[INSTANCE]] = tmap
		}
	}
	return ret
}

func ParseMatrixMetric(info *QueryResponseOfMatrix, key string) map[string][]RangeMetricItem {
	ret := make(map[string][]RangeMetricItem)
	if info == nil {
		return ret
	}
	for _, item := range info.Data.Result {
		for _, slice := range item.Values {
			var rangeItem RangeMetricItem
			rangeItem.Timestamp = slice[0].(float64)
			rangeItem.Value = slice[1].(string)
			ret[item.Metric[key]] = append(ret[item.Metric[key]], rangeItem)
		}
	}
	return ret
}

func QueryInstantMetric(name string, results *chan MetricResult) {
	var res QueryResponseOfVector
	err := core.GMetricClient.GetMetricFromPrometheus(
		INSTANT_METRIC_PATH, VECTOR_METRIC_QUERY_KEY, name, &res)
	*results <- MetricResult{
		Key:    name,
		Err:    err,
		Result: &res,
	}
}

func QueryRangeMetric(name string, results *chan MetricResult) {
	var res QueryResponseOfMatrix
	err := core.GMetricClient.GetMetricFromPrometheus(
		RANGE_METRIC_PATH, VECTOR_METRIC_QUERY_KEY, name, &res)
	*results <- MetricResult{
		Key:    name,
		Err:    err,
		Result: &res,
	}
}

func GetUtilization(name string, start, end, interval uint64) (map[string][]RangeMetricItem, error) {
	metricName := fmt.Sprintf("%s&start=%d&end=%d&step=%ds", name, start, end, interval)
	utilization := make(map[string][]RangeMetricItem)
	requestSize := 1
	results := make(chan MetricResult, requestSize)
	QueryRangeMetric(metricName, &results)
	count := 0
	for res := range results {
		if res.Err != nil {
			return nil, res.Err
		}
		utilization = ParseMatrixMetric(res.Result.(*QueryResponseOfMatrix), INSTANCE)
		count += 1
		if count >= requestSize {
			break
		}
	}
	return utilization, nil
}

func GetPerformance(name string, start, end, interval uint64) ([]Performance, error) {
	performance := []Performance{}
	retMap := make(map[float64]*Performance)

	// writeIOPS, writeBPS, readIOPS, readBPS
	requestSize := 4
	results := make(chan MetricResult, requestSize)
	writeIOPSName := fmt.Sprintf("%s%s&start=%d&end=%d&step=%ds", name, WRITE_IOPS, start, end, interval)
	writeBPSName := fmt.Sprintf("%s%s&start=%d&end=%d&step=%ds", name, WRITE_RATE, start, end, interval)
	readIOPSName := fmt.Sprintf("%s%s&start=%d&end=%d&step=%ds", name, READ_IOPS, start, end, interval)
	readBPSName := fmt.Sprintf("%s%s&start=%d&end=%d&step=%ds", name, READ_REAT, start, end, interval)

	go QueryRangeMetric(writeIOPSName, &results)
	go QueryRangeMetric(writeBPSName, &results)
	go QueryRangeMetric(readIOPSName, &results)
	go QueryRangeMetric(readBPSName, &results)

	count := 0
	for res := range results {
		if res.Err != nil {
			return nil, res.Err
		}
		ret := ParseMatrixMetric(res.Result.(*QueryResponseOfMatrix), INSTANCE)
		for _, v := range ret {
			for _, data := range v {
				if p, ok := retMap[data.Timestamp]; ok {
					switch res.Key.(string) {
					case writeIOPSName:
						p.WriteIOPS = data.Value
					case writeBPSName:
						p.WriteBPS = data.Value
					case readIOPSName:
						p.ReadIOPS = data.Value
					case readBPSName:
						p.ReadBPS = data.Value
					}
				} else {
					switch res.Key.(string) {
					case writeIOPSName:
						retMap[data.Timestamp] = &Performance{
							Timestamp: data.Timestamp,
							WriteIOPS: data.Value,
						}
					case writeBPSName:
						retMap[data.Timestamp] = &Performance{
							Timestamp: data.Timestamp,
							WriteBPS:  data.Value,
						}
					case readIOPSName:
						retMap[data.Timestamp] = &Performance{
							Timestamp: data.Timestamp,
							ReadIOPS:  data.Value,
						}
					case readBPSName:
						retMap[data.Timestamp] = &Performance{
							Timestamp: data.Timestamp,
							ReadBPS:   data.Value,
						}
					}
				}
			}
			break
		}
		count += 1
		if count >= requestSize {
			break
		}
	}

	for _, v := range retMap {
		performance = append(performance, *v)
	}
	return performance, nil
}

func GetUserPerformance(name string, start, end, interval uint64) ([]UserPerformance, error) {
	performance := []UserPerformance{}
	retMap := make(map[float64]*UserPerformance)

	// writeQPS, writeBPS, readQPS, readBPS
	requestSize := 4
	results := make(chan MetricResult, requestSize)
	writeQPSName := fmt.Sprintf("%s%s&start=%d&end=%d&step=%ds", name, WRITE_QPS, start, end, interval)
	writeBPSName := fmt.Sprintf("%s%s&start=%d&end=%d&step=%ds", name, WRITE_BPS, start, end, interval)
	readQPSName := fmt.Sprintf("%s%s&start=%d&end=%d&step=%ds", name, READ_QPS, start, end, interval)
	readBPSName := fmt.Sprintf("%s%s&start=%d&end=%d&step=%ds", name, READ_BPS, start, end, interval)

	go QueryRangeMetric(writeQPSName, &results)
	go QueryRangeMetric(writeBPSName, &results)
	go QueryRangeMetric(readQPSName, &results)
	go QueryRangeMetric(readBPSName, &results)

	count := 0
	for res := range results {
		if res.Err != nil {
			return nil, res.Err
		}
		ret := ParseMatrixMetric(res.Result.(*QueryResponseOfMatrix), INSTANCE)
		for _, v := range ret {
			for _, data := range v {
				if p, ok := retMap[data.Timestamp]; ok {
					switch res.Key.(string) {
					case writeQPSName:
						p.WriteQPS = data.Value
					case writeBPSName:
						p.WriteBPS = data.Value
					case readQPSName:
						p.ReadQPS = data.Value
					case readBPSName:
						p.ReadBPS = data.Value
					}
				} else {
					switch res.Key.(string) {
					case writeQPSName:
						retMap[data.Timestamp] = &UserPerformance{
							Timestamp: data.Timestamp,
							WriteQPS:  data.Value,
						}
					case writeBPSName:
						retMap[data.Timestamp] = &UserPerformance{
							Timestamp: data.Timestamp,
							WriteBPS:  data.Value,
						}
					case readQPSName:
						retMap[data.Timestamp] = &UserPerformance{
							Timestamp: data.Timestamp,
							ReadQPS:   data.Value,
						}
					case readBPSName:
						retMap[data.Timestamp] = &UserPerformance{
							Timestamp: data.Timestamp,
							ReadBPS:   data.Value,
						}
					}
				}
			}
			break
		}
		count += 1
		if count >= requestSize {
			break
		}
	}

	for _, v := range retMap {
		performance = append(performance, *v)
	}
	return performance, nil
}
