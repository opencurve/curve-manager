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
	"github.com/opencurve/curve-manager/internal/errno"
	metricomm "github.com/opencurve/curve-manager/internal/metrics/common"
	"github.com/opencurve/pigeon"
)

func ListHost(r *pigeon.Request, size, page uint32) (interface{}, errno.Errno) {
	hostInfos := []HostInfo{}
	// map[instance]*HostInfo
	hostsMap := make(map[string]*HostInfo)
	// 1. get host base info
	baseInfo, err := metricomm.GetHostsInfo()
	if err != nil {
		r.Logger().Error("ListHost metricomm.GetHostsInfo failed",
			pigeon.Field("error", err))
		return nil, errno.GET_HOST_INFO_FAILED
	}
	for k, info := range baseInfo {
		hostsMap[k] = &HostInfo{
			HostName: info.HostName,
			IP:       info.IP,
			Machine:  info.Machine,
			Release:  info.Release,
			System:   info.System,
			Version:  info.Version,
		}
	}

	// 2. get node cpu info
	cupInfo, err := metricomm.GetHostCPUInfo()
	if err != nil {
		r.Logger().Error("ListHost metricomm.GetHostCPUInfo failed",
			pigeon.Field("error", err))
		return nil, errno.GET_HOST_CPU_INFO_FAILED
	}
	for k, v := range cupInfo {
		if info, ok := hostsMap[k]; ok {
			(*info).CPUCores = v
		}
	}

	// 3. get memory info
	memInfo, err := metricomm.GetHostMemoryInfo()
	if err != nil {
		r.Logger().Error("ListHost metricomm.GetHostMemoryInfo failed",
			pigeon.Field("error", err))
		return nil, errno.GET_HOST_MEM_INFO_FAILED
	}
	for k, v := range memInfo {
		if info, ok := hostsMap[k]; ok {
			(*info).MemoryTotal = v
		}
	}

	// 4. get disk num
	diskNum, err := metricomm.GetHostDiskNum()
	if err != nil {
		r.Logger().Error("ListHost metricomm.GetHostDiskNum failed",
			pigeon.Field("error", err))
		return nil, errno.GET_HOST_DISK_NUM_FAILED
	}
	for k, v := range diskNum {
		if info, ok := hostsMap[k]; ok {
			(*info).DiskNum = v
		}
	}

	for _, v := range hostsMap {
		hostInfos = append(hostInfos, *v)
	}
	return hostInfos, errno.OK
}

func GetHostPerformance(r *pigeon.Request, hostname string) (interface{}, errno.Errno) {
	hostPerformance := HostPerformance{}
	instance, err := getInstanceByHostName(hostname)
	if err != nil {
		r.Logger().Error("GetHostPerformance getInstanceByHostName failed",
			pigeon.Field("hostname", hostname),
			pigeon.Field("error", err))
		return nil, errno.GET_INSTANCE_BY_HOSTNAME_FAILED
	}
	// 1. get cpu utilization
	cpuUtilization, err := metricomm.GetHostCPUUtilization(instance)
	if err != nil {
		r.Logger().Error("GetHostPerformance metricomm.GetHostCPUUtilization failed",
			pigeon.Field("instance", instance),
			pigeon.Field("error", err))
		return nil, errno.GET_HOST_CPU_UTILIZATION_FAILED
	}
	hostPerformance.CPUUtilization = cpuUtilization[instance]

	// 2. get memory utilization
	memUtilization, err := metricomm.GetHostMemUtilization(instance)
	if err != nil {
		r.Logger().Error("GetHostPerformance metricomm.GetHostMemUtilization failed",
			pigeon.Field("instance", instance),
			pigeon.Field("error", err))
		return nil, errno.GET_HOST_MEM_UTILIZATION_FAILED
	}
	hostPerformance.MemUtilization = memUtilization[instance]

	// 3. get disk performance
	diskPerformance, err := metricomm.GetDiskPerformance(instance)
	if err != nil {
		r.Logger().Error("GetHostPerformance metricomm.GetDiskPerformance failed",
			pigeon.Field("instance", instance),
			pigeon.Field("error", err))
		return nil, errno.GET_HOST_DISK_PERFORMANCE_FAILED
	}
	hostPerformance.DiskPerformance = diskPerformance

	// 4. get network traffic
	networkReceive, networkTransmit, err := metricomm.GetNetWorkTraffic(instance)
	if err != nil {
		r.Logger().Error("GetHostPerformance metricomm.GetNetWorkTraffic failed",
			pigeon.Field("instance", instance),
			pigeon.Field("error", err))
		return nil, errno.GET_HOST_NETWORK_TRAFFIC_FAILED
	}
	hostPerformance.NetWorkTraffic.NetWorkReceive = networkReceive
	hostPerformance.NetWorkTraffic.NetWorkTransmit = networkTransmit

	return hostPerformance, errno.OK
}
