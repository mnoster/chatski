package main

import (
	"github.com/gorilla/websocket"
	"github.com/mnoster/chat/trace"
	"log"
	"net/http"
)

type room struct {
	// Holds incoming messages that should be forwarded
	forward chan []byte
	// Channel for clients that want to join a room
	join chan *client
	// Channel for clients that want to leave the room
	leave chan *client
	// Holds all current clients in this room
	clients map[*client]bool
	// Tracer will receive trace information of activity in the room
	tracer trace.Tracer
}

// newRoom makes a new room ready to go.
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

func (r *room) run() {
	// This top level for loop indicated that  this method with run forever
	// until the program is terminated. . Since this code is run as a Go routine,
	// it will run in the backgroundm whic won't block the rest of the application.
	// If a msg is received on any of the join,leave,forward channels, then the
	// select statement will be triggered.
	for {
		// Select helps synchronize shared memory
		select {
		// If we received a msg on the join channel then update the r.clients map
		// to keep reference of the client that has joined the room.
		// Setting the value to true is just a handy low-memory way to store.
		case client := <-r.join:
			// Joining
			r.clients[client] = true
			r.tracer.Trace("New client joined")
		// If we receive a mssg on the leave channel then simply delete the client type
		// from the map, and close its send channel
		case client := <-r.leave:
			// Leaving
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left")
		// If we receive a msg on the forward channel then we iterate over all the clients
		// and send the mesg down each  client's send channel. Then the write method
		// of our client type will pick it up and send it down the socket to the browser.
		case msg := <-r.forward:
			// Forward msg to all clients
			for client := range r.clients {
				select {
				case client.send <- msg:
					r.tracer.Trace(" -- sent to client")
				// If the send channel is closed then we know the client is not receiving any
				//	more msgs, so remove client from the channel.
				default:
					// Failed to send
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

// In order to use websockets we must upgrade the HTTP connection using the websocker.Upgrader type
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

// ServeHTTP can now serve as a handler for room
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
