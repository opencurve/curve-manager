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
	"github.com/opencurve/curve-manager/internal/http/curvebs"
	"sort"

	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/errno"
	bshttp "github.com/opencurve/curve-manager/internal/http/curvebs"
	"github.com/opencurve/curve-manager/internal/metrics/bsmetric"
	metricomm "github.com/opencurve/curve-manager/internal/metrics/common"
	"github.com/opencurve/pigeon"
)

type Server struct {
	Id           uint32                `json:"id" binding:"required"`
	Hostname     string                `json:"hostname" binding:"required"`
	InternalIp   string                `json:"internalIp" binding:"required"`
	InternalPort uint32                `json:"internalPort" binding:"required"`
	ExternalIp   string                `json:"externalIp" binding:"required"`
	ExternalPort uint32                `json:"externalPort" binding:"required"`
	ChunkServers []curvebs.ChunkServer `json:"chunkservers" binding:"required"`
}

type Zone struct {
	Id      uint32   `json:"id" binding:"required"`
	Name    string   `json:"name" binding:"required"`
	Servers []Server `json:"servers" binding:"required"`
}

type Pool struct {
	Id             uint32 `json:"id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Type           string `json:"type" binding:"required"`
	Zones          []Zone `json:"zones" binding:"required"`
	physicalPoolId uint32
}

type Space struct {
	Total       uint64 `json:"total" binding:"required"`
	Alloc       uint64 `json:"alloc" binding:"required"`
	CanRecycled uint64 `json:"canRecycled" binding:"required"`
}

type PoolInfo struct {
	Id             uint32 `json:"id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	PhysicalPoolId uint32 `json:"physicalPoolId" binding:"required"`
	Type           string `json:"type" binding:"required"`
	CreateTime     string `json:"createTime" binding:"required"`
	AllocateStatus string `json:"allocateStatus" binding:"required"`
	ScanEnable     bool   `json:"scanEnable"`
	ServerNum      uint32 `json:"serverNum" binding:"required"`
	ChunkServerNum uint32 `json:"chunkServerNum" binding:"required"`
	CopysetNum     uint32 `json:"copysetNum" binding:"required"`
	Space          Space  `json:"space" binding:"required"`
}

type PoolInfoWithPerformance struct {
	Info        PoolInfo                `json:"info" binding:"required"`
	Performance []metricomm.Performance `json:"performance" binding:"required"`
}

// get chunkservers of server concurrently
func listChunkServer(pools *[]Pool, size int) error {
	ret := make(chan common.QueryResult, size)
	count := 0
	for pIndex, pool := range *pools {
		for zIndex, zone := range pool.Zones {
			for sIndex, server := range zone.Servers {
				go func(id uint32, addr *Server) {
					chunkservers, err := bshttp.GMdsClient.ListChunkServer(id)
					ret <- common.QueryResult{
						Key:    addr,
						Result: chunkservers,
						Err:    err,
					}
				}(server.Id, &(*pools)[pIndex].Zones[zIndex].Servers[sIndex])
			}
		}
	}
	for res := range ret {
		if res.Err != nil {
			return res.Err
		}
		for _, cs := range res.Result.([]curvebs.ChunkServer) {
			res.Key.(*Server).ChunkServers = append(res.Key.(*Server).ChunkServers, cs)
		}
		count += 1
		if count >= size {
			break
		}
	}
	return nil
}

// get servers of zones concurrently
func listZoneServer(pools *[]Pool, size int) error {
	ret := make(chan common.QueryResult, size)
	count := 0
	number := 0
	for pIndex, pool := range *pools {
		for zIndex, zone := range pool.Zones {
			go func(id uint32, addr *Zone) {
				servers, err := bshttp.GMdsClient.ListZoneServer(id)
				ret <- common.QueryResult{
					Key:    addr,
					Result: servers,
					Err:    err,
				}
			}(zone.Id, &(*pools)[pIndex].Zones[zIndex])
		}
	}
	for res := range ret {
		if res.Err != nil {
			return res.Err
		}
		for _, s := range res.Result.([]curvebs.Server) {
			var server Server
			server.Id = s.Id
			server.Hostname = s.HostName
			server.InternalIp = s.InternalIp
			server.InternalPort = s.InternalPort
			server.ExternalIp = s.ExternalIp
			server.ExternalPort = s.ExternalPort
			res.Key.(*Zone).Servers = append(res.Key.(*Zone).Servers, server)
			number = number + 1
		}
		count += 1
		if count >= size {
			break
		}
	}
	return listChunkServer(pools, number)
}

