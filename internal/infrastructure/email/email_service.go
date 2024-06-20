package email

import (
	"context"
	"fmt"
	"log"
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

func (s *EmailService) SendEmail(ctx context.Context, summary string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if os.Getenv("FAKE_EMAIL") == "true" {
		printFakeEmail(summary)
		return nil
	}

	from := os.Getenv("EMAIL_FROM")
	to := os.Getenv("EMAIL_TO")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	msg := buildEmailMessage(from, to, summary)

	err := s.sendMailFn(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func printFakeEmail(summary string) {
	log.Println("============================================================")
	log.Println("                   FAKE EMAIL ENABLED                       ")
	log.Println("============================================================")
	log.Println("To: recipient@example.com")
	log.Println("Subject: Monthly Transaction Summary")
	log.Println("------------------------------------------------------------")
	log.Println(summary)
	log.Println("============================================================")
}

func buildEmailMessage(from, to, summary string) string {
	return fmt.Sprintf("To: %s\r\nSubject: Monthly Transaction Summary\r\n\r\n%s\r\n", to, summary)
}
