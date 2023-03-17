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
	"github.com/opencurve/curve-manager/internal/storage"
	"github.com/opencurve/pigeon"
)

const (
	SYSTEM_LOG_EXPIRATION_DAYS         = "system.log.expiration.days"
	DEFAULT_SYSTEM_LOG_EXPIRATION_DAYS = 30
)

var (
	GSystemLogChann chan storage.SystemLog
)

func Init(cfg *pigeon.Configure, logger *pigeon.Logger) {
	expirationDays := cfg.GetConfig().GetInt(SYSTEM_LOG_EXPIRATION_DAYS)
	if expirationDays <= 0 {
		expirationDays = DEFAULT_SYSTEM_LOG_EXPIRATION_DAYS
	}
	// write system operation log
	GSystemLogChann = make(chan storage.SystemLog, 128)
	go writeSystemLog(logger)
	// clear expired logs
	go clearExpiredSystemLog(expirationDays, logger)
}
