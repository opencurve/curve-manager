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

package agent

import (
	"fmt"
	"sort"

	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/errno"
	metricomm "github.com/opencurve/curve-manager/internal/metrics/common"
	"github.com/opencurve/pigeon"
)

const (
	// request
	GET_HOST_INFO             = "GetHostInfo"
	GET_HOST_CPU_INFO         = "GetHostCPUInfo"
	GET_HOST_MEM_INFO         = "GetHostMemoryInfo"
	GET_HOST_DISK_NUM         = "GetHostDiskNum"
	GET_HOST_CPU_UTILIZATION  = "GetHostCPUUtilization"
	GET_HOST_MEM_UTILIZATION  = "GetHostMemUtilization"
	GET_HOST_DISK_PERFORMANCE = "GetDiskPerformance"
	GET_HOST_NETWORK_TRAFFIC  = "GetNetWorkTraffic"
)

type ListHostInfo struct {
	Total int        `json:"total" binding:"required"`
	Info  []HostInfo `json:"info" binding:"required"`
}

type HostInfo struct {
	HostName    string            `json:"hostName" binding:"required"`
	IP          string            `json:"ip" binding:"required"`
	Machine     string            `json:"machine" binding:"required"`
	Release     string            `json:"kernelRelease" binding:"required"`
	Version     string            `json:"kernelVersion" binding:"required"`
	System      string            `json:"operatingSystem" binding:"required"`
	CPUCores    metricomm.CPUInfo `json:"cpuCores" binding:"required"`
	DiskNum     uint32            `json:"diskNum" binding:"required"`
	MemoryTotal uint64            `json:"memory" binding:"required"`
}

type HostInfoWithPerformance struct {
	Host            HostInfo                           `json:"host" binding:"required"`
	CPUUtilization  []metricomm.RangeMetricItem        `json:"cpuUtilization" binding:"required"`
	MemUtilization  []metricomm.RangeMetricItem        `json:"memUtilization" binding:"required"`
	DiskPerformance map[string][]metricomm.Performance `json:"diskPerformance" binding:"required"`
	NetWorkTraffic  NetWorkTraffic                     `json:"networkTraffic" binding:"required"`
}

type NetWorkTraffic struct {
	NetWorkReceive  map[string][]metricomm.RangeMetricItem `json:"receive" binding:"required"`
	NetWorkTransmit map[string][]metricomm.RangeMetricItem `json:"transmit" binding:"required"`
}

type callback func(instance string) (interface{}, error)

func getHostErrnoByName(name string) errno.Errno {
	switch name {
	case GET_HOST_INFO:
		return errno.GET_HOST_INFO_FAILED
	case GET_HOST_CPU_INFO:
		return errno.GET_HOST_CPU_INFO_FAILED
	case GET_HOST_MEM_INFO:
		return errno.GET_HOST_MEM_INFO_FAILED
	case GET_HOST_DISK_NUM:
		return errno.GET_HOST_DISK_NUM_FAILED
	case GET_HOST_CPU_UTILIZATION:
		return errno.GET_HOST_CPU_UTILIZATION_FAILED
	case GET_HOST_MEM_UTILIZATION:
		return errno.GET_HOST_MEM_UTILIZATION_FAILED
	case GET_HOST_DISK_PERFORMANCE:
		return errno.GET_HOST_DISK_PERFORMANCE_FAILED
	case GET_HOST_NETWORK_TRAFFIC:
		return errno.GET_HOST_NETWORK_TRAFFIC_FAILED
	}
	return errno.UNKNOW_ERROR
}

func getInstanceByHostName(hostname string) (string, error) {
	if hostname == "" {
		return "", nil
	}
	baseInfo, err := metricomm.GetHostInfo("")
	if err != nil {
		return "", err
	}
	for k, info := range baseInfo.(map[string]metricomm.HostInfo) {
		if info.HostName == hostname {
			return k, nil
		}
	}
	return "", fmt.Errorf("hostname not exist, hostname = %s", hostname)
}

func getHostNameByInstance(instances []string) (map[string]string, error) {
	baseInfo, err := metricomm.GetHostInfo("")
	if err != nil {
		return nil, err
	}
	ret := make(map[string]string)
	for inst, info := range baseInfo.(map[string]metricomm.HostInfo) {
		ret[inst] = info.HostName
	}
	return ret, nil
}

func sortHost(hosts []HostInfo) {
	sort.Slice(hosts, func(i, j int) bool {
		return hosts[i].HostName < hosts[j].HostName
	})
}

