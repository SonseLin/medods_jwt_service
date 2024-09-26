package model

type SmtpServer struct {
	Host string
	Port string
}

func (s *SmtpServer) Address() string {
	return s.Host + ":" + s.Port
}
