package smtp

import (
	"bytes"
	"fmt"
	"net/smtp"

	"github.com/grez-lucas/boxer66-service/internal/config"
)

type SMTPService struct {
	cfg config.SMTPConfig
}

func NewSMTPService(cfg config.SMTPConfig) *SMTPService {
	return &SMTPService{
		cfg: cfg,
	}
}

func (s *SMTPService) SendVerificationEmail(to, verificationCode string) error {
	subject := "Boxer66 Verification Code"
	body := fmt.Sprintf(`
		<html>
		<head>
			<title>%s</title>
		</head>
		<body>
			<p> Hi there,</p>
			<p>Please use the verification code below:</p>
			<h3>%s</h3>
			<p>This code will expire in 60 minutes.</p>
			<p>Thanks,</p>
			<Boxer66 Team</p>
		</body>
		</html>
		`, subject, verificationCode)

	if err := s.SendEmail(to, subject, body); err != nil {
		return fmt.Errorf("failed to send email to recipient %s: %w", to, err)
	}
	return nil
}

func (s *SMTPService) SendEmail(to, subject, body string) error {
	var msg bytes.Buffer

	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("From: %s\r\n", s.cfg.User))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	auth := smtp.PlainAuth("", s.cfg.User, s.cfg.Password, s.cfg.Host)
	return smtp.SendMail(s.cfg.Host+":"+s.cfg.Port, auth, s.cfg.User, []string{to}, msg.Bytes())
}
