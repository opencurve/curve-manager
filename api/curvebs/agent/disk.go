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
	"sort"

	"github.com/opencurve/curve-manager/internal/errno"
	metricomm "github.com/opencurve/curve-manager/internal/metrics/common"
	"github.com/opencurve/pigeon"
)

func sortDisk(disks []DiskInfo) {
	sort.Slice(disks, func(i, j int) bool {
		if disks[i].HostName < disks[j].HostName {
			return true
		} else if disks[i].HostName == disks[j].HostName {
			return disks[i].Device < disks[j].Device
		}
		return false
	})
}

func ListDisk(r *pigeon.Request, size, page uint32, hostname string) (interface{}, errno.Errno) {
	disksInfo := []DiskInfo{}
	// map[instance]map[device]*DiskInfo
	retMap := make(map[string]map[string]*DiskInfo)
	instance, err := getInstanceByHostName(hostname)
	if err != nil {
		r.Logger().Error("ListDisk getInstanceByHostName failed",
			pigeon.Field("hostname", hostname),
			pigeon.Field("error", err))
		return nil, errno.GET_INSTANCE_BY_HOSTNAME_FAILED
	}
	// 1. get disk device list
	disks, err := metricomm.ListDisk(instance)
	if err != nil {
		r.Logger().Error("ListDisk metricomm.ListDisk failed",
			pigeon.Field("instance", instance),
			pigeon.Field("error", err))
		return nil, errno.LIST_DISK_FAILED
	}
	// nstance -> hostname
	insts := make([]string, len(disks))
	for k := range disks {
		insts = append(insts, k)
	}
	inst2host, err := getHostNameByInstance(insts)
	if err != nil {
		r.Logger().Error("ListDisk getHostNameByInstance failed",
			pigeon.Field("error", err))
		return nil, errno.GET_HOSTNAME_BY_INSTANCE_FAILED
	}
	for inst, devs := range disks {
		retMap[inst] = make(map[string]*DiskInfo)
		hostName := inst2host[inst]
		for _, dev := range devs {
			retMap[inst][dev] = &DiskInfo{
				HostName: hostName,
				Device:   dev,
			}
		}
	}
	// 2. get filesystem info
	fileSystemInfos, err := metricomm.GetFileSystemInfo(instance)
	if err != nil {
		r.Logger().Error("ListDisk metricomm.GetFileSystemInfo failed",
			pigeon.Field("instance", instance),
			pigeon.Field("error", err))
		return nil, errno.GET_FILESYSTEM_INFO_FAILED
	}
	for inst, infos := range fileSystemInfos {
		for dev, info := range infos {
			if _, ok := retMap[inst][dev]; ok {
				retMap[inst][dev].FileSystem = info.FsType
				retMap[inst][dev].MountPoint = info.MountPoint
				retMap[inst][dev].SpaceTotal = uint32(info.SpaceTotal)
				retMap[inst][dev].SpaceAvail = uint32(info.SpaceAvail)
			}
		}
	}

	for _, item := range retMap {
		for _, v := range item {
			disksInfo = append(disksInfo, *v)
		}
	}
	sortDisk(disksInfo)
	length := uint32(len(disksInfo))
	start := (page - 1) * size
	var end uint32
	if page*size > length {
		end = length
	} else {
		end = page * size
	}
	return disksInfo[start:end], errno.OK
}
