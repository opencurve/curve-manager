package storage

import (
	"crypto/md5"
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

type UserInfo struct {
	UserName   string `json:"userName" binding:"required"`
	Email      string `json:"email"`
	Permission int    `json:"permission" binding:"required"`
}

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
	if err = gStorage.execSQL(CREATE_USER_TABLE); err != nil {
		return err
	}

	// create admin user
	if err = createAdminUser(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
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

func createAdminUser() error {
	passwd := fmt.Sprintf("%x", md5.Sum([]byte(USER_ADMIN_PASSWORD)))
	return gStorage.execSQL(CREATE_ADMIN, USER_ADMIN_NAME, passwd, "", ADMIN_PERM)
}

func Login(name, passwd string) (interface{}, error) {
	rows, err := gStorage.db.Query(GET_USER, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var user UserInfo
	var passWord string
	if rows.Next() {
		err = rows.Scan(&user.UserName, &passWord, &user.Email, &user.Permission)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("user not exist")
	}

	if passWord == passwd {
		return &user, nil
	}
	return nil, fmt.Errorf(fmt.Sprintf("passwd mismatch, storedPasswd=%s, inPasswd=%s", passWord, passwd))
}

func CreateUser(name, passwd, email string, permission int) error {
	return gStorage.execSQL(CREATE_USER, name, passwd, email, permission)
}

func DeleteUser(name string) error {
	return gStorage.execSQL(DELETE_USER, name)
}

func ChangePassWord(name, passwd string) error {
	return gStorage.execSQL(UPDATE_PASSWORD, passwd, name)
}

func UpdateUserInfo(name, email string, permission int) error {
	return gStorage.execSQL(UPDATE_USER_INFO, email, permission, name)
}

func ListUser() (interface{}, error) {
	rows, err := gStorage.db.Query(LIST_USER, USER_ADMIN_NAME)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []UserInfo
	for rows.Next() {
		var user UserInfo
		err = rows.Scan(&user.UserName, &user.Email, &user.Permission)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return &users, nil
}
