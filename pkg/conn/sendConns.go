package conn

import (
	"sync"
)

type SendConns map[string]*SendConn

var mutex = sync.Mutex{}

func (r *SendConns) Add(conn *SendConn) error {
	if cFound := r.Get(conn.Id); cFound != nil {
		return ErrConnExists
	}
	mutex.Lock()
	(*r)[conn.Id] = conn
	mutex.Unlock()
	return nil
}

func (r *SendConns) Get(connId string) *SendConn {
	mutex.Lock()
	conn := (*r)[connId]
	mutex.Unlock()
	return conn
}
