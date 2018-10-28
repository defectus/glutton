package notifier

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"

	"github.com/pkg/errors"
	"github.com/defectus/glutton/pkg/iface"
)

// NilNotifier implements the glutton.PayloadNotifier interface but does nothing.
type NilNotifier struct{}

// SMTPNotifier implements the glutton.PayloadNotifier interface and delivers notifications over SMTP.
type SMTPNotifier struct {
	Server   string
	Port     string
	UseTLS   bool
	From     string
	Password string
	To       string
	Subject  string
}

// Notify does nothing.
func (n *NilNotifier) Notify(*iface.PayloadRecord) error {
	return nil
}

// Configure does nothing.
func (n *NilNotifier) Configure(*iface.Settings) error {
	return nil
}

// PayloadToSMTPMessage takes payload and format it to SMTP format.
func (s *SMTPNotifier) PayloadToSMTPMessage(payload *iface.PayloadRecord) []byte {
	smtpTemplateData := &struct {
		From    string
		To      string
		Subject string
		Body    string
	}{s.From, s.To, s.Subject, payload.String()}
	const emailTemplate = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

{{.Body}}

Sincerely,

{{.From}}
`
	var err error
	var doc bytes.Buffer

	t := template.New("emailTemplate")
	t, err = t.Parse(emailTemplate)
	if err != nil {
		log.Printf("error trying to parse mail template %+v", err)
	}
	err = t.Execute(&doc, smtpTemplateData)
	if err != nil {
		log.Printf("error trying to execute mail template %+v", err)
	}
	return doc.Bytes()
}

// Notify sends notification over SMTP.
func (s *SMTPNotifier) Notify(payload *iface.PayloadRecord) error {
	err := smtp.SendMail(s.Server+":"+s.Port,
		smtp.PlainAuth("", s.From, s.Password, s.Server),
		s.From, []string{s.To}, s.PayloadToSMTPMessage(payload))
	return errors.Wrapf(err, "error sending notification %+v", payload)
}

// Configure configures this notifier according to settings.
func (s *SMTPNotifier) Configure(settings *iface.Settings) error {
	s.Server = settings.SMTPServer
	s.Port = settings.SMTPPort
	s.UseTLS = settings.SMTPUseTLS
	s.From = settings.SMTPFrom
	s.Password = settings.SMTPPassword
	s.To = settings.SMTPTo
	s.Subject = "Notification from " + settings.Name
	return nil
}
