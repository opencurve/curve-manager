package curvebs

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/proto/nameserver2"
	"github.com/opencurve/curve-manager/internal/proto/topology"
	"github.com/opencurve/curve-manager/internal/proto/topology/statuscode"
	"github.com/opencurve/curve-manager/internal/rpc/baserpc"
	"github.com/opencurve/pigeon"
)

var (
	GMdsClient *mdsClient
)

type RpcResult common.QueryResult

const (
	CURVEBS_MDS_ADDRESS = "mds.address"

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

	// file type
	INODE_DIRECTORY         = "INODE_DIRECTORY"
	INODE_PAGEFILE          = "INODE_PAGEFILE"
	INODE_APPENDFILE        = "INODE_APPENDFILE"
	INODE_APPENDECFILE      = "INODE_APPENDECFILE"
	INODE_SNAPSHOT_PAGEFILE = "INODE_SNAPSHOT_PAGEFILE"

	// file status
	FILE_CREATED             = "kFileCreated"
	FILE_DELETING            = "kFileDeleting"
	FILE_CLONING             = "kFileCloning"
	FILE_CLONEMETA_INSTALLED = "kFileCloneMetaInstalled"
	FILE_CLONED              = "kFileCloned"
	FILE_BEIING_CLONED       = "kFileBeingCloned"

	// throttle type
	IOPS_TOTAL = "IOPS_TOTAL"
	IOPS_READ  = "IOPS_READ"
	IOPS_WRITE = "IOPS_WRITE"
	BPS_TOTAL  = "BPS_TOTAL"
	BPS_READ   = "BPS_READ"
	BPS_WRITE  = "BPS_WRITE"

	LIST_PHYSICAL_POOL_FUNC         = "ListPhysicalPool"
	LIST_LOGICAL_POOL_FUNC          = "ListLogicalPool"
	LIST_POOL_ZONE_FUNC             = "ListPoolZone"
	LIST_ZONE_SERVER_FUNC           = "ListZoneServer"
	LIST_CHUNKSERVER_FUNC           = "ListChunkServer"
	GET_CHUNKSERVER_IN_CLUSTER_FUNC = "GetChunkServerInCluster"
	GET_FILE_ALLOC_SIZE_FUNC        = "GetAllocatedSize"
	LIST_DIR_FUNC                   = "ListDir"
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
		addrs: strings.Split(addrs, common.CURVEBS_ADDRESS_DELIMITER),
	}
	return nil
}

