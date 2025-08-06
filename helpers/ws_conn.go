// Package helpers // ws_conn.go
// implement sync access to websocket connection
package helpers

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type SyncConn struct {
	conn *websocket.Conn
	m    sync.Mutex
}

func NewSyncConnection(conn *websocket.Conn) *SyncConn {
	return &SyncConn{
		conn: conn,
		m:    sync.Mutex{},
	}
}

func (c *SyncConn) WriteMessage(messageType int, data []byte) error {
	c.m.Lock()
	defer c.m.Unlock()
	return c.conn.WriteMessage(messageType, data)
}

func (c *SyncConn) WriteJSON(v interface{}) error {
	c.m.Lock()
	defer c.m.Unlock()
	return c.conn.WriteJSON(v)
}

func (c *SyncConn) Close() {
	c.m.Lock()
	defer c.m.Unlock()
	_ = c.conn.Close()
}

func (c *SyncConn) SetReadDeadline(t time.Time) error {
	c.m.Lock()
	defer c.m.Unlock()
	return c.conn.SetReadDeadline(t)
}

func (c *SyncConn) SetPongHandler(h func(appData string) error) {
	c.m.Lock()
	defer c.m.Unlock()
	c.conn.SetPongHandler(h)
}
