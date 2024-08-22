package notifications

import (
	"bytes"
	"net/http"
)

//var ErrNotSend = errors.New("não foi posssível enviar o email")

type NtfyMobilePushNotifications struct {
	host     string
	port     int
	username string
	password string
	from     string
	to       string
	subject  string
	message  string
}

func NewNtfyMobilePushNotifications() *NtfyMobilePushNotifications {
	return &NtfyMobilePushNotifications{}
}

func (m *NtfyMobilePushNotifications) Config(host string, port int, username, password string) {
	m.host = host
	m.port = port
	m.username = username
	m.password = password
}

func (m *NtfyMobilePushNotifications) SetFrom(from string) {
	m.from = from
}

func (m *NtfyMobilePushNotifications) SetTo(to string) {
	m.to = to
}

func (m *NtfyMobilePushNotifications) SetSubject(subject string) {
	m.subject = subject
}

func (m *NtfyMobilePushNotifications) SetMessage(message string) {
	m.message = message
}

func (m *NtfyMobilePushNotifications) Send() error {

	body := []byte(`{"title":"Contato NotaBaiana", "message":"Assunto: ` + m.subject + `"}`)
	_, err := http.Post("https://ntfy.sh/notabaiana", "application/json", bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return nil
}
