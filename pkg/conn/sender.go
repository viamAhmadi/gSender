package conn

import (
	"github.com/zeromq/goczmq"
)

type Sender struct {
	sock *goczmq.Sock
}

func NewSender(sock *goczmq.Sock) *Sender {
	b := Sender{sock: sock}
	return &b
}

func (s *Sender) Close() {
	if s.sock != nil {
		s.sock.Destroy()
		s.sock = nil
	}
}
