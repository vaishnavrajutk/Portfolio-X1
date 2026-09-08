package mailer

import (
	"fmt"
	"log"
	"net/smtp"

	"portfolio-backend/internal/config"
)

type Mailer interface {
	Send(name, email, message string) error
}

type SMTPMailer struct {
	host     string
	port     string
	username string
	password string
	from     string
	to       string
}

func NewSMTP(cfg config.Config) *SMTPMailer {
	return &SMTPMailer{
		host:     cfg.SMTPHost,
		port:     cfg.SMTPPort,
		username: cfg.SMTPUsername,
		password: cfg.SMTPPassword,
		from:     cfg.ContactFrom,
		to:       cfg.ContactTo,
	}
}

func (m *SMTPMailer) Send(name, email, message string) error {
	addr := fmt.Sprintf("%s:%s", m.host, m.port)
	auth := smtp.PlainAuth("", m.username, m.password, m.host)

	subject := fmt.Sprintf("Subject: Portfolio contact from %s\r\n", name)
	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nReply-To: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n", m.from, m.to, email)
	body := fmt.Sprintf("Name: %s\r\nEmail: %s\r\n\r\n%s\r\n", name, email, message)
	msg := []byte(subject + headers + "\r\n" + body)

	return smtp.SendMail(addr, auth, m.from, []string{m.to}, msg)
}

// ConsoleMailer just logs the message. Used when SMTP isn't configured so the
// contact form still works during local development.
type ConsoleMailer struct{}

func (ConsoleMailer) Send(name, email, message string) error {
	log.Printf("[contact] from %s <%s>: %s", name, email, message)
	return nil
}
