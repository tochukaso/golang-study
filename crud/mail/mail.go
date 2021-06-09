package mail

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"text/template"

	"github.com/tochukaso/golang-study/env"
	"github.com/tochukaso/golang-study/model"
)

func SendMail(from string, to []string, subject, mailTemplate string, variables interface{}) error {

	env := env.GetEnv()
	if env.SmtpHost == "" {
		log.Println("SMTPのホストが未設定のめメールは送信しません。")
		return nil
	}
	if env.SmtpHost == "localhost" {
		error := usePostfix(from, to, subject, mailTemplate, variables)
		log.Println("postfix err", error)
		return error
	}
	// Authentication.
	auth := authenticate()

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	var body bytes.Buffer
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", subject, mimeHeaders)))

	t, err := template.New("Mail").Parse(mailTemplate)
	if err != nil {
		log.Panic(err)
	}
	t.Execute(&body, variables)

	err = smtp.SendMail(env.SmtpHost+":"+env.SmtpPort, auth, env.SmtpFrom, to, body.Bytes())
	if err != nil {
		log.Println("sendmail error", err)
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

func SendProductRegisterMail(product model.Product) {
	mailTemplate := model.GetMailTemplateFromCode(int(model.ProductRegister))

	SendMail(mailTemplate.From,
		[]string{mailTemplate.Cc},
		mailTemplate.Subject,
		mailTemplate.MailDetail,
		product,
	)
}

func authenticate() smtp.Auth {
	env := env.GetEnv()
	return smtp.PlainAuth("", env.SmtpFrom, env.SmtpPassword, env.SmtpHost)
}

func usePostfix(from string, to []string, subject, mailTemplate string, variables interface{}) error {
	env := env.GetEnv()
	toHeader := strings.Join(to, ", ")
	header := make(map[string]string)
	header["To"] = toHeader
	header["From"] = from
	header["Subject"] = subject
	header["Content-Type"] = `text/html; charset="UTF-8"`
	msg := ""

	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	c, err := smtp.Dial(env.SmtpHost + ":" + env.SmtpPort)

	if err != nil {
		return err
	}

	defer c.Close()
	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}

	t, err := template.New("Mail").Parse(mailTemplate)
	if err != nil {
		log.Panic(err)
	}

	t.Execute(w, variables)

	err = w.Close()
	if err != nil {
		return err
	}

	err = c.Quit()
	if err != nil {
		return err
	}
	return nil
}
