package manager

import (
	"github.com/mcuadros/go-defaults"
	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/curve-manager/internal/agent"
	"github.com/opencurve/pigeon"
)

func ListVolume(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*ListVolumeRequest)
	defaults.SetDefaults(data)
	volumes, err := agent.ListVolume(data.Size, data.Page, data.Path, data.SortKey)
	if err != nil {
		r.Logger().Error("list volume failed",
			pigeon.Field("path", data.Path),
			pigeon.Field("size", data.Size),
			pigeon.Field("page", data.Page),
			pigeon.Field("sortkey", data.SortKey),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.LIST_VOLUME_FAILED)
	}
	return core.ExitSuccessWithData(r, volumes)
}
