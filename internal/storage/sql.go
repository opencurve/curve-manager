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

/*
 * +------------------------------+
 * |go        | sqlite3           |
 * |----------|-------------------|
 * |nil       | null              |
 * |int       | integer           |
 * |int64     | integer           |
 * |float64   | float             |
 * |bool      | integer           |
 * |[]byte    | blob              |
 * |string    | text              |
 * |time.Time | timestamp/datetime|
 * +------------------------------+
 */

package storage

var (
	// user table
	CREATE_USER_TABLE = `
		CREATE TABLE IF NOT EXISTS user (
			username TEXT NOT NULL PRIMARY KEY,
			password TEXT NOT NULL,
			email TEXT,
			permission INTEGER NOT NULL
		)
	`
	// system log table
	CREATE_SYSTEM_LOG_TABLE = `
		CREATE TABLE IF NOT EXISTS system_log (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp INTEGER,
			ip TEXT,
			user TEXT,
			module TEXT,
			method TEXT,
			error_code INTEGER,
			error_msg TEXT,
			content TEXT
		)
	`
	// alert conf table
	CREATE_ALERT_CONF_TABLE = `
		CREATE TABLE IF NOT EXISTS alert_conf (
			cluster INTEGER,
			name TEXT,
			level INTEGER,
			interval INTEGER,
			times INTERGER,
			enable INTEGER CHECK (enable IN (0, 1)),
			rule TEXT,
			desc TEXT,
			PRIMARY KEY (cluster, name)
		)
	`
	// alert user table
	CREATE_ALERT_USER_TABLE = `
		CREATE TABLE IF NOT EXISTS alert_user (
			cluster INTERGER,
			alert TEXT,
			user TEXT,
			UNIQUE (cluster, alert, user) ON CONFLICT IGNORE
		)
	`
	// alert table
	CREATE_ALERT_TABLE = `
		CREATE TABLE IF NOT EXISTS alert (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			cluster INTERGER,
			timestamp INTEGER,
			level INTEGER,
			name TEXT,
			duration INTEGER,
			summary TEXT
		)
	`
	// read alert table
	CREATE_READ_ALERT_TABLE = `
		CREATE TABLE IF NOT EXISTS read_alert (
			username TEXT PRIMARY KEY,
			id INTEGER
		)
	`

	// user
	CREATE_ADMIN           = `INSERT OR IGNORE INTO user(username, password, email, permission) VALUES(?, ?, ?, ?)`
	CREATE_USER            = `INSERT INTO user(username, password, email, permission) VALUES(?, ?, ?, ?)`
	DELETE_USER            = `DELETE FROM user WHERE username=?`
	GET_USER               = `SELECT * FROM user WHERE username=?`
	GET_USER_EMAIL         = `SELECT email FROM user WHERE username=?`
	GET_USER_PASSWORD      = `SELECT password FROM user WHERE username=?`
	LIST_USER              = `SELECT * FROM user WHERE username!=?`
	LIST_USER_WITH_EMAIL   = `SELECT username FROM user WHERE email!=?`
	UPDATE_USER_PASSWORD   = `UPDATE user SET password=? WHERE username=?`
	UPDATE_USER_EMAIL      = `UPDATE user SET email=? WHERE username=?`
	UPDATE_USER_PERMISSION = `UPDATE user SET permission=? WHERE username=?`

	// system log
	ADD_SYSTEM_LOG = `INSERT OR IGNORE INTO system_log(timestamp, ip, user, module, method, error_code, error_msg, content)
	 VALUES(?, ?, ?, ?, ?, ?, ?, ?)`
	GET_SYSTEM_LOG_NUM = `SELECT COUNT(*) FROM system_log WHERE timestamp>=? AND timestamp<=? AND
	 ip||user||module||method||error_code||error_msg||content LIKE ?`
	GET_SYSTEM_LOG = `SELECT * FROM system_log WHERE timestamp>=? AND timestamp<=? AND
	 ip||user||module||method||error_code||error_msg||content LIKE ? ORDER BY timestamp DESC LIMIT ? OFFSET ?`
	GET_SYSTEM_LOG_NUM_OF_USER = `SELECT COUNT(*) FROM system_log WHERE timestamp>=? AND timestamp<=? AND user=? AND
	 ip||user||module||method||error_code||error_msg||content LIKE ?`
	GET_SYSTEM_LOG_OF_USER = `SELECT * FROM system_log WHERE timestamp >= ? AND timestamp<=? AND user=? AND
	 ip||user||module||method||error_code||error_msg||content LIKE ? ORDER BY timestamp DESC LIMIT ? OFFSET ?`
	DELETE_SYSTEM_LOG = `DELETE FROM system_log WHERE timestamp<?`

	// alert conf
	ADD_ALERT_CONF    = `INSERT OR IGNORE INTO alert_conf(cluster, name, level, interval, times, enable, rule, desc) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`
	UPDATE_ALERT_CONF = `UPDATE alert_conf SET interval=?, times=?, enable=?, rule=? WHERE cluster=? AND name=?`
	GET_ALERT_CONF    = `SELECT * FROM alert_conf WHERE cluster=? ORDER BY name ASC`

	// alert user
	ADD_ALERT_USER    = `INSERT OR IGNORE INTO alert_user(cluster, alert, user) VALUES(?, ?, ?)`
	DELETE_ALERT_USER = `DELETE FROM alert_user WHERE cluster=? AND alert=? AND user=?`
	GET_ALERT_USER    = `SELECT user FROM alert_user WHERE cluster=? AND alert=?`

	// alert
	ADD_SYSTEM_ALERT            = `INSERT OR IGNORE INTO alert(cluster, timestamp, level, name, duration, summary) VALUES(?, ?, ?, ?, ?, ?)`
	GET_SYSTEM_ALERT_NUM        = `SELECT COUNT(*) FROM alert WHERE timestamp>=? AND timestamp<=? AND cluster=?`
	GET_SYSTEM_ALERT            = `SELECT * FROM alert WHERE timestamp>=? AND timestamp<=? AND cluster=?`
	DELETE_SYSTEM_ALERT         = `DELETE FROM alert WHERE cluster=? AND timestamp<?`
	GET_UNREAD_SYSTEM_ALERT_NUM = `SELECT COUNT(*) FROM alert WHERE id>? AND cluster=?`

	// read alert
	ADD_READ_ALERT_ID    = `INSERT OR IGNORE INTO read_alert(username, id) VALUES(?, ?)`
	GET_READ_ALERT_ID    = `SELECT id FROM read_alert WHERE username=?`
	UPDATE_READ_ALERT_ID = `UPDATE read_alert SET id=? WHERE username=?`
)
