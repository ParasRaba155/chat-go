package chat

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type room struct {
	// room-id
	id uuid.UUID
	// clients the map to keep track of connected clients
	clients map[*client]bool

	broadcast chan []byte

	register chan *client

	unregister chan *client
}

func NewRoom(id uuid.UUID) *room {
	return &room{
		id:         id,
		clients:    make(map[*client]bool, 2),
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
}

func (r *room) Run(l *zap.Logger) {
	for {
		select {
		case client := <-r.register:
			l.Info("registered the client", zap.String("client", client.senderMail))
			r.clients[client] = true
		case client := <-r.unregister:
			if _, ok := r.clients[client]; !ok {
				delete(r.clients, client)
				close(client.send)
			}
		case msg := <-r.broadcast:
			l.Debug("let us send some message")
			for client := range r.clients {
				select {
				case client.send <- msg:
					l.Info("msg sent", zap.ByteString("message", msg))
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		}
	}
}
