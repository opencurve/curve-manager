package manager

import (
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/pigeon"
)

const (
	MODULE = "manager"
)

var METHOD_REQUEST map[string]Request

type (
	HandlerFunc func(r *pigeon.Request, ctx *Context) bool

	Context struct {
		Data interface{}
	}

	Request struct {
		httpMethod string
		method     string
		vType      interface{}
		handler    HandlerFunc
	}
)

func init() {
	METHOD_REQUEST = map[string]Request{}
	for _, request := range requests {
		METHOD_REQUEST[request.method] = request
		core.BELONG[request.method] = MODULE
	}
}

type GetEtcdStatusRequest struct{}

type GetMdsStatusRequest struct{}

type GetSnapShotCloneServerStatusRequest struct{}

type ListTopologyRequest struct{}

type ListLogicalPoolRequest struct{}

type GetChunkServerStatusRequest struct{}

type GetClusterStatusRequest struct{}

type GetClusterSpaceRequest struct{}

type GetClusterPerformanceRequest struct{}

type ListVolumeRequest struct {
	Size    uint32 `json:"size" binding:"required"`
	Page    uint32 `json:"page" binding:"required"`
	Path    string `json:"path" default:"/"`
	SortKey string `json:"sortKey" default:"id"`
}

type ListSnapshotRequest struct {
	Size     uint32 `json:"size" binding:"required"`
	Page     uint32 `json:"page" binding:"required"`
	UUID     string `json:"uuid"`
	User     string `json:"user"`
	FileName string `json:"fileName"`
	Status   string `json:"status"`
}

var requests = []Request{
	{
		"GET",
		"status.etcd",
		GetEtcdStatusRequest{},
		GetEtcdStatus,
	},
	{
		"GET",
		"status.mds",
		GetMdsStatusRequest{},
		GetMdsStatus,
	},
	{
		"GET",
		"status.snapshotcloneserver",
		GetSnapShotCloneServerStatusRequest{},
		GetSnapShotCloneServerStatus,
	},
	{
		"GET",
		"status.chunkserver",
		GetChunkServerStatusRequest{},
		GetChunkServerStatus,
	},
	{
		"GET",
		"status.cluster",
		GetClusterStatusRequest{},
		GetClusterStatus,
	},
	{
		"GET",
		"space.cluster",
		GetClusterSpaceRequest{},
		GetClusterSpace,
	},
	{
		"GET",
		"performance.cluster",
		GetClusterPerformanceRequest{},
		GetClusterPerformance,
	},
	{
		"GET",
		"topo.list",
		ListTopologyRequest{},
		ListTopology,
	},
	{
		"GET",
		"topo.pool.list",
		ListLogicalPoolRequest{},
		ListLogicalPool,
	},
	{
		"POST",
		"volume.list",
		ListVolumeRequest{},
		ListVolume,
	},
	{
		"POST",
		"snapshot.list",
		ListSnapshotRequest{},
		ListSnapshot,
	},
}
