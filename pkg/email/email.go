package email

import (
	"bytes"
	"fmt"
	"net/smtp"
	"wormhole/pkg/config"
)

var auth smtp.Auth
var smtpAddr string

func Init() {
	initTemplates()

	auth = smtp.PlainAuth("", config.Get().EmailUsername, config.Get().EmailPassword, config.Get().EmailHost)
	smtpAddr = fmt.Sprintf("%s:%s", config.Get().EmailHost, config.Get().EmailPort)
}

func Send(to string, email *Template, data map[string]any) error {
	data["From"] = config.Get().EmailUsername
	data["To"] = to
	data["Subject"] = email.Subject
	data["ApiBase"] = config.Get().ApiBase

	var contentBuf bytes.Buffer
	err := email.Formatter.Execute(&contentBuf, data)
	if err != nil {
		return err
	}

	data["Content"] = contentBuf.String()

	var baseBuf bytes.Buffer
	err = baseTemplate.Execute(&baseBuf, data)
	if err != nil {
		return err
	}

	return smtp.SendMail(smtpAddr, auth, config.Get().EmailUsername, []string{to}, baseBuf.Bytes())
}
