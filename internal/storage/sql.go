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

const (
	USER_ADMIN_NAME     = "admin"
	USER_ADMIN_PASSWORD = "curve"

	ADMIN_PERM      = 0
	READ_PERM       = 1
	READ_WRITE_PERM = 2
)

var (
	// table user
	CREATE_USER_TABLE = `
		CREATE TABLE IF NOT EXISTS user (
			username TEXT NOT NULL PRIMARY KEY,
			password TEXT NOT NULL,
			email TEXT,
			permission INTERER NOT NULL
		)
	`
	CREATE_ADMIN = `INSERT OR IGNORE INTO user(username, password, email, permission) VALUES(?, ?, ?, ?)`

	CREATE_USER = `INSERT INTO user(username, password, email, permission) VALUES(?, ?, ?, ?)`

	DELETE_USER = `DELETE FROM user WHERE username = ?`

	GET_USER = `SELECT * FROM user WHERE username = ?`

	GET_USER_EMAIL = `SELECT email FROM user WHERE username = ?`

	LIST_USER = `SELECT username, email, permission FROM user WHERE username != ?`

	UPDATE_PASSWORD = `UPDATE user SET password = ? WHERE username = ?`

	UPDATE_USER_INFO = `UPDATE user SET email = ?, permission = ? WHERE username = ?`
)
