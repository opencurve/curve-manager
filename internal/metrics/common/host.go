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
* Created Date: 2023-02-15
* Author: wanghai (SeanHai)
 */

package common

import (
	"fmt"
	"strconv"
	"time"

	"github.com/opencurve/curve-manager/internal/common"
)

const (
	// node info
	NODE_UNAME_INFO      = "node_uname_info"
	NODE_UNAME_NODE_NAME = "nodename"
	NODE_UNAME_MACHINE   = "machine"
	NODE_UNAME_RELEASE   = "release"
	NODE_UNAME_SYSNAME   = "sysname"
	NODE_UNAME_VERSION   = "version"

	// cpu info
	NODE_CPU_INFO        = "node_cpu_info"
	NODE_CPU_MODLE       = "model_name"
	NODE_CPU_UTILIZATION = "node_cpu_seconds_total"
	NODE_CPU_IDLE        = "idle"

	// memory info
	NODE_MEMORY_TOTAL_BYTES = "node_memory_MemTotal_bytes"

	// network info
	NODE_NETWORK_DEVICE_FILTER        = "tap.*|veth.*|br.*|docker.*|virbr*|lo*"
	NODE_NETWORK_RECEIVE_BYTES_TOTAL  = "node_network_receive_bytes_total"
	NODE_NETWORK_TRANSMIT_BYTES_TOTAL = "node_network_transmit_bytes_total"
)

type HostInfo struct {
	HostName string `json:"hostname"`
	IP       string `json:"ip"`
	Machine  string `json:"machine"`
	Release  string `json:"kernel-release"`
	Version  string `json:"kernel-version"`
	System   string `json:"operating-system"`
}

type CPUInfo struct {
	TotalNum uint32            `json:"totalNum"`
	Models   map[string]uint32 `json:"cpuModles"`
}

type NetworkTraffic struct {
	Receive  map[string][]RangeMetricItem
	Transmit map[string][]RangeMetricItem
}

func GetHostInfo(instance string) (interface{}, error) {
	hostsInfo := make(map[string]HostInfo)
	requestSize := 1
	results := make(chan MetricResult, requestSize)
	metricName := NODE_UNAME_INFO
	if instance != "" {
		metricName = fmt.Sprintf("%s{instance=%q}", metricName, instance)
	}
	QueryInstantMetric(metricName, &results)
	count := 0
	for res := range results {
		if res.Err != nil {
			return nil, res.Err
		}
		ret := ParseVectorMetric(res.Result.(*QueryResponseOfVector), false)
		for k, v := range ret {
			info := HostInfo{}
			info.HostName = v[NODE_UNAME_NODE_NAME]
			info.IP, _ = common.GetIPFromEndpoint(v[INSTANCE])
			info.Machine = v[NODE_UNAME_MACHINE]
			info.Release = v[NODE_UNAME_RELEASE]
			info.System = v[NODE_UNAME_SYSNAME]
			info.Version = v[NODE_UNAME_VERSION]
			hostsInfo[k] = info
		}
		count += 1
		if count >= requestSize {
			break
		}
	}
	return hostsInfo, nil
}

func GetHostCPUInfo(instance string) (interface{}, error) {
	cpuInfoMap := make(map[string]CPUInfo)
	requestSize := 1
	results := make(chan MetricResult, requestSize)
	metricName := NODE_CPU_INFO
	if instance != "" {
		metricName = fmt.Sprintf("%s{instance=%q}", metricName, instance)
	}
	QueryInstantMetric(metricName, &results)
	count := 0
	for res := range results {
		if res.Err != nil {
			return nil, res.Err
		}
		info := res.Result.(*QueryResponseOfVector)
		for _, item := range info.Data.Result {
			instance := item.Metric[INSTANCE]
			model := item.Metric[NODE_CPU_MODLE]
			if cpu, ok := cpuInfoMap[instance]; ok {
				cpu.TotalNum++
				cpu.Models[model]++
				cpuInfoMap[instance] = cpu
			} else {
				cpu := CPUInfo{}
				cpu.TotalNum = 1
				cpu.Models = make(map[string]uint32)
				cpu.Models[model] = 1
				cpuInfoMap[instance] = cpu
			}
		}
		count += 1
		if count >= requestSize {
			break
		}
	}
	return cpuInfoMap, nil
}

func GetHostMemoryInfo(instance string) (interface{}, error) {
	memoryInfo := make(map[string]uint64)
	requestSize := 1
	results := make(chan MetricResult, requestSize)
	metricName := NODE_MEMORY_TOTAL_BYTES
	if instance != "" {
		metricName = fmt.Sprintf("%s{instance=%q}", metricName, instance)
	}
	QueryInstantMetric(metricName, &results)
	count := 0
	for res := range results {
		if res.Err != nil {
			return nil, res.Err
		}
		ret := ParseVectorMetric(res.Result.(*QueryResponseOfVector), true)
		for k, v := range ret {
			memoryInfo[k], _ = strconv.ParseUint(v["value"], 10, 64)
			memoryInfo[k] /= common.GiB
		}
		count += 1
		if count >= requestSize {
			break
		}
	}
	return memoryInfo, nil
}

func GetHostDiskNum(instance string) (interface{}, error) {
	diskNum := make(map[string]uint32)
	requestSize := 1
	results := make(chan MetricResult, requestSize)
	QueryInstantMetric(NODE_DISK_NUMBER, &results)
	count := 0
	for res := range results {
		if res.Err != nil {
			return nil, res.Err
		}
		ret := ParseVectorMetric(res.Result.(*QueryResponseOfVector), true)
		for k, v := range ret {
			num, err := strconv.ParseUint(v["value"], 10, 32)
			if err == nil {
				diskNum[k] = uint32(num)
			} else {
				return nil, err
			}
		}
		count += 1
		if count >= requestSize {
			break
		}
	}
	return diskNum, nil
}

/*
* return: reveive, transmit, error
* map[string][]RangeMetricItem: key: network device, value: performance in different timestamp
 */
func GetNetWorkTraffic(instance string) (interface{}, error) {
	networkTraffic := NetworkTraffic{}

	// receive, transmit
	requestSize := 2
	results := make(chan MetricResult, requestSize)
	end := time.Now().Unix()
	start := end - DEFAULT_RANGE
	receiveName := fmt.Sprintf("%s&start=%d&end=%d&step=%ds",
		GetNodeNetWorkReveiveName(NODE_NETWORK_RECEIVE_BYTES_TOTAL, instance), start, end, DEFAULT_STEP)
	transmitName := fmt.Sprintf("%s&start=%d&end=%d&step=%ds",
		GetNodeNetWorkReveiveName(NODE_NETWORK_TRANSMIT_BYTES_TOTAL, instance), start, end, DEFAULT_STEP)

	go QueryRangeMetric(receiveName, &results)
	go QueryRangeMetric(transmitName, &results)

	count := 0
	for res := range results {
		if res.Err != nil {
			return networkTraffic, res.Err
		}
		ret := ParseMatrixMetric(res.Result.(*QueryResponseOfMatrix), DEVICE)
		switch res.Key.(string) {
		case receiveName:
			networkTraffic.Receive = ret
		case transmitName:
			networkTraffic.Transmit = ret
		}
		count += 1
		if count >= requestSize {
			break
		}
	}
	return networkTraffic, nil
}

func GetHostCPUUtilization(instance string) (interface{}, error) {
	return GetUtilization(GetNodeCPUUtilizationName(instance))
}

func GetHostMemUtilization(instance string) (interface{}, error) {
	return GetUtilization(GetNodeMemUtilizationName(instance))
}
