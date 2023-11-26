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
* Project: Curve-Go-RPC
* Created Date: 2023-03-03
* Author: wanghai (SeanHai)
 */

package curvebs

import (
	"fmt"
	"github.com/SeanHai/curve-go-rpc/rpc/baserpc"
	"github.com/opencurve/curve-manager/internal/http/common"
	"strconv"
	"time"

	"github.com/SeanHai/curve-go-rpc/curvebs_proto/proto/topology"
	"github.com/SeanHai/curve-go-rpc/curvebs_proto/proto/topology/statuscode"
)

const (
	// invalid type
	INVALID = "INVALID"

	// logical pool type
	PAGEFILE_TYPE     = "PAGEFILE"
	APPENDFILE_TYPE   = "APPENDFILE"
	APPENDECFILE_TYPE = "APPENDECFILE"

	// logical pool allocate status
	ALLOW_STATUS = "ALLOW"
	DENY_STATUS  = "DENY"

	// chunkserver status
	READWRITE_STATUS = "READWRITE"
	PENDDING_STATUS  = "PENDDING"
	RETIRED_STATUS   = "RETIRED"

	// chunkserver disk status
	DISKNORMAL_STATUS = "DISKNORMAL"
	DISKERROR_STATUS  = "DISKERROR"

	// chunkserver online status
	ONLINE_STATUS   = "ONLINE"
	OFFLINE_STATUS  = "OFFLINE"
	UNSTABLE_STATUS = "UNSTABLE"

	// apis
	LIST_PHYSICAL_POOL_FUNC          = "ListPhysicalPool"
	LIST_LOGICAL_POOL_FUNC           = "ListLogicalPool"
	LIST_POOL_ZONE_FUNC              = "ListPoolZone"
	LIST_ZONE_SERVER_FUNC            = "ListZoneServer"
	LIST_CHUNKSERVER_FUNC            = "ListChunkServer"
	GET_CHUNKSERVER_IN_CLUSTER_FUNC  = "GetChunkServerInCluster"
	GET_COPYSET_IN_CHUNKSERVER_FUNC  = "GetCopySetsInChunkServer"
	GET_CHUNKSERVER_LIST_IN_COPYSETS = "GetChunkServerListInCopySets"
	GET_COPYSETS_IN_CLUSTER          = "GetCopySetsInCluster"
	GET_LOGICAL_POOL                 = "GetLogicalPool"

	//http path
	LIST_PHYSICAL_POOL_FUNC_HTTP          = "host.list"
	LIST_LOGICAL_POOL_FUNC_HTTP           = "ListLogicalPool"
	LIST_POOL_ZONE_FUNC_HTTP              = "ListPoolZone"
	LIST_ZONE_SERVER_FUNC_HTTP            = "ListZoneServer"
	LIST_CHUNKSERVER_FUNC_HTTP            = "ListChunkServer"
	GET_CHUNKSERVER_IN_CLUSTER_FUNC_HTTP  = "GetChunkServerInCluster"
	GET_COPYSET_IN_CHUNKSERVER_FUNC_HTTP  = "GetCopySetsInChunkServer"
	GET_CHUNKSERVER_LIST_IN_COPYSETS_HTTP = "GetChunkServerListInCopySets"
	GET_COPYSETS_IN_CLUSTER_HTTP          = "GetCopySetsInCluster"
	GET_LOGICAL_POOL_HTTP                 = "GetLogicalPool"
)

type PhysicalPool struct {
	Id   uint32 `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc"`
}

