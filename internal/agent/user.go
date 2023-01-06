package agent

import "github.com/opencurve/curve-manager/internal/storage"

func Login(name, passwd string) (interface{}, error) {
	return storage.Login(name, passwd)
}

func CreateUser(name, passwd, email string, permission int) error {
	return storage.CreateUser(name, passwd, email, permission)
}

func DeleteUser(name string) error {
	return storage.DeleteUser(name)
}

func ChangePassWord(name, passwd string) error {
	return storage.ChangePassWord(name, passwd)
}

func UpdateUserInfo(name, email string, permission int) error {
	return storage.UpdateUserInfo(name, email, permission)
}

func ListUser() (interface{}, error) {
	return storage.ListUser()
}