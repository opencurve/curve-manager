package manager

import (
	"reflect"

	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/pigeon"
)

func Entrypoint(r *pigeon.Request) bool {
	if r.GetModuleCtx(MODULE) == nil {
		return r.NextHandler()
	}

	if r.Method != pigeon.HTTP_METHOD_GET &&
		r.Method != pigeon.HTTP_METHOD_POST {
		return core.Exit(r, errno.UNSUPPORT_HTTP_METHOD)
	}

	request, ok := METHOD_REQUEST[r.Args["method"]]
	if !ok {
		return core.Exit(r, errno.UNSUPPORT_METHOD_ARGUMENT)
	} else if request.httpMethod != r.Method {
		return core.Exit(r, errno.HTTP_METHOD_MISMATCHED)
	}

	vType := reflect.TypeOf(request.vType)
	data := reflect.New(vType).Interface()
	if err := r.BindBody(data); err != nil {
		r.Logger().Error("bad request form param",
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return core.Exit(r, errno.BAD_REQUEST_FORM_PARAM)
	}

	if !core.AccessAllowed(r, data) {
		return core.Exit(r, errno.REQUEST_IS_DENIED_FOR_SIGNATURE)
	}
	return request.handler(r, &Context{data})
}
