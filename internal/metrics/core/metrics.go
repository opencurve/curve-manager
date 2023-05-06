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
	"fmt"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/opencurve/curve-manager/internal/common"
)

const (
	CURVEBS_MONITOR_ADDRESS              = "monitor.prometheus.address"
	CURVEBS_ETCD_ADDRESS                 = "etcd.address"
	CURVEBS_MDS_DUMMY_ADDRESS            = "mds.dummy.address"
	CURVEBS_SNAPSHOT_CLONE_DUMMY_ADDRESS = "snapshot.clone.dummy.address"
)

type metricClient struct {
	client                       *resty.Client
	PromeAddr                    string
	EtcdAddr                     []string
	MdsDummyAddr                 []string
	SnapShotCloneServerDummyAddr []string
}

var (
	GMetricClient *metricClient
)

func Init(cfg map[string]interface{}) {
	GMetricClient = &metricClient{}
	GMetricClient.client = resty.NewWithClient(common.GetHttpClient())

	GMetricClient.PromeAddr = cfg[CURVEBS_MONITOR_ADDRESS].(string)
	GMetricClient.EtcdAddr = strings.Split(cfg[CURVEBS_ETCD_ADDRESS].(string), common.CURVEBS_ADDRESS_DELIMITER)
	GMetricClient.MdsDummyAddr = strings.Split(cfg[CURVEBS_MDS_DUMMY_ADDRESS].(string), common.CURVEBS_ADDRESS_DELIMITER)
	GMetricClient.SnapShotCloneServerDummyAddr = strings.Split(cfg[CURVEBS_SNAPSHOT_CLONE_DUMMY_ADDRESS].(string),
		common.CURVEBS_ADDRESS_DELIMITER)
}

func (cli *metricClient) GetMetricFromService(host, path string) (interface{}, error) {
	url := (&url.URL{
		Scheme: "http",
		Host:   host,
		Path:   path,
	}).String()

	resp, err := cli.client.R().
		SetHeader("Connection", "Keep-Alive").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "curl/7.52.1").
		Execute("GET", url)
	if err != nil {
		return nil, fmt.Errorf("get bvar metric failed, addr: %s, err: %v", host, err)
	} else if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get bvar metric failed, addr: %s, status = %s", host, resp.Status())
	}
	return resp.String(), nil
}

func (cli *metricClient) GetMetricFromPrometheus(path, queryKey, queryValue string, ret interface{}) error {
	url := (&url.URL{
		Scheme:   "http",
		Host:     cli.PromeAddr,
		Path:     path,
		RawQuery: fmt.Sprintf("%s=%s", queryKey, queryValue),
	}).String()

	resp, err := cli.client.R().
		SetHeader("Connection", "Keep-Alive").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Curve-Manager").
		SetResult(ret).
		Execute("GET", url)
	if err != nil {
		return fmt.Errorf("get prometheus metric failed: %v", err)
	} else if resp.StatusCode() != 200 {
		return fmt.Errorf("get prometheus metric failed, status = %s, url = %v",
			resp.Status(), url)
	}
	return nil
}
