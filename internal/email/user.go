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

package email

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/opencurve/pigeon"
)

const (
	SMTP_ADDRESS                 = "smtp.address"
	SMTP_PORT                    = "smtp.port"
	SMTP_USERNAME                = "smtp.username"
	SMTP_PASSWORD                = "smtp.password"
	EMAIL_SUBJECT_RESET_PASSWORD = "Curve-Manager Reset Password"
	EMAIL_SUBJECT_ALERT          = "Curve-Manager Alert"
)

var (
	smtp_address  string
	smtp_port     string
	smtp_username string
	smtp_password string
)

func Init(cfg *pigeon.Configure) {
	smtp_address = cfg.GetConfig().GetString(SMTP_ADDRESS)
	smtp_port = cfg.GetConfig().GetString(SMTP_PORT)
	smtp_username = cfg.GetConfig().GetString(SMTP_USERNAME)
	smtp_password = cfg.GetConfig().GetString(SMTP_PASSWORD)
}

func SendNewPassWord(name, to, passwd string) error {
	if smtp_address == "" || smtp_port == "" || smtp_username == "" || smtp_password == "" {
		return fmt.Errorf("the manage email info not set")
	}
	content := fmt.Sprintf("The password of Curve-Manager has been reset successfully.\n"+
		"UserName: %s\nNewPassWord: %s\n", name, passwd)

	e := email.Email{
		From:    smtp_username,
		To:      []string{to},
		Subject: EMAIL_SUBJECT_RESET_PASSWORD,
		Text:    []byte(content),
	}
	return e.Send(fmt.Sprintf("%v:%v", smtp_address, smtp_port),
		smtp.PlainAuth("", smtp_username, smtp_password, smtp_address))
}

func SendAlert2Users(content string, tos []string) error {
	if smtp_address == "" || smtp_port == "" || smtp_username == "" || smtp_password == "" {
		return fmt.Errorf("the manage email info not set")
	}
	var errors error
	for _, to := range tos {
		e := email.Email{
			From:    smtp_username,
			To:      []string{to},
			Subject: EMAIL_SUBJECT_ALERT,
			Text:    []byte(content),
		}
		err := e.Send(fmt.Sprintf("%v:%v", smtp_address, smtp_port),
			smtp.PlainAuth("", smtp_username, smtp_password, smtp_address))
		if err != nil {
			errors = fmt.Errorf("dest: %s, error: %s, %s", to, err.Error(), errors.Error())
		}
	}
	return errors
}
