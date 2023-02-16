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

package bsmetric

import (
	"fmt"

	comm "github.com/opencurve/curve-manager/internal/metrics/common"
	"github.com/opencurve/curve-manager/internal/metrics/core"
)

const (
	SNAPSHOT_CLONE_STATUS           = "snapshotcloneserver_status"
	SNAPSHOT_CLONE_CONF_LISTEN_ADDR = "snapshotcloneserver_config_server_address"
	SNAPSHOT_CLONE_LEADER           = "active"
)

type ServiceStatus comm.ServiceStatus

func GetSnapShotCloneServerStatus() ([]ServiceStatus, string) {
	ret := []ServiceStatus{}
	size := len(core.GMetricClient.SnapShotCloneServerDummyAddr)
	if size == 0 {
		return ret, "no snapshotclone service address found"
	}
	results := make(chan comm.MetricResult, size)
	names := fmt.Sprintf("%s,%s,%s", comm.CURVEBS_VERSION, SNAPSHOT_CLONE_STATUS, SNAPSHOT_CLONE_CONF_LISTEN_ADDR)
	comm.GetBvarMetric(core.GMetricClient.SnapShotCloneServerDummyAddr, names, &results)

	count := 0
	var errors string
	for res := range results {
		if res.Err == nil {
			addr := ""
			v, e := comm.ParseBvarMetric(res.Result.(string))
			if e != nil {
				errors = fmt.Sprintf("%s; %s:%s", errors, res.Key, e.Error())
			} else {
				addr = comm.GetBvarConfMetricValue((*v)[SNAPSHOT_CLONE_CONF_LISTEN_ADDR])
			}
			ret = append(ret, ServiceStatus{
				Address: addr,
				Version: (*v)[comm.CURVEBS_VERSION],
				Online:  true,
				Leader:  (*v)[SNAPSHOT_CLONE_STATUS] == SNAPSHOT_CLONE_LEADER,
			})
		} else {
			errors = fmt.Sprintf("%s; %s:%s", errors, res.Key, res.Err.Error())
			ret = append(ret, ServiceStatus{
				Address: res.Key.(string),
				Version: "",
				Leader:  false,
				Online:  false,
			})
		}
		count += 1
		if count >= size {
			break
		}
	}
	return ret, errors
}
