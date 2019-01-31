package main

import (
	"github.com/gorilla/websocket"
)

// Client represents a single chatting user

type client struct {
	// socket is a web socket for this client
	socket *websocket.Conn
	// send is the channel that messages are sent
	send chan []byte
	// room is the client's room
	room *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
