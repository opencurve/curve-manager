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

package snapshotclone

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/opencurve/curve-manager/internal/common"
)

type snapshotCloneClient struct {
	client                 *resty.Client
	SnapShotCloneProxyAddr []string
}

var (
	GSnapshotCloneClient *snapshotCloneClient
)

const (
	CURVEBS_SNAPSHOT_CLONE_PROXY_ADDRESS = "snapshot.clone.proxy.address"
)

func Init(cfg map[string]string) {
	GSnapshotCloneClient = &snapshotCloneClient{}
	GSnapshotCloneClient.client = resty.NewWithClient(common.GetHttpClient())
	GSnapshotCloneClient.SnapShotCloneProxyAddr = strings.Split(cfg[CURVEBS_SNAPSHOT_CLONE_PROXY_ADDRESS],
		common.CURVEBS_ADDRESS_DELIMITER)
}

func (cli *snapshotCloneClient) sendHttp2SnapshotClone(queryParam string) (string, error) {
	if cli.SnapShotCloneProxyAddr == nil || len(cli.SnapShotCloneProxyAddr) == 0 {
		return "", fmt.Errorf("no snapshot proxy address found")
	}
	url := (&url.URL{
		Scheme:   "http",
		Host:     cli.SnapShotCloneProxyAddr[0],
		Path:     SERVICE_ROUTER,
		RawQuery: queryParam,
	}).String()

	resp, err := cli.client.R().
		SetHeader("Connection", "Keep-Alive").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Curve-Manager").
		Execute("GET", url)
	if err != nil {
		return "", fmt.Errorf("sendHttp2SnapshotClone failed: %v", err)
	}
	return resp.String(), nil
}
