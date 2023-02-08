package common

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/metrics/core"
)

const (
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

	DEFAULT_RANGE = 180
	DEFAULT_STEP  = 15
	WRITE_IOPS    = "write_iops"
	WRITE_RATE    = "write_rate"
	READ_IOPS     = "read_iops"
	READ_REAT     = "read_rate"
)

type MetricResult common.QueryResult

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
	Timestamp float64
	value     string
}

type Performance struct {
	Timestamp float64 `json:"timestamp" binding:"required"`
	WriteIOPS string  `json:"writeIOPS" binding:"required"`
	WriteBPS  string  `json:"writeBPS" binding:"required"`
	ReadIOPS  string  `json:"readIOPS" binding:"required"`
	ReadBPS   string  `json:"readBPS" binding:"required"`
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
func ParseRaftStatusMetric(value string) ([]map[string]string, error) {
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
				tmap[common.RAFT_STATUS_KEY_GROUPID] = line[start:end]
				continue
			}
			hit := strings.Count(line, ": ")
			if hit == 1 {
				c := strings.Split(line, ": ")
				if len(c) != 2 {
					return nil, fmt.Errorf(fmt.Sprintf("format error: %s", line))
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
					return nil, fmt.Errorf(fmt.Sprintf("format error: %s", line))
				}
				for _, i := range v {
					j := strings.Split(i, ": ")
					if len(j) != 2 {
						return nil, fmt.Errorf(fmt.Sprintf("format error: %s", i))
					}
					tmap[j[0]] = j[1]
				}
			} else if strings.Contains(line, common.RAFT_STATUS_KEY_STORAGE) {
				storageItems := strings.Split(line, "\n")
				for _, sitem := range storageItems {
					sitemArr := strings.Split(sitem, ": ")
					if len(sitemArr) != 2 {
						return nil, fmt.Errorf(fmt.Sprintf("format error: %s", sitem))
					}
					tmap[sitemArr[0]] = sitemArr[1]
				} 
			} else {
				return nil, fmt.Errorf(fmt.Sprintf("format error: %s", line))
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

// eg: LogicalPool1 -> logical_pool1
func FormatToMetricName(name string) string {
	var target []string
	pos := 0
	for index, ch := range name {
		if ch >= 65 && ch <= 90 && index != 0 {
			target = append(target, strings.ToLower(name[pos:index]))
			pos = index
		}
	}
	target = append(target, strings.ToLower(name[pos:]))
	return strings.Join(target, "_")
}

func ParseVectorMetric(info *QueryResponseOfVector, key string, isValue bool) *map[string]string {
	ret := make(map[string]string)
	if info == nil {
		return &ret
	}
	for _, item := range info.Data.Result {
		if !isValue {
			ret[item.Metric["instance"]] = item.Metric[key]
		} else {
			ret[item.Metric["instance"]] = item.Value[1].(string)
		}
	}
	return &ret
}

func ParseMatrixMetric(info *QueryResponseOfMatrix) *map[string][]RangeMetricItem {
	ret := make(map[string][]RangeMetricItem)
	if info == nil {
		return &ret
	}
	for _, item := range info.Data.Result {
		for _, slice := range item.Values {
			var rangeItem RangeMetricItem
			rangeItem.Timestamp = slice[0].(float64)
			rangeItem.value = slice[1].(string)
			ret[item.Metric["instance"]] = append(ret[item.Metric["instance"]], rangeItem)
		}
	}
	return &ret
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

func GetPerformance(name string) ([]Performance, error) {
	performance := []Performance{}
	retMap := make(map[float64]*Performance)

	// writeIOPS, writeBPS, readIOPS, readBPS
	rquestSize := 4
	results := make(chan MetricResult, rquestSize)
	end := time.Now().Unix()
	start := end - DEFAULT_RANGE
	writeIOPSName := fmt.Sprintf("%s%s&start=%d&end=%d&step=%ds", name, WRITE_IOPS, start, end, DEFAULT_STEP)
	writeBPSName := fmt.Sprintf("%s%s&start=%d&end=%d&step=%ds", name, WRITE_RATE, start, end, DEFAULT_STEP)
	readIOPSName := fmt.Sprintf("%s%s&start=%d&end=%d&step=%ds", name, READ_IOPS, start, end, DEFAULT_STEP)
	readBPSName := fmt.Sprintf("%s%s&start=%d&end=%d&step=%ds", name, READ_REAT, start, end, DEFAULT_STEP)

	go QueryRangeMetric(writeIOPSName, &results)
	go QueryRangeMetric(writeBPSName, &results)
	go QueryRangeMetric(readIOPSName, &results)
	go QueryRangeMetric(readBPSName, &results)

	count := 0
	for res := range results {
		if res.Err != nil {
			return nil, res.Err
		}
		ret := ParseMatrixMetric(res.Result.(*QueryResponseOfMatrix))
		for _, v := range *ret {
			for _, data := range v {
				if p, ok := retMap[data.Timestamp]; ok {
					switch res.Key.(string) {
					case writeIOPSName:
						p.WriteIOPS = data.value
					case writeBPSName:
						p.WriteBPS = data.value
					case readIOPSName:
						p.ReadIOPS = data.value
					case readBPSName:
						p.ReadBPS = data.value
					}
				} else {
					switch res.Key.(string) {
					case writeIOPSName:
						retMap[data.Timestamp] = &Performance{
							Timestamp: data.Timestamp,
							WriteIOPS: data.value,
						}
					case writeBPSName:
						retMap[data.Timestamp] = &Performance{
							Timestamp: data.Timestamp,
							WriteBPS:  data.value,
						}
					case readIOPSName:
						retMap[data.Timestamp] = &Performance{
							Timestamp: data.Timestamp,
							ReadIOPS:  data.value,
						}
					case readBPSName:
						retMap[data.Timestamp] = &Performance{
							Timestamp: data.Timestamp,
							ReadBPS:   data.value,
						}
					}
				}
			}
			break
		}
		count += 1
		if count >= rquestSize {
			break
		}
	}

	for _, v := range retMap {
		performance = append(performance, *v)
	}
	return performance, nil
}
