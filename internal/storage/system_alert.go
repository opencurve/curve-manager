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
* Created Date: 2023-03-21
* Author: wanghai (SeanHai)
 */

package storage

import (
	"fmt"
	"time"

	"github.com/opencurve/curve-manager/internal/common"
)

const (
	ALERT_CRITICAL = 1
	ALERT_WARNING  = 2
)

type SystemAlert struct {
	id          int64
	TimeMs      int64  `json:"-"`
	Time        string `json:"time"`
	Level       int    `json:"-"`
	LevelStr    string `json:"level"`
	Module      string `json:"module"`
	DurationSec uint32 `json:"duration"`
	Summary     string `json:"summary"`
}

type SystemAlertInfo struct {
	Total int           `json:"total"`
	Info  []SystemAlert `json:"info"`
}

func getLevelStr(l int) string {
	switch l {
	case ALERT_CRITICAL:
		return "critical"
	case ALERT_WARNING:
		return "warning"
	default:
		return "invalid"
	}
}

func AddSystemAlert(alert *SystemAlert) error {
	return gStorage.execSQL(ADD_SYSTEM_ALERT, alert.TimeMs, alert.Level, alert.Module, alert.DurationSec, alert.Summary)
}

func DeleteSystemAlert(expirationMs int64) error {
	return gStorage.execSQL(DELETE_SYSTEM_ALERT, expirationMs)
}

func GetSystemAlert(start, end int64, limit, offset uint32, filter string) (SystemAlertInfo, error) {
	alerts := SystemAlertInfo{}
	filter = fmt.Sprintf("%%%s%%", filter)
	// get total
	rows, err := gStorage.db.Query(GET_SYSTEM_ALERT_NUM, start, end, filter)
	if err != nil {
		return alerts, err
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&alerts.Total)
	}

	if alerts.Total > int(offset) {
		rows.Close()
		rows, err = gStorage.db.Query(GET_SYSTEM_ALERT, start, end, filter, limit, offset)
		if err != nil {
			return alerts, err
		}
		for rows.Next() {
			var alert SystemAlert
			err = rows.Scan(&alert.TimeMs, &alert.Level, &alert.Module, &alert.DurationSec, &alert.Summary)
			if err != nil {
				return alerts, err
			}
			sec := alert.TimeMs / 1000
			t := time.Unix(sec, alert.TimeMs%1000*int64(time.Millisecond))
			alert.Time = t.Format(common.TIME_MS_FORMAT)
			alert.LevelStr = getLevelStr(alert.Level)
			alerts.Info = append(alerts.Info, alert)
		}
	}
	return alerts, nil
}

func UpdateReadAlertId(id int64) error {
	return gStorage.execSQL(UPDATE_USER_SYSTEM_ALERT_ID, id)
}

func GetNotReadAlertNum(userName string) (int64, error) {
	var maxId int64
	var readId int64
	rows, err := gStorage.db.Query(GET_LAST_SYSTEM_ALERT_ID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&maxId)
		if err != nil {
			return 0, err
		}
	}

	rows.Close()
	rows, err = gStorage.db.Query(GET_USER_SYSTEM_ALERT_ID, userName)
	if err != nil {
		return 0, err
	}
	if rows.Next() {
		err = rows.Scan(&readId)
		if err != nil {
			return 0, err
		}
	}
	return maxId - readId, nil
}
