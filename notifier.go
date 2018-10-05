package glutton

type NilNotifier struct{}

type SMTPNotifier struct {
	Server   string
	Port     string
	UseTLS   bool
	User     string
	Password string
}

func (n *NilNotifier) Notify(*PayloadRecord) error {
	return nil
}

func (n *NilNotifier) Configure(*Settings) error {
	return nil
}

func (s *SMTPNotifier) Notify(*PayloadRecord) error {
	return nil
}

func (s *SMTPNotifier) Configure(settings *Settings) error {
	s.Server = settings.SMTPServer
	s.Port = settings.SMTPPort
	s.UseTLS = settings.SMTPUseTLS
	s.User = settings.SMTPUser
	s.Password = settings.SMTPPassword
	return nil
}
