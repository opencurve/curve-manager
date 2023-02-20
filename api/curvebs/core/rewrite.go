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

	method := r.Args[METHOD]
	module, ok := BELONG[method]
	if ok {
		r.SetModuleCtx(module, true)
	} else {
		return Exit(r, errno.UNSUPPORT_METHOD_ARGUMENT)
	}
	return r.NextHandler()
}
