package email

import (
	"context"
	"fmt"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
)

var (
	smtpAddr     = os.Getenv("TRAVELREYS_SMTP_ADDR")
	smtpUsername = os.Getenv("TRAVELREYS_SMTP_USERNAME")
	smtpPassword = os.Getenv("TRAVELREYS_SMTP_PASSWORD")
)

type Service interface {
	SendMail(ctx context.Context, to, from, subj, body string) error
}

type service struct {
	addr     string
	username string
	password string
}

func NewDefaultService() Service {
	return &service{smtpAddr, smtpUsername, smtpPassword}
}

func (svc service) SendMail(ctx context.Context, to, from, subj, body string) error {
	toAddr := mail.Address{Address: to}
	fromAddr := mail.Address{Address: from}
	fromHeader := fromAddr.String()

	header := map[string]string{}
	header["To"] = strings.Join([]string{toAddr.String()}, ",")
	header["From"] = fromHeader
	header["Subject"] = subj
	header["Content-Type"] = `text/html; charset="UTF-8"`

	msg := ""
	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + body

	bMsg := []byte(msg)
	// Send using local postfix service
	c, err := smtp.Dial(smtpAddr)
	if err != nil {
		return err
	}
	defer c.Close()

	if err = c.Mail(fromHeader); err != nil {
		return err
	}
	if err = c.Rcpt(to); err != nil {
		return err
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err = w.Write(bMsg); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}

	return c.Quit()
}