func ListHost(r *pigeon.Request, size, page uint32) (interface{}, errno.Errno) {
	instance := ""
	hostInfos := []HostInfo{}
	// map[instance]*HostInfo
	hostsMap := make(map[string]*HostInfo)
	requests := []callback{
		metricomm.GetHostInfo,
		metricomm.GetHostCPUInfo,
		metricomm.GetHostMemoryInfo,
		metricomm.GetHostDiskNum,
	}
	// TODO: improve with reflect func name
	requestName := []string{
		GET_HOST_INFO,
		GET_HOST_CPU_INFO,
		GET_HOST_MEM_INFO,
		GET_HOST_DISK_NUM,
	}
	requestSize := len(requests)
	ret := make(chan common.QueryResult, requestSize)
	for index, fn := range requests {
		go func(key string, fn callback) {
			info, err := fn(instance)
			ret <- common.QueryResult{
				Key:    key,
				Err:    err,
				Result: info,
			}
		}(requestName[index], fn)
	}
	count := 0
	for res := range ret {
		if res.Err != nil {
			r.Logger().Error("ListHost failed",
				pigeon.Field("step", res.Key.(string)),
				pigeon.Field("error", res.Err),
				pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
			return nil, getHostErrnoByName(res.Key.(string))
		}
		switch res.Key.(string) {
		case GET_HOST_INFO:
			for k, info := range res.Result.(map[string]metricomm.HostInfo) {
				if _, ok := hostsMap[k]; ok {
					hostsMap[k].HostName = info.HostName
					hostsMap[k].IP = info.IP
					hostsMap[k].Machine = info.Machine
					hostsMap[k].Release = info.Release
					hostsMap[k].System = info.System
					hostsMap[k].Version = info.Version
				} else {
					hostsMap[k] = &HostInfo{
						HostName: info.HostName,
						IP:       info.IP,
						Machine:  info.Machine,
						Release:  info.Release,
						System:   info.System,
						Version:  info.Version,
					}
				}
			}
		case GET_HOST_CPU_INFO:
			cpuInfo := res.Result.(map[string]metricomm.CPUInfo)
			for k, v := range cpuInfo {
				if info, ok := hostsMap[k]; ok {
					(*info).CPUCores = v
				} else {
					hostsMap[k] = &HostInfo{
						CPUCores: v,
					}
				}
			}
		case GET_HOST_MEM_INFO:
			memInfo := res.Result.(map[string]uint64)
			for k, v := range memInfo {
				if info, ok := hostsMap[k]; ok {
					(*info).MemoryTotal = v
				} else {
					hostsMap[k] = &HostInfo{
						MemoryTotal: v,
					}
				}
			}
		case GET_HOST_DISK_NUM:
			diskNum := res.Result.(map[string]uint32)
			for k, v := range diskNum {
				if info, ok := hostsMap[k]; ok {
					(*info).DiskNum = v
				} else {
					hostsMap[k] = &HostInfo{
						DiskNum: v,
					}
				}
			}
		}
		count += 1
		if count >= requestSize {
			break
		}
	}

	for _, v := range hostsMap {
		hostInfos = append(hostInfos, *v)
	}
	sortHost(hostInfos)
	length := uint32(len(hostInfos))
	start := (page - 1) * size
	end := common.MinUint32(page*size, length)
	listHostInfo := ListHostInfo{
		Info: []HostInfo{},
	}
	listHostInfo.Total = len(hostInfos)
	if start >= length {
		return listHostInfo, errno.OK
	}
	listHostInfo.Info = hostInfos[start:end]
	return listHostInfo, errno.OK
}

func GetHost(r *pigeon.Request, hostname string) (interface{}, errno.Errno) {
	hostPerformance := HostInfoWithPerformance{}
	instance, err := getInstanceByHostName(hostname)
	if err != nil {
		r.Logger().Error("GetHostPerformance getInstanceByHostName failed",
			pigeon.Field("hostname", hostname),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_INSTANCE_BY_HOSTNAME_FAILED
	}

	requests := []callback{
		metricomm.GetHostInfo,
		metricomm.GetHostCPUInfo,
		metricomm.GetHostMemoryInfo,
		metricomm.GetHostDiskNum,
		metricomm.GetHostCPUUtilization,
		metricomm.GetHostMemUtilization,
		metricomm.GetDiskPerformance,
		metricomm.GetNetWorkTraffic,
	}

	// TODO: improve with reflect func name
	requestName := []string{
		GET_HOST_INFO,
		GET_HOST_CPU_INFO,
		GET_HOST_MEM_INFO,
		GET_HOST_DISK_NUM,
		GET_HOST_CPU_UTILIZATION,
		GET_HOST_MEM_UTILIZATION,
		GET_HOST_DISK_PERFORMANCE,
		GET_HOST_NETWORK_TRAFFIC,
	}
	requestSize := len(requests)
	ret := make(chan common.QueryResult, requestSize)
	for index, fn := range requests {
		go func(key string, fn callback) {
			info, err := fn(instance)
			ret <- common.QueryResult{
				Key:    key,
				Err:    err,
				Result: info,
			}
		}(requestName[index], fn)
	}

	count := 0
	for res := range ret {
		if res.Err != nil {
			r.Logger().Error("GetHost failed",
				pigeon.Field("step", res.Key.(string)),
				pigeon.Field("instance", instance),
				pigeon.Field("error", res.Err),
				pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
			return nil, getHostErrnoByName(res.Key.(string))
		}
		switch res.Key.(string) {
		case GET_HOST_INFO:
			info := res.Result.(map[string]metricomm.HostInfo)[instance]
			hostPerformance.Host.HostName = info.HostName
			hostPerformance.Host.IP = info.IP
			hostPerformance.Host.Machine = info.Machine
			hostPerformance.Host.Release = info.Release
			hostPerformance.Host.System = info.System
			hostPerformance.Host.Version = info.Version
		case GET_HOST_CPU_INFO:
			cpuInfo := res.Result.(map[string]metricomm.CPUInfo)
			hostPerformance.Host.CPUCores = cpuInfo[instance]
		case GET_HOST_MEM_INFO:
			memInfo := res.Result.(map[string]uint64)
			hostPerformance.Host.MemoryTotal = memInfo[instance]
		case GET_HOST_DISK_NUM:
			diskNum := res.Result.(map[string]uint32)
			hostPerformance.Host.DiskNum = diskNum[instance]
		case GET_HOST_CPU_UTILIZATION:
			cpuUtilization := res.Result.(map[string][]metricomm.RangeMetricItem)
			hostPerformance.CPUUtilization = cpuUtilization[instance]
		case GET_HOST_MEM_UTILIZATION:
			memUtilization := res.Result.(map[string][]metricomm.RangeMetricItem)
			hostPerformance.MemUtilization = memUtilization[instance]
		case GET_HOST_DISK_PERFORMANCE:
			diskPerformance := res.Result.(map[string][]metricomm.Performance)
			hostPerformance.DiskPerformance = diskPerformance
		case GET_HOST_NETWORK_TRAFFIC:
			networkTraffic := res.Result.(metricomm.NetworkTraffic)
			hostPerformance.NetWorkTraffic.NetWorkReceive = networkTraffic.Receive
			hostPerformance.NetWorkTraffic.NetWorkTransmit = networkTraffic.Transmit
		}
		count += 1
		if count >= requestSize {
			break
		}
	}
	// ensure performance data is time sequence
	sort.Slice(hostPerformance.CPUUtilization, func(i, j int) bool {
		return hostPerformance.CPUUtilization[i].Timestamp < hostPerformance.CPUUtilization[j].Timestamp
	})
	sort.Slice(hostPerformance.MemUtilization, func(i, j int) bool {
		return hostPerformance.MemUtilization[i].Timestamp < hostPerformance.MemUtilization[j].Timestamp
	})
	for key := range hostPerformance.DiskPerformance {
		sort.Slice(hostPerformance.DiskPerformance[key], func(i, j int) bool {
			return hostPerformance.DiskPerformance[key][i].Timestamp < hostPerformance.DiskPerformance[key][j].Timestamp
		})
	}
	for key := range hostPerformance.NetWorkTraffic.NetWorkReceive {
		sort.Slice(hostPerformance.NetWorkTraffic.NetWorkReceive[key], func(i, j int) bool {
			return hostPerformance.NetWorkTraffic.NetWorkReceive[key][i].Timestamp <
			hostPerformance.NetWorkTraffic.NetWorkReceive[key][j].Timestamp
		})
	}
	for key := range hostPerformance.NetWorkTraffic.NetWorkTransmit {
		sort.Slice(hostPerformance.NetWorkTraffic.NetWorkTransmit[key], func(i, j int) bool {
			return hostPerformance.NetWorkTraffic.NetWorkTransmit[key][i].Timestamp <
			hostPerformance.NetWorkTraffic.NetWorkTransmit[key][j].Timestamp
		})
	}
	return hostPerformance, errno.OK
}
