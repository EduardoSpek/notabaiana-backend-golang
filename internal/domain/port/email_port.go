package port

type EmailPort interface {
	Config(host string, port int, username, password string)
	SetFrom(from string)
	SetTo(to string)
	SetSubject(subject string)
	SetMessage(message string)
	Send() error
}
