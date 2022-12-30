package manager

import (
	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/curve-manager/internal/metrics/bsmetric"
	"github.com/opencurve/pigeon"
)

func GetEtcdStatus(r *pigeon.Request, ctx *Context) bool {
	status, err := bsmetric.GetEtcdStatus()
	if err != nil {
		r.Logger().Error("GetEtcdStatus failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.GET_ETCD_STATUS_FAILED)
	}
	return core.ExitSuccessWithData(r, status)
}

func GetMdsStatus(r *pigeon.Request, ctx *Context) bool {
	status, err := bsmetric.GetMdsStatus()
	if err != "" {
		r.Logger().Error("GetMdsStatus failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.GET_MDS_STATUS_FAILED)
	}
	return core.ExitSuccessWithData(r, status)
}

func GetSnapShotCloneServerStatus(r *pigeon.Request, ctx *Context) bool {
	status, err := bsmetric.GetSnapShotCloneServerStatus()
	if err != "" {
		r.Logger().Error("GetSnapShotCloneServerStatus failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.GET_SNAPSHOT_CLONE_STATUS_FAILED)
	}
	return core.ExitSuccessWithData(r, status)
}
