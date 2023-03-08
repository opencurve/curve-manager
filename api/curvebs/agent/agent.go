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
* Created Date: 2023-03-07
* Author: wanghai (SeanHai)
 */

package agent

import (
	"fmt"
	"strings"

	bsrpc "github.com/SeanHai/curve-go-rpc/rpc/curvebs"
	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/pigeon"
)

var (
	GMdsClient *bsrpc.MdsClient
)

const (
	CURVEBS_MDS_ADDRESS = "mds.address"

	DEFAULT_RPC_TIMEOUT_MS  = 500
	DEFAULT_RPC_RETRY_TIMES = 3
)

func Init(cfg *pigeon.Configure) error {
	addrs := cfg.GetConfig().GetString(CURVEBS_MDS_ADDRESS)
	if len(addrs) == 0 {
		return fmt.Errorf("no cluster mds address found")
	}
	GMdsClient = bsrpc.NewMdsClient(bsrpc.MdsClientOption{
		TimeoutMs:  DEFAULT_RPC_TIMEOUT_MS,
		RetryTimes: DEFAULT_RPC_RETRY_TIMES,
		Addrs:      strings.Split(addrs, common.CURVEBS_ADDRESS_DELIMITER),
	})
	return nil
}
