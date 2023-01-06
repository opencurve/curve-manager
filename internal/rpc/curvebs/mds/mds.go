package mds

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/opencurve/curve-manager/internal/proto/topology"
	"github.com/opencurve/curve-manager/internal/proto/topology/statuscode"
	"github.com/opencurve/curve-manager/internal/rpc/baserpc"
	comm "github.com/opencurve/curve-manager/internal/rpc/curvebs/common"
	"github.com/opencurve/pigeon"
)

var (
	GMdsClient *mdsClient
)

const (
	CURVEBS_MDS_ADDRESS           = "mds.address"
	CURVEBS_MDS_ADDRESS_DELIMITER = ","

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

	LIST_PHYSICAL_POOL_NAME = "ListPhysicalPool"
	LIST_LOGICAL_POOL_NAME  = "ListLogicalPool"
	LIST_POOL_ZONE_NAME     = "ListPoolZone"
	LIST_ZONE_SERVER_NAME   = "ListZoneServer"
	LIST_CHUNKSERVER_NAME   = "ListChunkServer"
)

type mdsClient struct {
	addrs []string
}

func Init(cfg *pigeon.Configure) error {
	addrs := cfg.GetConfig().GetString(CURVEBS_MDS_ADDRESS)
	if len(addrs) == 0 {
		return fmt.Errorf("no cluster mds address found")
	}
	GMdsClient = &mdsClient{
		addrs: strings.Split(addrs, CURVEBS_MDS_ADDRESS_DELIMITER),
	}
	return nil
}

// list physical pool
func (cli *mdsClient) ListPhysicalPool() (interface{}, error) {
	Rpc := &ListPhysicalPoolRpc{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, LIST_PHYSICAL_POOL_NAME)
	Rpc.Request = &topology.ListPhysicalPoolRequest{}
	ret := baserpc.GBaseClient.SendRpc(Rpc.ctx, Rpc)
	if ret.Err != nil {
		return nil, ret.Err
	}

	response := ret.Result.(*topology.ListPhysicalPoolResponse)
	statusCode := response.GetStatusCode()
	if statusCode != int32(statuscode.TopoStatusCode_Success) {
		return nil, fmt.Errorf(statuscode.TopoStatusCode_name[statusCode])
	}

	var infos []comm.PhysicalPool
	for _, pool := range response.GetPhysicalPoolInfos() {
		info := comm.PhysicalPool{}
		info.Id = pool.GetPhysicalPoolID()
		info.Name = pool.GetPhysicalPoolName()
		info.Desc = pool.GetDesc()
		infos = append(infos, info)
	}
	return infos, nil
}

// list logical pool
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

func (cli *mdsClient) ListLogicalPool() (interface{}, error) {
	// list physical pool and get pool id
	physicalPools, err := cli.ListPhysicalPool()
	if err != nil {
		return nil, err
	}
	size := len(physicalPools.([]comm.PhysicalPool))
	results := make(chan baserpc.RpcResult, size)
	for _, pool := range physicalPools.([]comm.PhysicalPool) {
		go func(id uint32) {
			Rpc := &ListLogicalPoolRpc{}
			Rpc.ctx = baserpc.NewRpcContext(cli.addrs, LIST_LOGICAL_POOL_NAME)
			Rpc.Request = &topology.ListLogicalPoolRequest{
				PhysicalPoolID: &id,
			}
			ret := baserpc.GBaseClient.SendRpc(Rpc.ctx, Rpc)
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
					var pools []comm.LogicalPool
					for _, pool := range response.GetLogicalPoolInfos() {
						info := comm.LogicalPool{}
						info.Id = pool.GetLogicalPoolID()
						info.Name = pool.GetLogicalPoolName()
						info.PhysicalPoolId = pool.GetPhysicalPoolID()
						info.Type = getLogicalPoolType(pool.GetType())
						info.CreateTime = time.Unix(int64(pool.GetCreateTime()), 0).UTC().String()
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

	var pools []comm.LogicalPool
	count := 0
	for res := range results {
		if res.Err != nil {
			return nil, fmt.Errorf("physical pool id: %d; %v", res.Key, res.Err)
		}
		pools = append(pools, (*res.Result.(*[]comm.LogicalPool))...)
		count = count + 1
		if count >= size {
			break
		}
	}
	return pools, nil
}

// list zones of physical pool
func (cli *mdsClient) ListPoolZone(poolId uint32) (interface{}, error) {
	Rpc := &ListPoolZonesRpc{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, LIST_POOL_ZONE_NAME)
	Rpc.Request = &topology.ListPoolZoneRequest{
		PhysicalPoolID: &poolId,
	}
	ret := baserpc.GBaseClient.SendRpc(Rpc.ctx, Rpc)
	if ret.Err != nil {
		return nil, ret.Err
	}

	response := ret.Result.(*topology.ListPoolZoneResponse)
	statusCode := response.GetStatusCode()
	if statusCode != int32(statuscode.TopoStatusCode_Success) {
		return nil, fmt.Errorf(statuscode.TopoStatusCode_name[statusCode])
	}

	var infos []comm.Zone
	for _, zone := range response.GetZones() {
		info := comm.Zone{}
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
func (cli *mdsClient) ListZoneServer(zoneId uint32) (interface{}, error) {
	Rpc := &ListZoneServer{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, LIST_ZONE_SERVER_NAME)
	Rpc.Request = &topology.ListZoneServerRequest{
		ZoneID: &zoneId,
	}
	ret := baserpc.GBaseClient.SendRpc(Rpc.ctx, Rpc)
	if ret.Err != nil {
		return nil, ret.Err
	}

	response := ret.Result.(*topology.ListZoneServerResponse)
	statusCode := response.GetStatusCode()
	if statusCode != int32(statuscode.TopoStatusCode_Success) {
		return nil, fmt.Errorf(statuscode.TopoStatusCode_name[statusCode])
	}

	var infos []comm.Server
	for _, server := range response.GetServerInfo() {
		info := comm.Server{}
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
		return READWRITE_STATUS
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

func (cli *mdsClient) ListChunkServer(serverId uint32) (interface{}, error) {
	Rpc := &ListChunkServer{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, LIST_CHUNKSERVER_NAME)
	Rpc.Request = &topology.ListChunkServerRequest{
		ServerID: &serverId,
	}
	ret := baserpc.GBaseClient.SendRpc(Rpc.ctx, Rpc)
	if ret.Err != nil {
		return nil, ret.Err
	}

	response := ret.Result.(*topology.ListChunkServerResponse)
	statusCode := response.GetStatusCode()
	if statusCode != int32(statuscode.TopoStatusCode_Success) {
		return nil, fmt.Errorf(statuscode.TopoStatusCode_name[statusCode])
	}

	var infos []comm.ChunkServer
	for _, cs := range response.GetChunkServerInfos() {
		info := comm.ChunkServer{}
		info.Id = cs.GetChunkServerID()
		info.DiskType = cs.GetDiskType()
		info.HostIp = cs.GetHostIp()
		info.Port = cs.GetPort()
		info.Status = getChunkServerStatus(cs.GetStatus())
		info.DiskStatus = getDiskStatus(cs.GetDiskStatus())
		info.OnlineStatus = getOnlineStatus(cs.GetOnlineState())
		info.MountPoint = cs.GetMountPoint()
		info.DiskCapacity = strconv.FormatUint(cs.GetDiskCapacity(), 10)
		info.DiskUsed = strconv.FormatUint(cs.GetDiskUsed(), 10)
		info.ExternalIp = cs.GetExternalIp()
		infos = append(infos, info)
	}
	return infos, nil
}
