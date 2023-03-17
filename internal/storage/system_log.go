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
	"time"

	"github.com/opencurve/curve-manager/internal/common"
)

type SystemLog struct {
	id        int64
	TimeMs    int64  `json:"-"`
	IP        string `json:"ip"`
	User      string `json:"user"`
	Module    string `json:"module"`
	Method    string `json:"method"`
	ErrorCode int    `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	Content   string `json:"content"`
	Time      string `json:"time"`
}

type SystemLogInfo struct {
	Total int         `json:"total"`
	Info  []SystemLog `json:"info"`
}

func AddSystemLog(log *SystemLog) error {
	return gStorage.execSQL(ADD_SYSTEM_LOG, log.TimeMs, log.IP, log.User, log.Module, log.Method, log.ErrorCode,
		log.ErrorMsg, log.Content)
}

func DeleteSystemLog(expirationMs int64) error {
	return gStorage.execSQL(DELETE_SYSTEM_LOG, expirationMs)
}

func GetSystemLog(start, end int64, limit, offset uint32) (SystemLogInfo, error) {
	logs := SystemLogInfo{}
	// get total
	num, err := gStorage.db.Query(GET_SYSTEM_LOG_NUM, start, end)
	if err != nil {
		return logs, err
	}
	defer num.Close()
	for num.Next() {
		num.Scan(&logs.Total)
	}

	if logs.Total > int(offset) {
		rows, err := gStorage.db.Query(GET_SYSTEM_LOG, start, end, limit, offset)
		if err != nil {
			return logs, err
		}
		defer rows.Close()
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
