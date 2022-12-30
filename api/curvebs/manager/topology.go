package manager

import (
	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/curve-manager/internal/rpc/curvebs/mds"
	"github.com/opencurve/pigeon"
)

func ListLogicalPool(r *pigeon.Request, ctx *Context) bool {
	pools, err := mds.GMdsClient.ListLogicalPool()
	if err != nil {
		r.Logger().Error("ListLogicalPool failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.LIST_POOL_FAILED)
	}
	return core.ExitSuccessWithData(r, pools)
}
