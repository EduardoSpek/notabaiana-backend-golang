package notifications

import (
	"net/smtp"
	"os"
)

//var ErrNotSend = errors.New("não foi posssível enviar o email")

type GmailSMTP struct {
	host     string
	port     int
	username string
	password string
	from     string
	to       string
	subject  string
	message  string
}

func NewGmailSMTP() *GmailSMTP {
	return &GmailSMTP{}
}

func (m *GmailSMTP) Config(host string, port int, username, password string) {
	m.host = host
	m.port = port
	m.username = username
	m.password = password
}

func (m *GmailSMTP) SetFrom(from string) {
	m.from = from
}

func (m *GmailSMTP) SetTo(to string) {
	m.to = to
}

func (m *GmailSMTP) SetSubject(subject string) {
	m.subject = subject
}

func (m *GmailSMTP) SetMessage(message string) {
	m.message = message
}

func (m *GmailSMTP) Send() error {

	auth := smtp.PlainAuth(
		"",
		os.Getenv("TO_EMAILSERVER"),
		os.Getenv("PASSWORD_GMAIL_APP"),
		"smtp.gmail.com",
	)

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		os.Getenv("TO_EMAILSERVER"),
		[]string{os.Getenv("TO_EMAILSERVER")},
		[]byte("Subject: "+m.subject+"\n"+headers+"\n\n"+m.message),
	)

	return nil
}
