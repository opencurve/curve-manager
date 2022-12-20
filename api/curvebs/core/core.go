package core

import (
	"strconv"

	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/pigeon"
)

func Exit(r *pigeon.Request, code errno.Errno) bool {
	r.SendJSON(pigeon.JSON{
		"errorCode": strconv.Itoa(code.Code()),
		"errorMsg":  code.Description(),
	})
	if code != errno.OK {
		r.HeadersOut[comm.HEADER_ERROR_CODE] = strconv.Itoa(code.Code())
	}

	return r.Exit(code.HTTPCode())
}

func Default(r *pigeon.Request) bool {
	r.Logger().Warn("unupport request uri", pigeon.Field("uri", r.Uri))
	return Exit(r, errno.UNSUPPORT_REQUEST_URI)
}

func ExitSuccessWithData(r *pigeon.Request, data interface{}) bool {
	r.SendJSON(pigeon.JSON{
		"data":      data,
		"errorCode": "0",
		"errorMsg":  "success",
	})
	return r.Exit(200)
}