package conn

import (
	"errors"
	gConn "github.com/viamAhmadi/gReceiver2/pkg/conn"
	"github.com/zeromq/goczmq"
)

const YES = 1
const NO = 0

var ErrConvertToModel = errors.New("convert error")
var ErrConnExists = errors.New("connection exists")
var ErrDealer = errors.New("dealer was nil")

type SendConn struct {
	Sender      *Sender
	Destination string
	Count       int
	Id          string // 20 char
	IsOpen      int    // open 1, closed 0
	Successful  int    //  1 - 0
	Factor      *gConn.Factor
	FactorCh    chan *gConn.Factor
	SendCh      chan gConn.Send
	ErrorCh     chan Error
}

// id = from rec conns get id from the connection
// for send  generate random connection id
func NewSendConn(des, id string, count int) *SendConn {
	return &SendConn{
		Destination: des,
		Count:       count,
		Id:          id,
		Factor:      nil,
		FactorCh:    make(chan *gConn.Factor),
		SendCh:      make(chan gConn.Send),
		ErrorCh:     make(chan Error),
		Successful:  NO,
	}
}

func (c *SendConn) OpenConnection(addr string) error {
	dealer, err := goczmq.NewDealer(addr)
	if err != nil {
		return err
	}
	c.Sender = NewSender(dealer)
	c.IsOpen = YES
	return dealer.SendFrame(SerializeSendConn(c.Destination, c.Count, c.Id), goczmq.FlagNone)
}

func (c *SendConn) Send(msg *gConn.Message) error {
	if c.Sender.sock == nil {
		return ErrDealer
	}
	return c.Sender.sock.SendFrame(*gConn.SerializeMessage(msg), goczmq.FlagNone)
}

func (c *SendConn) Close() {
	c.IsOpen = NO
	c.Sender.Close()
	close(c.FactorCh)
	close(c.SendCh)
	close(c.ErrorCh)
}

func (c *SendConn) Receiving() error {
	for {
		if c.Sender.sock == nil {
			return nil
		}
		msg, err := c.Sender.sock.RecvMessage()
		if err != nil {
			return err
		}

		valStr := string((msg)[0][0])

		if valStr == "s" {
			s, err := gConn.ConvertToSend(msg[0])
			if err != nil {
				return err
			}
			c.SendCh <- s
		} else if valStr == "f" {
			f, err := gConn.ConvertToFactor(&(msg)[0])
			if err != nil {
				return err
			}
			c.FactorCh <- f
		}
	}
}
