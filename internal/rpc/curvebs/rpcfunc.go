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

package curvebs

import (
	"context"

	"github.com/opencurve/curve-manager/internal/proto/nameserver2"
	"github.com/opencurve/curve-manager/internal/proto/topology"
	"github.com/opencurve/curve-manager/internal/rpc/baserpc"
	"google.golang.org/grpc"
)

// topology
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

// get copysets on chunkserver
type GetCopySetsInChunkServer struct {
	ctx     *baserpc.RpcContext
	client  topology.TopologyServiceClient
	Request *topology.GetCopySetsInChunkServerRequest
}

func (rpc *GetCopySetsInChunkServer) NewRpcClient(cc grpc.ClientConnInterface) {
	rpc.client = topology.NewTopologyServiceClient(cc)
}

func (rpc *GetCopySetsInChunkServer) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return rpc.client.GetCopySetsInChunkServer(ctx, rpc.Request, opt...)
}

// get chunkserver list in copysets
type GetChunkServerListInCopySets struct {
	ctx     *baserpc.RpcContext
	client  topology.TopologyServiceClient
	Request *topology.GetChunkServerListInCopySetsRequest
}

func (rpc *GetChunkServerListInCopySets) NewRpcClient(cc grpc.ClientConnInterface) {
	rpc.client = topology.NewTopologyServiceClient(cc)
}

func (rpc *GetChunkServerListInCopySets) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return rpc.client.GetChunkServerListInCopySets(ctx, rpc.Request, opt...)
}

// get copysets in cluster
type GetCopySetsInCluster struct {
	ctx     *baserpc.RpcContext
	client  topology.TopologyServiceClient
	Request *topology.GetCopySetsInClusterRequest
}

func (rpc *GetCopySetsInCluster) NewRpcClient(cc grpc.ClientConnInterface) {
	rpc.client = topology.NewTopologyServiceClient(cc)
}

func (rpc *GetCopySetsInCluster) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return rpc.client.GetCopySetsInCluster(ctx, rpc.Request, opt...)
}

// get logical pool
type GetLogicalPool struct {
	ctx     *baserpc.RpcContext
	client  topology.TopologyServiceClient
	Request *topology.GetLogicalPoolRequest
}

func (rpc *GetLogicalPool) NewRpcClient(cc grpc.ClientConnInterface) {
	rpc.client = topology.NewTopologyServiceClient(cc)
}

func (rpc *GetLogicalPool) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return rpc.client.GetLogicalPool(ctx, rpc.Request, opt...)
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

// get volume
type GetFileInfo struct {
	ctx     *baserpc.RpcContext
	client  nameserver2.CurveFSServiceClient
	Request *nameserver2.GetFileInfoRequest
}

func (rpc *GetFileInfo) NewRpcClient(cc grpc.ClientConnInterface) {
	rpc.client = nameserver2.NewCurveFSServiceClient(cc)
}

func (rpc *GetFileInfo) Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error) {
	return rpc.client.GetFileInfo(ctx, rpc.Request, opt...)
}
