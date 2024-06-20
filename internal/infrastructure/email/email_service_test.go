package email

import (
	"context"
	"net/smtp"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSMTPClient is a mock of the smtp client
type MockSMTPClient struct {
	mock.Mock
}

// Mock the SendMail function
func (m *MockSMTPClient) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	args := m.Called(addr, a, from, to, msg)
	return args.Error(0)
}

func TestSendEmail_FakeEmail(t *testing.T) {
	os.Setenv("FAKE_EMAIL", "true")
	os.Setenv("EMAIL_FROM", "test-email@example.com")
	os.Setenv("EMAIL_TO", "test-recipient@example.com")
	os.Setenv("EMAIL_PASSWORD", "test-email-password")
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_PORT", "587")

	service := NewEmailService()
	summary := "This is a test summary"

	err := service.SendEmail(context.Background(), summary)
	assert.NoError(t, err)
}

func TestSendEmail_RealEmail(t *testing.T) {
	os.Setenv("FAKE_EMAIL", "false")
	os.Setenv("EMAIL_FROM", "test-email@example.com")
	os.Setenv("EMAIL_TO", "test-recipient@example.com")
	os.Setenv("EMAIL_PASSWORD", "test-email-password")
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_PORT", "587")

	service := NewEmailService()
	summary := "This is a test summary"

	mockClient := new(MockSMTPClient)
	mockClient.On(
		"SendMail",
		"smtp.example.com:587",
		smtp.PlainAuth("", "test-email@example.com", "test-email-password", "smtp.example.com"),
		"test-email@example.com",
		[]string{"test-recipient@example.com"},
		[]byte(buildEmailMessage("test-email@example.com", "test-recipient@example.com", summary)),
	).Return(nil)

	// Inject the mock sendMail function
	service.sendMailFn = mockClient.SendMail

	err := service.SendEmail(context.Background(), summary)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestPrintFakeEmail(t *testing.T) {
	summary := "This is a test summary"
	printFakeEmail(summary)
}

func TestBuildEmailMessage(t *testing.T) {
	from := "test-email@example.com"
	to := "test-recipient@example.com"
	summary := "This is a test summary"
	expected := "To: test-recipient@example.com\r\nSubject: Monthly Transaction Summary\r\n\r\nThis is a test summary\r\n"
	msg := buildEmailMessage(from, to, summary)
	assert.Equal(t, expected, msg)
}
