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

package agent

import (
	"fmt"
	"sort"
	"time"

	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/curve-manager/internal/metrics/bsmetric"
	comm "github.com/opencurve/curve-manager/internal/metrics/common"
	bsrpc "github.com/opencurve/curve-manager/internal/rpc/curvebs"
	"github.com/opencurve/pigeon"
)

const (
	ORDER_BY_ID              = "id"
	ORDER_BY_CTIME           = "ctime"
	ORDER_BY_LENGTH          = "length"
	ORDER_DIRECTION_INCREASE = 1
	ORDER_DIRECTION_DECREASE = -1
)

func sortFile(files []bsrpc.FileInfo, orderKey string, direction int) {
	sort.Slice(files, func(i, j int) bool {
		switch orderKey {
		case ORDER_BY_CTIME:
			itime, _ := time.Parse(common.TIME_FORMAT, files[i].Ctime)
			jtime, _ := time.Parse(common.TIME_FORMAT, files[j].Ctime)
			if direction == ORDER_DIRECTION_INCREASE {
				return itime.Unix() < jtime.Unix()
			}
			return itime.Unix() > jtime.Unix()
		case ORDER_BY_LENGTH:
			if direction == ORDER_DIRECTION_INCREASE {
				return files[i].Length < files[j].Length
			}
			return files[i].Length > files[j].Length
		}
		if direction == ORDER_DIRECTION_INCREASE {
			return files[i].Id < files[j].Id
		}
		return files[i].Id > files[j].Id
	})
}

func getVolumeAllocSize(path string, volumes *[]VolumeInfo) error {
	size := len(*volumes)
	ret := make(chan common.QueryResult, size)
	for index, volume := range *volumes {
		go func(vname string, addr *VolumeInfo) {
			_, poolSize, err := bsrpc.GMdsClient.GetFileAllocatedSize(vname)
			ret <- common.QueryResult{
				Key:    addr,
				Result: poolSize,
				Err:    err,
			}
		}(path+volume.Info.FileName, &(*volumes)[index])
	}
	count := 0
	for res := range ret {
		if res.Err != nil {
			return res.Err
		}
		var totalAlloc uint64
		for id, pSize := range res.Result.(map[uint32]uint64) {
			vPool := VolumePoolInfo{}
			vPool.Id = id
			vPool.Alloc = uint32(pSize)
			res.Key.(*VolumeInfo).Pools = append(res.Key.(*VolumeInfo).Pools, vPool)
			totalAlloc += pSize
		}
		res.Key.(*VolumeInfo).Info.AllocateSize = totalAlloc
		count += 1
		if count >= size {
			break
		}
	}
	return getVolumePoolInfo(volumes)
}

func getVolumePoolInfo(volumes *[]VolumeInfo) error {
	pools, err := bsrpc.GMdsClient.ListLogicalPool()
	if err != nil {
		return fmt.Errorf("getVolumePoolInfo failed, %s", err)
	}

	poolMap := make(map[uint32]*bsrpc.LogicalPool)
	for index, pool := range pools {
		poolMap[pool.Id] = &pools[index]
	}

	for i, vInfo := range *volumes {
		for j, _ := range vInfo.Pools {
			id := (*volumes)[i].Pools[j].Id
			(*volumes)[i].Pools[j].Name = *&poolMap[id].Name
			(*volumes)[i].Pools[j].Type = *&poolMap[id].Type
		}
	}
	return nil
}

func getVolumePerformance(path string, volumes *[]VolumeInfo) error {
	size := len(*volumes)
	ret := make(chan common.QueryResult, size)
	for index, volume := range *volumes {
		go func(vname string, addr *VolumeInfo) {
			performances, err := bsmetric.GetVolumePerformance(vname)
			ret <- common.QueryResult{
				Key:    addr,
				Result: performances,
				Err:    err,
			}
		}(path+volume.Info.FileName, &(*volumes)[index])
	}
	count := 0
	for res := range ret {
		if res.Err != nil {
			return res.Err
		}
		res.Key.(*VolumeInfo).Performance = res.Result.([]comm.UserPerformance)
		count += 1
		if count >= size {
			break
		}
	}
	return nil
}

func ListVolume(r *pigeon.Request, size, page uint32, path, key string, direction int) (interface{}, errno.Errno) {
	// get root auth info
	authInfo, err := bsmetric.GetAuthInfoOfRoot()
	if err != "" {
		r.Logger().Error("ListVolume bsmetric.GetAuthInfoOfRoot failed",
			pigeon.Field("path", path),
			pigeon.Field("size", size),
			pigeon.Field("page", page),
			pigeon.Field("sortkey", key),
			pigeon.Field("error", err))
		return nil, errno.GET_ROOT_AUTH_FAILED
	}

	// create signature
	date := time.Now().UnixMicro()
	str2sig := common.GetString2Signature(date, authInfo.UserName)
	sig := common.CalcString2Signature(str2sig, authInfo.PassWord)

	fileInfos, e := bsrpc.GMdsClient.ListDir(path, authInfo.UserName, sig, uint64(date))
	if e != nil {
		r.Logger().Error("ListVolume bsrpc.ListDir failed",
			pigeon.Field("path", path),
			pigeon.Field("size", size),
			pigeon.Field("page", page),
			pigeon.Field("sortkey", key),
			pigeon.Field("error", err))
		return nil, errno.LIST_VOLUME_FAILED
	}

	if len(fileInfos) == 0 {
		return []VolumeInfo{}, errno.OK
	}

	sortFile(fileInfos, key, direction)
	length := uint32(len(fileInfos))
	start := (page - 1) * size
	var end uint32
	if page*size > length {
		end = length
	} else {
		end = page * size
	}

	var volumes []VolumeInfo
	for _, info := range fileInfos[start:end] {
		vInfo := VolumeInfo{}
		vInfo.Info = info
		vInfo.Pools = []VolumePoolInfo{}
		vInfo.Performance = []comm.UserPerformance{}
		volumes = append(volumes, vInfo)
	}
	e = getVolumeAllocSize(path, &volumes)
	if e != nil {
		r.Logger().Error("ListVolume getVolumeAllocSize failed",
			pigeon.Field("error", err))
		return nil, errno.GET_VOLUME_ALLOC_SIZE_FAILED
	}

	// get performance of the volume
	e = getVolumePerformance(path, &volumes)
	if e != nil {
		r.Logger().Error("ListVolume getVolumePerformance failed",
			pigeon.Field("error", err))
		return nil, errno.GET_VOLUME_PERFORMANCE_FAILED
	}
	return volumes, errno.OK
}
