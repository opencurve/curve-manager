package mds

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/opencurve/curve-manager/internal/metrics/bsmetric"
	"github.com/opencurve/curve-manager/internal/proto/topology"
	"github.com/opencurve/curve-manager/internal/proto/topology/statuscode"
	"github.com/opencurve/curve-manager/internal/rpc/baserpc"
	comm "github.com/opencurve/curve-manager/internal/rpc/curvebs/common"
	"github.com/opencurve/pigeon"
	"google.golang.org/grpc"
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

	LIST_PHYSICAL_POOL_NAME = "ListPhysicalPool"
	LIST_LOGICAL_POOL_NAME  = "ListLogicalPool"
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
type ListPhysicalPoolRpc struct {
	ctx            *baserpc.RpcContext
	Request        *topology.ListPhysicalPoolRequest
	topologyClient topology.TopologyServiceClient
}

func (lpRpc *ListPhysicalPoolRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	lpRpc.topologyClient = topology.NewTopologyServiceClient(cc)
}

func (lpRpc *ListPhysicalPoolRpc) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return lpRpc.topologyClient.ListPhysicalPool(ctx, lpRpc.Request, opt...)
}

func (cli *mdsClient) listPhysicalPool() (interface{}, error) {
	listPoolRpc := &ListPhysicalPoolRpc{}
	listPoolRpc.ctx = baserpc.NewRpcContext(cli.addrs, LIST_PHYSICAL_POOL_NAME)
	listPoolRpc.Request = &topology.ListPhysicalPoolRequest{}
	ret := baserpc.GBaseClient.SendRpc(listPoolRpc.ctx, listPoolRpc)
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
type ListLogicalPoolRpc struct {
	ctx            *baserpc.RpcContext
	Request        *topology.ListLogicalPoolRequest
	topologyClient topology.TopologyServiceClient
}

func (lpRpc *ListLogicalPoolRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	lpRpc.topologyClient = topology.NewTopologyServiceClient(cc)
}

func (lpRpc *ListLogicalPoolRpc) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return lpRpc.topologyClient.ListLogicalPool(ctx, lpRpc.Request, opt...)
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

func (cli *mdsClient) ListLogicalPool() (interface{}, error) {
	// list physical pool and get pool id
	physicalPools, err := cli.listPhysicalPool()
	if err != nil {
		return nil, err
	}
	size := len(physicalPools.([]comm.PhysicalPool))
	results := make(chan baserpc.RpcResult, size)
	for _, pool := range physicalPools.([]comm.PhysicalPool) {
		go func(id uint32) {
			listPoolRpc := &ListLogicalPoolRpc{}
			listPoolRpc.ctx = baserpc.NewRpcContext(cli.addrs, LIST_LOGICAL_POOL_NAME)
			listPoolRpc.Request = &topology.ListLogicalPoolRequest{
				PhysicalPoolID: &id,
			}
			ret := baserpc.GBaseClient.SendRpc(listPoolRpc.ctx, listPoolRpc)
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
						// get space
						space, err := bsmetric.GetPoolSpace(pool.GetLogicalPoolName())
						if err != nil {
							results <- baserpc.RpcResult{
								Key: id,
								Err: fmt.Errorf("get space failed, poolName: %s, err: %v", pool.GetLogicalPoolName(), err),
								Result: nil,
							}
						} else {
							info := comm.LogicalPool{}
							info.Id = pool.GetLogicalPoolID()
							info.Name = pool.GetLogicalPoolName()
							info.PhysicalPoolId = pool.GetPhysicalPoolID()
							info.Type = getLogicalPoolType(pool.GetType())
							info.CreateTime = time.Unix(int64(pool.GetCreateTime()), 0).UTC().String()
							info.AllocateStatus = getLogicalPoolAllocateStatus(pool.GetAllocateStatus())
							info.ScanEnable = pool.GetScanEnable()
							info.TotalSpace = space.Total
							info.UsedSpace = space.Used
							pools = append(pools, info)
						}
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

// list chunkserver
type ListChunkServerRpc struct {
	ctx            *baserpc.RpcContext
	Request        *topology.ListChunkServerRequest
	topologyClient topology.TopologyServiceClient
}

func (lpRpc *ListChunkServerRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	lpRpc.topologyClient = topology.NewTopologyServiceClient(cc)
}

func (lpRpc *ListChunkServerRpc) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return lpRpc.topologyClient.ListChunkServer(ctx, lpRpc.Request, opt...)
}

func (cli *mdsClient) ListChunkServer() (interface{}, error) {
	listChunkServerRpc := &ListChunkServerRpc{}
	listChunkServerRpc.ctx = baserpc.NewRpcContext(cli.addrs, LIST_CHUNKSERVER_NAME)
	listChunkServerRpc.Request = &topology.ListChunkServerRequest{}

	return nil, nil
}