type LogicalPool struct {
	Id             uint32 `json:"id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	PhysicalPoolId uint32 `json:"physicalPoolId" binding:"required"`
	Type           string `json:"type" binding:"required"`
	CreateTime     string `json:"createTime" binding:"required"`
	AllocateStatus string `json:"allocateStatus" binding:"required"`
	ScanEnable     bool   `json:"scanEnable"`
}

type Zone struct {
	Id               uint32 `json:"id" binding:"required"`
	Name             string `json:"name" binding:"required"`
	PhysicalPoolId   uint32 `json:"physicalPoolId" binding:"required"`
	PhysicalPoolName string `json:"physicalName" binding:"required"`
	Desc             string `json:"desc"`
}

type Server struct {
	Id               uint32 `json:"id" binding:"required"`
	HostName         string `json:"hostName" binding:"required"`
	InternalIp       string `json:"internalIp" binding:"required"`
	InternalPort     uint32 `json:"internalPort" binding:"required"`
	ExternalIp       string `json:"externalIp" binding:"required"`
	ExternalPort     uint32 `json:"externalPort" binding:"required"`
	ZoneId           uint32 `json:"zoneId" binding:"required"`
	ZoneName         string `json:"zoneName" binding:"required"`
	PhysicalPoolId   uint32 `json:"physicalPoolId" binding:"required"`
	PhysicalPoolName string `json:"physicalName" binding:"required"`
	Desc             string `json:"desc"`
}

type ChunkServer struct {
	Id           uint32 `json:"id" binding:"required"`
	DiskType     string `json:"diskType" binding:"required"`
	HostIp       string `json:"hostIp" binding:"required"`
	Port         uint32 `json:"port" binding:"required"`
	Status       string `json:"status" binding:"required"`
	DiskStatus   string `json:"diskStatus" binding:"required"`
	OnlineStatus string `json:"onlineStatus" binding:"required"`
	MountPoint   string `json:"mountPoint" binding:"required"`
	DiskCapacity string `json:"diskCapacity" binding:"required"`
	DiskUsed     string `json:"diskUsed" binding:"required"`
	ExternalIp   string `json:"externalIp"`
}

type CopySetInfo struct {
	LogicalPoolId      uint32 `json:"logicalPoolId" binding:"required"`
	CopysetId          uint32 `json:"copysetId" binding:"required"`
	Scanning           bool   `json:"scanning"`
	LastScanSec        uint64 `json:"lastScanSec"`
	LastScanConsistent bool   `json:"lastScanConsistent"`
}

type ChunkServerLocation struct {
	ChunkServerId uint32 `json:"chunkServerId" binding:"required"`
	HostIp        string `json:"hostIp" binding:"required"`
	Port          uint32 `josn:"port" binding:"required"`
	ExternalIp    string `json:"externalIp"`
}
type CopySetServerInfo struct {
	CopysetId uint32                `json:"copysetId" binding:"required"`
	CsLocs    []ChunkServerLocation `json:"csLocs" binding:"required"`
}

func (cli *MdsClient) ListPhysicalPool_http() ([]PhysicalPool, error) {
	var host = cli.addrs
	var path = LIST_PHYSICAL_POOL_FUNC_HTTP
	ret := cli.baseClient_http.SendHTTP(host, path)
	if ret.Err != nil {
		return nil, ret.Err
	}
	response := ret.Result.(*topology.ListPhysicalPoolResponse)
	statusCode := response.GetStatusCode()
	if statusCode != int32(statuscode.TopoStatusCode_Success) {
		return nil, fmt.Errorf(statuscode.TopoStatusCode_name[statusCode])
	}

	var infos []PhysicalPool
	for _, pool := range response.GetPhysicalPoolInfos() {
		info := PhysicalPool{}
		info.Id = pool.GetPhysicalPoolID()
		info.Name = pool.GetPhysicalPoolName()
		info.Desc = pool.GetDesc()
		infos = append(infos, info)
	}
	return infos, nil
}

func getLogicalPoolType(t topology.LogicalPoolType) string {
	switch t {
	case topology.LogicalPoolType_PAGEFILE:
		return PAGEFILE_TYPE
	case topology.LogicalPoolType_APPENDFILE:
		return APPENDFILE_TYPE
	case topology.LogicalPoolType_APPENDECFILE:
		return APPENDECFILE_TYPE
	default:
		return INVALID
	}
}

func getLogicalPoolAllocateStatus(s topology.AllocateStatus) string {
	switch s {
	case topology.AllocateStatus_ALLOW:
		return ALLOW_STATUS
	case topology.AllocateStatus_DENY:
		return DENY_STATUS
	default:
		return INVALID
	}
}

func (cli *MdsClient) ListLogicalPool_http() ([]LogicalPool, error) {
	// list physical pool and get pool id
	physicalPools, err := cli.ListPhysicalPool_http()
	if err != nil {
		return nil, err
	}
	size := len(physicalPools)
	results := make(chan baserpc.RpcResult, size)
	for _, pool := range physicalPools {
		go func(id uint32) {
			var host = cli.addrs
			var path = LIST_LOGICAL_POOL_FUNC_HTTP
			path = fmt.Sprintf("%s %s %s", path, "PhysicalPoolId=", &id)

			ret := cli.baseClient_http.SendHTTP(host, path)
			if ret.Err != nil {
				results <- baserpc.RpcResult{
					Key:    id,
					Err:    fmt.Errorf("%s: %v", ret.Key, ret.Err),
					Result: nil,
				}
			} else {
				response := ret.Result.(*topology.ListLogicalPoolResponse)
				statusCode := response.GetStatusCode()
				if statusCode != int32(statuscode.TopoStatusCode_Success) {
					results <- baserpc.RpcResult{
						Key:    id,
						Err:    fmt.Errorf("%s", statuscode.TopoStatusCode_name[statusCode]),
						Result: nil,
					}
				} else {
					var pools []LogicalPool
					for _, pool := range response.GetLogicalPoolInfos() {
						info := LogicalPool{}
						info.Id = pool.GetLogicalPoolID()
						info.Name = pool.GetLogicalPoolName()
						info.PhysicalPoolId = pool.GetPhysicalPoolID()
						info.Type = getLogicalPoolType(pool.GetType())
						info.CreateTime = time.Unix(int64(pool.GetCreateTime()), 0).Format(common.TIME_FORMAT)
						info.AllocateStatus = getLogicalPoolAllocateStatus(pool.GetAllocateStatus())
						info.ScanEnable = pool.GetScanEnable()
						pools = append(pools, info)
					}
					results <- baserpc.RpcResult{
						Key:    id,
						Err:    nil,
						Result: &pools,
					}
				}
			}
		}(pool.Id)
	}

	pools := []LogicalPool{}
	count := 0
	for res := range results {
		if res.Err != nil {
			return nil, fmt.Errorf("physical pool id: %d; %v", res.Key, res.Err)
		}
		pools = append(pools, (*res.Result.(*[]LogicalPool))...)
		count++
		if count >= size {
			break
		}
	}
	return pools, nil
}

func (cli *MdsClient) GetLogicalPool_http(poolId uint32) (LogicalPool, error) {
	info := LogicalPool{}
	var host = cli.addrs
	var path = GET_LOGICAL_POOL_HTTP
	path = fmt.Sprintf("%s %s %s", path, "LogicPoolId=", &poolId)

	ret := cli.baseClient_http.SendHTTP(host, path)

	if ret.Err != nil {
		return info, ret.Err
	}
	response := ret.Result.(*topology.GetLogicalPoolResponse)
	statusCode := response.GetStatusCode()
	if statusCode != int32(statuscode.TopoStatusCode_Success) {
		return info, fmt.Errorf(statuscode.TopoStatusCode_name[statusCode])
	}
	pool := response.GetLogicalPoolInfo()
	info.Id = pool.GetLogicalPoolID()
	info.Name = pool.GetLogicalPoolName()
	info.PhysicalPoolId = pool.GetPhysicalPoolID()
	info.Type = getLogicalPoolType(pool.GetType())
	info.CreateTime = time.Unix(int64(pool.GetCreateTime()), 0).Format(common.TIME_FORMAT)
	info.AllocateStatus = getLogicalPoolAllocateStatus(pool.GetAllocateStatus())
	info.ScanEnable = pool.GetScanEnable()
	return info, nil
}

// list zones of physical pool

func (cli *MdsClient) ListPoolZone_http(poolId uint32) ([]Zone, error) {
	var host = cli.addrs
	var path = LIST_POOL_ZONE_FUNC_HTTP
	path = fmt.Sprintf("%s %s %s", path, "PhysicalPoolId=", &poolId)

	ret := cli.baseClient_http.SendHTTP(host, path)

	if ret.Err != nil {
		return nil, ret.Err
	}
	response := ret.Result.(*topology.ListPoolZoneResponse)
	statusCode := response.GetStatusCode()
	if statusCode != int32(statuscode.TopoStatusCode_Success) {
		return nil, fmt.Errorf(statuscode.TopoStatusCode_name[statusCode])
	}

	infos := []Zone{}
	for _, zone := range response.GetZones() {
		info := Zone{}
		info.Id = zone.GetZoneID()
		info.Name = zone.GetZoneName()
		info.PhysicalPoolId = zone.GetPhysicalPoolID()
		info.PhysicalPoolName = zone.GetPhysicalPoolName()
		info.Desc = zone.GetDesc()
		infos = append(infos, info)
	}
	return infos, nil
}

// list servers of zone

func (cli *MdsClient) ListZoneServer_http(zoneId uint32) ([]Server, error) {

	var host = cli.addrs
	var path = LIST_ZONE_SERVER_FUNC_HTTP
	path = fmt.Sprintf("%s %s %s", path, "ZoneId=", &zoneId)

	ret := cli.baseClient_http.SendHTTP(host, path)
	if ret.Err != nil {
		return nil, ret.Err
	}
	response := ret.Result.(*topology.ListZoneServerResponse)
	statusCode := response.GetStatusCode()
	if statusCode != int32(statuscode.TopoStatusCode_Success) {
		return nil, fmt.Errorf(statuscode.TopoStatusCode_name[statusCode])
	}

	infos := []Server{}
	for _, server := range response.GetServerInfo() {
		info := Server{}
		info.Id = server.GetServerID()
		info.HostName = server.GetHostName()
		info.InternalIp = server.GetInternalIp()
		info.InternalPort = server.GetInternalPort()
		info.ExternalIp = server.GetExternalIp()
		info.ExternalPort = server.GetExternalPort()
		info.ZoneId = server.GetZoneID()
		info.ZoneName = server.GetZoneName()
		info.PhysicalPoolId = server.GetPhysicalPoolID()
		info.PhysicalPoolName = server.GetPhysicalPoolName()
		info.Desc = server.GetDesc()
		infos = append(infos, info)
	}
	return infos, nil
}

// list chunkservers of server
func getChunkServerStatus(s topology.ChunkServerStatus) string {
	switch s {
	case topology.ChunkServerStatus_READWRITE:
		return READWRITE_STATUS
	case topology.ChunkServerStatus_PENDDING:
		return PENDDING_STATUS
	case topology.ChunkServerStatus_RETIRED:
		return RETIRED_STATUS
	default:
		return INVALID
	}
}

func getDiskStatus(s topology.DiskState) string {
	switch s {
	case topology.DiskState_DISKNORMAL:
		return DISKNORMAL_STATUS
	case topology.DiskState_DISKERROR:
		return DISKERROR_STATUS
	default:
		return INVALID
	}
}

func getOnlineStatus(s topology.OnlineState) string {
	switch s {
	case topology.OnlineState_ONLINE:
		return ONLINE_STATUS
	case topology.OnlineState_OFFLINE:
		return OFFLINE_STATUS
	case topology.OnlineState_UNSTABLE:
		return UNSTABLE_STATUS
	default:
		return INVALID
	}
}

func (cli *MdsClient) ListChunkServer_http(serverId uint32) ([]ChunkServer, error) {
	var host = cli.addrs
	var path = LIST_CHUNKSERVER_FUNC_HTTP
	path = fmt.Sprintf("%s %s %s", path, "ServerId=", &serverId)
	ret := cli.baseClient_http.SendHTTP(host, path)

	if ret.Err != nil {
		return nil, ret.Err
	}
	response := ret.Result.(*topology.ListChunkServerResponse)
	statusCode := response.GetStatusCode()
	if statusCode != int32(statuscode.TopoStatusCode_Success) {
		return nil, fmt.Errorf(statuscode.TopoStatusCode_name[statusCode])
	}

	infos := []ChunkServer{}
	for _, cs := range response.GetChunkServerInfos() {
		if cs.GetStatus() == topology.ChunkServerStatus_RETIRED {
			continue
		}
		info := ChunkServer{}
		info.Id = cs.GetChunkServerID()
		info.DiskType = cs.GetDiskType()
		info.HostIp = cs.GetHostIp()
		info.Port = cs.GetPort()
		info.Status = getChunkServerStatus(cs.GetStatus())
		info.DiskStatus = getDiskStatus(cs.GetDiskStatus())
		info.OnlineStatus = getOnlineStatus(cs.GetOnlineState())
		info.MountPoint = cs.GetMountPoint()
		info.DiskCapacity = strconv.FormatUint(cs.GetDiskCapacity()/common.GiB, 10)
		info.DiskUsed = strconv.FormatUint(cs.GetDiskUsed()/common.GiB, 10)
		info.ExternalIp = cs.GetExternalIp()
		infos = append(infos, info)
	}
	return infos, nil
}

func (cli *MdsClient) GetChunkServerInCluster_http() ([]ChunkServer, error) {
	var host = cli.addrs
	var path = GET_CHUNKSERVER_IN_CLUSTER_FUNC_HTTP
	ret := cli.baseClient_http.SendHTTP(host, path)
	if ret.Err != nil {
		return nil, ret.Err
	}
	response := ret.Result.(*topology.GetChunkServerInClusterResponse)
	statusCode := response.GetStatusCode()
	if statusCode != int32(statuscode.TopoStatusCode_Success) {
		return nil, fmt.Errorf(statuscode.TopoStatusCode_name[statusCode])
	}
	infos := []ChunkServer{}
	for _, cs := range response.GetChunkServerInfos() {
		if cs.GetStatus() == topology.ChunkServerStatus_RETIRED {
			continue
		}
		info := ChunkServer{}
		info.Id = cs.GetChunkServerID()
		info.DiskType = cs.GetDiskType()
		info.HostIp = cs.GetHostIp()
		info.Port = cs.GetPort()
		info.Status = getChunkServerStatus(cs.GetStatus())
		info.DiskStatus = getDiskStatus(cs.GetDiskStatus())
		info.OnlineStatus = getOnlineStatus(cs.GetOnlineState())
		info.MountPoint = cs.GetMountPoint()
		info.DiskCapacity = strconv.FormatUint(cs.GetDiskCapacity()/common.GiB, 10)
		info.DiskUsed = strconv.FormatUint(cs.GetDiskUsed()/common.GiB, 10)
		info.ExternalIp = cs.GetExternalIp()
		infos = append(infos, info)
	}
	return infos, nil
}

func (cli *MdsClient) GetCopySetsInChunkServer_http(ip string, port uint32) ([]CopySetInfo, error) {
	var host = cli.addrs
	var path = GET_COPYSET_IN_CHUNKSERVER_FUNC_HTTP

	ret := cli.baseClient_http.SendHTTP(host, path)
	if ret.Err != nil {
		return nil, ret.Err
	}
	response := ret.Result.(*topology.GetCopySetsInChunkServerResponse)
	statusCode := response.GetStatusCode()
	if statusCode != int32(statuscode.TopoStatusCode_Success) {
		return nil, fmt.Errorf(statuscode.TopoStatusCode_name[statusCode])
	}
	infos := []CopySetInfo{}
	for _, cs := range response.GetCopysetInfos() {
		info := CopySetInfo{}
		info.LogicalPoolId = cs.GetLogicalPoolId()
		info.CopysetId = cs.GetCopysetId()
		info.Scanning = cs.GetScaning()
		info.LastScanSec = cs.GetLastScanSec()
		info.LastScanConsistent = cs.GetLastScanConsistent()
		infos = append(infos, info)
	}
	return infos, nil
}

func (cli *MdsClient) GetChunkServerListInCopySets_http(logicalPoolId uint32, copysetIds []uint32) ([]CopySetServerInfo, error) {
	var host = cli.addrs
	var path = GET_CHUNKSERVER_LIST_IN_COPYSETS_HTTP
	path = fmt.Sprintf("%s %s %s", path, "LogicPoolId=", &logicalPoolId)
	ret := cli.baseClient_http.SendHTTP(host, path)
	if ret.Err != nil {
		return nil, ret.Err
	}
	response := ret.Result.(*topology.GetChunkServerListInCopySetsResponse)
	statusCode := response.GetStatusCode()
	if statusCode != int32(statuscode.TopoStatusCode_Success) {
		return nil, fmt.Errorf(statuscode.TopoStatusCode_name[statusCode])
	}
	infos := []CopySetServerInfo{}
	for _, csInfo := range response.GetCsInfo() {
		info := CopySetServerInfo{}
		info.CopysetId = csInfo.GetCopysetId()
		for _, locs := range csInfo.GetCsLocs() {
			var l ChunkServerLocation
			l.ChunkServerId = locs.GetChunkServerID()
			l.HostIp = locs.GetHostIp()
			l.Port = locs.GetPort()
			l.ExternalIp = locs.GetExternalIp()
			info.CsLocs = append(info.CsLocs, l)
		}
		infos = append(infos, info)
	}
	return infos, nil
}

func (cli *MdsClient) GetCopySetsInCluster_http() ([]CopySetInfo, error) {
	var host = cli.addrs
	var path = LIST_PHYSICAL_POOL_FUNC_HTTP
	ret := cli.baseClient_http.SendHTTP(host, path)
	if ret.Err != nil {
		return nil, ret.Err
	}
	response := ret.Result.(*topology.GetCopySetsInClusterResponse)
	statusCode := response.GetStatusCode()
	if statusCode != int32(statuscode.TopoStatusCode_Success) {
		return nil, fmt.Errorf(statuscode.TopoStatusCode_name[statusCode])
	}
	infos := []CopySetInfo{}
	for _, csInfo := range response.GetCopysetInfos() {
		info := CopySetInfo{}
		info.CopysetId = csInfo.GetCopysetId()
		info.LastScanConsistent = csInfo.GetLastScanConsistent()
		info.LastScanSec = csInfo.GetLastScanSec()
		info.LogicalPoolId = csInfo.GetLogicalPoolId()
		info.Scanning = csInfo.GetScaning()
		infos = append(infos, info)
	}
	return infos, nil
}
