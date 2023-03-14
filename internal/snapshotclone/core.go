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
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/pigeon"
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

func Init(cfg *pigeon.Configure) {
	GSnapshotCloneClient = &snapshotCloneClient{}
	// init http client
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxConnsPerHost:       100,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
	}

	GSnapshotCloneClient.client = resty.NewWithClient(httpClient)
	addr := cfg.GetConfig().GetString(CURVEBS_SNAPSHOT_CLONE_PROXY_ADDRESS)
	if len(addr) != 0 {
		GSnapshotCloneClient.SnapShotCloneProxyAddr = strings.Split(addr, common.CURVEBS_ADDRESS_DELIMITER)
	}
}

func (cli *snapshotCloneClient) sendHttp2SnapshotClone(queryParam string) (string, error) {
	if cli.SnapShotCloneProxyAddr == nil {
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
