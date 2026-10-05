package client

import (
	"log"
	"time"

	"pouic/internal/websocket/command"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space = []byte{' '}
)

type Hub interface {
	AddClient(c Client)
	RemoveClient(id string)
	AddCommand(command command.Command)
}

type Parser interface {
	Parse(b []byte, senderId string) (command.Command, error)
}

type Client struct {
	hub Hub
	id string
	conn *websocket.Conn
	parser Parser
	toSend chan []byte
}

func NewClient(hub Hub, parser Parser, id string, conn *websocket.Conn) *Client {
	c := Client{
		hub: hub,
		id: id,
		conn: conn,
		parser: parser,
		toSend: make(chan []byte, 256),
	}
	c.hub.AddClient(c)
	return &c
}

// Get the id of the player
func (c Client) GetId() string {
	return c.id
}

// Start goroutine for write and read client channels
func (c *Client) Go() {
	go c.readPump()
	go c.writePump()
}

// Close the toSend channel
func (c *Client) Close() {
	close(c.toSend)
}

// Read from the client fd and pass the binaries to
// the hub's command worlflow.
func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
		c.hub.RemoveClient(c.id)
	}()
	
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("socket connection error: %v", err)
			}
			break
		}

		command, err := c.parser.Parse(message, c.id)
		if err != nil {
			log.Printf("command parsing error: %v. Received: %s", err, message)
			// TODO jsonified emit for error
			// c.toSend <- ErrorEmit(err)
		} else {
			c.hub.AddCommand(command)
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.toSend:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Poor all the queued messages to send.
			// WHATIF: toSend queue is closed during the pooring ?
			n := len(c.toSend)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.toSend)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}