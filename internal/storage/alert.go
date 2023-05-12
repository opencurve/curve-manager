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

	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/email"
)

type Alert struct {
	Id          int64  `json:"id"`
	ClusterId   int    `json:"-"`
	TimeMs      int64  `json:"-"`
	Time        string `json:"time"`
	Level       int    `json:"-"`
	LevelStr    string `json:"level"`
	Name        string `json:"name"`
	DurationSec uint32 `json:"duration"`
	Summary     string `json:"summary"`
}

type AlertInfo struct {
	Total int64   `json:"total"`
	Info  []Alert `json:"info"`
}

func AddAlert(alert *Alert) error {
	return gStorage.execSQL(ADD_SYSTEM_ALERT, alert.ClusterId, alert.TimeMs, alert.Level, alert.Name, alert.DurationSec, alert.Summary)
}

func DeleteAlert(clusterId int, expirationMs int64) error {
	return gStorage.execSQL(DELETE_SYSTEM_ALERT, clusterId, expirationMs)
}

func SendAlert(clusterId int, alert *Alert) error {
	users, err := GetAlertUser(clusterId, alert.Name)
	if err != nil {
		return err
	}
	emails := []string{}
	var errors error
	for _, user := range users {
		email, err := GetUserEmail(user)
		if err != nil {
			errors = fmt.Errorf("user: %s, error: %s, %s", user, err.Error(), errors.Error())
		} else {
			emails = append(emails, email)
		}
	}
	if len(emails) != 0 {
		content := fmt.Sprintf("Alert Info:\nCluster ID: %d\nName: %s\nLevel: %s\nDuration Second: %d\nSummary: %s\nTime: %s\n",
			alert.ClusterId, alert.Name, getLevelStr(alert.Level), alert.DurationSec, alert.Summary, common.Mill2TimeStr(alert.TimeMs))
		err = email.SendAlert2Users(content, emails)
		if err != nil {
			errors = fmt.Errorf("send email error: %s, %s", err.Error(), errors.Error())
		}
	}
	return errors
}

func GetAlert(clusterId int, start, end int64, limit, offset uint32, name, level, filter string) (AlertInfo, error) {
	alerts := AlertInfo{}
	filter = fmt.Sprintf("%%%s%%", filter)
	numSql := GET_SYSTEM_ALERT_NUM
	contSql := GET_SYSTEM_ALERT
	l := 0
	params := []interface{}{start, end, clusterId}
	if name != "" {
		numSql += " AND name = ?"
		contSql += " AND name = ?"
		params = append(params, name)
	}
	if level != "" {
		switch level {
		case CRITICAL:
			l = 1
		case WARNING:
			l = 2
		}
		numSql += " AND level = ?"
		contSql += " AND level = ?"
		params = append(params, l)
	}
	if filter != "" {
		numSql += " AND summary like ?"
		contSql += " AND summary like ?"
		params = append(params, filter)
	}
	// get total
	rows, err := gStorage.querySQL(numSql, params...)
	if err != nil {
		return alerts, err
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&alerts.Total)
	}

	if alerts.Total > int64(offset) {
		rows.Close()
		contSql += " ORDER BY timestamp DESC LIMIT ? OFFSET ?"
		params = append(params, limit, offset)
		rows, err = gStorage.querySQL(contSql, params...)
		if err != nil {
			return alerts, err
		}
		for rows.Next() {
			var alert Alert
			err = rows.Scan(&alert.Id, &alert.ClusterId, &alert.TimeMs, &alert.Level, &alert.Name, &alert.DurationSec, &alert.Summary)
			if err != nil {
				return alerts, err
			}
			alert.Time = common.Mill2TimeStr(alert.TimeMs)
			alert.LevelStr = getLevelStr(alert.Level)
			alerts.Info = append(alerts.Info, alert)
		}
	}
	return alerts, nil
}

func UpdateReadAlertId(id int64, name string) error {
	return gStorage.execSQL(UPDATE_READ_ALERT_ID, id, name)
}

func GetUnreadAlertNum(clusterId int, readId int64) (int64, error) {
	var ret int64 = 0
	rows, err := gStorage.querySQL(GET_UNREAD_SYSTEM_ALERT_NUM, readId, clusterId)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&ret)
		if err != nil {
			return 0, err
		}
	}
	return ret, nil
}

func GetReadAlertId(userName string) (int64, error) {
	var readId int64 = -1
	rows, err := gStorage.querySQL(GET_READ_ALERT_ID, userName)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&readId)
		if err != nil {
			return 0, err
		}
	}
	return readId, nil
}

func AddReadAlertId(userName string) error {
	return gStorage.execSQL(ADD_READ_ALERT_ID, userName, 0)
}
