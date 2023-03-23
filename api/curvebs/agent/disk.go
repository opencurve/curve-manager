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

	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/errno"
	metricomm "github.com/opencurve/curve-manager/internal/metrics/common"
	"github.com/opencurve/pigeon"
)

const (
	// node disk info
	DEVICE = "device"
	MODEL  = "model"

	// request
	GET_DISK_FILESYSTEM_INFO         = "GetDiskFileSystemInfo"
	GET_DISK_TYPE                    = "GetDiskType"
	GET_DISK_WRITE_CACHE_ENABLE_FLAG = "GetDiskWriteCacheEnableFlag"
)

type ListDiskInfo struct {
	Total int        `json:"total" binding:"required"`
	Info  []DiskInfo `json:"info" binding:"required"`
}

type DiskInfo struct {
	HostName   string `json:"hostName" binding:"required"`
	Device     string `json:"device" binding:"required"`
	DiskType   string `json:"diskType" binding:"required"`
	Model      string `json:"model" binding:"required"`
	WriteCache string `json:"writeCache"`
	MountPoint string `json:"mountPoint"`
	SpaceTotal uint32 `json:"spaceTotal"`
	SpaceUsed  uint32 `json:"spaceUsed"`
}

func getDiskErrnoByName(name string) errno.Errno {
	switch name {
	case GET_DISK_FILESYSTEM_INFO:
		return errno.GET_DISK_FILESYSTEM_INFO_FAILED
	case GET_DISK_TYPE:
		return errno.GET_DISK_TYPE_FAILED
	case GET_DISK_WRITE_CACHE_ENABLE_FLAG:
		return errno.GET_DISK_WRITE_CACHE_FAILED
	}
	return errno.UNKNOW_ERROR
}

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
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return disksInfo, errno.OK
	}
	// 1. get disk device list
	disks, err := metricomm.ListDiskInfo(instance)
	if err != nil {
		r.Logger().Error("ListDisk metricomm.ListDisk failed",
			pigeon.Field("instance", instance),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.LIST_DISK_INFO_FAILED
	}
	// instance -> hostname
	insts := make([]string, len(disks.(map[string][]map[string]string)))
	for k := range disks.(map[string][]map[string]string) {
		insts = append(insts, k)
	}
	inst2host, err := getHostNameByInstance(insts)
	if err != nil {
		r.Logger().Error("ListDisk getHostNameByInstance failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_HOSTNAME_BY_INSTANCE_FAILED
	}
	for inst, devs := range disks.(map[string][]map[string]string) {
		retMap[inst] = make(map[string]*DiskInfo)
		hostName := inst2host[inst]
		for _, dev := range devs {
			retMap[inst][dev[DEVICE]] = &DiskInfo{
				HostName: hostName,
				Device:   dev[DEVICE],
				Model:    dev[MODEL],
			}
		}
	}

	requests := []callback{
		metricomm.GetDiskFileSystemInfo,
		metricomm.GetDiskType,
		metricomm.GetDiskWriteCacheEnableFlag,
	}
	// TODO: improve with reflect func name
	requestName := []string{
		GET_DISK_FILESYSTEM_INFO,
		GET_DISK_TYPE,
		GET_DISK_WRITE_CACHE_ENABLE_FLAG,
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
			r.Logger().Error("ListDisk failed",
				pigeon.Field("step", res.Key.(string)),
				pigeon.Field("error", res.Err),
				pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
			return nil, getDiskErrnoByName(res.Key.(string))
		}
		switch res.Key.(string) {
		case GET_DISK_FILESYSTEM_INFO:
			fileSystemInfos := res.Result.(map[string]map[string]metricomm.FileSystemInfo)
			for inst, infos := range fileSystemInfos {
				if _, ok := retMap[inst]; ok {
					for dev, info := range infos {
						if _, ok := retMap[inst][dev]; ok {
							retMap[inst][dev].MountPoint = info.MountPoint
							retMap[inst][dev].SpaceTotal = uint32(info.SpaceTotal)
							retMap[inst][dev].SpaceUsed = uint32(info.SpaceTotal) - uint32(info.SpaceAvail)
						}
					}
				}
			}
		case GET_DISK_TYPE:
			diskTypes := res.Result.(map[string]map[string]string)
			for inst, infos := range diskTypes {
				if _, ok := retMap[inst]; ok {
					for dev, v := range infos {
						if _, ok := retMap[inst][dev]; ok {
							retMap[inst][dev].DiskType = v
						}
					}
				}
			}
		case GET_DISK_WRITE_CACHE_ENABLE_FLAG:
			wrietCaches := res.Result.(map[string]map[string]string)
			for inst, infos := range wrietCaches {
				if _, ok := retMap[inst]; ok {
					for dev, v := range infos {
						if _, ok := retMap[inst][dev]; ok {
							retMap[inst][dev].WriteCache = v
						}
					}
				}
			}
		}
		count += 1
		if count >= requestSize {
			break
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
	end := common.MinUint32(page*size, length)
	listDiskInfo := ListDiskInfo{
		Info: []DiskInfo{},
	}
	listDiskInfo.Total = len(disksInfo)
	if start >= length {
		return listDiskInfo, errno.OK
	}
	listDiskInfo.Info = disksInfo[start:end]
	return listDiskInfo, errno.OK
}
