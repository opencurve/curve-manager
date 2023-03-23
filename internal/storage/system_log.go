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

package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/opencurve/curve-manager/internal/common"
)

type SystemLog struct {
	id        int64
	TimeMs    int64  `json:"-"`
	Time      string `json:"time"`
	IP        string `json:"ip"`
	User      string `json:"user"`
	Module    string `json:"module"`
	Method    string `json:"method"`
	ErrorCode int    `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	Content   string `json:"content"`
}

type SystemLogInfo struct {
	Total int64         `json:"total"`
	Info  []SystemLog `json:"info"`
}

func AddSystemLog(log *SystemLog) error {
	return gStorage.execSQL(ADD_SYSTEM_LOG, log.TimeMs, log.IP, log.User, log.Module, log.Method, log.ErrorCode,
		log.ErrorMsg, log.Content)
}

func DeleteSystemLog(expirationMs int64) error {
	return gStorage.execSQL(DELETE_SYSTEM_LOG, expirationMs)
}

func GetSystemLog(start, end int64, limit, offset uint32, filter, userName string) (SystemLogInfo, error) {
	logs := SystemLogInfo{}
	filter = fmt.Sprintf("%%%s%%", filter)
	// get total
	// admin can get all system operation log, others can only get themself's
	var rows *sql.Rows
	var err error
	if userName == USER_ADMIN_NAME {
		rows, err = gStorage.querySQL(GET_SYSTEM_LOG_NUM, start, end, filter)
	} else {
		rows, err = gStorage.querySQL(GET_SYSTEM_LOG_NUM_OF_USER, start, end, userName, filter)
	}
	if err != nil {
		return logs, err
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&logs.Total)
	}

	if logs.Total > int64(offset) {
		rows.Close()
		if userName == USER_ADMIN_NAME {
			rows, err = gStorage.querySQL(GET_SYSTEM_LOG, start, end, filter, limit, offset)
		} else {
			rows, err = gStorage.querySQL(GET_SYSTEM_LOG_OF_USER, start, end, userName, filter, limit, offset)
		}
		if err != nil {
			return logs, err
		}
		for rows.Next() {
			var log SystemLog
			err = rows.Scan(&log.id, &log.TimeMs, &log.IP, &log.User, &log.Module, &log.Method, &log.ErrorCode,
				&log.ErrorMsg, &log.Content)
			if err != nil {
				return logs, err
			}
			sec := log.TimeMs / 1000
			t := time.Unix(sec, log.TimeMs%1000*int64(time.Millisecond))
			log.Time = t.Format(common.TIME_MS_FORMAT)
			logs.Info = append(logs.Info, log)
		}
	}
	return logs, nil
}
