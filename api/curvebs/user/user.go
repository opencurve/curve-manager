package user

import (
	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/curve-manager/internal/storage"
	"github.com/opencurve/pigeon"
)

func CreateUser(r *pigeon.Request, ctx *Context) bool {
	data := ctx.Data.(*CreateUserRequest)
	err := storage.InsertUser(data.UserName, data.PassWord, data.Permission)
	if err != nil {
		r.Logger().Error("create user failed",
			pigeon.Field("userName", data.UserName),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.CREATE_USER_FAILED)
	}
	return core.Exit(r, errno.OK)
}
