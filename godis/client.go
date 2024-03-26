package godis

import (
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	addr string
	path string
	conn *websocket.Conn
}

func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
		path: "/cmd",
	}
}

func (c *Client) Dial() error {
	u := url.URL{Scheme: "ws", Host: c.addr, Path: c.path}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) Hangup() error {
	return c.conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func (c *Client) SendCommand(op string, args ...string) (*Reply, error) {
	start := time.Now()
	cmd := MakeCommand(op, args...)
	err := c.conn.WriteJSON(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("-> sent: %v", cmd)

	var rep Reply
	err = c.conn.ReadJSON(&rep)
	if err != nil {
		return nil, err
	}
	log.Printf("<- recv: %v (%v)", rep, time.Since(start))
	return &rep, nil
}
