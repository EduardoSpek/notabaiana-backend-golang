package notifications

import (
	"errors"

	"gopkg.in/gomail.v2"
)

var ErrNotSend = errors.New("não foi posssível enviar o email")

type MailtrapEmailServer struct {
	host     string
	port     int
	username string
	password string
	from     string
	to       string
	subject  string
	message  string
}

func NewMailtrapEmailServer() *MailtrapEmailServer {
	return &MailtrapEmailServer{}
}

func (m *MailtrapEmailServer) Config(host string, port int, username, password string) {
	m.host = host
	m.port = port
	m.username = username
	m.password = password
}

func (m *MailtrapEmailServer) SetFrom(from string) {
	m.from = from
}

func (m *MailtrapEmailServer) SetTo(to string) {
	m.to = to
}

func (m *MailtrapEmailServer) SetSubject(subject string) {
	m.subject = subject
}

func (m *MailtrapEmailServer) SetMessage(message string) {
	m.message = message
}

func (m *MailtrapEmailServer) Send() error {

	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", m.to)
	msg.SetHeader("Subject", m.subject)
	msg.SetBody("text/html", m.message)

	dialer := gomail.NewDialer(m.host, m.port, m.username, m.password)

	if err := dialer.DialAndSend(msg); err != nil {
		return ErrNotSend

	}

	return nil
}
