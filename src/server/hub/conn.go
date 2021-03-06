package hub

import (
	"github.com/garyburd/go-websocket/websocket"
	"io/ioutil"
	"time"
	"fmt"
)

const (
	// Time allowed to write a message to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next message from the client.
	readWait = 60 * time.Second

	// Send pings to client with this period. Must be less than readWait.
	pingPeriod = (readWait * 9) / 10

	// Maximum message size allowed from client.
	maxMessageSize = 512
)

// connection is an middleman between the websocket connection and the hub.
type Connection struct {
	// The websocket connection.
	Ws *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	//hub
	H Hub
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Connection) ReadPump() {
	defer func() {
		c.H.Unregister <- c
		c.Ws.Close()
	}()
	c.Ws.SetReadLimit(maxMessageSize)
	c.Ws.SetReadDeadline(time.Now().Add(readWait))
	for {
		op, r, err := c.Ws.NextReader()
		if err != nil {
			break
		}
		switch op {
		case websocket.OpPong:
			c.Ws.SetReadDeadline(time.Now().Add(readWait))
		case websocket.OpText:
			message, err := ioutil.ReadAll(r)
			if err != nil {
				fmt.Println("ERROR!")
				break
			}
			c.H.Broadcast <- message
		}
	}
}

// write writes a message with the given opCode and payload.
func (c *Connection) Write(opCode int, payload []byte) error {
	c.Ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Ws.WriteMessage(opCode, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Connection) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Write(websocket.OpClose, []byte{})
				return
			}
			if err := c.Write(websocket.OpText, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.Write(websocket.OpPing, []byte{}); err != nil {
				return
			}
		}
	}
}

// DocumentConnection is an middleman between the websocket connection and the DocumentHub.
type DocumentConnection struct {
	// The websocket connection.
	Ws *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan Message

	//hub
	H DocumentHub
}

type Message struct {
	Conn *DocumentConnection
	M []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (c *DocumentConnection) ReadPump() {
	defer func() {
		c.H.Unregister <- c
		c.Ws.Close()
	}()
	c.Ws.SetReadLimit(maxMessageSize)
	c.Ws.SetReadDeadline(time.Now().Add(readWait))
	for {
		op, r, err := c.Ws.NextReader()
		if err != nil {
			break
		}
		switch op {
		case websocket.OpPong:
			c.Ws.SetReadDeadline(time.Now().Add(readWait))
		case websocket.OpText:
			message, err := ioutil.ReadAll(r)
			if err != nil {
				fmt.Println("ERROR!")
				break
			}
			mess := Message{M:message, Conn:c}
			c.H.Broadcast <- mess
		}
	}
}

// write writes a message with the given opCode and payload.
func (c *DocumentConnection) Write(opCode int, payload []byte) error {
	c.Ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Ws.WriteMessage(opCode, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *DocumentConnection) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Write(websocket.OpClose, []byte{})
				return
			}
			if err := c.Write(websocket.OpText, message.M); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.Write(websocket.OpPing, []byte{}); err != nil {
				return
			}
		}
	}
}

// connection is an middleman between the websocket connection and the hub.
type ChatConnection struct {
	// The websocket connection.
	Ws *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	//hub
	H ChatHub
}

// readPump pumps messages from the websocket connection to the hub.
func (c *ChatConnection) ReadPump() {
	defer func() {
		c.H.Unregister <- c
		c.Ws.Close()
	}()
	c.Ws.SetReadLimit(maxMessageSize)
	c.Ws.SetReadDeadline(time.Now().Add(readWait))
	for {
		op, r, err := c.Ws.NextReader()
		if err != nil {
			break
		}
		switch op {
		case websocket.OpPong:
			c.Ws.SetReadDeadline(time.Now().Add(readWait))
		case websocket.OpText:
			message, err := ioutil.ReadAll(r)
			if err != nil {
				fmt.Println("ERROR!")
				break
			}
			c.H.Broadcast <- message
		}
	}
}

// write writes a message with the given opCode and payload.
func (c *ChatConnection) Write(opCode int, payload []byte) error {
	c.Ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Ws.WriteMessage(opCode, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *ChatConnection) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Write(websocket.OpClose, []byte{})
				return
			}
			if err := c.Write(websocket.OpText, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.Write(websocket.OpPing, []byte{}); err != nil {
				return
			}
		}
	}
}