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
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/pigeon"
)

const (
	CURVEADM_SERVICE_ADDRESS = "curveadm.service.address"

	CLUSTER_SERVICES_ADDRESS = "cluster.service.addr"

	METHOD_CLUSTER_DEPLOY = "deploy.cluster.deploy"
)

var (
	curveadm_service_addr = ""
)

type AdmHttpResponse struct {
	ErrorCode string      `json:"errorCode"`
	ErrorMsg  string      `json:"errorMsg"`
	Data      interface{} `json:"data"`
}

func ProxyPass(r *pigeon.Request, body interface{}, method string) bool {
	args := fmt.Sprintf("method=%s", method)
	if method == METHOD_CLUSTER_DEPLOY {
		defer InitClients()
	}
	return r.ProxyPass(curveadm_service_addr, r.WithURI("/"), r.WithArgs(args), r.WithScheme("http"), r.WithBody(body))
}

func GetCurrentClusterServicesAddr() (map[string]interface{}, error) {
	httpClient := common.GetHttpClient()
	url := (&url.URL{
		Scheme:   "http",
		Host:     curveadm_service_addr,
		Path:     "/",
		RawQuery: fmt.Sprintf("%s=%s", "method", CLUSTER_SERVICES_ADDRESS),
	}).String()

	resp, err := resty.NewWithClient(httpClient).R().
		SetHeader("Connection", "Keep-Alive").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Curve-Manager").
		Execute("GET", url)
	if err != nil {
		return nil, fmt.Errorf("getClusterServicesAddr failed: %v", err)
	}

	respStruct := AdmHttpResponse{}
	err = json.Unmarshal([]byte(resp.String()), &respStruct)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal getClusterServicesAddr response failed, resp = %s, err = %v", resp.String(), err)
	}
	return respStruct.Data.(map[string]interface{}), nil
}
