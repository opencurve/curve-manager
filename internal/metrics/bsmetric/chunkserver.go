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

import metricomm "github.com/opencurve/curve-manager/internal/metrics/common"

// @return map[version]number
func GetChunkServerVersion(endpoints *[]string) (map[string]int, error) {
	size := len(*endpoints)
	results := make(chan metricomm.MetricResult, size)
	metricomm.GetBvarMetric(*endpoints, metricomm.CURVEBS_VERSION, &results)

	count := 0
	ret := make(map[string]int)
	for res := range results {
		if res.Err == nil {
			v, e := metricomm.ParseBvarMetric(res.Result.(string))
			if e != nil {
				return nil, e
			} else {
				ret[(*v)[metricomm.CURVEBS_VERSION]] += 1
			}
		} else {
			return nil, res.Err
		}
		count += 1
		if count >= size {
			break
		}
	}
	return ret, nil
}

// @return key: chunkserver's addr, value: copysets' raft status
func GetCopysetRaftStatus(endpoints *[]string) (map[string][]map[string]string, error) {
	size := len(*endpoints)
	results := make(chan metricomm.MetricResult, size)
	metricomm.GetRaftStatusMetric(*endpoints, &results)

	count := 0
	ret := map[string][]map[string]string{}
	for res := range results {
		if res.Err == nil {
			v, e := metricomm.ParseRaftStatusMetric(res.Result.(string))
			if e != nil {
				return nil, e
			} else {
				ret[res.Key.(string)] = v
			}
		} else {
			ret[res.Key.(string)] = nil
		}
		count += 1
		if count >= size {
			break
		}
	}
	return ret, nil
}
