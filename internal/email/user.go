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
	EMAIL_ADDRESS = "email.addr"
	EMAIL_AUTH    = "email.auth"

	EMAIL_SERVER_163             = "smtp.163.com"
	EMAIL_SERVER_163_ENDPOINT    = "smtp.163.com:25"
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
	content := fmt.Sprintf("The password of Curve-Manager has been reset successfully.\n"+
		"UserName: %s\nNewPassWord: %s\n", name, passwd)

	e := email.Email{
		From:    email_addr,
		To:      []string{to},
		Subject: EMAIL_SUBJECT_RESET_PASSWORD,
		Text:    []byte(content),
	}
	return e.Send(EMAIL_SERVER_163_ENDPOINT, smtp.PlainAuth("", email_addr, email_auth, EMAIL_SERVER_163))
}
