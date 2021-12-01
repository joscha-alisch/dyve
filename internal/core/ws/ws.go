package ws

import (
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Handler interface {
}

func New(w http.ResponseWriter, r *http.Request, disableOriginCheck bool) (*Connection, error) {
	upgrader := websocket.Upgrader{}

	if disableOriginCheck {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return &Connection{
		w:        w,
		r:        r,
		conn:     conn,
		handlers: make(map[string]func()),
	}, nil
}

type Connection struct {
	w        http.ResponseWriter
	r        *http.Request
	conn     *websocket.Conn
	errCount int
	handlers map[string]func()
}

type OnCommand func(cmd string, c *Connection)

func (c *Connection) On(cmd string, cb func()) {
	c.handlers[cmd] = cb
}

func (c *Connection) Run() error {
	for {
		t, bytes, err := c.conn.ReadMessage()
		if websocket.IsCloseError(err) {
			return nil
		}
		if err != nil {
			c.errCount++
			if c.errCount > 10 {
				log.Error().Err(err).Msg("too many errors when listening on websocket, closing connection now")
				break
			}
			continue
		}
		c.errCount = 0

		if t == websocket.CloseMessage {
			break
		}
		if t != websocket.TextMessage {
			continue
		}

		if handler := c.handlers[string(bytes)]; handler != nil {
			handler()
		}
	}
	return c.conn.Close()
}

func (c *Connection) Send(msg interface{}) error {
	return c.conn.WriteJSON(msg)
}

func (c *Connection) Close() error {
	return c.conn.Close()
}
