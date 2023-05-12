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
	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/api/curvebs/agent"
	"github.com/opencurve/curve-manager/api/curvebs/core"
	"github.com/opencurve/pigeon"
)

const (
	ADM_HOST_LIST          = "host.list"
	ADM_HOST_COMMIT        = "host.commit"
	ADM_DISK_LIST          = "disk.list"
	ADM_DISK_COMMIT        = "disk.commit"
	ADM_DISK_FORMAT_STATUS = "disk.format.status"
	ADM_DISK_FORMAT        = "disk.format"
	ADM_CONFIG_SHOW        = "config.show"
	ADM_CONFIG_COMMIT      = "config.commit"
	ADM_CLUSTER_LIST       = "cluster.list"
	ADM_CLUSTER_CHECKOUT   = "cluster.checkout"
	ADM_CLUSTER_ADD        = "cluster.add"
	ADM_CLUSTER_DEPLOY     = "cluster.deploy"
	ADM_UNSUPPORT          = "unsupport"
)

func getAdmMethod(method string) string {
	switch method {
	case core.DEPLOY_HOST_LIST:
		return ADM_HOST_LIST
	case core.DEPLOY_HOST_COMMIT:
		return ADM_HOST_COMMIT
	case core.DEPLOY_DISK_LIST:
		return ADM_DISK_LIST
	case core.DEPLOY_DISK_COMMIT:
		return ADM_DISK_COMMIT
	case core.DEPLOY_DISK_FORMAT_STATUS:
		return ADM_DISK_FORMAT_STATUS
	case core.DEPLOY_DISK_FORMAT:
		return ADM_DISK_FORMAT
	case core.DEPLOY_CONFIG_SHOW:
		return ADM_CONFIG_SHOW
	case core.DEPLOY_CONFIG_COMMIT:
		return ADM_CONFIG_COMMIT
	case core.DEPLOY_CLUSTER_ADD:
		return ADM_CLUSTER_ADD
	case core.DEPLOY_CLUSTER_DEPLOY:
		return ADM_CLUSTER_DEPLOY
	case core.DEPLOY_CLUSTER_LIST:
		return ADM_CLUSTER_LIST
	case core.DEPLOY_CLUSTER_CHECKOUT:
		return ADM_CLUSTER_CHECKOUT
	default:
		return ADM_UNSUPPORT
	}
}

func DealDeploy(r *pigeon.Request, ctx *Context) bool {
	ret := agent.ProxyPass(r, ctx.Data, getAdmMethod(r.Args[core.METHOD]))
	if core.NeedRecordLog(r) {
		agent.WriteSystemLog(r.Context.ClientIP(), r.HeadersIn[comm.HEADER_LOG_USER], core.BELONG[r.Args[core.METHOD]],
			r.Args[core.METHOD], "proxy", r.HeadersIn[comm.HEADER_LOG_CONTENT], r.Status)
	}
	return ret
}