// get zones of pools concurrently
func listPoolZone(pools *[]Pool) error {
	size := len(*pools)
	ret := make(chan common.QueryResult, size)
	count := 0
	number := 0
	for index, pool := range *pools {
		go func(id uint32, addr *Pool) {
			zones, err := bshttp.GMdsClient.ListPoolZone(id)
			ret <- common.QueryResult{
				Key:    addr,
				Result: zones,
				Err:    err,
			}
		}(pool.physicalPoolId, &(*pools)[index])
	}
	for res := range ret {
		if res.Err != nil {
			return res.Err
		}
		for _, z := range res.Result.([]curvebs.Zone) {
			var zone Zone
			zone.Id = z.Id
			zone.Name = z.Name
			res.Key.(*Pool).Zones = append(res.Key.(*Pool).Zones, zone)
			number = number + 1
		}
		count += 1
		if count >= size {
			break
		}
	}
	return listZoneServer(pools, number)
}

func getPoolSpace(pools *[]PoolInfo) error {
	// get can be recycled space
	_, recycledSize, err := bshttp.GMdsClient.GetFileAllocatedSize(RECYCLEBIN_DIR)
	if err != nil {
		return err
	}

	// get capacity and used space
	size := len(*pools)
	ret := make(chan common.QueryResult, size)
	for index, pool := range *pools {
		go func(name string, index int) {
			space, err := bsmetric.GetPoolSpace(name)
			ret <- common.QueryResult{
				Key:    index,
				Result: space,
				Err:    err,
			}
		}(pool.Name, index)
	}
	count := 0
	for res := range ret {
		if res.Err != nil {
			return res.Err
		}
		(*pools)[res.Key.(int)].Space.Total = res.Result.(*metricomm.Space).Total
		(*pools)[res.Key.(int)].Space.Alloc = res.Result.(*metricomm.Space).Used
		(*pools)[res.Key.(int)].Space.CanRecycled = recycledSize[(*pools)[res.Key.(int)].Id]
		count += 1
		if count >= size {
			break
		}
	}
	return nil
}

func getPoolItemNum(pools *[]PoolInfo) error {
	size := len(*pools)
	ret := make(chan common.QueryResult, size)
	for index, pool := range *pools {
		go func(name string, index int) {
			number, err := bsmetric.GetPoolItemNum(name)
			ret <- common.QueryResult{
				Key:    index,
				Result: number,
				Err:    err,
			}
		}(pool.Name, index)
	}
	count := 0
	for res := range ret {
		if res.Err != nil {
			return res.Err
		}
		(*pools)[res.Key.(int)].ServerNum = res.Result.(*bsmetric.PoolItemNum).ServerNum
		(*pools)[res.Key.(int)].ChunkServerNum = res.Result.(*bsmetric.PoolItemNum).ChunkServerNum
		(*pools)[res.Key.(int)].CopysetNum = res.Result.(*bsmetric.PoolItemNum).CopysetNum
		count += 1
		if count >= size {
			break
		}
	}
	return nil
}

func getPoolPerformance(pools *[]PoolInfoWithPerformance, start, end, interval uint64) error {
	size := len(*pools)
	ret := make(chan common.QueryResult, size)
	for index, pool := range *pools {
		go func(name string, index int) {
			performance, err := bsmetric.GetPoolPerformance(name, start, end, interval)
			ret <- common.QueryResult{
				Key:    index,
				Result: performance,
				Err:    err,
			}
		}(pool.Info.Name, index)
	}
	count := 0
	for res := range ret {
		if res.Err != nil {
			return res.Err
		}
		(*pools)[res.Key.(int)].Performance = append((*pools)[res.Key.(int)].Performance, res.Result.([]metricomm.Performance)...)
		count += 1
		if count >= size {
			break
		}
	}
	return nil
}

func sortLogicalPool(pools []PoolInfo) {
	sort.Slice(pools, func(i, j int) bool {
		return pools[i].Id < pools[j].Id
	})
}

