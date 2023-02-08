package agent

import (
	"github.com/opencurve/curve-manager/internal/common"
	"github.com/opencurve/curve-manager/internal/storage"
	"github.com/opencurve/curve-manager/internal/email"
)

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

func ResetPassWord(name string) error {
	emailAddr, err := storage.GetUserEmail(name)
	if err != nil {
		return err
	}
	passwd := storage.GetNewPassWord()
	err = ChangePassWord(name, common.GetMd5Sum32Little(passwd))
	if err != nil {
		return err
	}

	err = email.SendNewPassWord(name, emailAddr, passwd)
	return err
}

func UpdateUserInfo(name, email string, permission int) error {
	return storage.UpdateUserInfo(name, email, permission)
}

func ListUser() (interface{}, error) {
	return storage.ListUser()
}
