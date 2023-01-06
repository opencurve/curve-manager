package mds

import (
	"context"

	"github.com/opencurve/curve-manager/internal/proto/topology"
	"github.com/opencurve/curve-manager/internal/rpc/baserpc"
	"google.golang.org/grpc"
)

// list physical pool
type ListPhysicalPoolRpc struct {
	ctx            *baserpc.RpcContext
	Request        *topology.ListPhysicalPoolRequest
	topologyClient topology.TopologyServiceClient
}

func (rpc *ListPhysicalPoolRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	rpc.topologyClient = topology.NewTopologyServiceClient(cc)
}

func (rpc *ListPhysicalPoolRpc) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return rpc.topologyClient.ListPhysicalPool(ctx, rpc.Request, opt...)
}

// list logical pool
type ListLogicalPoolRpc struct {
	ctx            *baserpc.RpcContext
	Request        *topology.ListLogicalPoolRequest
	topologyClient topology.TopologyServiceClient
}

func (rpc *ListLogicalPoolRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	rpc.topologyClient = topology.NewTopologyServiceClient(cc)
}

func (rpc *ListLogicalPoolRpc) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return rpc.topologyClient.ListLogicalPool(ctx, rpc.Request, opt...)
}

// list zones of logical pool
type ListPoolZonesRpc struct {
	ctx            *baserpc.RpcContext
	Request        *topology.ListPoolZoneRequest
	topologyClient topology.TopologyServiceClient
}

func (rpc *ListPoolZonesRpc) NewRpcClient(cc grpc.ClientConnInterface) {
	rpc.topologyClient = topology.NewTopologyServiceClient(cc)
}

func (rpc *ListPoolZonesRpc) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return rpc.topologyClient.ListPoolZone(ctx, rpc.Request, opt...)
}

// list servers of zone
type ListZoneServer struct {
	ctx            *baserpc.RpcContext
	Request        *topology.ListZoneServerRequest
	topologyClient topology.TopologyServiceClient
}

func (rpc *ListZoneServer) NewRpcClient(cc grpc.ClientConnInterface) {
	rpc.topologyClient = topology.NewTopologyServiceClient(cc)
}

func (rpc *ListZoneServer) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return rpc.topologyClient.ListZoneServer(ctx, rpc.Request, opt...)
}

// list chunkservers of server
type ListChunkServer struct {
	ctx            *baserpc.RpcContext
	Request        *topology.ListChunkServerRequest
	topologyClient topology.TopologyServiceClient
}

func (rpc *ListChunkServer) NewRpcClient(cc grpc.ClientConnInterface) {
	rpc.topologyClient = topology.NewTopologyServiceClient(cc)
}

func (rpc *ListChunkServer) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return rpc.topologyClient.ListChunkServer(ctx, rpc.Request, opt...)
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