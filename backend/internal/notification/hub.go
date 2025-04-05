package notification

import (
	"log"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	conn *websocket.Conn

	send chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("New client registered!")
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client unregistered!")
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func NewWs(hub *Hub, c *websocket.Conn) {
	client := &Client{
		conn: c,
		send: make(chan []byte, 256),
	}

	hub.register <- client
	go client.writePump(hub)

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}

	hub.unregister <- client
}

func (c *Client) writePump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		message, ok := <-c.send
		if !ok {
			// Connection closed
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Failed to write message - %s", err.Error())
			return
		}
	}
}
