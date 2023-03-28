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
* Created Date: 2023-03-30
* Author: wanghai (SeanHai)
 */

package storage

import "fmt"

func AddAlertUser(alertName string, users []string) error {
	for _, user := range users {
		err := gStorage.execSQL(ADD_ALERT_USER, alertName, user)
		if err != nil {
			return fmt.Errorf("user: %s, error: %s", user, err)
		}
	}
	return nil
}

func DeleteAlertUser(alertName string, users []string) error {
	for _, user := range users {
		err := gStorage.execSQL(DELETE_ALERT_USER, alertName, user)
		if err != nil {
			return fmt.Errorf("user: %s, error: %s", user, err)
		}
	}
	return nil
}

func GetAlertUser(alertName string) ([]string, error) {
	users := []string{}
	rows, err := gStorage.querySQL(GET_ALERT_USER, alertName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user string
		err = rows.Scan(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
