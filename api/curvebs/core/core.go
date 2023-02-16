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
* Created Date: 2023-02-11
* Author: wanghai (SeanHai)
 */

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
