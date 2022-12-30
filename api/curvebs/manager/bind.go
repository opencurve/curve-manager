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

type GetEtcdStatusRequest struct {
}

type GetMdsStatusRequest struct {
}

type GetSnapShotCloneServerStatusRequest struct {
}

type ListLogicalPoolRequest struct {
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
		"topo.list.pool",
		ListLogicalPoolRequest{},
		ListLogicalPool,
	},
}
