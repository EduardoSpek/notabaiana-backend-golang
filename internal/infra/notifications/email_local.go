package notifications

import (
	"gopkg.in/gomail.v2"
)

//var ErrNotSend = errors.New("não foi posssível enviar o email")

type LocalEmailServer struct {
	host     string
	port     int
	username string
	password string
	from     string
	to       string
	subject  string
	message  string
}

func NewLocalEmailServer() *LocalEmailServer {
	return &LocalEmailServer{}
}

func (m *LocalEmailServer) Config(host string, port int, username, password string) {
	m.host = host
	m.port = port
	m.username = username
	m.password = password
}

func (m *LocalEmailServer) SetFrom(from string) {
	m.from = from
}

func (m *LocalEmailServer) SetTo(to string) {
	m.to = to
}

func (m *LocalEmailServer) SetSubject(subject string) {
	m.subject = subject
}

func (m *LocalEmailServer) SetMessage(message string) {
	m.message = message
}

func (m *LocalEmailServer) Send() error {

	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", m.to)
	msg.SetHeader("Subject", m.subject)
	msg.SetBody("text/html", m.message)

	dialer := gomail.NewDialer("localhost", 1025, m.username, m.password)

	if err := dialer.DialAndSend(msg); err != nil {
		return ErrNotSend

	}

	return nil
}
