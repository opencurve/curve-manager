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
	// system alert table
	CREATE_SYSTEM_ALERT_TABLE = `
		CREATE TABLE IF NOT EXISTS system_alert (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp INTEGER,
			level INTEGER,
			module TEXT,
			duration INTEGER,
			summary TEXT
		)
	`
	// user read system alert table
	CREATE_USER_SYSTEM_LOG_TABLE = `
		CREATE TABLE IF NOT EXISTS user_system_alert (
			username TEXT PRIMARY KEY,
			id INTEGER
		)
	`

	// user
	CREATE_ADMIN           = `INSERT OR IGNORE INTO user(username, password, email, permission) VALUES(?, ?, ?, ?)`
	CREATE_USER            = `INSERT INTO user(username, password, email, permission) VALUES(?, ?, ?, ?)`
	DELETE_USER            = `DELETE FROM user WHERE username = ?`
	GET_USER               = `SELECT * FROM user WHERE username = ?`
	GET_USER_EMAIL         = `SELECT email FROM user WHERE username = ?`
	GET_USER_PASSWORD      = `SELECT password FROM user WHERE username = ?`
	LIST_USER              = `SELECT * FROM user WHERE username != ?`
	UPDATE_USER_PASSWORD   = `UPDATE user SET password = ? WHERE username = ?`
	UPDATE_USER_EMAIL      = `UPDATE user SET email = ? WHERE username = ?`
	UPDATE_USER_PERMISSION = `UPDATE user SET permission = ? WHERE username = ?`

	// system log
	ADD_SYSTEM_LOG = `INSERT INTO system_log(timestamp, ip, user, module, method, error_code, error_msg, content)
	 VALUES(?, ?, ?, ?, ?, ?, ?, ?)`
	GET_SYSTEM_LOG_NUM = `SELECT COUNT(*) FROM system_log WHERE timestamp >= ? AND timestamp <= ? AND
	 ip||user||module||method||error_code||error_msg||content LIKE ?`
	GET_SYSTEM_LOG = `SELECT * FROM system_log WHERE timestamp >= ? AND timestamp <= ? AND
	 ip||user||module||method||error_code||error_msg||content LIKE ? ORDER BY timestamp DESC LIMIT ? OFFSET ?`
	GET_SYSTEM_LOG_NUM_OF_USER = `SELECT COUNT(*) FROM system_log WHERE timestamp >= ? AND timestamp <= ? AND user = ? AND
	 ip||user||module||method||error_code||error_msg||content LIKE ?`
	GET_SYSTEM_LOG_OF_USER = `SELECT * FROM system_log WHERE timestamp >= ? AND timestamp <= ? AND user = ? AND
	 ip||user||module||method||error_code||error_msg||content LIKE ? ORDER BY timestamp DESC LIMIT ? OFFSET ?`
	DELETE_SYSTEM_LOG = `DELETE FROM system_log WHERE timestamp < ?`

	// system alert
	ADD_SYSTEM_ALERT     = `INSERT INTO system_alert(timestamp, level, module, duration, summary) VALUES(?, ?, ?, ?, ?)`
	GET_SYSTEM_ALERT_NUM = `SELECT COUNT(*) FROM system_alert WHERE timestamp >= ? AND timestamp <= ? AND
	 level||module||duration||summary LIKE ?`
	GET_SYSTEM_ALERT = `SELECT * FROM system_alert WHERE timestamp >= ? AND timestamp <= ? AND
	 level||module||duration||summary LIKE ? ORDER BY timestamp DESC LIMIT ? OFFSET ?`
	DELETE_SYSTEM_ALERT      = `DELETE FROM system_alert WHERE timestamp < ?`
	GET_LAST_SYSTEM_ALERT_ID = `SELECT MAX(id) from system_alert`

	ADD_READ_SYSTEM_ALERT_ID    = `INSERT INTO user_system_alert(username, id) VALUES(?, ?)`
	GET_READ_SYSTEM_ALERT_ID    = `SELECT id FROM user_system_alert WHERE username = ?`
	UPDATE_READ_SYSTEM_ALERT_ID = `UPDATE user_system_alert SET id = ? WHERE username = ?`
)
