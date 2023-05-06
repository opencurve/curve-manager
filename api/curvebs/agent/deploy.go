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
* Created Date: 2023-04-14
* Author: wanghai (SeanHai)
 */

package agent

import (
	"fmt"

	"github.com/opencurve/pigeon"
)

const (
	CURVEADM_SERVICE_ADDRESS = "curveadm.service.address"
)

var (
	curveadm_service_addr = ""
)

func ProxyPass(r *pigeon.Request, body interface{}, method string) bool {
	args := fmt.Sprintf("method=%s", method)
	return r.ProxyPass(curveadm_service_addr, r.WithURI("/"), r.WithArgs(args), r.WithScheme("http"), r.WithBody(body))
}
