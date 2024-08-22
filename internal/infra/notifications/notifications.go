package notifications

import "github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"

//var ErrNotSend = errors.New("não foi posssível enviar o email")

type Notifications struct {
	Services []port.EmailPort
}

func NewNotifications(services []port.EmailPort) *Notifications {
	return &Notifications{Services: services}
}

func (m *Notifications) Config(host string, port int, username, password string) {

	for i := range m.Services {
		m.Services[i].Config(host, port, username, password)
	}
}

func (m *Notifications) SetFrom(from string) {
	for i := range m.Services {
		m.Services[i].SetFrom(from)
	}
}

func (m *Notifications) SetTo(to string) {
	for i := range m.Services {
		m.Services[i].SetTo(to)
	}
}

func (m *Notifications) SetSubject(subject string) {
	for i := range m.Services {
		m.Services[i].SetSubject(subject)
	}
}

func (m *Notifications) SetMessage(message string) {
	for i := range m.Services {
		m.Services[i].SetMessage(message)
	}

}

func (m *Notifications) Send() error {

	for i := range m.Services {

		go func(i int) error {
			err := m.Services[i].Send()
			if err != nil {
				return err
			}

			return nil
		}(i)

	}

	return nil
}
