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
* Created Date: 2023-03-24
* Author: wanghai (SeanHai)
 */

package storage

const (
	ALERT_CRITICAL = 1
	ALERT_WARNING  = 2

	CRITICAL = "critical"
	WARNING  = "warning"
	INVALID  = "invalid"
)

func getLevelStr(l int) string {
	switch l {
	case ALERT_CRITICAL:
		return CRITICAL
	case ALERT_WARNING:
		return WARNING
	default:
		return INVALID
	}
}

type AlertConf struct {
	id         int64
	Name       string
	Level      int
	LevelStr   string
	Interval   uint32
	Times      uint32
	Enable     int
	EnableBool bool
	Rule       string
	Desc       string
}

func AddAlertConf(conf *AlertConf) error {
	return gStorage.execSQL(ADD_ALERT_CONF, conf.Name, conf.Level, conf.Interval, conf.Times, conf.Enable, conf.Rule, conf.Desc)
}

func UpdateAlertConf(conf *AlertConf) error {
	return gStorage.execSQL(UPDATE_ALERT_CONF, conf.Interval, conf.Times, conf.Enable, conf.Rule, conf.Name)
}

func GetAlertConf() ([]AlertConf, error) {
	info := []AlertConf{}
	rows, err := gStorage.querySQL(GET_ALERT_CONF)
	if err != nil {
		return info, err
	}
	defer rows.Close()
	for rows.Next() {
		item := AlertConf{}
		err := rows.Scan(&item.id, &item.Name, &item.Level, &item.Interval, &item.Times, &item.Enable, &item.Rule, &item.Desc)
		if err != nil {
			return nil, err
		}
		item.LevelStr = getLevelStr(item.Level)
		item.EnableBool = item.Enable == 1
		info = append(info, item)
	}
	return info, nil
}
