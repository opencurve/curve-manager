package storage

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/opencurve/pigeon"
)

var (
	gStorage *Storage
)

const (
	SQLITE_DB_FILE = "db.sqlite.filepath"
)

type Storage struct {
	db    *sql.DB
	mutex *sync.Mutex
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
	gStorage = &Storage{db: db, mutex: &sync.Mutex{}}

	// init user table
	if err := gStorage.execSQL(CREATE_USER_TABLE); err != nil {
		return err
	}

	return nil
}

func (s *Storage) execSQL(query string, args ...interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(args...)
	return err
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func InsertUser(name, passwd string, permission int) error {
	return gStorage.execSQL(INSERT_USER, name, passwd, permission)
}
