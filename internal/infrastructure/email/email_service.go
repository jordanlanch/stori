package email

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"net/smtp"
	"os"
	"sync"
)

type EmailService struct {
	mu         sync.Mutex
	sendMailFn func(string, smtp.Auth, string, []string, []byte) error
}

func NewEmailService() *EmailService {
	return &EmailService{
		sendMailFn: smtp.SendMail,
	}
}

func (s *EmailService) SendEmail(ctx context.Context, templatePath string, data interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	from := os.Getenv("EMAIL_FROM")
	to := os.Getenv("EMAIL_TO")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if from == "" || to == "" || password == "" || smtpHost == "" || smtpPort == "" {
		return errors.New("missing required environment variables for email configuration")
	}

	auth := smtp.PlainAuth("", from, password, smtpHost)
	msg, err := buildEmailMessage(templatePath, data)
	if err != nil {
		return err
	}

	// Guardar el mensaje como un archivo HTML
	err = os.WriteFile("./internal/infrastructure/email/output_email/email_output.html", msg, 0644)
	if err != nil {
		return err
	}

	if s.sendMailFn == nil {
		return errors.New("sendMailFn is not initialized")
	}

	err = s.sendMailFn(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return err
	}

	return nil
}

func buildEmailMessage(templatePath string, data interface{}) ([]byte, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return nil, err
	}

	return body.Bytes(), nil
}
