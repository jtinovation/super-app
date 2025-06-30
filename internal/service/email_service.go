package service

import (
	"jti-super-app-go/config"
	"jti-super-app-go/pkg/helper"

	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendEmail(to, subject, body string) error
	SendEmailWithTemplate(to, subject string, templateFileName string, data interface{}) error
}

type emailService struct {
	dialer *gomail.Dialer
	sender string
}

func NewEmailService(cfg config.EmailConfig) EmailService {
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Password)
	return &emailService{
		dialer: dialer,
		sender: cfg.SenderEmail,
	}
}

func (s *emailService) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.sender)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if err := s.dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func (s *emailService) SendEmailWithTemplate(to, subject string, templateFileName string, data interface{}) error {
	body, err := helper.ParseTemplate(templateFileName, data)
	if err != nil {
		return err
	}

	return s.SendEmail(to, subject, body)
}
