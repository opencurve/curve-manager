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
* Created Date: 2023-03-16
* Author: wanghai (SeanHai)
 */

package agent

import (
	"time"

	comm "github.com/opencurve/curve-manager/api/common"
	"github.com/opencurve/curve-manager/internal/errno"
	"github.com/opencurve/curve-manager/internal/storage"
	"github.com/opencurve/pigeon"
)

const (
	WRITE_SYSTEM_LOG_INTERVAL = 1 * time.Second
	CLEAR_SYSTEM_LOG_INTERVAL = 1 * time.Hour
)

func clearExpiredSystemLog(expirationDays int, logger *pigeon.Logger) {
	timer := time.NewTimer(CLEAR_SYSTEM_LOG_INTERVAL)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			err := storage.DeleteSystemLog(time.Now().AddDate(0, 0, -expirationDays).UnixMilli())
			if err != nil {
				logger.Error("clear expired system log failed",
					pigeon.Field("error", err))
			}
			timer.Reset(CLEAR_SYSTEM_LOG_INTERVAL)
		}
	}
}

func writeSystemLog(logger *pigeon.Logger) {
	timer := time.NewTimer(WRITE_SYSTEM_LOG_INTERVAL)
	defer timer.Stop()
	for {
		select {
		case data, ok := <-systemLogChann:
			if ok {
				err := storage.AddSystemLog(&data)
				if err != nil {
					logger.Error("write system log failed",
						pigeon.Field("error", err))
				}
			}
		case <-timer.C:
			timer.Reset(WRITE_SYSTEM_LOG_INTERVAL)
		}
	}
}

func WriteSystemLog(ip, user, module, method, error_msg, content string, error_code int) {
	timeMs := time.Now().UnixMilli()
	logItem := storage.SystemLog{
		TimeMs:    timeMs,
		IP:        ip,
		User:      user,
		Module:    module,
		Method:    method,
		ErrorCode: error_code,
		ErrorMsg:  error_msg,
		Content:   content,
	}
	systemLogChann <- logItem
}

func GetSysLog(r *pigeon.Request, start, end int64, page, size uint32, filter string) (interface{}, errno.Errno) {
	if start == 0 && end == 0 {
		end = time.Now().UnixMilli()
	}
	userName := storage.GetLoginUserByToken(r.HeadersIn[comm.HEADER_AUTH_TOKEN])
	info, err := storage.GetSystemLog(start, end, size, (page-1)*size, filter, userName)
	if err != nil {
		r.Logger().Error("GetSysLog failed",
			pigeon.Field("start", start),
			pigeon.Field("end", end),
			pigeon.Field("filter", filter),
			pigeon.Field("page", page),
			pigeon.Field("size", size),
			pigeon.Field("error", err),
			pigeon.Field("requestId", r.HeadersIn[comm.HEADER_REQUEST_ID]))
		return nil, errno.GET_SYSTEM_LOG_FAILED
	}
	return info, errno.OK
}
