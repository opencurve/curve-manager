package curvebs

import (
	"context"

	"github.com/opencurve/curve-manager/internal/proto/nameserver2"
	"github.com/opencurve/curve-manager/internal/proto/topology"
	"github.com/opencurve/curve-manager/internal/rpc/baserpc"
	"google.golang.org/grpc"
)

// list physical pool
type ListPhysicalPoolRpc struct {
	ctx     *baserpc.RpcContext
	client  topology.TopologyServiceClient
	Request *topology.ListPhysicalPoolRequest
}

func (rpc *ListPhysicalPoolRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	rpc.client = topology.NewTopologyServiceClient(cc)
}

func (rpc *ListPhysicalPoolRpc) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return rpc.client.ListPhysicalPool(ctx, rpc.Request, opt...)
}

// list logical pool
type ListLogicalPoolRpc struct {
	ctx     *baserpc.RpcContext
	client  topology.TopologyServiceClient
	Request *topology.ListLogicalPoolRequest
}

func (rpc *ListLogicalPoolRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	rpc.client = topology.NewTopologyServiceClient(cc)
}

func (rpc *ListLogicalPoolRpc) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return rpc.client.ListLogicalPool(ctx, rpc.Request, opt...)
}

// list zones of logical pool
type ListPoolZonesRpc struct {
	ctx     *baserpc.RpcContext
	client  topology.TopologyServiceClient
	Request *topology.ListPoolZoneRequest
}

func (rpc *ListPoolZonesRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	rpc.client = topology.NewTopologyServiceClient(cc)
}

func (rpc *ListPoolZonesRpc) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return rpc.client.ListPoolZone(ctx, rpc.Request, opt...)
}

// list servers of zone
type ListZoneServer struct {
	ctx     *baserpc.RpcContext
	client  topology.TopologyServiceClient
	Request *topology.ListZoneServerRequest
}

func (rpc *ListZoneServer) NewRpcClient(cc grpc.ClientConnInterface) {
	rpc.client = topology.NewTopologyServiceClient(cc)
}

func (rpc *ListZoneServer) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return rpc.client.ListZoneServer(ctx, rpc.Request, opt...)
}

// list chunkservers of server
type ListChunkServer struct {
	ctx     *baserpc.RpcContext
	client  topology.TopologyServiceClient
	Request *topology.ListChunkServerRequest
}

func (rpc *ListChunkServer) NewRpcClient(cc grpc.ClientConnInterface) {
	rpc.client = topology.NewTopologyServiceClient(cc)
}

func (rpc *ListChunkServer) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return rpc.client.ListChunkServer(ctx, rpc.Request, opt...)
}

// get chunkserver in cluster
type GetChunkServerInCluster struct {
	ctx     *baserpc.RpcContext
	client  topology.TopologyServiceClient
	Request *topology.GetChunkServerInClusterRequest
}

func (rpc *GetChunkServerInCluster) NewRpcClient(cc grpc.ClientConnInterface) {
	rpc.client = topology.NewTopologyServiceClient(cc)
}

func (rpc *GetChunkServerInCluster) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return rpc.client.GetChunkServerInCluster(ctx, rpc.Request, opt...)
}

// nameserver2
// get file(include dir) allocated space
type GetFileAllocatedSize struct {
	ctx     *baserpc.RpcContext
	client  nameserver2.CurveFSServiceClient
	Request *nameserver2.GetAllocatedSizeRequest
}

func (rpc *GetFileAllocatedSize) NewRpcClient(cc grpc.ClientConnInterface) {
	rpc.client = nameserver2.NewCurveFSServiceClient(cc)
}

func (rpc *GetFileAllocatedSize) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return rpc.client.GetAllocatedSize(ctx, rpc.Request, opt...)
}

// list volume in dir
type ListDir struct {
	ctx     *baserpc.RpcContext
	client  nameserver2.CurveFSServiceClient
	Request *nameserver2.ListDirRequest
}

func (rpc *ListDir) NewRpcClient(cc grpc.ClientConnInterface) {
	rpc.client = nameserver2.NewCurveFSServiceClient(cc)
}

func (rpc *ListDir) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return rpc.client.ListDir(ctx, rpc.Request, opt...)
}