func sortTopology(pools []Pool) {
	sort.Slice(pools, func(i, j int) bool {
		return pools[i].Id < pools[j].Id
	})
	for index := range pools {
		sort.Slice(pools[index].Zones, func(i, j int) bool {
			return pools[index].Zones[i].Id < pools[index].Zones[j].Id
		})
	}
	for zindex := range pools {
		for sindex := range pools[zindex].Zones {
			sort.Slice(pools[zindex].Zones[sindex].Servers, func(i, j int) bool {
				return pools[zindex].Zones[sindex].Servers[i].Id < pools[zindex].Zones[sindex].Servers[j].Id
			})
		}
	}

}

func ListLogicalPool(r *pigeon.Request) (interface{}, errno.Errno) {
	result := []PoolInfo{}
	// get info from mds
	pools, err := bshttp.GMdsClient.ListLogicalPool()
	if err != nil {
		r.Logger().Error("ListLogicalPool bsrpc.ListLogicalPool failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.LIST_POOL_FAILED
	}

	for _, pool := range pools {
		var info PoolInfo
		info.Id = pool.Id
		info.Name = pool.Name
		info.PhysicalPoolId = pool.PhysicalPoolId
		info.Type = pool.Type
		info.CreateTime = pool.CreateTime
		info.AllocateStatus = pool.AllocateStatus
		info.ScanEnable = pool.ScanEnable
		result = append(result, info)
	}

	// get info from monitor
	err = getPoolItemNum(&result)
	if err != nil {
		r.Logger().Error("ListLogicalPool getPoolItemNum failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_POOL_ITEM_NUMBER_FAILED
	}

	err = getPoolSpace(&result)
	if err != nil {
		r.Logger().Error("ListLogicalPool getPoolSpace failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_POOL_SPACE_FAILED
	}
	sortLogicalPool(result)
	return &result, errno.OK
}

func GetLogicalPool(r *pigeon.Request, poolId uint32, start, end, interval uint64) (interface{}, errno.Errno) {
	pool, err := bshttp.GMdsClient.GetLogicalPool(poolId)
	if err != nil {
		r.Logger().Error("GetLogicalPool bsrpc.GetLogicalPool failed",
			pigeon.Field("poolId", poolId),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_POOL_FAILED
	}

	var poolInfo PoolInfoWithPerformance
	var info PoolInfo
	info.Id = pool.Id
	info.Name = pool.Name
	info.PhysicalPoolId = pool.PhysicalPoolId
	info.Type = pool.Type
	info.CreateTime = pool.CreateTime
	info.AllocateStatus = pool.AllocateStatus
	info.ScanEnable = pool.ScanEnable
	tmp := []PoolInfo{info}
	// get info from monitor
	err = getPoolItemNum(&tmp)
	if err != nil {
		r.Logger().Error("GetLogicalPool getPoolItemNum failed",
			pigeon.Field("poolId", poolId),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_POOL_ITEM_NUMBER_FAILED
	}

	err = getPoolSpace(&tmp)
	if err != nil {
		r.Logger().Error("GetLogicalPool getPoolSpace failed",
			pigeon.Field("poolId", poolId),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_POOL_SPACE_FAILED
	}

	poolInfo.Info = tmp[0]
	poolInfo.Performance = []metricomm.Performance{}
	result := []PoolInfoWithPerformance{poolInfo}
	err = getPoolPerformance(&result, start, end, interval)
	if err != nil {
		r.Logger().Error("GetLogicalPool getPoolPerformance failed",
			pigeon.Field("poolId", poolId),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_POOL_PERFORMANCE_FAILED
	}
	// ensure performance data is time sequence
	sort.Slice(result[0].Performance, func(i, j int) bool {
		return result[0].Performance[i].Timestamp < result[0].Performance[j].Timestamp
	})
	return result[0], errno.OK
}

func ListTopology(r *pigeon.Request) (interface{}, errno.Errno) {
	result := []Pool{}
	logicalPools, err := bshttp.GMdsClient.ListLogicalPool()
	if err != nil {
		r.Logger().Error("ListTopology bsrpc.ListLogicalPool failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.LIST_POOL_FAILED
	}
	for _, lp := range logicalPools {
		var pool Pool
		pool.Id = lp.Id
		pool.physicalPoolId = lp.PhysicalPoolId
		pool.Name = lp.Name
		pool.Type = lp.Type
		result = append(result, pool)
	}
	err = listPoolZone(&result)
	if err != nil {
		r.Logger().Error("ListTopology listPoolZone failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.LIST_POOL_ZONE_FAILED
	}
	sortTopology(result)
	return &result, errno.OK
}
