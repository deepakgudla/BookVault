package notifications

import (
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strconv"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type SimpleEmail struct {
	To      string
	Subject string
	Body    string
}

type EmailNotifier struct {
	config *SMTPConfig
}

func NewEmailNotifier(config *SMTPConfig) *EmailNotifier {
	return &EmailNotifier{
		config: config,
	}
}

func (e *EmailNotifier) SendSimpleEmail(email *SimpleEmail) error {
	addr := net.JoinHostPort(e.config.Host, strconv.Itoa(e.config.Port))

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, e.config.Host)
	if err != nil {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close SMTP connection: %v", err)
		}
		return err
	}

	defer func() {
		if err := client.Quit(); err != nil {
			log.Printf("failed to quit SMTP client: %v", err)
		}
	}()

	if e.config.Username != "" && e.config.Password != "" {
		auth := smtp.PlainAuth("", e.config.Username, e.config.Password, e.config.Host)
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	if err := client.Mail(e.config.From); err != nil {
		return err
	}

	if err := client.Rcpt(email.To); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		e.config.From,
		email.To,
		email.Subject,
		email.Body,
	)

	_, err = w.Write([]byte(msg))
	if err != nil {
		if closeErr := w.Close(); closeErr != nil {
			log.Printf("failed to close email writer: %v", closeErr)
		}
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func (e *EmailNotifier) SendLoginNotification(userEmail, userName string) error {
	email := &SimpleEmail{
		To:      userEmail,
		Subject: "Login Notification",
		Body: fmt.Sprintf(
			"Hello %s,\n\nYou have successfully logged into your account.\nIf this was not you, please contact support immediately.\n\nBest Regards,\nMahishmathi Samrajyam",
			userName,
		),
	}

	return e.SendSimpleEmail(email)
}
