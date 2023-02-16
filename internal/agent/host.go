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
	metricomm "github.com/opencurve/curve-manager/internal/metrics/common"
)

func ListHost(size, page uint32) (interface{}, error) {
	hostInfos := []HostInfo{}
	hostsMap := make(map[string]*HostInfo)
	// 1. get host base info
	baseInfo, err := metricomm.GetHostsInfo()
	if err != nil {
		return nil, err
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
		return nil, err
	}
	for k, v := range cupInfo {
		if info, ok := hostsMap[k]; ok {
			(*info).CPUCores = v
		}
	}

	// 3. get memory info
	memInfo, err := metricomm.GetHostMemoryInfo()
	if err != nil {
		return nil, err
	}
	for k, v := range memInfo {
		if info, ok := hostsMap[k]; ok {
			(*info).MemoryTotal = v
		}
	}

	// 4. get disk num
	diskNum, err := metricomm.GetHostDiskNum()
	if err != nil {
		return nil, err
	}
	for k, v := range diskNum {
		if info, ok := hostsMap[k]; ok {
			(*info).DiskNum = v
		}
	}

	for _, v := range hostsMap {
		hostInfos = append(hostInfos, *v)
	}
	return hostInfos, nil
}

func GetHostPerformance(hostname string) (interface{}, error) {
	hostPerformance := HostPerformance{}
	instance, err := getInstanceByHostName(hostname)
	if err != nil {
		return nil, err
	}
	// 1. get cpu utilization
	cpuUtilization, err := metricomm.GetHostCPUUtilization(instance)
	if err != nil {
		return nil, err
	}
	hostPerformance.CPUUtilization = cpuUtilization[instance]

	// 2. get memory utilization
	memUtilization, err := metricomm.GetHostMemUtilization(instance)
	if err != nil {
		return nil, err
	}
	hostPerformance.MemUtilization = memUtilization[instance]

	// 3. get disk performance
	diskPerformance, err := metricomm.GetDiskPerformance(instance)
	if err != nil {
		return nil, err
	}
	hostPerformance.DiskPerformance = diskPerformance

	// 4. get network traffic
	networkReceive, networkTransmit, err := metricomm.GetNetWorkTraffic(instance)
	if err != nil {
		return nil, err
	}
	hostPerformance.NetWorkTraffic.NetWorkReceive = networkReceive
	hostPerformance.NetWorkTraffic.NetWorkTransmit = networkTransmit

	return hostPerformance, nil
}
