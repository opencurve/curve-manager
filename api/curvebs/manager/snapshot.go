package manager

import (
	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/curve-manager/internal/agent"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/pigeon"
)

func ListSnapshot(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*ListSnapshotRequest)
	snapshots, err := agent.GetSnapshot(data.Size, data.Page, data.UUID, data.User, data.FileName, data.Status)
	if err != nil {
		r.Logger().Error("list snapshot failed",
			pigeon.Field("size", data.Size),
			pigeon.Field("page", data.Page),
			pigeon.Field("UUID", data.UUID),
			pigeon.Field("user", data.User),
			pigeon.Field("fileName", data.FileName),
			pigeon.Field("status", data.Status),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.LIST_SNAPSHOT_FAILED)
	}
	return core.ExitSuccessWithData(r, snapshots)
}
