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
	"github.com/opencurve/curve-manager/internal/email"
	metrics "github.com/opencurve/curve-manager/internal/metrics/core"
	bsrpc "github.com/opencurve/curve-manager/internal/rpc/curvebs"
	"github.com/opencurve/curve-manager/internal/snapshotclone"
	"github.com/opencurve/curve-manager/internal/storage"
	"github.com/opencurve/pigeon"
)

func initClients(cfg *pigeon.Configure) error {
	// init mds rpc client
	err := bsrpc.Init(cfg)
	if err != nil {
		return err
	}

	// init metric client
	err = metrics.Init(cfg)
	if err != nil {
		return err
	}

	// init snapshot clone client
	snapshotclone.Init(cfg)

	// init email client which used to reset password and some system notifications
	email.Init(cfg)
	return nil
}

func Init(cfg *pigeon.Configure) error {
	// init storage
	err := storage.Init(cfg)
	if err != nil {
		return err
	}

	return initClients(cfg)
}