// list physical pool
func (cli *mdsClient) ListPhysicalPool() ([]PhysicalPool, error) {
	Rpc := &ListPhysicalPoolRpc{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, LIST_PHYSICAL_POOL_FUNC)
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

func (cli *mdsClient) ListLogicalPool() ([]LogicalPool, error) {
	// list physical pool and get pool id
	physicalPools, err := cli.ListPhysicalPool()
	if err != nil {
		return nil, err
	}
	size := len(physicalPools)
	results := make(chan RpcResult, size)
	for _, pool := range physicalPools {
		go func(id uint32) {
			Rpc := &ListLogicalPoolRpc{}
			Rpc.ctx = baserpc.NewRpcContext(cli.addrs, LIST_LOGICAL_POOL_FUNC)
			Rpc.Request = &topology.ListLogicalPoolRequest{
				PhysicalPoolID: &id,
			}
			ret := baserpc.GBaseClient.SendRpc(Rpc.ctx, Rpc)
			if ret.Err != nil {
				results <- RpcResult{
					Key:    id,
					Err:    fmt.Errorf("%s: %v", ret.Key, ret.Err),
					Result: nil,
				}
			} else {
				response := ret.Result.(*topology.ListLogicalPoolResponse)
				statusCode := response.GetStatusCode()
				if statusCode != int32(statuscode.TopoStatusCode_Success) {
					results <- RpcResult{
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
					results <- RpcResult{
						Key:    id,
						Err:    nil,
						Result: &pools,
					}
				}
			}
		}(pool.Id)
	}

	var pools []LogicalPool
	count := 0
	for res := range results {
		if res.Err != nil {
			return nil, fmt.Errorf("physical pool id: %d; %v", res.Key, res.Err)
		}
		pools = append(pools, (*res.Result.(*[]LogicalPool))...)
		count += 1
		if count >= size {
			break
		}
	}
	return pools, nil
}

// list zones of physical pool
func (cli *mdsClient) ListPoolZone(poolId uint32) ([]Zone, error) {
	Rpc := &ListPoolZonesRpc{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, LIST_POOL_ZONE_FUNC)
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

	var infos []Zone
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
func (cli *mdsClient) ListZoneServer(zoneId uint32) ([]Server, error) {
	Rpc := &ListZoneServer{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, LIST_ZONE_SERVER_FUNC)
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

	var infos []Server
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

func (cli *mdsClient) ListChunkServer(serverId uint32) ([]ChunkServer, error) {
	Rpc := &ListChunkServer{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, LIST_CHUNKSERVER_FUNC)
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

	var infos []ChunkServer
	for _, cs := range response.GetChunkServerInfos() {
		info := ChunkServer{}
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

func (cli *mdsClient) GetChunkServerInCluster() ([]ChunkServer, error) {
	Rpc := &GetChunkServerInCluster{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, GET_CHUNKSERVER_IN_CLUSTER_FUNC)
	Rpc.Request = &topology.GetChunkServerInClusterRequest{}
	ret := baserpc.GBaseClient.SendRpc(Rpc.ctx, Rpc)
	if ret.Err != nil {
		return nil, ret.Err
	}

	response := ret.Result.(*topology.GetChunkServerInClusterResponse)
	statusCode := response.GetStatusCode()
	if statusCode != int32(statuscode.TopoStatusCode_Success) {
		return nil, fmt.Errorf(statuscode.TopoStatusCode_name[statusCode])
	}
	var infos []ChunkServer
	for _, cs := range response.GetChunkServerInfos() {
		info := ChunkServer{}
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

func (cli *mdsClient) GetFileAllocatedSize(filename string) (uint64, *map[uint32]uint64, error) {
	Rpc := &GetFileAllocatedSize{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, GET_FILE_ALLOC_SIZE_FUNC)
	Rpc.Request = &nameserver2.GetAllocatedSizeRequest{
		FileName: &filename,
	}
	ret := baserpc.GBaseClient.SendRpc(Rpc.ctx, Rpc)
	if ret.Err != nil {
		return 0, nil, ret.Err
	}

	response := ret.Result.(*nameserver2.GetAllocatedSizeResponse)
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return 0, nil, fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	infos := make(map[uint32]uint64)
	for k, v := range response.GetAllocSizeMap() {
		infos[k] = v
	}
	return response.GetAllocatedSize(), &infos, nil
}

func getFileType(t nameserver2.FileType) string {
	switch t {
	case nameserver2.FileType_INODE_DIRECTORY:
		return INODE_DIRECTORY
	case nameserver2.FileType_INODE_PAGEFILE:
		return INODE_PAGEFILE
	case nameserver2.FileType_INODE_APPENDFILE:
		return INODE_APPENDFILE
	case nameserver2.FileType_INODE_APPENDECFILE:
		return INODE_APPENDECFILE
	case nameserver2.FileType_INODE_SNAPSHOT_PAGEFILE:
		return INODE_SNAPSHOT_PAGEFILE
	default:
		return INVALID
	}
}

func getFileStatus(s nameserver2.FileStatus) string {
	switch s {
	case nameserver2.FileStatus_kFileCreated:
		return FILE_CREATED
	case nameserver2.FileStatus_kFileDeleting:
		return FILE_DELETING
	case nameserver2.FileStatus_kFileCloning:
		return FILE_CLONING
	case nameserver2.FileStatus_kFileCloneMetaInstalled:
		return FILE_CLONEMETA_INSTALLED
	case nameserver2.FileStatus_kFileCloned:
		return FILE_CLONED
	case nameserver2.FileStatus_kFileBeingCloned:
		return FILE_BEIING_CLONED
	default:
		return INVALID
	}
}

func getThrottleType(t nameserver2.ThrottleType) string {
	switch t {
	case nameserver2.ThrottleType_IOPS_TOTAL:
		return IOPS_TOTAL
	case nameserver2.ThrottleType_IOPS_READ:
		return IOPS_READ
	case nameserver2.ThrottleType_IOPS_WRITE:
		return IOPS_WRITE
	case nameserver2.ThrottleType_BPS_TOTAL:
		return BPS_TOTAL
	case nameserver2.ThrottleType_BPS_READ:
		return BPS_READ
	case nameserver2.ThrottleType_BPS_WRITE:
		return BPS_WRITE
	default:
		return INVALID
	}
}

func (cli *mdsClient) ListDir(filename, owner, sig string, date uint64) ([]FileInfo, error) {
	Rpc := &ListDir{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, LIST_DIR_FUNC)
	Rpc.Request = &nameserver2.ListDirRequest{
		FileName:  &filename,
		Owner:     &owner,
		Date:      &date,
		Signature: &sig,
	}
	ret := baserpc.GBaseClient.SendRpc(Rpc.ctx, Rpc)
	if ret.Err != nil {
		return nil, ret.Err
	}

	response := ret.Result.(*nameserver2.ListDirResponse)
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return nil, fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	var infos []FileInfo
	for _, v := range response.GetFileInfo() {
		var info FileInfo
		info.Id = v.GetId()
		info.FileName = v.GetFileName()
		info.ParentId = v.GetParentId()
		info.FileType = getFileType(v.GetFileType())
		info.Owner = v.GetOwner()
		info.ChunkSize = v.GetChunkSize()
		info.SegmentSize = v.GetSegmentSize()
		info.Length = v.GetLength() / common.GB
		info.Ctime = time.Unix(int64(v.GetCtime()/1000000), 0).Format(common.TIME_FORMAT)
		info.SeqNum = v.GetSeqNum()
		info.FileStatus = getFileStatus(v.GetFileStatus())
		info.OriginalFullPathName = v.GetOriginalFullPathName()
		info.CloneSource = v.GetCloneSource()
		info.CloneLength = v.GetCloneLength()
		info.StripeUnit = v.GetStripeUnit()
		info.StripeCount = v.GetStripeCount()
		for _, p := range v.GetThrottleParams().GetThrottleParams() {
			var param ThrottleParams
			param.Type = getThrottleType(p.GetType())
			param.Limit = p.GetLimit()
			param.Burst = p.GetBurst()
			param.BurstLength = p.GetBurstLength()
			info.ThrottleParams = append(info.ThrottleParams, param)
		}
		info.Epoch = v.GetEpoch()
		infos = append(infos, info)
	}
	return infos, nil
}
