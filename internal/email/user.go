package email

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/opencurve/pigeon"
)

const (
	EMAIL_ADDRESS = "email.addr"
	EMAIL_AUTH    = "email.auth"

	EMAIL_SERVER_163 = "smtp.163.com"
	EMAIL_SERVER_163_ENDPOINT = "smtp.163.com:25"
	EMAIL_SUBJECT_RESET_PASSWORD = "Curve-Manager Reset Password"
)

var (
	email_addr string
	email_auth string
)

func Init(cfg *pigeon.Configure) {
	email_addr = cfg.GetConfig().GetString(EMAIL_ADDRESS)
	email_auth = cfg.GetConfig().GetString(EMAIL_AUTH)
}

func SendNewPassWord(name, to, passwd string) error {
	if email_addr == "" || email_auth == "" {
		return fmt.Errorf("the manage email info not set")
	}
	content := fmt.Sprintf("The password of Curve-Manager has been reset successfully.\n" +
	"UserName: %s\nNewPassWord: %s\n", name, passwd)

	e := email.Email{
		From: email_addr,
		To: []string{to},
		Subject: EMAIL_SUBJECT_RESET_PASSWORD,
		Text: []byte(content),
	}
	return e.Send(EMAIL_SERVER_163_ENDPOINT, smtp.PlainAuth("", email_addr, email_auth, EMAIL_SERVER_163))
}
