package manager

import (
	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/curve-manager/internal/agent"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/pigeon"
)

func ListTopology(r *pigeon.Request, ctx *Context) bool {
	topo, err := agent.ListTopology()
	if err != nil {
		r.Logger().Error("ListTopology failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.LIST_TOPO_FAILED)
	}
	return core.ExitSuccessWithData(r, topo)
}

func ListLogicalPool(r *pigeon.Request, ctx *Context) bool {
	pools, err := agent.ListLogicalPool()
	if err != nil {
		r.Logger().Error("ListLogicalPool failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.LIST_POOL_FAILED)
	}
	return core.ExitSuccessWithData(r, pools)
}
