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

package storage

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/opencurve/pigeon"
)

var (
	gStorage *storage
)

const (
	SQLITE_DB_FILE = "db.sqlite.filepath"
)

type storage struct {
	db      *sql.DB
	dbMutex *sync.RWMutex
	// token->sessionItem
	session map[string]sessionItem
	// user->token
	loginOnce map[string]string
	mutex     *sync.Mutex
	// only one person with write permission is allowed to login
	loginedWriteUser string
}

func Init(cfg *pigeon.Configure) error {
	dbfile := cfg.GetConfig().GetString(SQLITE_DB_FILE)
	if len(dbfile) == 0 {
		return fmt.Errorf("no sqlite db file found")
	}

	// new db
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return err
	}
	gStorage = &storage{db: db, dbMutex: &sync.RWMutex{}, session: make(map[string]sessionItem), loginedWriteUser: "",
		loginOnce: make(map[string]string), mutex: &sync.Mutex{}}

	// init user table
	if err = gStorage.execSQL(CREATE_USER_TABLE); err != nil {
		return err
	}
	// create admin user
	if err = createAdminUser(); err != nil {
		return err
	}

	// create system operation log table
	if err = gStorage.execSQL(CREATE_SYSTEM_LOG_TABLE); err != nil {
		return err
	}

	// create system alert conf table
	if err = gStorage.execSQL(CREATE_ALERT_CONF_TABLE); err != nil {
		return err
	}

	// create alert user table
	if err = gStorage.execSQL(CREATE_ALERT_USER_TABLE); err != nil {
		return err
	}

	// create system alert table
	if err = gStorage.execSQL(CREATE_ALERT_TABLE); err != nil {
		return err
	}

	// create user system alert table
	if err = gStorage.execSQL(CREATE_READ_ALERT_TABLE); err != nil {
		return err
	}
	return nil
}

func (s *storage) Close() error {
	return s.db.Close()
}

func (s *storage) execSQL(query string, args ...interface{}) error {
	s.dbMutex.Lock()
	defer s.dbMutex.Unlock()
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(args...)
	return err
}

func (s *storage) querySQL(query string, args ...interface{}) (*sql.Rows, error) {
	s.dbMutex.RLock()
	defer s.dbMutex.RUnlock()
	return s.db.Query(query, args...)
}
