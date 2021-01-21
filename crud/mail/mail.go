package mail

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"text/template"

	"omori.jp/env"
	"omori.jp/model"
)

func SendMail(from string, to []string, subject, mailTemplate string, variables interface{}) error {

	env := env.GetEnv()
	// Authentication.
	auth := authenticate()

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	var body bytes.Buffer
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", subject, mimeHeaders)))

	t, err := template.New("Mail").Parse(mailTemplate)
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(&body, variables)

	err = smtp.SendMail(env.SmtpHost+":"+env.SmtpPort, auth, env.SmtpFrom, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func SendUserRegisterMail(user model.User) {
	mailTemplate := model.GetMailTemplateFromCode(int(model.UserRegister))

	SendMail(mailTemplate.From,
		[]string{user.Mail},
		mailTemplate.Subject,
		mailTemplate.MailDetail,
		user,
	)
}

func authenticate() smtp.Auth {
	env := env.GetEnv()
	return smtp.PlainAuth("", env.SmtpFrom, env.SmtpPassword, env.SmtpHost)
}
