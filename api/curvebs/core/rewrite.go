package core

import (
	"github.com/google/uuid"
	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/pigeon"
)

var (
	BELONG map[string]string
)

func init() {
	BELONG = map[string]string{}
}

func Rewrite(r *pigeon.Request) bool {
	// request id
	requestId := uuid.New().String()
	r.HeadersIn[comm.HEADER_REQUEST_ID] = requestId
	r.HeadersOut[comm.HEADER_REQUEST_ID] = requestId
	r.Var.LogAttach = requestId
	// console version
	r.HeadersOut[comm.HEADER_CURVE_CONSOLE_VERSION] = comm.VERSION

	method := r.Args["method"]
	module, ok := BELONG[method]
	if ok {
		r.SetModuleCtx(module, true)
	} else {
		return Exit(r, errno.UNSUPPORT_METHOD_ARGUMENT)
	}
	return r.NextHandler()
}
