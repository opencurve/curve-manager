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
	MDS_STATUS             = "mds_status"
	MDS_LEADER             = "leader"
	MDS_FOLLOWER           = "follower"
	MDS_CONF_LISTEN_ADDR   = "mds_config_mds_listen_addr"
	MDS_CONF_AUTH_USERNAME = "mds_config_mds_auth_root_user_name"
	MDS_CONF_AUTH_PASSWORD = "mds_config_mds_auth_root_password"
)

func GetMdsStatus() ([]ServiceStatus, string) {
	size := len(core.GMetricClient.MdsDummyAddr)
	results := make(chan comm.MetricResult, size)
	names := fmt.Sprintf("%s,%s,%s", comm.CURVEBS_VERSION, MDS_STATUS, MDS_CONF_LISTEN_ADDR)
	comm.GetBvarMetric(core.GMetricClient.MdsDummyAddr, names, &results)

	count := 0
	var errors string
	ret := []ServiceStatus{}
	for res := range results {
		if res.Err == nil {
			addr := ""
			v, e := comm.ParseBvarMetric(res.Result.(string))
			if e != nil {
				errors = fmt.Sprintf("%s; %s:%s", errors, res.Key, e.Error())
			} else {
				addr = comm.GetBvarConfMetricValue((*v)[MDS_CONF_LISTEN_ADDR])
			}
			ret = append(ret, ServiceStatus{
				Address: addr,
				Version: (*v)[comm.CURVEBS_VERSION],
				Online:  true,
				Leader:  (*v)[MDS_STATUS] == MDS_LEADER,
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

func GetAuthInfoOfRoot() (string, string, string) {
	size := len(core.GMetricClient.MdsDummyAddr)
	results := make(chan comm.MetricResult, size)
	names := fmt.Sprintf("%s,%s", MDS_CONF_AUTH_USERNAME, MDS_CONF_AUTH_PASSWORD)
	comm.GetBvarMetric(core.GMetricClient.MdsDummyAddr, names, &results)

	count := 0
	var errors string
	var userName, passWord string
	for res := range results {
		if res.Err == nil {
			v, e := comm.ParseBvarMetric(res.Result.(string))
			if e != nil {
				errors = fmt.Sprintf("%s; %s:%s", errors, res.Key, e.Error())
			} else {
				userName = comm.GetBvarConfMetricValue((*v)[MDS_CONF_AUTH_USERNAME])
				passWord = comm.GetBvarConfMetricValue((*v)[MDS_CONF_AUTH_PASSWORD])
				return userName, passWord, ""
			}
		} else {
			errors = fmt.Sprintf("%s; %s:%s", errors, res.Key, res.Err.Error())
		}
		count += 1
		if count >= size {
			break
		}
	}
	return "", "", errors
}
