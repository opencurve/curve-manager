package manager

import (
	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/curve-manager/internal/agent"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/pigeon"
)

func GetEtcdStatus(r *pigeon.Request, ctx *Context) bool {
	status, err := agent.GetEtcdStatus()
	if err != "" {
		r.Logger().Error("GetEtcdStatus failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.GET_ETCD_STATUS_FAILED)
	}
	return core.ExitSuccessWithData(r, status)
}

func GetMdsStatus(r *pigeon.Request, ctx *Context) bool {
	status, err := agent.GetMdsStatus()
	if err != "" {
		r.Logger().Error("GetMdsStatus failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.GET_MDS_STATUS_FAILED)
	}
	return core.ExitSuccessWithData(r, status)
}

func GetSnapShotCloneServerStatus(r *pigeon.Request, ctx *Context) bool {
	status, err := agent.GetSnapShotCloneServerStatus()
	if err != "" {
		r.Logger().Error("GetSnapShotCloneServerStatus failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.GET_SNAPSHOT_CLONE_STATUS_FAILED)
	}
	return core.ExitSuccessWithData(r, status)
}

func GetChunkServerStatus(r *pigeon.Request, ctx *Context) bool {
	status, err := agent.GetChunkServerStatus()
	if err != nil {
		r.Logger().Error("GetChunkServerStatus failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.GET_CHUNKSERVER_STATUS_FAILED)
	}
	return core.ExitSuccessWithData(r, status)
}

func GetClusterStatus(r *pigeon.Request, ctx *Context) bool {
	status, err := agent.GetClusterStatus()
	if err != nil {
		r.Logger().Error("GetClusterStatus failed",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
	}
	return core.ExitSuccessWithData(r, status)
}
