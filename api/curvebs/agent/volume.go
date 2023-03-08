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
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/SeanHai/curve-go-rpc/rpc/curvebs"
	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/curve-manager/internal/metrics/bsmetric"
	metricomm "github.com/opencurve/curve-manager/internal/metrics/common"
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

type AuthInfo struct {
	userName  string
	passWord  string
	signatrue string
	date      uint64
}

type VolumePoolInfo struct {
	Id    uint32 `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Type  string `json:"type" binding:"required"`
	Alloc uint32 `json:"alloc" binding:"required"`
}
type VolumeInfo struct {
	Info        curvebs.FileInfo            `json:"info" binding:"required"`
	Pools       []VolumePoolInfo            `json:"pools"`
	Performance []metricomm.UserPerformance `json:"performance" binding:"required"`
}

type ListVolumeInfo struct {
	Total int                `json:"total" binding:"required"`
	Info  []curvebs.FileInfo `json:"info" binding:"required"`
}

func getUpPath(dir string) string {
	return dir[:strings.LastIndex(dir, "/")]
}

func getString2Signature(date int64, owner string) string {
	return fmt.Sprintf("%d:%s", date, owner)
}

func calcString2Signature(in string, secretKet string) string {
	h := hmac.New(sha256.New, []byte(secretKet))
	h.Write([]byte(in))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func getAuthInfoOfRoot() (*AuthInfo, string) {
	// get root username and password
	userName, passWord, err := bsmetric.GetAuthInfoOfRoot()
	if err != "" {
		return nil, err
	}

	// create signature
	date := time.Now().UnixMicro()
	str2sig := getString2Signature(date, userName)
	sig := calcString2Signature(str2sig, passWord)
	return &AuthInfo{
		userName:  userName,
		passWord:  passWord,
		signatrue: sig,
		date:      uint64(date),
	}, ""
}

func sortFile(files []curvebs.FileInfo, orderKey string, direction int) {
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

func getVolumeAllocSize(dir string, volumes *[]VolumeInfo) error {
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
		}(path.Join(dir, volume.Info.FileName), &(*volumes)[index])
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

	poolMap := make(map[uint32]*curvebs.LogicalPool)
	for index, pool := range pools {
		poolMap[pool.Id] = &pools[index]
	}

	for i, vInfo := range *volumes {
		for j := range vInfo.Pools {
			id := (*volumes)[i].Pools[j].Id
			(*volumes)[i].Pools[j].Name = *&poolMap[id].Name
			(*volumes)[i].Pools[j].Type = *&poolMap[id].Type
		}
	}
	return nil
}

func getVolumePerformance(dir string, volumes *[]VolumeInfo) error {
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
		}(path.Join(dir, volume.Info.FileName), &(*volumes)[index])
	}
	count := 0
	for res := range ret {
		if res.Err != nil {
			return res.Err
		}
		res.Key.(*VolumeInfo).Performance = res.Result.([]metricomm.UserPerformance)
		count += 1
		if count >= size {
			break
		}
	}
	return nil
}

func getVolumeSpaceSize(dir string, size int, volumes *[]VolumeInfo) error {
	if size <= 0 {
		return nil
	}
	ret := make(chan common.QueryResult, size)
	for index, volume := range *volumes {
		if volume.Info.FileType == curvebs.INODE_DIRECTORY {
			go func(vname string, addr *VolumeInfo) {
				size, err := bsrpc.GMdsClient.GetFileSize(vname)
				ret <- common.QueryResult{
					Key:    addr,
					Result: size,
					Err:    err,
				}
			}(path.Join(dir, volume.Info.FileName), &(*volumes)[index])
		}
	}
	count := 0
	for res := range ret {
		if res.Err != nil {
			return res.Err
		}
		res.Key.(*VolumeInfo).Info.Length = res.Result.(uint64)
		count += 1
		if count >= size {
			break
		}
	}
	return nil
}

func ListVolume(r *pigeon.Request, size, page uint32, path, key string, direction int) (interface{}, errno.Errno) {
	listVolumeInfo := ListVolumeInfo{
		Info: []curvebs.FileInfo{},
	}
	authInfo, err := getAuthInfoOfRoot()
	if err != "" {
		r.Logger().Error("ListVolume getAuthInfoOfRoot failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_ROOT_AUTH_FAILED
	}
	fileInfos, e := bsrpc.GMdsClient.ListDir(path, authInfo.userName, authInfo.signatrue, authInfo.date)
	if e != nil {
		r.Logger().Error("ListVolume bsrpc.ListDir failed",
			pigeon.Field("path", path),
			pigeon.Field("size", size),
			pigeon.Field("page", page),
			pigeon.Field("sortkey", key),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.LIST_VOLUME_FAILED
	}
	listVolumeInfo.Total = len(fileInfos)
	if len(fileInfos) == 0 {
		return listVolumeInfo, errno.OK
	}
	sortFile(fileInfos, key, direction)
	length := uint32(len(fileInfos))
	start := (page - 1) * size
	end := common.MinUint32(page*size, length)
	if start >= length {
		return listVolumeInfo, errno.OK
	}

	tmpVolumes := []VolumeInfo{}
	dirSize := 0
	for _, v := range fileInfos[start:end] {
		tmpVolumes = append(tmpVolumes, VolumeInfo{
			Info: v,
		})
		if v.FileType == curvebs.INODE_DIRECTORY {
			dirSize++
		}
	}
	e = getVolumeSpaceSize(path, dirSize, &tmpVolumes)
	if e != nil {
		r.Logger().Error("ListVolume getVolumeSpaceSize failed",
			pigeon.Field("error", e),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return listVolumeInfo, errno.GET_VOLUME_SIZE_FAILED
	}

	e = getVolumeAllocSize(path, &tmpVolumes)
	if e != nil {
		r.Logger().Error("ListVolume getVolumeAllocSize failed",
			pigeon.Field("error", e),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return listVolumeInfo, errno.GET_VOLUME_ALLOC_SIZE_FAILED
	}

	for _, v := range tmpVolumes {
		listVolumeInfo.Info = append(listVolumeInfo.Info, v.Info)
	}
	return listVolumeInfo, errno.OK
}

func GetVolume(r *pigeon.Request, volumeName string) (interface{}, errno.Errno) {
	authInfo, err := getAuthInfoOfRoot()
	if err != "" {
		r.Logger().Error("GetVolume getAuthInfoOfRoot failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_ROOT_AUTH_FAILED
	}
	fileInfo, e := bsrpc.GMdsClient.GetFileInfo(volumeName, authInfo.userName, authInfo.signatrue, authInfo.date)
	if e != nil {
		r.Logger().Error("GetVolume bsrpc.GetFileInfo failed",
			pigeon.Field("fileName", volumeName),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_VOLUME_INFO_FAILED
	}

	volume := VolumeInfo{}
	volume.Info = fileInfo
	volume.Pools = []VolumePoolInfo{}
	volume.Performance = []metricomm.UserPerformance{}

	path := getUpPath(volumeName)
	volumes := []VolumeInfo{volume}
	e = getVolumeAllocSize(path, &volumes)
	if e != nil {
		r.Logger().Error("GetVolume getVolumeAllocSize failed",
			pigeon.Field("fileName", volumeName),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_VOLUME_ALLOC_SIZE_FAILED
	}

	// get performance of the volume
	e = getVolumePerformance(path, &volumes)
	if e != nil {
		r.Logger().Error("GetVolume getVolumePerformance failed",
			pigeon.Field("fileName", volumeName),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_VOLUME_PERFORMANCE_FAILED
	}
	// ensure performance data is time sequence
	sort.Slice(volumes[0].Performance, func(i, j int) bool {
		return volumes[0].Performance[i].Timestamp < volumes[0].Performance[i].Timestamp
	})
	return volumes[0], errno.OK
}
