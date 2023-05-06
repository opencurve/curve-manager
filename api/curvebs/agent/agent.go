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

	SYSTEM_ALERT_EXPIRATION_DAYS         = "system.alert.expiration.days"
	DEFAULT_SYSTEM_ALERT_EXPIRATION_DAYS = 30
)

var (
	systemLogChann chan storage.Log
)

func Init(cfg *pigeon.Configure, logger *pigeon.Logger) error {
	logExpirationDays := cfg.GetConfig().GetInt(SYSTEM_LOG_EXPIRATION_DAYS)
	if logExpirationDays <= 0 {
		logExpirationDays = DEFAULT_SYSTEM_LOG_EXPIRATION_DAYS
	}

	alertExpirationDays := cfg.GetConfig().GetInt(SYSTEM_ALERT_EXPIRATION_DAYS)
	if alertExpirationDays <= 0 {
		alertExpirationDays = DEFAULT_SYSTEM_ALERT_EXPIRATION_DAYS
	}

	curveadm_service_addr = cfg.GetConfig().GetString(CURVEADM_SERVICE_ADDRESS)

	// write system operation log
	systemLogChann = make(chan storage.Log, 128)
	go writeSystemLog(logger)
	// clear expired logs
	go clearExpiredSystemLog(logExpirationDays, logger)
	// start system alerts
	return initAlerts(alertExpirationDays, logger)
}
