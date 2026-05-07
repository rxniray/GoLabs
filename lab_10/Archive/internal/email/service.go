package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/smtp"

	"lab10/internal/config"
)

//go:embed templates/*.html
var templateFS embed.FS

type Service struct {
	cfg  *config.Config
	tmpl *template.Template
}

func NewService(cfg *config.Config) (*Service, error) {
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parse email templates: %w", err)
	}
	return &Service{cfg: cfg, tmpl: tmpl}, nil
}

func (s *Service) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)
	msg := buildMessage(s.cfg.SMTPFrom, to, subject, body)

	var auth smtp.Auth
	if s.cfg.SMTPUser != "" {
		auth = smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPHost)
	}

	if err := smtp.SendMail(addr, auth, s.cfg.SMTPFrom, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("smtp send: %w", err)
	}
	return nil
}

func (s *Service) SendTemplate(to, templateName string, data any) error {
	var buf bytes.Buffer
	if err := s.tmpl.ExecuteTemplate(&buf, templateName, data); err != nil {
		return fmt.Errorf("execute template %q: %w", templateName, err)
	}
	return s.Send(to, subjectFor(templateName), buf.String())
}

func buildMessage(from, to, subject, body string) string {
	return fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		from, to, subject, body,
	)
}

func subjectFor(templateName string) string {
	subjects := map[string]string{
		"welcome.html":        "Welcome!",
		"password_reset.html": "Reset your password",
	}
	if s, ok := subjects[templateName]; ok {
		return s
	}
	return "Notification"
}
