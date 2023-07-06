/*
*  Copyright (c) 2023 NetEase Inc.
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
 */

/*
* Project: Curve-Manager
* Created Date: 2023-04-13
* Author: wanghai (SeanHai)
 */

package deploy

import (
	"encoding/json"
	"reflect"

	"github.com/mcuadros/go-defaults"
	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/curve-manager/internal/storage"
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

	defaults.SetDefaults(data)
	if core.NeedRecordLog(r) {
		c, _ := json.Marshal(data)
		r.HeadersIn[comm.HEADER_LOG_CONTENT] = string(c)
		r.HeadersIn[comm.HEADER_LOG_USER] = storage.GetLoginUserByToken(r.HeadersIn[comm.HEADER_AUTH_TOKEN])
	}

	if e := core.AccessAllowed(r, data); e != errno.OK {
		return core.Exit(r, e)
	}

	return request.handler(r, &Context{data})
}
